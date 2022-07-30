package tslog

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	I           *log.Logger // information
	E           *log.Logger // error
	infoPrefix  string      = "INFO: "
	errorPrefix string      = "ERROR: "
)

func init() {
	Reset()
}

func Reset() {
	if err := setLog(); err != nil {
		setStdout()
		E.Printf("%v; switching log to stdout", err)
	}
}

func setLog() error {
	filename, isset := os.LookupEnv("TS_LOGFILE")
	if !isset {
		return fmt.Errorf("env variable $TS_LOGFILE not set")
	}
	if filename == "" {
		return fmt.Errorf("empty name in $TS_LOGFILE")
	}
	if filename[len(filename)-1:] == "/" {
		return fmt.Errorf("no file, but only directory in $TS_LOGFILE = %v", filename)
	}
	if filename == "stdout" {
		setStdout()
	} else if filename == "discard" {
		noLogger()
	} else {
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("log file: %w; $TS_SSL_LOGFILE = %v", err, filename)
		}
		setLogger(f)
	}
	return nil
}

func setLogger(f *os.File) {
	if f == nil {
		panic(fmt.Errorf("nil pointer detected"))
	}
	var flags int = log.Ldate | log.Ltime | log.Lshortfile
	I = log.New(f, infoPrefix, flags)
	E = log.New(f, errorPrefix, flags)
}

func setStdout() {
	setLogger(os.Stdout)
}

func noLogger() {
	I = log.New(io.Discard, "", 0)
	E = log.New(io.Discard, "", 0)
}
