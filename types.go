package main

import "time"

// TestEvent is a single line from `go test -json` output.
// See: https://pkg.go.dev/cmd/test2json#hdr-Output_Format
type TestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`  // run, pass, fail, output, skip, pause, cont
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Output  string    `json:"Output"`
	Elapsed float64   `json:"Elapsed"` // seconds
}

// TestResult holds the final state of a single test.
type TestResult struct {
	Name    string
	Status  string        // pass, fail, skip
	Duration time.Duration
	Output  []string      // captured output lines
}

// Package holds all test results for one package.
type Package struct {
	Name     string
	Status   string        // pass, fail, skip
	Duration time.Duration
	Tests    []TestResult
	Output   []string      // package-level output
}
