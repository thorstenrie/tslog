// Copyright (c) 2023 thorstenrie
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tslog

// Import standard library packages, tserr and tsfio.
import (
	"fmt" // fmt
	"log" // log
	"os"  // os

	"github.com/thorstenrie/tserr" // tserr
	"github.com/thorstenrie/tsfio" // tsfio
)

// Logger contains a log.logger for logging and the minimum level for logging.
// The minimum level for logging is set with SetLevel.
type Logger struct {
	minLvl int         // minimum level for logging
	logger *log.Logger // for logging
}

// New creates a new logger with default minimum level Info for logging. To alter
// the minimum level for logging use SetLevel. Logging is set to Stdout. To
// change logging output use SetOutput.
func New() *Logger {
	return &Logger{minLvl: defaultMinLvl, logger: log.New(os.Stdout, "", 0)}
}

// SetLevel sets the logging level. All levels equal or higher than the set level
// are logged. All log messages with levels below the set level are discarded.
// SetLevel returns an error for undefined levels, otherwise nil. If the provided
// level is lower than Trace level, the lowest level, the minimum level is set
// to Trace. If the provided level is higher than Fatal level, the highest level,
// the minimum level is set to Fatal.
func (l *Logger) SetLevel(level int) error {
	// Initially set error e to nil
	var e error = nil
	// If level is lower than Trace, the lowest level, return an error and
	// set the minimum level to Trace.
	if level < TraceLevel {
		// Set error to not existent
		e = tserr.NotExistent(fmt.Sprintf("log level %d", level))
		// Set minimum level to Trace level
		l.minLvl = TraceLevel
		// If level is higher than Fatal, the highest level, return an error and
		// set the minimum level to Fatal.
	} else if level > FatalLevel {
		// Set error to not existent
		e = tserr.NotExistent(fmt.Sprintf("log level %d", level))
		// Set minimum level to Fatal level
		l.minLvl = FatalLevel
	} else {
		// Set minimum level to provided level
		l.minLvl = level
	}
	// Return e
	return e
}

// SetOutput sets the logging output to fn. Special loggers are
// 'stdout' for logging to Stdout (default)
// 'discard' for no logging
// 'tmp' for logging to tslog_* in the temporary directory
// If SetOuput returns an error, logging is set to Stdout
func (l *Logger) SetOutput(fn tsfio.Filename) error {
	// Handle special loggers
	switch fn {
	case DiscardLogger:
		// discard, no logging
		l.noLogger()
		// Return nil
		return nil
	case StdoutLogger:
		// Logging to Stdout
		l.setStdout()
		// Return nil
		return nil
	case TmpLogger:
		// Define pattern for the temporary file
		p := fmt.Sprintf("%v_*", defaultPattern)
		// Create temporary file for logging
		f, err := os.CreateTemp(os.TempDir(), p)
		// If it fails, return an error
		if err != nil {
			// Set logging output to Stdout
			l.setStdout()
			// Return error
			return tserr.Op(&tserr.OpArgs{Op: "create temp file", Fn: p, Err: err})
		}
		// Activate logging to f
		l.logger.SetOutput(f)
		// Return nil
		return nil
	}

	// Check filename using tsfio.CheckFile
	if err := tsfio.CheckFile(fn); err != nil {
		// If the check fails, set logging output to Stdout and return an error
		l.setStdout()
		// Return error
		return tserr.Check(&tserr.CheckArgs{F: string(fn), Err: err})
	}

	// Open file with filename fn
	f, e := tsfio.OpenFile(fn)
	// If OpenFile fails, set logging output to Stdout and return an error
	if e != nil {
		// Set logging output to Stdout
		l.setStdout()
		// Return error
		return tserr.Op(&tserr.OpArgs{Op: "open file", Fn: string(fn), Err: e})
	}

	// Set Ouptut to file f
	l.logger.SetOutput(f)

	// Return nil
	return nil
}

// Trace logs a message at Trace level. It returns an error if JSON encoding of msg fails.
func (l *Logger) Trace(msg string) error {
	return l.tryLog(TraceLevel, msg)
}

// Debug logs a message at Debug level. It returns an error if JSON encoding of msg fails.
func (l *Logger) Debug(msg string) error {
	return l.tryLog(DebugLevel, msg)
}

// Info logs a message at Info level. It returns an error if JSON encoding of msg fails.
func (l *Logger) Info(msg string) error {
	return l.tryLog(InfoLevel, msg)
}

// Warn logs a message at Warn level. It returns an error if JSON encoding of msg fails.
func (l *Logger) Warn(msg string) error {
	return l.tryLog(WarnLevel, msg)
}

// Error logs error err at Error level. It returns an error if JSON encoding of msg fails.
func (l *Logger) Error(err error) error {
	return l.tryLog(ErrorLevel, err.Error())
}

// Fatal logs error err at Fatal level. It returns an error if JSON encoding of msg fails.
func (l *Logger) Fatal(err error) error {
	return l.tryLog(FatalLevel, err.Error())
}
