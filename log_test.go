// Copyright (c) 2023 thorstenrie
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tslog

// Import standard library packages, tserr and tsfio.
import (
	"encoding/json" // json
	"errors"
	"fmt" // io

	"testing" // testing
	"time"    // time

	"github.com/thorstenrie/tserr" // tserr
	"github.com/thorstenrie/tsfio" // tsfio
)

// A testcase serves input data for tests. A testcases contains the level and message.
type testcase struct {
	level int    // Log level
	in    string // Log message
}

// A testingtype interface implements Error and Fatal for T, B and F.
// The interface enables generic functions for all test types T, B and F.
type testingtype interface {
	*testing.T | *testing.B | *testing.F // Interface constraint to T, B and F
	Error(a ...any)                      // Record formated output followed by Fail
	Fatal(a ...any)                      // Record formated output followed by FailNow
}

// Slice of testcases
var (
	testcases = []*testcase{
		{TraceLevel, "test"},
		{DebugLevel, " "},
		{InfoLevel, "Hello World!"},
		{WarnLevel, "Warning!"},
		{ErrorLevel, "!12345"},
		{FatalLevel, "\n"},
	}
)

// TestNotSet performs logging of all testcases with default settings.
// Expected result is logging to Stdout.
func TestNotSet(t *testing.T) {
	testLogAll(t, testcases)
}

// TestDefaultLog retrieves default global pre-defined standard logger
// and performs logging of all testcases.
func TestDefaultLog(t *testing.T) {
	// Retrieve the global pre-defined standard logger in l
	l := Default()
	// Perform logging of all testcases with l
	testLoggerAll(t, testcases, l)
}

// TestStdout performs logging with the default logger set to stdout.
// Expected result is logging to Stdout.
func TestStdout(t *testing.T) {
	// Set ouput of the default logger to Stdout
	SetOutput(StdoutLogger)
	// Perform logging of testcases
	testLogAll(t, testcases)
}

// TestDiscard performs logging with the default logger set to discard.
// Expected result is no logging.
func TestDiscard(t *testing.T) {
	// Set ouput of the default logger to discard
	SetOutput(DiscardLogger)
	// Perform logging of testcases
	testLogAll(t, testcases)
}

// TestTmp performs logging with the default logger set to a temporary file.
// Expected result is logging in a temporary file in the temporary directory.
func TestTmp(t *testing.T) {
	// Set output of the default logger to a temporary file
	SetOutput(TmpLogger)
	// Perform logging of testcases
	testLogAll(t, testcases)
}

// TestDir sets output of the default logger to a directory. It is expected to
// return an error. The test fails if no error is returned.
func TestDir(t *testing.T) {
	// Create temporary directory d
	d := tmpDir(t)
	// Set output of the default logger to d
	if err := SetOutput(tsfio.Filename(d)); err == nil {
		// Record an error if SetOutput returns nil instead of an error
		t.Error(tserr.NilFailed("Set output to temp directory"))
	}
	// Remove the temporary directory d
	rm(t, d)
}

// TestLogger performs logging with a newly created logger with output set to a temporary file.
// It logs all testcases to the created logger and evaluates the output in the temporary file.
// It records an error if a performed operation reports an error or if the text in the
// temporary output file does not match the expected result based on the testcases.
func TestLogger(t *testing.T) {
	// Create the temporary file fn
	fn := tmp(t)
	// Create new logger lg
	lg := New()
	// Set output to temporary file fn
	if err := lg.SetOutput(fn); err != nil {
		// Record an error, if SetOutput fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Set output", Fn: string(fn), Err: err}))
	}
	// Set logging level to Trace
	if err := lg.SetLevel(TraceLevel); err != nil {
		// Record an error, if SetLevel fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: fmt.Sprintf("Set level to %d for", TraceLevel), Fn: string(fn), Err: err}))
	}
	// Log all testcases using logger lg
	testLoggerAll(t, testcases, lg)
	// Evaluate logging in output file fn
	evaluate(t, fn)
}

// TestLog performs logging with the default predefined standard logger with output set to a temporary file.
// It logs all testcases to the default logger and evaluates the output in the temporary file.
// It records an error if a performed operation reports an error or if the text in the
// temporary output file does not match the expected result based on the testcases.
func TestLog(t *testing.T) {
	// Create the temporary file fn
	fn := tmp(t)
	// Set output to temporary file fn
	if err := SetOutput(fn); err != nil {
		// Record an error, if SetOutput fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Set output", Fn: string(fn), Err: err}))
	}
	// Set logging level to Trace
	if err := SetLevel(TraceLevel); err != nil {
		// Record an error, if SetLevel fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: fmt.Sprintf("Set level to %d", TraceLevel), Fn: string(fn), Err: err}))
	}
	// Log all testcases using the default predefined standard logger
	testLogAll(t, testcases)
	// Evaluate logging in output file fn
	evaluate(t, fn)
}

// TestSetLevelErr sets the log level one below Trace level and one above Fatal level.
// It expects to receive an error, when calling SetLevel. The test fails if SetLevel
// returns nil.
func TestSetLevelErr(t *testing.T) {
	// Set log level minus one below Trace level
	if err := SetLevel(TraceLevel - 1); err == nil {
		// Record an error if SetLevel returns nil
		t.Error(tserr.NilFailed("Set level"))
	}
	// Set log level plus one above Fatal level
	if err := SetLevel(FatalLevel + 1); err == nil {
		// Record an error if SetLevel returns nil
		t.Error(tserr.NilFailed("Set level"))
	}
}

// TestSetLevelTrace tests log messages at all log levels to be logged
// if log level is set to Trace. It fails if an operation fails or if a
// messaged is logged other than Trace level.
func TestSetLevelTrace(t *testing.T) {
	testLevel(t, testTrace)
}

// TestSetLevelFatal tests log messages at all log levels to be logged
// if log level is set to Fatal. It fails if an operation fails or if a
// messaged is not logged.
func TestSetLevelFatal(t *testing.T) {
	testLevel(t, testFatal)
}

// testLevel iterates all log level from Trace level to Fatal level and calls testfunc tf.
func testLevel(t *testing.T, tf testfunc) {
	// Create an array with all log levels from Trace level to Fatal level
	lvls := [6]int{TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel}
	// Create the temporary file fn
	fn := tmp(t)
	// Set log output to the temporary file fn
	if err := SetOutput(fn); err != nil {
		// Record an error if SetOutput fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Set output", Fn: string(fn), Err: err}))
	}
	// Iterate all log levels
	for _, v := range lvls {
		// Call testfunc tf for each log level and providing fn
		tf(t, v, fn)
	}
	// Remove the temporary file fn
	rm(t, fn)
}

// testTrace implements testfunc. It sets log level to v, logs a testcase at Trace level
// and evaluates the output in file fn.
func testTrace(t *testing.T, v int, fn tsfio.Filename) {
	// Set log level to v
	if err := SetLevel(v); err != nil {
		// Record an error, if SetLevel fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Set level", Fn: fmt.Sprint(v), Err: err}))
	}
	// Create testcase with log level Trace
	tc := testcase{level: TraceLevel, in: "test"}
	// Log testcase on log level Trace
	if err := Trace(tc.in); err != nil {
		// Record an error, if Trace fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Trace for level", Fn: fmt.Sprint(v), Err: err}))
	}
	// Read contents of file fn
	in, e := tsfio.ReadFile(fn)
	// Record an error, if ReadFile fails
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Read file", Fn: string(fn), Err: e}))
	}
	// Reset file fn
	if err := tsfio.ResetFile(fn); err != nil {
		// Record an error, if ResetFile fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "ResetFile", Fn: string(fn), Err: err}))
	}
	// Evaluate log message from fn, in case v equals Trace level
	if v == TraceLevel {
		testMessage(t, in, &tc)
	} else {
		// Check fn for its length, in case v equals a higher than Trace level
		si := size(t, fn)
		// Record an error, if length of fn is higher than zero
		if si > 0 {
			t.Error(tserr.Equal(&tserr.EqualArgs{Var: "Size of log file", Actual: si, Want: 0}))
		}
	}
}

// testFatal implements testfunc. It sets log level to v, logs a testcase at Fatal level
// and evaluates the output in file fn.
func testFatal(t *testing.T, v int, fn tsfio.Filename) {
	// Set log level to v
	if err := SetLevel(v); err != nil {
		// Record an error, if SetLevel fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Set level", Fn: fmt.Sprint(v), Err: err}))
	}
	// Create testcase with log level Fatal
	tc := testcase{level: FatalLevel, in: "test"}
	// Log testcase on log level Fatal
	if err := Fatal(errors.New(tc.in)); err != nil {
		// Record an error, if Trace fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Fatal for level", Fn: fmt.Sprint(v), Err: err}))
	}
	// Reset file fn
	in, e := tsfio.ReadFile(fn)
	// Record an error, if ReadFile fails
	if e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Read file", Fn: string(fn), Err: e}))
	}
	// Reset file fn
	if err := tsfio.ResetFile(fn); err != nil {
		// Record an error, if ResetFile fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "ResetFile", Fn: string(fn), Err: err}))
	}
	// Evaluate log message from fn
	testMessage(t, in, &tc)
}

// testMessage checks the prefix and the contents of the log message in.
// The expected prefix and the expected contents is compared to the actual log message.
// It panics if t is nil. The execution stops if want or in are nil. The test fails
// if Unmarchal fails, the actual prefix does not match the expected prefix or if the
// expected message does not equal the actual message.
func testMessage(t *testing.T, in []byte, want *testcase) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Execution stops if want or in are nil
	if (want == nil) || (in == nil) {
		t.Fatal(tserr.NilPtr())
	}
	// Retrieve wanted log level as string
	wantl, err := level(want.level)
	// Record error if the log level does not exit
	if err != nil {
		t.Error(tserr.NotExistent(fmt.Sprintf("log level %d", want.level)))
	}
	// Unmarshal log message in
	var lmsg logwrap
	if err := json.Unmarshal(in, &lmsg); err != nil {
		// Record an error if Unmarshal fails
		t.Error(tserr.Op(&tserr.OpArgs{Op: "json unmarshal", Fn: string(in), Err: err}))
	}
	// Record an error if the expected log level does not equal the actual log level
	if lmsg.L.Lvl != wantl {
		t.Error(tserr.NotEqualStr(&tserr.NotEqualStrArgs{X: wantl, Y: lmsg.L.Lvl}))
	}
	// Record an error if the expected log message does not equal the actual log message
	if lmsg.L.Msg != want.in {
		t.Error(tserr.NotEqualStr(&tserr.NotEqualStrArgs{X: want.in, Y: lmsg.L.Msg}))
	}
	// Record an error if the timestamp of the log message cannot be parsed
	if _, err := time.Parse(timeLayout, lmsg.L.Now); err != nil {
		t.Error(tserr.Check(&tserr.CheckArgs{F: lmsg.L.Now, Err: err}))
	}
}

// testLogAll logs all testcases in tc. It panics, if t is nil. It records an error if
// tc is nil or if logging of a testcase in tc fails.
func testLogAll(t *testing.T, tc []*testcase) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Record an error if tc is nil
	if tc == nil {
		t.Error(tserr.NilPtr())
	}
	// Iterate all testcases in tc
	for i := range tc {
		// Log all testcases
		if e := testLog(tc[i]); e != nil {
			// Record an error if logging fails
			t.Error(tserr.Op(&tserr.OpArgs{Op: "test log", Fn: fmt.Sprint(tc[i]), Err: e}))
		}
	}
}

// testLog logs testcase tc. It returns an error fi tc is nil or if the log level
// of the testcase does not exist.
func testLog(tc *testcase) error {
	// Return an error if tc is nil
	if tc == nil {
		return tserr.NilPtr()
	}
	// Log according to the defined log level in testcase tc
	switch tc.level {
	case TraceLevel:
		return Trace(tc.in)
	case DebugLevel:
		return Debug(tc.in)
	case InfoLevel:
		return Info(tc.in)
	case WarnLevel:
		return Warn(tc.in)
	case ErrorLevel:
		return Error(errors.New(tc.in))
	case FatalLevel:
		return Fatal(errors.New(tc.in))
	}
	// Return an error if the log level does not exist
	return tserr.NotExistent(fmt.Sprintf("%d", tc.level))
}

// testLoggerAll logs all testcases in tc using the logger l. It panics
// if t is nil. It stops execution if l or tc are nil. It records an error
// if the testcase is nil or if logging fails.
func testLoggerAll(t *testing.T, tc []*testcase, l *Logger) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Stop execution if l or tc are nil
	if (l == nil) || (tc == nil) {
		t.Fatal(tserr.NilPtr())
	}
	// Iterate all testcases in tc
	for i := range tc {
		// Record an error if the testcase is nil
		if tc[i] == nil {
			t.Error(tserr.NilPtr())
		} else {
			// Log testcase with logger l
			if e := testLogger(tc[i], l); e != nil {
				// Record an error if logging fails
				t.Error(tserr.Op(&tserr.OpArgs{Op: "test log", Fn: fmt.Sprint(tc[i]), Err: e}))
			}
		}
	}
}

// testLogger logs the testcase using logger l. It returns an error
// if l or tc are nil. It also returns an error if the log level
// in testcase tc does not exist.
func testLogger(tc *testcase, l *Logger) error {
	// Return an error if l or tc are nil
	if (l == nil) || (tc == nil) {
		return tserr.NilPtr()
	}
	// Log testcase according to the defined log level
	// in the testcase
	switch tc.level {
	case TraceLevel:
		return l.Trace(tc.in)
	case DebugLevel:
		return l.Debug(tc.in)
	case InfoLevel:
		return l.Info(tc.in)
	case WarnLevel:
		return l.Warn(tc.in)
	case ErrorLevel:
		return l.Error(errors.New(tc.in))
	case FatalLevel:
		return l.Fatal(errors.New(tc.in))
	}
	// Return an error if the log level in the testcase does not exist.
	return tserr.NotExistent(fmt.Sprintf("%d", tc.level))
}
