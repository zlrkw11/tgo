package main

import (
	"fmt"
	"testing"
	"time"
)

// --- basic pass/fail ---

func TestPass(t *testing.T) {
	if 1+1 != 2 {
		t.Fatal("math is broken")
	}
}

func TestFail(t *testing.T) {
	t.Fatal("this test always fails")
}

// --- skip ---

func TestSkipped(t *testing.T) {
	t.Skip("skipping: not implemented yet")
}

// --- multiple errors in one test ---

func TestMultipleErrors(t *testing.T) {
	t.Error("first error")
	t.Error("second error")
	t.Error("third error")
}

// --- nested subtests ---

func TestSubtests(t *testing.T) {
	t.Run("passing", func(t *testing.T) {
		// ok
	})
	t.Run("failing", func(t *testing.T) {
		t.Fatal("subtest failed")
	})
	t.Run("skipped", func(t *testing.T) {
		t.Skip("skipped subtest")
	})
}

// --- long error message ---

func TestLongError(t *testing.T) {
	t.Fatalf("something went wrong: got %q, want %q",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	)
}

// --- multiline error ---

func TestMultilineError(t *testing.T) {
	t.Fatal("line 1\nline 2\nline 3")
}

// --- panic recovery ---

func TestPanic(t *testing.T) {
	panic("unexpected panic!")
}

// --- slow test ---

func TestSlow(t *testing.T) {
	time.Sleep(2 * time.Second)
}

// --- test with verbose output ---

func TestWithLogs(t *testing.T) {
	t.Log("info: setting up")
	t.Log("info: doing work")
	t.Log("info: tearing down")
	t.Fatal("failed after logging")
}

// --- error with formatting ---

func TestFormattedError(t *testing.T) {
	got := map[string]int{"a": 1, "b": 2}
	want := map[string]int{"a": 1, "b": 3}
	t.Fatalf("mismatch:\n  got:  %v\n  want: %v", got, want)
}

// --- nil pointer (runtime error) ---

func TestNilPointer(t *testing.T) {
	var s *string
	fmt.Println(*s)
}

// --- empty test (should pass) ---

func TestEmpty(t *testing.T) {
}
