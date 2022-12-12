// Package tslog implements logging that tries to keep it simple.
//
// The tslog package provides one log.Logger for informational logging (I) and one
// log.Logger for error logging (E). Both log.Logger write into the configured
// io.Writer (a file, a tmp file, Stdout or discard). The io.Writer is configured
// using the environment variable "TS_LOGFILE" during the initial startup of the app.
//
// Set TS_LOGFILE to
// 'stdout' for logging to Stdout (default)
// 'discard' for no logging
// 'tmp' for logging to tslog_* in temporary directory
// <filename> for logging to <filename>
//
// Copyright (c) 2022 thorstenrie
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tslog

// Import standard library packages.
import (
	// fmt
	"io"  // io
	"log" // log
	"os"  // os

	"github.com/thorstenrie/tserr" // tserr
	"github.com/thorstenrie/tsfio" // tsfio
)

// Global informational logger and error logger provided.
var (
	I *log.Logger // information
	E *log.Logger // error
)

// Prefixes of the loggers set as constants.
const (
	infoPrefix  string = "INFO: "  // information prefix
	errorPrefix string = "ERROR: " // error prefix
)

// Strings for special loggers.
const (
	stdoutLogger  string = "stdout"  // Stdout
	discardLogger string = "discard" // discard, no logging
	tmpLogger     string = "tmp"     // temporary file
)

// Flags for logging properties
const (
	flags int = log.Ldate | log.Ltime | log.Lshortfile
)

// init initializes global loggers.
func init() {
	initialize()
}

// initialize sets global loggers according to env variable TS_LOGFILE.
// On error, the function falls back to Stdout.
func initialize() {
	if err := setLog(); err != nil {
		setStdout()
		I.Printf("%v; using default log stdout", err)
	}
}

// setLog interpretes env variable TS_LOGFILE and sets global loggers.
func setLog() error {

	// read env variable TS_LOGFILE
	fn, isset := os.LookupEnv("TS_LOGFILE")

	// error handling
	// return error, if not set
	if !isset {
		return tserr.NotSet("env variable $TS_LOGFILE")
	}

	// Handle special loggers
	switch fn {
	case discardLogger:
		// discard, no logging
		noLogger()
		// Return nil
		return nil
	case stdoutLogger:
		// Logging to stdout
		setStdout()
		// Return nil
		return nil
	case tmpLogger:
		// Create temporary file for logging
		f, err := os.CreateTemp(os.TempDir(), "tslog_*")
		// If it fails, return an error
		if err != nil {
			return tserr.Op(&tserr.OpArgs{Op: "create temp file", Fn: "tslog_*", Err: err})
		}
		// Activate file logging
		setLogger(f)
		// Return nil
		return nil
	}

	// Use type tsfio.Filename
	filename := tsfio.Filename(fn)

	// Check filename using tsfio.CheckFile
	if err := tsfio.CheckFile(filename); err != nil {
		// If the check fails, return an error
		return tserr.Check(&tserr.CheckArgs{F: string(filename), Err: err})
	}

	// Set file
	f, e := tsfio.OpenFile(tsfio.Filename(filename))
	// If OpenFile fails, return an error
	if e != nil {
		return tserr.Op(&tserr.OpArgs{Op: "open file", Fn: string(filename), Err: e})
	}

	// Activate file logging
	setLogger(f)

	// Return nil
	return nil
}

// setLogger initializes global loggers with f.
func setLogger(f *os.File) {
	if f == nil {
		// If f is nil, fall back to logging to Stdout
		f = os.Stdout
	}
	// Set loggers I and E to f
	I = log.New(f, infoPrefix, flags)
	E = log.New(f, errorPrefix, flags)

}

// setStdout set global loggers to Stdout.
func setStdout() {
	setLogger(os.Stdout)
}

// noLogger set global loggers to discard logging.
func noLogger() {
	I = log.New(io.Discard, "", 0)
	E = log.New(io.Discard, "", 0)
}
