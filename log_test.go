package tslog

// Import standard library packages.
import (
	"io"      // io
	"os"      // os
	"testing" // testing
	"time"    // time
)

// A testcase serves input data for tests. Prefix and in are defined separately.
// For valid prefixes, global constants infoPrefix and errorPrefix can be used.
type testcase struct {
	prefix, in string
}

// A testcheck holds an actual output log message and the wanted result
type testcheck struct {
	in   string   // actual output log message
	want testcase // wanted result (normally the input testcase)
}

// A testingtype interface implements Errorf for T, B and F
// The interface enables generic functions for all test types T, B and F
type testingtype interface {
	*testing.T | *testing.B | *testing.F
	Errorf(format string, a ...any)
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
	// Set env variable TS_LOGFILE to an empty string
	if err := setEnv(""); err != nil {
		t.Errorf("set empty env TS_LOGFILE failed: %v", err)
	}
	// Reconfigure logging
	Reset()
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestNotSet performs logging with the env variable TS_LOGFILE being unset.
// Expected result is fallback logging to Stdout.
func TestNotSet(t *testing.T) {
	// Unset env variable TS_LOGFILE
	if err := os.Unsetenv("TS_LOGFILE"); err != nil {
		t.Errorf("unset env TS_LOGFILE failed: %v", err)
	}
	// Reconfigure logging
	Reset()
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestDirectory1 performs logging with the env variable TS_LOGFILE set to a directory.
// Expected result is fallback logging to Stdout.
func TestDirectory1(t *testing.T) {
	// Set env variable TS_LOGFILE to temp directory
	if err := setEnv(os.TempDir()); err != nil {
		t.Errorf("set env TS_LOGFILE = /tmp/ failed: %v", err)
	}
	// Reconfigure logging
	Reset()
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestDirectory2 performs logging with the env variable TS_LOGFILE set to a directory.
// Expected result is fallback logging to Stdout.
func TestDirectory2(t *testing.T) {
	// Set env variable TS_LOGFILE to temp directory plus /
	if err := setEnv(os.TempDir() + string(os.PathSeparator)); err != nil {
		t.Errorf("set env TS_LOGFILE = %v failed: %v", os.TempDir(), err)
	}
	// Reconfigure logging
	Reset()
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestStdout performs logging with the env variable TS_LOGFILE set to stdout.
// Expected result is logging to Stdout.
func TestStdout(t *testing.T) {
	// Set env variable TS_LOGFILE to stdout
	if err := setEnv(stdoutLogger); err != nil {
		t.Errorf("set env TS_LOGFILE = stdout failed: %v", err)
	}
	// Reconfigure logging
	Reset()
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestTmp performs logging with the env variable TS_LOGFILE set to stdout.
// Expected result is logging to a temp file in the temp directory.
func TestTmp(t *testing.T) {
	// Set env variable TS_LOGFILE to tmp
	if err := setEnv(tmpLogger); err != nil {
		t.Errorf("set env TS_LOGFILE = stdout failed: %v", err)
	}
	// Reconfigure logging
	Reset()
	// Perform logging of testcases
	testLogAll(testcases)
}

// TestDiscard performs logging with the env variable TS_LOGFILE set to discard.
// Expected result is no logging.
func TestDiscard(t *testing.T) {
	// Set env variable TS_LOGFILE to discard
	if err := setEnv(discardLogger); err != nil {
		t.Errorf("set env TS_LOGFILE = discard failed: %v", err)
	}
	// Reconfigure logging
	Reset()
	// Perform logging of testcases
	testLogAll(testcases)
}

// BenchmarkLog performs a benchmark logging into a temp file in temp directory.
func BenchmarkLog(b *testing.B) {
	// Create temp file, set env variable and close the file
	tmpLog(b).Close()
	// Reconfigure logging
	Reset()
	// Reset benchmark timer
	b.ResetTimer()
	// Run benchmark with all testcases in each iteration
	for i := 0; i < b.N; i++ {
		testLogAll(testcases)
	}
}

func tmpLog[T testingtype](tt T) *os.File {
	f, err := os.CreateTemp(os.TempDir(), "tslog_test_*")
	if err != nil {
		f.Close()
		tt.Errorf("creating %v failed: %v", f.Name(), err)
		return os.Stdout
	}
	if err := setEnv(f.Name()); err != nil {
		tt.Errorf("set env TS_LOGFILE = %v failed: %v", f.Name(), err)
	}
	return f
}

func FuzzInfo(f *testing.F) {
	for _, tc := range testcases {
		f.Add(tc.in)
	}
	f.Fuzz(func(t *testing.T, a string) {
		tc := testcase{prefix: infoPrefix, in: a}
		testWrapper(t, tc, testLength)
		testWrapper(t, tc, testPrefix)
		testWrapper(t, tc, testMessage)
	})
}

func testWrapper(t *testing.T, tc testcase, tf testfunc) {
	f := tmpLog(t)
	testLog(tc)
	var want testcase
	in, err := io.ReadAll(f)
	if err != nil {
		t.Errorf("open %v failed: %v", f.Name(), err)
		return
	}
	tf(t, &testcheck{in: string(in), want: want})
	if err := f.Close(); err != nil {
		t.Errorf("closing %v failed: %v", f.Name(), err)
	}
}

func testLength(t *testing.T, tc *testcheck) {
	if t == nil {
		E.Fatalln("nil pointer")
	}
	if tc == nil {
		t.Errorf("nil pointer")
		return
	}
	minl := len(tc.want.prefix) + len(tc.want.in) + len(time.Now().Format("2009/01/23 01:23:23")) + 2 /*spaces*/ + 2 /*colons*/
	actl := len(tc.in)
	if actl < minl {
		t.Errorf("minimum length %d expected, but length is %d", minl, actl)
	}
}

func testPrefix(t *testing.T, tc *testcheck) {
	if t == nil {
		E.Fatalln("nil pointer")
	}
	if tc == nil {
		t.Errorf("nil pointer")
		return
	}
	minl := len(tc.want.prefix)
	actl := len(tc.in)
	if actl < minl {
		t.Errorf("log message length %d shorter than length %d of prefix %v", actl, minl, tc.want.prefix)
		return
	}
	actp := tc.in[0:minl]
	if actp != tc.want.prefix {
		t.Errorf("expected prefix %v but got %v", tc.want.prefix, actp)
	}
}

func testMessage(t *testing.T, tc *testcheck) {
	if t == nil {
		E.Fatalln("nil pointer")
	}
	if tc == nil {
		t.Errorf("nil pointer")
		return
	}
	minl := len(tc.want.in)
	actl := len(tc.in)
	if actl < minl {
		t.Errorf("log message length %d shorter than length %d of message %v", actl, minl, tc.want.in)
		return
	}
	actm := tc.in[len(tc.in)-minl:]
	if actm != tc.want.in {
		t.Errorf("expected log message %v but got %v", tc.want.in, actm)
	}
}

func TestLogLength(t *testing.T) {
	for tc := range testcases {
		testWrapper(t, testcases[tc], testLength)
	}
}

func TestLogPrefix(t *testing.T) {
	for tc := range testcases {
		testWrapper(t, testcases[tc], testPrefix)
	}
}

func TestLogMessage(t *testing.T) {
	for tc := range testcases {
		testWrapper(t, testcases[tc], testMessage)
	}
}

func testLog(tc testcase) {
	if tc.prefix == infoPrefix {
		I.Print(tc.in)
	} else if tc.prefix == errorPrefix {
		E.Print(tc.in)
	} else {
		E.Printf("expected prefix %v or %v, but got prefix %v for log message %v", infoPrefix, errorPrefix, tc.prefix, tc.in)
	}
}

func testLogAll(tc []testcase) {
	for i := range tc {
		testLog(tc[i])
	}
}

func setEnv(fn string) error {
	err := os.Setenv("TS_LOGFILE", fn)
	Reset()
	return err
}
