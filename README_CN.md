<p align="center">
  <h1 align="center">⚡ TGO</h1>
  <p align="center">
    <strong>一个漂亮的、交互式的 <code>go test</code> 终端 UI</strong>
  </p>
  <p align="center">
    <a href="https://github.com/zlrkw11/tgo/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
    <a href="https://goreportcard.com/report/github.com/zlrkw11/tgo"><img src="https://goreportcard.com/badge/github.com/zlrkw11/tgo" alt="Go Report Card"></a>
    <a href="https://pkg.go.dev/github.com/zlrkw11/tgo"><img src="https://pkg.go.dev/badge/github.com/zlrkw11/tgo.svg" alt="Go Reference"></a>
  </p>
  <p align="center">
    <a href="README.md">English</a> | 中文
  </p>
</p>

<br>

> 别再盯着密密麻麻的 `go test` 输出了。**tgo** 给你一个实时的、可交互的终端界面，带通过/失败高亮、可展开的包、进度条和内联错误信息。

<br>

## 功能特性

- **实时流式输出** — 测试通过或失败时即时显示
- **交互式导航** — 用键盘浏览包和测试用例
- **可展开的包** — 点进任何包查看每个测试的详细结果
- **内联错误信息** — 失败的测试输出直接显示在对应位置
- **进度条** — 一眼看清通过/失败比例
- **精美 UI** — 圆角边框、彩色状态图标、高亮光标

```
╭──────────────────────────────────────────────────────────╮
│                         ⚡ TGO                            │
│                  TEST UI built for Go                    │
│                                                          │
│  3 tests  ████████████████░░░░  2 passed  1 failed       │
│  ────────────────────────────────────────────────────    │
│  ▸ ✕ github.com/you/pkg             2/3   120ms          │
│      ✓ TestAdd                             0ms           │
│      ✓ TestSubtract                        0ms           │
│      ✕ TestFail                            0ms           │
│          sample_test.go:18: this test always fails       │
│  ────────────────────────────────────────────────────    │
│  done                                                    │
│                                                          │
│  ↑↓ navigate  ⏎ expand  q quit                           │
╰──────────────────────────────────────────────────────────╯
```

<br>

## 安装

**Go 1.22+**

```bash
go install github.com/zlrkw11/tgo@latest
```

**从源码构建**

```bash
git clone https://github.com/zlrkw11/tgo.git
cd tgo
go build -o tgo .
```

<br>

## 使用方法

```bash
# 运行当前项目的所有测试
tgo ./...

# 运行指定包的测试
tgo ./pkg/auth/...

# 不带参数默认运行 ./...
tgo
```

### 快捷键

| 按键 | 操作 |
|------|------|
| `↑` `k` | 向上移动 |
| `↓` `j` | 向下移动 |
| `Enter` `空格` | 展开 / 收起包 |
| `q` `Ctrl+C` | 退出 |

<br>

## 工作原理

tgo 底层运行 `go test -json`，实时解析结构化的 JSON 输出。每个测试事件通过 channel 流入 [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI，渲染成可交互、可导航的测试结果视图。

```
go test -json ./...  →  解析事件流  →  Bubble Tea TUI
```

<br>

## 技术栈

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI 框架
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — 样式渲染

<br>

## 开源协议

[MIT](LICENSE)
