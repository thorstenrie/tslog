// Copyright (c) 2022 thorstenrie
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tslog

// Import standard library packages.
import (
	"io"      // io
	"os"      // os
	"testing" // testing
	"time"    // time

	"github.com/thorstenrie/tserr" // tserr
)

// A testcase serves input data for tests. Prefix and in are defined separately.
// For valid prefixes, global constants infoPrefix and errorPrefix can be used.
type testcase struct {
	prefix, in string
}

// A testcheck holds an actual output log message and the wanted result.
type testcheck struct {
	in   string   // actual output log message
	want testcase // wanted result (normally the input testcase)
}

// A testingtype interface implements Errorf for T, B and F.
// The interface enables generic functions for all test types T, B and F.
type testingtype interface {
	*testing.T | *testing.B | *testing.F
	Error(a ...any)
	Fatal(a ...any)
	Errorf(format string, a ...any)
	Fatalf(format string, a ...any)
}

// A testfunc is a function testing different dimensions of a testcheck.
type testfunc func(*testing.T, *testcheck)

// Slice of testcases.
var (
	testcases = []testcase{
		{errorPrefix, "test"},
		{infoPrefix, " "},
		{errorPrefix, "Hello World!"},
		{infoPrefix, "!12345"},
		{errorPrefix, "\n"},
	}
)

// TestEmpty performs logging with the env variable TS_LOGFILE set empty.
// Expected result is fallback logging to Stdout.
func TestEmpty(t *testing.T) {
	// Set env variable TS_LOGFILE to an empty string and reconfigure logging
	setEnv(t, "")
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestNotSet performs logging with the env variable TS_LOGFILE being unset.
// Expected result is fallback logging to Stdout.
func TestNotSet(t *testing.T) {
	// Unset env variable TS_LOGFILE
	if err := os.Unsetenv("TS_LOGFILE"); err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "unset env", Fn: "TS_LOGFILE", Err: err}))
	}
	// Re-initialize logging
	initialize()
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestDirectory1 performs logging with the env variable TS_LOGFILE set to a directory.
// Expected result is fallback logging to Stdout.
func TestDirectory1(t *testing.T) {
	// Set env variable TS_LOGFILE to temp directory and re-initialize logging
	setEnv(t, os.TempDir())
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestDirectory2 performs logging with the env variable TS_LOGFILE set to a directory.
// Expected result is fallback logging to Stdout.
func TestDirectory2(t *testing.T) {
	// Set env variable TS_LOGFILE to temp directory plus / and re-initialize logging
	setEnv(t, os.TempDir()+string(os.PathSeparator))
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestStdout performs logging with the env variable TS_LOGFILE set to stdout.
// Expected result is logging to Stdout.
func TestStdout(t *testing.T) {
	// Set env variable TS_LOGFILE to stdout and re-initialize logging
	setEnv(t, stdoutLogger)
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestTmp performs logging with the env variable TS_LOGFILE set to stdout.
// Expected result is logging to a temp file in the temp directory.
func TestTmp(t *testing.T) {
	// Set env variable TS_LOGFILE to tmp and re-initialize logging
	setEnv(t, tmpLogger)
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestDiscard performs logging with the env variable TS_LOGFILE set to discard.
// Expected result is no logging.
func TestDiscard(t *testing.T) {
	// Set env variable TS_LOGFILE to discard and re-initialize logging
	setEnv(t, discardLogger)
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestLogLength checks the length of all testcases in the log file.
// Since log.Lshortfile is not securly known during runtime, it only
// checks for the minimal length without log.Lshortfile.
// Note: Hard-coding log.Lshortfile in test functions would break
// tests if the source filename is changed
func TestLogLength(t *testing.T) {
	for tc := range testcases {
		testWrapper(t, testcases[tc], testLength)
	}
}

// TestLogPrefix checks the prefix of all testcases in the log file.
func TestLogPrefix(t *testing.T) {
	for tc := range testcases {
		testWrapper(t, testcases[tc], testPrefix)
	}
}

// TestLogMessage checks the contents of all testcases in the log file.
func TestLogMessage(t *testing.T) {
	for tc := range testcases {
		testWrapper(t, testcases[tc], testMessage)
	}
}

// BenchmarkLog performs a benchmark logging into a temp file in temp directory.
func BenchmarkLog(b *testing.B) {
	// Create temp file, set env variable, close the file and reconfigure logging
	tmpLog(b).Close()
	// Reset benchmark timer
	b.ResetTimer()
	// Run benchmark with all testcases in each iteration
	for i := 0; i < b.N; i++ {
		testLogAll(testcases)
	}
}

// FuzzInfo conducts fuzzing on log messages and checks for
// errors. The checks include the length of the log message,
// the prefix and the correct logging of the fuzzed message.
func FuzzInfo(f *testing.F) {
	// Addition of testcases to the seed corpus
	for _, tc := range testcases {
		f.Add(tc.in)
	}
	// Fuzz target
	f.Fuzz(func(t *testing.T, a string) {
		// Create testcase as informational log with fuzzing applied
		// to the log message
		tc := testcase{prefix: infoPrefix, in: a}
		// Test log message length
		testWrapper(t, tc, testLength)
		// Test prefix
		testWrapper(t, tc, testPrefix)
		// Test log message text
		testWrapper(t, tc, testMessage)
	})
}

// tmpLog creates a temp log file tslog_test_* in the temp directory.
// The env variable TS_LOGFILE is set accordingly.
// tmpLog returns the temp file. In case of errors tmpLog returns Stdout.
func tmpLog[T testingtype](tt T) *os.File {
	// Create temp log file tslog_test_* in the temp directory
	f, err := os.CreateTemp(os.TempDir(), "tslog_test_*")
	// In case of an error fall back to Stdout for logging
	if err != nil {
		f.Close()
		tt.Error(tserr.Op(&tserr.OpArgs{Op: "create", Fn: f.Name(), Err: err}))
		return os.Stdout
	}
	// Set TS_LOGFILE to temp log file tslog_test_* and re-initialize logging
	setEnv(tt, f.Name())
	// Return temp log file tslog_test_*
	return f
}

// testWrapper logs a testcase into a temp file and checks the
// result with tf.
func testWrapper(t *testing.T, tc testcase, tf testfunc) {
	// Create temp file, set env variable and reconfigure logging
	f := tmpLog(t)
	// Log testcase
	testLog(tc)

	var want testcase
	// Read log file
	in, err := io.ReadAll(f)
	if err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "ReadAll", Fn: f.Name(), Err: err}))
		return
	}
	// Check log file with tf
	tf(t, &testcheck{in: string(in), want: want})
	// Close temp log file
	if err := f.Close(); err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Close", Fn: f.Name(), Err: err}))
	}
}

// testLength checks the length of a log message.
// The minimum expected length of the log message is compared to
// the actual length of the log message.
// testLength implements testfunc.
func testLength(t *testing.T, tc *testcheck) {
	if t == nil {
		panic("nil pointer")
	}
	if tc == nil {
		t.Fatal(tserr.NilPtr())
	}
	// Calculates minimum length
	// Note: length of log.Lshortfile not known
	minl := len(tc.want.prefix) +
		len(tc.want.in) +
		len(time.Now().Format("2009/01/23 01:23:23")) +
		2 /*spaces*/ +
		2 /*colons*/
	// Get actual length of log message
	actl := len(tc.in)
	// Error in case actual length is lower than the calculated minimum length
	if actl < minl {
		t.Error(tserr.Higher(&tserr.HigherArgs{Var: "length of log message", Actual: int64(actl), LowerBound: int64(minl)}))
	}
}

// testPrefix checks the prefix of a log message
// The expected prefix is compared to the actual prefix.
// testPrefix implements testfunc.
func testPrefix(t *testing.T, tc *testcheck) {
	if t == nil {
		panic("nil pointer")
	}
	if tc == nil {
		t.Fatal(tserr.NilPtr())
	}
	// Check if the actual log message length is at least the prefix length
	minl := len(tc.want.prefix)
	actl := len(tc.in)
	if actl < minl {
		t.Errorf("log message length %d shorter than length %d of prefix %v", actl, minl, tc.want.prefix)
		return
	}
	// Get the actual prefix of the log message
	actp := tc.in[0:minl]
	// Error in case the actual prefix does not match the expected prefix
	if actp != tc.want.prefix {
		t.Error(tserr.NotEqualStr(&tserr.NotEqualStrArgs{X: tc.want.prefix, Y: actp}))
	}
}

// testMessage checks the contents of a log message.
// The expected contents is compared to the actual log message.
// testMessage implements testfunc.
func testMessage(t *testing.T, tc *testcheck) {
	if t == nil {
		panic("nil pointer")
	}
	if tc == nil {
		t.Fatal(tserr.NilPtr())
	}
	// Check if the actual log message length is at least the expected contents length
	minl := len(tc.want.in)
	actl := len(tc.in)
	if actl < minl {
		t.Fatal(tserr.Higher(&tserr.HigherArgs{Var: "length of log message", Actual: int64(actl), LowerBound: int64(minl)}))
	}
	// Get the actual log message without prefix and flags
	actm := tc.in[len(tc.in)-minl:]
	// Error in case the actual log message does not match the expected contents
	if actm != tc.want.in {
		t.Error(tserr.NotEqualStr(&tserr.NotEqualStrArgs{X: tc.want.in, Y: actm}))
	}
}

// testLog logs the testcase into the log file.
func testLog(tc testcase) {
	if tc.prefix == infoPrefix {
		I.Print(tc.in)
	} else if tc.prefix == errorPrefix {
		E.Print(tc.in)
	} else {
		E.Printf("expected prefix %v or %v, but got prefix %v for log message %v", infoPrefix, errorPrefix, tc.prefix, tc.in)
	}
}

// testLogAll logs all testcases into the log file.
func testLogAll(tc []testcase) {
	for i := range tc {
		testLog(tc[i])
	}
}

// setEnv sets env variable TS_LOGFILE to fn and re-initialize loggers.
func setEnv[T testingtype](tt T, fn string) {
	if err := os.Setenv("TS_LOGFILE", fn); err != nil {
		tt.Fatal(tserr.Op(&tserr.OpArgs{Op: "set env", Fn: "TS_LOGFILE", Err: err}))
	}
	initialize()
}
