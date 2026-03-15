package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"./..."}
	}

	ch := make(chan TestEvent)

	err := RunTests(args, ch)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	model := initialModel(ch)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	_ = fmt.Println
	_ = os.Args
	_ = tea.NewProgram
}
