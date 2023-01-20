// Copyright (c) 2023 thorstenrie
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tslog

// Import standard library packages and tserr.
import (
	"encoding/json" // json
	"fmt"           // fmt
	"io"            // io
	"os"            // os
	"time"          // time

	"github.com/thorstenrie/tserr" // tserr
)

// Struct logmsg contains the content of the log message.
// - Lvl: log level as string
// - Msg: log message as string
// - Now: timestamp as string
type logmsg struct {
	Lvl string `json:"level"`   // level
	Msg string `json:"message"` // message
	Now string `json:"time"`    // timestamp
}

// Struct logwrap is the JSON root element holding the log message.
type logwrap struct {
	L logmsg `json:"log"` // JSON root element
}

// setStdout sets logging to Stdout.
func (l *Logger) setStdout() {
	l.logger.SetOutput(os.Stdout)
}

// noLogger sets logging to discard logging.
func (l *Logger) noLogger() {
	l.logger.SetOutput(io.Discard)
}

// trylog logs message msg, if lvl is equal to or higher than the
// minimum log level. It returns an error if JSON encoding of msg fails.
func (l *Logger) tryLog(lvl int, msg string) error {
	// Log message if lvl is equal to or higher than the minimum log level
	if lvl >= l.minLvl {
		// Format log message in JSON format
		j, e := jsonFormat(lvl, msg)
		// Log JSON encoded log message using the logger
		l.logger.Println(string(j))
		// Return an error from JSON encoding, if any
		return e
	}
	// Return nil
	return nil
}

// jsonFormat encodes lvl and msg into a JSON log message. It returns the
// JSON encoded log message or an error, if any. If JSON encoding fails,
// it returns nil and an error.
func jsonFormat(lvl int, msg string) ([]byte, error) {
	// Retrieve string representation for log level lvl
	ls, errl := level(lvl)
	// Return nil and an error for invalid log levels
	if errl != nil {
		return nil, errl
	}
	// data holds the log message
	data := logmsg{Lvl: ls, Msg: msg, Now: time.Now().Format(timeLayout)}
	// wrap holds the log message and the JSON root element
	wrap := logwrap{L: data}
	// Retrieve the JSON encoding of wrap
	j, errj := json.Marshal(&wrap)
	// Return nil and an error, if JSON encoding fails
	if errj != nil {
		return nil, tserr.Op(&tserr.OpArgs{Op: "JSON Marshal", Fn: msg, Err: errj})
	}
	// Return the JSON encoded log message and nil
	return j, nil
}

// level returns the string representation of lvl. It returns "error" and an error,
// if lvl is non existent.
func level(lvl int) (string, error) {
	switch lvl {
	case TraceLevel:
		return "trace", nil
	case DebugLevel:
		return "debug", nil
	case InfoLevel:
		return "info", nil
	case WarnLevel:
		return "warn", nil
	case ErrorLevel:
		return "error", nil
	case FatalLevel:
		return "fatal", nil
	default:
		return "error", tserr.NotExistent(fmt.Sprintf("log level %d", lvl))
	}
}
