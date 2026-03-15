package main

import "testing"

func TestAdd(t *testing.T) {
	if 1+1 != 2 {
		t.Fatal("math is broken")
	}
}

func TestSubtract(t *testing.T) {
	if 5-3 != 2 {
		t.Fatal("math is broken")
	}
}

func TestFail(t *testing.T) {
	t.Fatal("this test always fails")
}
