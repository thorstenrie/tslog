// Copyright (c) 2023 thorstenrie
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tslog

// Import standard library packages, tserr and tsfio.
import (
	"bufio"   // bufio
	"bytes"   // bytes
	"os"      // os
	"testing" // testing

	"github.com/thorstenrie/tserr" // tserr
	"github.com/thorstenrie/tsfio" // tsfio
)

// A testfunc is a function testing logging into a file.
type testfunc func(*testing.T, int, tsfio.Filename)

// Interface fio is constrained to type tsfio.Filename and tsfio.Directory
type fio interface {
	tsfio.Filename | tsfio.Directory
}

// tmpLog creates a temp log file tslog_test_* in the temp directory.
// It returns the temp filename. In case of errors tmpLog returns Stdout.
func tmp[T testingtype](tt T) tsfio.Filename {
	// Panic if tt is nil
	if tt == nil {
		panic("nil pointer")
	}
	// Create temp log file tslog_test_* in the temp directory
	f, err := os.CreateTemp(os.TempDir(), "tslog_test_*")
	// In case of an error fall back to Stdout for logging
	if err != nil {
		// Close temporary file
		f.Close()
		// Record error
		tt.Error(tserr.Op(&tserr.OpArgs{Op: "create", Fn: f.Name(), Err: err}))
		// Return Stdout
		return StdoutLogger
	}
	// Retrieve filename of temporary file f
	fn := tsfio.Filename(f.Name())
	// Close temporary file f
	if err := f.Close(); err != nil {
		// Record an error, if close fails
		tt.Error(tserr.Op(&tserr.OpArgs{Op: "close", Fn: string(fn), Err: err}))
	}
	// Return filename of temp log file tslog_test_*
	return fn
}

// tmpDir creates a new temporary directory in the default directory for temporary files
// with the prefix tslog_testdir_ and a random string to the end. In case of an error
// the execution stops.
func tmpDir[T testingtype](tt T) tsfio.Directory {
	// Panic if tt is nil
	if tt == nil {
		panic("nil pointer")
	}
	// Create the temporary directory
	d, err := os.MkdirTemp("", "tslog_testdir_*")
	// Stop execution in case of an error
	if err != nil {
		tt.Fatal(tserr.Op(&tserr.OpArgs{Op: "create temp dir", Fn: d, Err: err}))
	}
	// Return the temporary Directory
	return tsfio.Directory(d)
}

// rm removes the file named Filename a or empty Directory a. In case of an error
// execution stops.
func rm[T fio](t *testing.T, a T) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Remove file or empty directory
	if err := os.Remove(string(a)); err != nil {
		// Stop execution in case of an error
		t.Fatal(tserr.Op(&tserr.OpArgs{Op: "Remove", Fn: string(a), Err: err}))
	}
}

// size returns the length of regular file fn.
func size[T testingtype](tt T, fn tsfio.Filename) int64 {
	// Panic if tt is nil
	if tt == nil {
		panic("nil pointer")
	}
	// Retrieve length of file fn
	fs, e := tsfio.FileSize(fn)
	// Record an error, if FileSize fails.
	if e != nil {
		tt.Error(tserr.Op(&tserr.OpArgs{Op: "FileSize", Fn: string(fn), Err: e}))
	}
	// Return the file size
	return fs
}

// scanner returns a *bufio.Scanner on the contents of file fn. The returned Scanner
// returns each line of the contents of file fn. It panics if tt is nil. Execution
// stops if reading the file fn fails.
func scanner[T testingtype](tt T, fn tsfio.Filename) *bufio.Scanner {
	// Panic if tt is nil
	if tt == nil {
		panic("nil pointer")
	}
	// Read file fn
	in, err := tsfio.ReadFile(fn)
	// Execution stops, if ReadFile fails.
	if err != nil {
		tt.Fatal(tserr.Op(&tserr.OpArgs{Op: "Read file", Fn: string(fn), Err: err}))
	}
	// Create new buffer on text in file fn
	bf := bytes.NewBuffer(in)
	// Create new scanner on buffer with text in file fn
	fs := bufio.NewScanner(bf)
	// Set split function to scan lines of the text in file fn
	fs.Split(bufio.ScanLines)
	// Return scanner
	return fs
}

// evaluate compares the logging result in fn with the testcases. It panics if
// t is nil. It records an error if a performed operation reports an error or if the text in the
// temporary output file does not match the expected result based on the testcases.
func evaluate(t *testing.T, fn tsfio.Filename) {
	// Panic if t is nil
	if t == nil {
		panic("nil pointer")
	}
	// Create scanner fs on logging output file fn
	fs := scanner(t, fn)
	// Remove logging output file fn
	rm(t, fn)
	// Iterate over fs line by line
	var i, m int = 0, len(testcases)
	for ; fs.Scan() && i < m; i++ {
		// Evaluate log file with testcases
		testMessage(t, fs.Bytes(), testcases[i])
	}
	// Record an error if scanner returns an error
	if err := fs.Err(); err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Scan", Fn: string(fn), Err: err}))
	}
	// Record an error if no. lines in logging output file does not equal no. testcases
	if i != m {
		t.Error(tserr.Equal(&tserr.EqualArgs{Var: "No. lines", Actual: int64(i), Want: int64(m)}))
	}
}
