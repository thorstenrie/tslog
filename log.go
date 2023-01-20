// Package tslog implements logging that tries to keep it simple.
//
// The tslog package is a logging interface in Go that tries to keep it simple.
// It provides log levels Trace, Debug, Info, Warn, Error and Fatal.
// The log messages are formatted in JSON format to enable parsing.
// The predefined default logger is set to log to Stdout on Info level. A new
// logger instance can be created with New(). The output of a logger can be set
// to a specific file, a temporary file, to Stdout and to discard.
// All function calls return an error, if any.
//
// Copyright (c) 2023 thorstenrie
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tslog

// Import tsfio.
import (
	"github.com/thorstenrie/tsfio" // tsfio
)

// Strings for special loggers
const (
	StdoutLogger  tsfio.Filename = tsfio.Filename("stdout")  // Stdout
	DiscardLogger tsfio.Filename = tsfio.Filename("discard") // discard, no logging
	TmpLogger     tsfio.Filename = tsfio.Filename("tmp")     // temporary file
)

// Enum for log levels.
const (
	// Trace: log the execution of code of the app
	TraceLevel int = 1
	// Debug: log detailed events for debugging of the app
	DebugLevel int = 2
	// Info: log an event under normal conditions of the app
	InfoLevel int = 3
	// Warn: log an unintended event, which is tried to be recovered and potentially
	// impacting execution of the app
	WarnLevel int = 4
	// Error: log an unexpected event with at least one function of the app being not operable
	ErrorLevel int = 5
	// Fatal: log an unexpected critical event forcing a shutdown of the app
	FatalLevel int = 6
)

// Defaults for logging
const (
	// Layout for timestamp in the log message
	timeLayout string = "2006-01-02 15:04:05 -0700 MST"
	// Root element for JSON format
	defaultPattern string = "tslog"
	// Default log level is InfoLevel
	defaultMinLvl int = InfoLevel
)

// Global logger to provide a predefined standard logger
var (
	globalLogger *Logger = New()
)

// Default returns the global predefined standard logger
func Default() *Logger {
	return globalLogger
}

// SetLevel sets the logging level. All levels equal or higher than the set level
// are logged. All log messages with levels below the set level are discarded.
// SetLevel returns an error for undefined levels, otherwise nil.
func SetLevel(level int) error {
	return globalLogger.SetLevel(level)
}

// SetOutput sets the logging output to fn. Special loggers are
// 'stdout' for logging to Stdout (default)
// 'discard' for no logging
// 'tmp' for logging to tslog_* in the temporary directory
// If SetOuput returns an error, logging is set to Stdout
func SetOutput(fn tsfio.Filename) error {
	return globalLogger.SetOutput(fn)
}

// Trace logs a message at Trace level on the global predefined standard logger.
// It returns an error if JSON encoding of msg fails.
func Trace(msg string) error {
	return globalLogger.Trace(msg)
}

// Debug logs a message at Debug level on the global predefined standard logger.
// It returns an error if JSON encoding of msg fails.
func Debug(msg string) error {
	return globalLogger.Debug(msg)
}

// Info logs a message at Info level on the global predefined standard logger.
// It returns an error if JSON encoding of msg fails.
func Info(msg string) error {
	return globalLogger.Info(msg)
}

// Warn logs a message at Warn level on the global predefined standard logger.
// It returns an error if JSON encoding of msg fails.
func Warn(msg string) error {
	return globalLogger.Warn(msg)
}

// Error logs error err at Error level on the global predefined standard logger.
// It returns an error if JSON encoding of msg fails.
func Error(err error) error {
	return globalLogger.Error(err)
}

// Fatal logs error err at Fatal level on the global predefined standard logger.
// It returns an error if JSON encoding of msg fails.
func Fatal(err error) error {
	return globalLogger.Fatal(err)
}
