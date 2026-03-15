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

	// package list
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
		s += line + "\n"
		row++

		if m.expanded[i] {
			for _, t := range pkg.Tests {
				testCursor := "    "
				if row == m.cursor {
					testCursor = "  ▸ "
				}

				testLine := fmt.Sprintf("%s%s %-34s %dms",
					testCursor,
					statusIcon(t.Status),
					t.Name,
					t.Duration.Milliseconds(),
				)
				if row == m.cursor {
					testLine = hlStyle.Render(testLine)
				}
				s += testLine + "\n"
				row++

				if t.Status == "fail" && len(t.Output) > 0 {
					for _, line := range t.Output {
						trimmed := strings.TrimRight(line, "\n")
						if trimmed != "" && !strings.HasPrefix(trimmed, "=== RUN") && !strings.HasPrefix(trimmed, "--- FAIL") {
							s += failStyle.Render("        "+trimmed) + "\n"
						}
					}
				}
			}
		}
	}

	// footer
	s += "  " + divider + "\n"
	if m.done {
		s += dimStyle.Render("  done") + "\n"
	} else {
		s += dimStyle.Render("  running...") + "\n"
	}
	s += "\n"
	s += "  " + pillStyle.Render("↑↓ navigate") + " " +
		pillStyle.Render("⏎ expand") + " " +
		pillStyle.Render("q quit")

	return "\n" + boxStyle.Render(s) + "\n"
}
