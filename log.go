// Package tslog implements easy-to-use logging that tries to keep it simple.
//
// The tslog package provides one log.Logger for informational logging (I) and one
// log.Logger for error logging (E). Both log.Logger write into the configured
// io.Writer (a file, a tmp file, Stdout or discard). The io.Writer is configured
// using the environment variable "TS_LOGFILE"
//
// Set TS_LOGFILE to
// 'stdout' for logging to Stdout (default)
// 'discard' for no logging
// 'tmp' for logging to tslog_* in temporary directory
// <filename> for logging to <filename>
package tslog

// Import standard library packages.
import (
	"fmt" // fmt
	"io"  // io
	"log" // log
	"os"  // os
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

// init calls Reset to initialize global loggers.
func init() {
	Reset()
}

// Reset sets global loggers according to env variable TS_LOGFILE.
//
// On error, the function falls back to Stdout.
func Reset() {
	if err := setLog(); err != nil {
		setStdout()
		E.Printf("%v; switching log to stdout", err)
	}
}

// setLog interpretes env variable TS_LOGFILE and sets global loggers.
func setLog() error {

	// read env variable TS_LOGFILE
	filename, isset := os.LookupEnv("TS_LOGFILE")

	// error handling
	// return error, if not set
	if !isset {
		return fmt.Errorf("env variable $TS_LOGFILE not set")
	}
	// return error if empty
	if filename == "" {
		return fmt.Errorf("empty name in $TS_LOGFILE")
	}
	// return error if directory
	if l := filename[len(filename)-1:]; (l == "/") || (l == "\\") {
		return fmt.Errorf("no file, but only directory in $TS_LOGFILE = %v", filename)
	}

	// handle "discard" and return
	if filename == "discard" {
		noLogger()
		return nil
	}

	// handle "stdout" and return
	if filename == "stdout" {
		setStdout()
		return nil
	}

	// file ptr and error
	var (
		f   *os.File
		err error
	)

	// set file
	if filename == "tmp" {
		f, err = os.CreateTemp(os.TempDir(), "tslog_*")
	} else {
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	// handle file errors
	if err != nil {
		f.Close()
		return fmt.Errorf("log file: %w; $TS_SSL_LOGFILE = %v", err, filename)
	}

	// activate file logging and return
	setLogger(f)
	return nil
}

// setLogger initializes global loggers with f.
func setLogger(f *os.File) {
	if f == nil {
		panic(fmt.Errorf("nil pointer"))
	}
	var flags int = log.Ldate | log.Ltime | log.Lshortfile
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
