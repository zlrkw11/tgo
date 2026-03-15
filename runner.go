package main

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"time"
)

// RunTests executes `go test -json` and sends events to the channel.
// The channel is closed when the command finishes.
func RunTests(args []string, events chan<- TestEvent) error {
	cmdArgs := append([]string{"test", "-json"}, args...)
	cmd := exec.Command("go", cmdArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		defer close(events)
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var ev TestEvent
			if json.Unmarshal(scanner.Bytes(), &ev) == nil {
				events <- ev
			}
		}
		cmd.Wait()
	}()

	return nil
}

// ParseEvents processes a stream of TestEvents into Packages.
func ParseEvents(events []TestEvent) []Package {
	pkgMap := make(map[string]*Package)
	testMap := make(map[string]*TestResult) // key: "pkg/TestName"

	for _, ev := range events {
		// ensure package exists
		pkg, ok := pkgMap[ev.Package]
		if !ok {
			pkg = &Package{Name: ev.Package}
			pkgMap[ev.Package] = pkg
		}

		if ev.Test == "" {
			// package-level event
			switch ev.Action {
			case "output":
				pkg.Output = append(pkg.Output, ev.Output)
			case "pass", "fail", "skip":
				pkg.Status = ev.Action
				pkg.Duration = time.Duration(ev.Elapsed * float64(time.Second))
			}
			continue
		}

		// test-level event
		key := ev.Package + "/" + ev.Test
		tr, ok := testMap[key]
		if !ok {
			tr = &TestResult{Name: ev.Test}
			testMap[key] = tr
			pkg.Tests = append(pkg.Tests, *tr)
		}

		switch ev.Action {
		case "output":
			tr.Output = append(tr.Output, ev.Output)
		case "pass", "fail", "skip":
			tr.Status = ev.Action
			tr.Duration = time.Duration(ev.Elapsed * float64(time.Second))
		}

		// update the test in the package slice
		for i := range pkg.Tests {
			if pkg.Tests[i].Name == ev.Test {
				pkg.Tests[i] = *tr
			}
		}
	}

	// collect into slice
	var packages []Package
	for _, pkg := range pkgMap {
		packages = append(packages, *pkg)
	}
	return packages
}
