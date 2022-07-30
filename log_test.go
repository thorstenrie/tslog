package tslog

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
)

var (
	infoText  string = "info test"
	errorText string = "error test"
)

func TestEmpty(t *testing.T) {
	if err := setEnv(""); err != nil {
		t.Errorf("set empty env TS_LOGFILE failed: %v", err)
	}
	testLog()
}

func TestNotSet(t *testing.T) {
	if err := os.Unsetenv("TS_LOGFILE"); err != nil {
		t.Errorf("unset env TS_LOGFILE failed: %v", err)
	}
	Reset()
	testLog()
}

func TestDirectory1(t *testing.T) {
	if err := setEnv("/tmp/"); err != nil {
		t.Errorf("set env TS_LOGFILE = /tmp/ failed: %v", err)
	}
	testLog()
}

func TestDirectory2(t *testing.T) {
	if err := setEnv("/tmp"); err != nil {
		t.Errorf("set env TS_LOGFILE = /tmp failed: %v", err)
	}
	testLog()
}

func TestStdout(t *testing.T) {
	if err := setEnv("stdout"); err != nil {
		t.Errorf("set env TS_LOGFILE = stdout failed: %v", err)
	}
	testLog()
}

func TestDiscard(t *testing.T) {
	if err := setEnv("discard"); err != nil {
		t.Errorf("set env TS_LOGFILE = discard failed: %v", err)
	}
	testLog()
}

func TestInit(t *testing.T) {
	testLog()
}

func TestFile(t *testing.T) {
	fn := filepath.Join(os.TempDir(), "tslog_test.txt")
	if _, err := os.Create(fn); err != nil {
		t.Errorf("creating %v failed: %v", fn, err)
	}
	if err := setEnv(fn); err != nil {
		t.Errorf("set env TS_LOGFILE = %v failed: %v", fn, err)
	}
	testLog()
	f, err := os.Open(fn)
	if err != nil {
		t.Errorf("open logfile %v failed: %v", fn, err)
	}
	scanner := bufio.NewScanner(f)
	i := 0
	for ; scanner.Scan(); i++ {
		if ln := scanner.Text(); len(ln) > 0 {
			if len(ln) < mLen() {
				f.Close()
				t.Errorf("Line %v in logfile %v too short", ln, fn)
			} else if (ln[0:len(infoPrefix)] != infoPrefix) && (ln[0:len(errorPrefix)] != errorPrefix) {
				f.Close()
				t.Errorf("Prefix not found; %v or %v expected", infoPrefix, errorPrefix)
			} else if (ln[len(ln)-len(infoText):] != infoText) && (ln[len(ln)-len(errorText):] != errorText) {
				f.Close()
				t.Errorf("Text not found; %v or %v expected", infoText, errorText)
			}
		}
	}
	f.Close()
}

func testLog() {
	I.Print(infoText)
	E.Print(errorText)
}

func setEnv(fn string) error {
	err := os.Setenv("TS_LOGFILE", fn)
	Reset()
	return err
}

func mLen() int {
	v := []int{len(infoPrefix), len(errorPrefix), len(infoText), len(errorText)}
	min := v[0]
	for i := range v {
		if i < min {
			min = i
		}
	}
	return min
}
