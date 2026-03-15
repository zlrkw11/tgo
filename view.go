package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// 样式定义 — 你可以调整颜色
var (
	passStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))  // 绿色
	failStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // 红色
	skipStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // 黄色
	dimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("245")) // 灰色
)

// statusIcon 根据状态返回图标
// STEP 1: 实现这个函数
// 提示:
//   "pass" → passStyle.Render("✅")
//   "fail" → failStyle.Render("❌")
//   "skip" → skipStyle.Render("⏭")
//   其他   → dimStyle.Render("⏳")  (还在跑)
func statusIcon(status string) string {
	// TODO: your code here
	return ""
}

// View renders the entire TUI. This is called by bubbletea on every update.
func (m model) View() string {
	s := ""

	// STEP 2: 标题
	// 提示: s += "\n  tgo — test runner\n\n"

	// STEP 3: 遍历 m.packages，渲染每个包
	// 对于每个 package:
	//   row := 0  (用来追踪当前行号，和 m.cursor 比较来高亮)
	//   for i, pkg := range m.packages {
	//
	//     // 统计通过的测试数
	//     passed := 0
	//     for _, t := range pkg.Tests {
	//         if t.Status == "pass" { passed++ }
	//     }
	//
	//     // 渲染包那一行: icon + 包名 + passed/total + 耗时
	//     // 例如: "✅ pkg/auth  4/4  120ms"
	//     line := fmt.Sprintf("  %s %-40s %d/%d  %dms",
	//         statusIcon(pkg.Status),
	//         pkg.Name,
	//         passed,
	//         len(pkg.Tests),
	//         pkg.Duration.Milliseconds(),
	//     )
	//
	//     // 如果 row == m.cursor，加高亮背景
	//     // 提示: lipgloss.NewStyle().Background(lipgloss.Color("236")).Render(line)
	//
	//     s += line + "\n"
	//     row++
	//
	//     // 如果这个包被展开了 (m.expanded[i] == true)，渲染每个测试
	//     if m.expanded[i] {
	//         for _, t := range pkg.Tests {
	//             testLine := fmt.Sprintf("    %s %s  %dms",
	//                 statusIcon(t.Status),
	//                 t.Name,
	//                 t.Duration.Milliseconds(),
	//             )
	//             // 同样检查 row == m.cursor 来高亮
	//             s += testLine + "\n"
	//             row++
	//         }
	//     }
	//   }

	// STEP 4: 底部状态栏
	// 如果 m.done: s += dimStyle.Render("  done. press q to quit")
	// 否则:       s += dimStyle.Render("  running...")
	// 最后加:     s += "\n  ↑↓ navigate  enter expand  q quit\n"

	_ = fmt.Sprintf
	_ = lipgloss.NewStyle
	_ = passStyle
	_ = failStyle
	_ = skipStyle
	_ = dimStyle

	return s
}
