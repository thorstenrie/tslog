package tslog

import (
	"io"
	"os"
	"testing"
	"time"
)

var (
	testcases = []testcase{
		{errorPrefix, "test"},
		{infoPrefix, " "},
		{errorPrefix, "Hello World!"},
		{infoPrefix, "!12345"},
	}
)

type testcase struct {
	prefix, in string
}

type testStruct struct {
	in   string
	want testcase
}

type testFunc func(*testing.T, *testStruct)

func TestEmpty(t *testing.T) {
	if err := setEnv(""); err != nil {
		t.Errorf("set empty env TS_LOGFILE failed: %v", err)
	}
	testLogAll(testcases)
}

func TestNotSet(t *testing.T) {
	if err := os.Unsetenv("TS_LOGFILE"); err != nil {
		t.Errorf("unset env TS_LOGFILE failed: %v", err)
	}
	Reset()
	testLogAll(testcases)
}

func TestDirectory1(t *testing.T) {
	if err := setEnv("/tmp/"); err != nil {
		t.Errorf("set env TS_LOGFILE = /tmp/ failed: %v", err)
	}
	testLogAll(testcases)
}

func TestDirectory2(t *testing.T) {
	if err := setEnv("/tmp"); err != nil {
		t.Errorf("set env TS_LOGFILE = /tmp failed: %v", err)
	}
	testLogAll(testcases)
}

func TestStdout(t *testing.T) {
	if err := setEnv("stdout"); err != nil {
		t.Errorf("set env TS_LOGFILE = stdout failed: %v", err)
	}
	testLogAll(testcases)
}

func TestDiscard(t *testing.T) {
	if err := setEnv("discard"); err != nil {
		t.Errorf("set env TS_LOGFILE = discard failed: %v", err)
	}
	testLogAll(testcases)
}

func TestInit(t *testing.T) {
	for tc := range testcases {
		testLog(testcases[tc])
	}
}

func BenchmarkLog(b *testing.B) {
	tmpLog(b).Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testLogAll(testcases)
	}
}

type testingType interface {
	*testing.T | *testing.B | *testing.F
	Errorf(format string, a ...any)
}

func tmpLog[T testingType](tt T) *os.File {
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

func testWrapper(t *testing.T, tc testcase, tf testFunc) {
	f := tmpLog(t)
	testLog(tc)
	var want testcase
	in, err := io.ReadAll(f)
	if err != nil {
		t.Errorf("open %v failed: %v", f.Name(), err)
		return
	}
	tf(t, &testStruct{in: string(in), want: want})
	if err := f.Close(); err != nil {
		t.Errorf("closing %v failed: %v", f.Name(), err)
	}
}

func testLength(t *testing.T, tc *testStruct) {
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

func testPrefix(t *testing.T, tc *testStruct) {
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

func testMessage(t *testing.T, tc *testStruct) {
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
