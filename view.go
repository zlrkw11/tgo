package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	passStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	failStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	skipStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	dimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	hlStyle   = lipgloss.NewStyle().Background(lipgloss.Color("236"))
	boxStyle  = lipgloss.NewStyle().
			Border(lipgloss.MarkdownBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(0, 0)
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Align(lipgloss.Center).
			Width(56)
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

func (m model) View() string {
	logo := `
 ████████╗ ██████╗  ██████╗
 ╚══██╔══╝██╔════╝ ██╔═══██╗
    ██║   ██║  ███╗██║   ██║
    ██║   ██║   ██║██║   ██║
    ██║   ╚██████╔╝╚██████╔╝
    ╚═╝    ╚═════╝  ╚═════╝`
	s := titleStyle.Render(logo) + "\n"
	s += dimStyle.Align(lipgloss.Center).Foreground(lipgloss.Color("39")).Width(56).Render("TEST UI built for Go") + "\n\n"

	row := 0
	for i, pkg := range m.packages {
		passed := 0
		for _, t := range pkg.Tests {
			if t.Status == "pass" {
				passed++
			}
		}

		line := fmt.Sprintf("  %s %-40s %d/%d  %dms",
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
				testLine := fmt.Sprintf("    %s %-38s %dms",
					statusIcon(t.Status),
					t.Name,
					t.Duration.Milliseconds(),
				)
				if row == m.cursor {
					testLine = hlStyle.Render(testLine)
				}
				s += testLine + "\n"
				row++

				// show error output for failed tests
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

	s += "\n"
	if m.done {
		s += dimStyle.Render("  done. press q to quit")
	} else {
		s += dimStyle.Render("  running...")
	}
	s += "\n  [↑↓] navigate | [enter] expand/close | [q] quit"

	return "\n" + boxStyle.Render(s) + "\n"
}
