package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// suppress unused imports — remove these when you use them in RerunTest
var _ = fmt.Sprintf
var _ = strings.ReplaceAll

// RunTests executes `go test -json` and sends events to the channel.
func RunTests(ctx context.Context, args []string, events chan<- TestEvent) error {
	cmdArgs := append([]string{"test", "-json"}, args...)
	cmd := exec.CommandContext(ctx, "go", cmdArgs...)

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

func RerunTest(pkgName, testName string) (TestResult, error) {
	runPattern := "^" + strings.ReplaceAll(testName, "/", "$/^") + "$"
	cmd := exec.Command("go", "test", "-json", "-run", runPattern, pkgName)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return TestResult{}, err
	}

	if err := cmd.Start(); err != nil {
		return TestResult{}, err
	}

	var result TestResult
	result.Name = testName
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		var ev TestEvent
		if json.Unmarshal(scanner.Bytes(), &ev) != nil {
			continue
		}

		if ev.Test != testName {
			continue
		}
		switch ev.Action {
		case "output":
			result.Output = append(result.Output, ev.Output)
		case "pass", "fail", "skip":
			result.Status = ev.Action
			result.Duration = time.Duration(ev.Elapsed * float64(time.Second))
		}
	}

	cmd.Wait()

	if result.Status == "" {
		return result, fmt.Errorf("test %s not found in output", testName)
	}
	return result, nil
}

// ParseEvents processes a stream of TestEvents into Packages.
func ParseEvents(events []TestEvent) []Package {
	pkgMap := make(map[string]*Package)
	testMap := make(map[string]*TestResult)

	for _, ev := range events {
		pkg, ok := pkgMap[ev.Package]
		if !ok {
			pkg = &Package{Name: ev.Package}
			pkgMap[ev.Package] = pkg
		}

		if ev.Test == "" {
			switch ev.Action {
			case "output":
				pkg.Output = append(pkg.Output, ev.Output)
			case "pass", "fail", "skip":
				pkg.Status = ev.Action
				pkg.Duration = time.Duration(ev.Elapsed * float64(time.Second))
			}
			continue
		}

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

		for i := range pkg.Tests {
			if pkg.Tests[i].Name == ev.Test {
				pkg.Tests[i] = *tr
			}
		}
	}

	var packages []Package
	for _, pkg := range pkgMap {
		packages = append(packages, *pkg)
	}
	return packages
}
