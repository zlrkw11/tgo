package main

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	packages   []Package
	events     []TestEvent
	eventCh    chan TestEvent
	cursor     int
	expanded   map[int]bool
	showErrors map[string]bool
	done       bool
	err        error
	height     int
	offset     int

	// rerun
	rerunning map[string]bool // "pkg/test" -> rerunning

	// watch mode
	watchMode bool
	testArgs  []string
	watchCh   chan fileChangedMsg
	runCount  int

	// trends
	history TestHistory
	trends  map[string]Trend // "pkg/test" -> trend
}

type eventMsg TestEvent
type doneMsg struct{}

type rerunResultMsg struct {
	pkgName  string
	testName string
	result   TestResult
	err      error
}

type historySavedMsg struct{}

func initialModel(eventCh chan TestEvent, testArgs []string, watchMode bool, history TestHistory) model {
	m := model{
		eventCh:    eventCh,
		testArgs:   testArgs,
		watchMode:  watchMode,
		expanded:   make(map[int]bool),
		showErrors: make(map[string]bool),
		rerunning:  make(map[string]bool),
		history:    history,
		trends:     make(map[string]Trend),
		height:     24,
		runCount:   1,
	}
	if watchMode {
		m.watchCh = make(chan fileChangedMsg, 1)
	}
	return m
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

func waitForFileChange(ch chan fileChangedMsg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

func rerunTestCmd(pkgName, testName string) tea.Cmd {
	return func() tea.Msg {
		result, err := RerunTest(pkgName, testName)
		return rerunResultMsg{
			pkgName:  pkgName,
			testName: testName,
			result:   result,
			err:      err,
		}
	}
}

func saveHistoryCmd(h TestHistory) tea.Cmd {
	return func() tea.Msg {
		SaveHistory(h)
		return historySavedMsg{}
	}
}

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{waitForEvent(m.eventCh)}
	if m.watchMode && m.watchCh != nil {
		cmds = append(cmds, waitForFileChange(m.watchCh))
	}
	return tea.Batch(cmds...)
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
		// record history and compute trends
		m.history = RecordRun(m.history, m.packages)
		m.computeTrends()
		cmds := []tea.Cmd{saveHistoryCmd(m.history)}
		// in watch mode, keep listening for file changes
		if m.watchMode && m.watchCh != nil {
			cmds = append(cmds, waitForFileChange(m.watchCh))
		}
		return m, tea.Batch(cmds...)

	case fileChangedMsg:
		// restart test run
		m.done = false
		m.events = nil
		m.packages = nil
		m.expanded = make(map[int]bool)
		m.showErrors = make(map[string]bool)
		m.rerunning = make(map[string]bool)
		m.trends = make(map[string]Trend)
		m.cursor = 0
		m.offset = 0
		m.runCount++
		m.eventCh = make(chan TestEvent)
		ctx := context.Background()
		if err := RunTests(ctx, m.testArgs, m.eventCh); err != nil {
			m.err = err
			return m, nil
		}
		return m, waitForEvent(m.eventCh)

	case rerunResultMsg:
		key := msg.pkgName + "/" + msg.testName
		delete(m.rerunning, key)
		if msg.err == nil {
			m.updateTestResult(msg.pkgName, msg.testName, msg.result)
		}
		return m, nil

	case historySavedMsg:
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
				key := testKey(pkgIdx, testIdx)
				m.showErrors[key] = !m.showErrors[key]
			} else if pkgIdx >= 0 {
				m.expanded[pkgIdx] = !m.expanded[pkgIdx]
			}
		case "r":
			pkgIdx, testIdx := m.itemAtCursor()
			if pkgIdx >= 0 && testIdx >= 0 && testIdx < len(m.packages[pkgIdx].Tests) {
				pkg := m.packages[pkgIdx]
				t := pkg.Tests[testIdx]
				key := pkg.Name + "/" + t.Name
				if !m.rerunning[key] {
					m.rerunning[key] = true
					m.adjustScroll()
					return m, rerunTestCmd(pkg.Name, t.Name)
				}
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

// ---------------------------------------------------------------------------
// TODO: 实现 updateTestResult — 把 rerun 的结果更新到 model 里
// ---------------------------------------------------------------------------
//  1. 遍历 m.packages，找到 Name == pkgName 的包
//  2. 在那个包里遍历 Tests，找到 Name == testName 的测试
//  3. 用 result 替换那个测试: m.packages[i].Tests[j] = result
//  4. 重新计算包的状态:
//     - 遍历所有 tests，如果没有 fail 的 → pkg.Status = "pass"
//     - 否则 pkg.Status = "fail"
func (m *model) updateTestResult(pkgName, testName string, result TestResult) {
	allPass := true
	for i := range m.packages {
		if m.packages[i].Name != pkgName {
			continue
		}
		for j := range m.packages[i].Tests {
			if m.packages[i].Tests[j].Name != testName {
				continue
			}
			m.packages[i].Tests[j] = result
		}
		for _, t := range m.packages[i].Tests {
			if t.Status == "fail" {
				allPass = false
				break
			}
		}
		if allPass {
			m.packages[i].Status = "pass"
		}
		if !allPass{
			m.packages[i].Status = "fail"
		}
		return
	}


}

func (m *model) computeTrends() {
	m.trends = make(map[string]Trend)
	for _, pkg := range m.packages {
		for _, t := range pkg.Tests {
			trend := GetTrend(m.history, pkg.Name, t.Name, t.Duration)
			if trend.Direction != "" {
				key := pkg.Name + "/" + t.Name
				m.trends[key] = trend
			}
		}
	}
}

func testKey(pkgIdx, testIdx int) string {
	return fmt.Sprintf("%d/%d", pkgIdx, testIdx)
}

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
