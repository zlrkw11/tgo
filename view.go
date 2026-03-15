package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	passStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	failStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	skipStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	dimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	hlStyle   = lipgloss.NewStyle().Background(lipgloss.Color("236"))
	boxStyle  = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Align(lipgloss.Center).
			Width(56)
	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Align(lipgloss.Center).
			Width(56)
	pillStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("237")).
			Padding(0, 1)
	divider = dimStyle.Render(strings.Repeat("─", 52))
)

func statusIcon(status string) string {
	switch status {
	case "pass":
		return passStyle.Render("✓")
	case "fail":
		return failStyle.Render("✕")
	case "skip":
		return skipStyle.Render("⏭")
	default:
		return dimStyle.Render("⏳")
	}
}

func progressBar(passed, failed, total int, width int) string {
	if total == 0 {
		return dimStyle.Render(strings.Repeat("░", width))
	}
	passWidth := passed * width / total
	failWidth := failed * width / total
	remaining := width - passWidth - failWidth

	bar := passStyle.Render(strings.Repeat("█", passWidth)) +
		failStyle.Render(strings.Repeat("█", failWidth)) +
		dimStyle.Render(strings.Repeat("░", remaining))
	return bar
}

// firstErrorLine extracts the first meaningful error message from test output.
func firstErrorLine(output []string) string {
	for _, l := range output {
		trimmed := strings.TrimRight(l, "\n")
		trimmed = strings.TrimSpace(trimmed)
		if trimmed == "" || strings.HasPrefix(trimmed, "=== RUN") || strings.HasPrefix(trimmed, "--- FAIL") {
			continue
		}
		// truncate if too long
		if len(trimmed) > 50 {
			trimmed = trimmed[:47] + "..."
		}
		return trimmed
	}
	return ""
}

// rowGroup is a navigable row + any extra display lines (e.g. error output).
type rowGroup struct {
	lines []string
}

// buildRowGroups collects all visible content grouped by navigable row.
func (m model) buildRowGroups() []rowGroup {
	var groups []rowGroup

	row := 0
	for i, pkg := range m.packages {
		passed := 0
		for _, t := range pkg.Tests {
			if t.Status == "pass" {
				passed++
			}
		}

		cursor := "  "
		if row == m.cursor {
			cursor = "▸ "
		}

		line := fmt.Sprintf("%s%s %-36s %d/%d  %dms",
			cursor,
			statusIcon(pkg.Status),
			pkg.Name,
			passed,
			len(pkg.Tests),
			pkg.Duration.Milliseconds(),
		)
		if row == m.cursor {
			line = hlStyle.Render(line)
		}
		groups = append(groups, rowGroup{lines: []string{line}})
		row++

		if m.expanded[i] {
			for j, t := range pkg.Tests {
				testCursor := "    "
				if row == m.cursor {
					testCursor = "  ▸ "
				}

				// show spinner if this test is being rerun
				rerunKey := pkg.Name + "/" + t.Name
				icon := statusIcon(t.Status)
				if m.rerunning[rerunKey] {
					icon = skipStyle.Render("⟳")
				}

				testLine := fmt.Sprintf("%s%s %-34s %dms",
					testCursor,
					icon,
					t.Name,
					t.Duration.Milliseconds(),
				)

				// append short error summary for failed tests
				if t.Status == "fail" && !m.showErrors[testKey(i, j)] {
					if msg := firstErrorLine(t.Output); msg != "" {
						testLine += "  " + dimStyle.Render(msg)
					}
				}

				if row == m.cursor {
					testLine = hlStyle.Render(testLine)
				}

				g := rowGroup{lines: []string{testLine}}

				// only show errors when user explicitly toggles this test
				if m.showErrors[testKey(i, j)] && t.Status == "fail" && len(t.Output) > 0 {
					for _, l := range t.Output {
						trimmed := strings.TrimRight(l, "\n")
						if trimmed != "" && !strings.HasPrefix(trimmed, "=== RUN") && !strings.HasPrefix(trimmed, "--- FAIL") {
							g.lines = append(g.lines, failStyle.Render("        "+trimmed))
						}
					}
				}

				groups = append(groups, g)
				row++
			}
		}
	}

	return groups
}

func (m model) View() string {
	s := titleStyle.Render("⚡ TGO") + "\n"
	s += subtitleStyle.Render("TEST UI built for Go") + "\n\n"

	// summary stats
	totalTests := 0
	totalPassed := 0
	totalFailed := 0
	for _, pkg := range m.packages {
		for _, t := range pkg.Tests {
			totalTests++
			switch t.Status {
			case "pass":
				totalPassed++
			case "fail":
				totalFailed++
			}
		}
	}

	summary := fmt.Sprintf("  %d tests  ", totalTests)
	summary += progressBar(totalPassed, totalFailed, totalTests, 20)
	summary += "  " + passStyle.Render(fmt.Sprintf("%d passed", totalPassed))
	if totalFailed > 0 {
		summary += "  " + failStyle.Render(fmt.Sprintf("%d failed", totalFailed))
	}
	s += summary + "\n"
	s += "  " + divider + "\n"

	// scrollable package/test list
	groups := m.buildRowGroups()
	viewable := m.viewableRows()
	end := m.offset + viewable
	if end > len(groups) {
		end = len(groups)
	}

	if m.offset > 0 {
		s += dimStyle.Render(fmt.Sprintf("  ↑ %d more above", m.offset)) + "\n"
	}

	for _, g := range groups[m.offset:end] {
		for _, line := range g.lines {
			s += line + "\n"
		}
	}

	remaining := len(groups) - end
	if remaining > 0 {
		s += dimStyle.Render(fmt.Sprintf("  ↓ %d more below", remaining)) + "\n"
	}

	// footer
	s += "  " + divider + "\n"
	if m.rerunErr != "" {
		s += failStyle.Render("  rerun error: "+m.rerunErr) + "\n"
	} else if m.done {
		s += dimStyle.Render("  done") + "\n"
	} else {
		s += dimStyle.Render("  running...") + "\n"
	}
	s += "\n"
	s += "  " + pillStyle.Render("↑↓ navigate") + " " +
		pillStyle.Render("g/G top/bottom") + " " +
		pillStyle.Render("⏎ expand") + " " +
		pillStyle.Render("r rerun") + " " +
		pillStyle.Render("q quit")

	rendered := boxStyle.Render(s)

	// place at top of terminal
	return lipgloss.Place(0, m.height, lipgloss.Left, lipgloss.Top, rendered)
}
