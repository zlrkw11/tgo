package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	packages    []Package
	events      []TestEvent
	eventCh     chan TestEvent
	cursor      int
	expanded    map[int]bool   // package index -> expanded
	showErrors  map[string]bool // "pkgIdx/testIdx" -> show errors
	done        bool
	err         error
	height      int
	offset      int
}

type eventMsg TestEvent
type doneMsg struct{}

func initialModel(eventCh chan TestEvent) model {
	return model{
		eventCh:    eventCh,
		expanded:   make(map[int]bool),
		showErrors: make(map[string]bool),
		height:     24,
	}
}

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

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.adjustScroll()
		return m, nil

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
			pkgIdx, testIdx := m.itemAtCursor()
			if testIdx >= 0 {
				// on a test row: toggle error output
				key := testKey(pkgIdx, testIdx)
				m.showErrors[key] = !m.showErrors[key]
			} else if pkgIdx >= 0 {
				// on a package row: toggle expand
				m.expanded[pkgIdx] = !m.expanded[pkgIdx]
			}
		}

		total := m.totalRows()
		if total > 0 && m.cursor >= total {
			m.cursor = total - 1
		}
	}

	m.adjustScroll()
	return m, nil
}

func testKey(pkgIdx, testIdx int) string {
	return string(rune('0'+pkgIdx)) + "/" + string(rune('0'+testIdx))
}

// itemAtCursor returns (pkgIdx, testIdx). testIdx is -1 if cursor is on a package row.
func (m model) itemAtCursor() (int, int) {
	row := 0
	for i, pkg := range m.packages {
		if row == m.cursor {
			return i, -1
		}
		row++
		if m.expanded[i] {
			for j := range pkg.Tests {
				if row == m.cursor {
					return i, j
				}
				row++
			}
		}
	}
	return -1, -1
}

func (m *model) adjustScroll() {
	viewable := m.viewableRows()
	if viewable < 1 {
		viewable = 1
	}

	if m.cursor < m.offset {
		m.offset = m.cursor
	}

	for m.displayLinesBetween(m.offset, m.cursor) >= viewable && m.offset < m.cursor {
		m.offset++
	}
}

func (m model) displayLinesBetween(startRow, endRow int) int {
	groups := m.rowGroupSizes()
	count := 0
	for i := startRow; i <= endRow && i < len(groups); i++ {
		count += groups[i]
	}
	return count
}

func (m model) rowGroupSizes() []int {
	var sizes []int

	for i, pkg := range m.packages {
		sizes = append(sizes, 1)

		if m.expanded[i] {
			for j, t := range pkg.Tests {
				lines := 1
				if m.showErrors[testKey(i, j)] && t.Status == "fail" {
					for _, o := range t.Output {
						if o != "" && len(o) > 1 {
							lines++
						}
					}
				}
				sizes = append(sizes, lines)
			}
		}
	}

	return sizes
}

func (m model) viewableRows() int {
	headerLines := 8
	footerLines := 6
	viewable := m.height - headerLines - footerLines
	if viewable < 3 {
		viewable = 3
	}
	return viewable
}

func (m model) totalRows() int {
	count := 0
	for i, pkg := range m.packages {
		count++
		if m.expanded[i] {
			count += len(pkg.Tests)
		}
	}
	return count
}

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
