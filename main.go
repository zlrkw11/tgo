package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	args := os.Args[1:]
	watchMode := false

	// parse --watch flag
	var testArgs []string
	for _, a := range args {
		if a == "--watch" || a == "-w" {
			watchMode = true
		} else {
			testArgs = append(testArgs, a)
		}
	}
	if len(testArgs) == 0 {
		testArgs = []string{"./..."}
	}

	ch := make(chan TestEvent)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := RunTests(ctx, testArgs, ch)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// load history
	history, _ := LoadHistory()

	m := initialModel(ch, testArgs, watchMode, history)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
