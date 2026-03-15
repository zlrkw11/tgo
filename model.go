package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// model is the bubbletea Model that holds all TUI state.
type model struct {
	packages []Package
	events   []TestEvent
	eventCh  chan TestEvent
	cursor   int          // which row is selected
	expanded map[int]bool // which packages are expanded
	done     bool         // test run finished
	err      error
}

// eventMsg wraps a TestEvent for the bubbletea message loop.
type eventMsg TestEvent

// doneMsg signals that the test run is complete.
type doneMsg struct{}

func initialModel(eventCh chan TestEvent) model {
	return model{
		eventCh:  eventCh,
		expanded: make(map[int]bool),
	}
}

// waitForEvent returns a command that waits for the next test event.
func waitForEvent(ch chan TestEvent) tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return doneMsg{}
		}
		return eventMsg(ev)
	}
}

func (m model) Init() tea.Cmd {
	return waitForEvent(m.eventCh)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case eventMsg:
		m.events = append(m.events, TestEvent(msg))
		m.packages = ParseEvents(m.events)
		return m, waitForEvent(m.eventCh)

	case doneMsg:
		m.done = true
		m.packages = ParseEvents(m.events)
		return m, nil

	case tea.KeyMsg:
		totalRows := m.totalRows()
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < totalRows-1 {
				m.cursor++
			}
		case "enter", " ":
			// toggle expand on package rows
			pkgIdx := m.packageIndexAtCursor()
			if pkgIdx >= 0 {
				m.expanded[pkgIdx] = !m.expanded[pkgIdx]
			}
		}
	}

	return m, nil
}

// totalRows returns the total number of visible rows.
func (m model) totalRows() int {
	count := 0
	for i, pkg := range m.packages {
		count++ // package row
		if m.expanded[i] {
			count += len(pkg.Tests)
		}
	}
	return count
}

// packageIndexAtCursor returns the package index if cursor is on a package row, or -1 if on a test row.
func (m model) packageIndexAtCursor() int {
	row := 0
	for i, pkg := range m.packages {
		if row == m.cursor {
			return i
		}
		row++
		if m.expanded[i] {
			row += len(pkg.Tests)
		}
	}
	return -1
}
