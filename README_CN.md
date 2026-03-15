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

## 为什么选 tgo？

<table>
<tr>
<td><h3>重跑失败测试</h3>按 <code>r</code> 即可重跑单个失败测试——不用重启，不用重新输命令。修改、重跑、循环。</td>
<td><h3>Watch 模式</h3><code>tgo --watch</code> 监听 <code>.go</code> 文件变更，保存后自动重跑测试。Go 版的 Jest 体验。</td>
<td><h3>耗时趋势</h3>tgo 跨运行追踪测试速度。如果测试慢了 1.5 倍，会标记 <code>▲ +45ms</code> 提醒你。</td>
</tr>
</table>

<br>

## 功能特性

- **实时流式输出** — 测试通过或失败时即时显示
- **交互式导航** — 用键盘浏览包和测试用例
- **可展开的包** — 点进任何包查看每个测试的详细结果
- **错误摘要** — 失败测试显示一行错误预览，按 Enter 查看完整详情
- **重跑失败测试** — 按 `r` 即可重跑选中的测试，无需重启
- **Watch 模式** — `.go` 文件变更时自动重跑测试
- **耗时趋势** — 跨运行记录测试耗时，标记变慢的测试
- **进度条** — 一眼看清通过/失败比例
- **精美 UI** — 圆角边框、彩色状态图标、高亮光标

<p align="center">
  <img src="doc/demo.gif" alt="tgo demo" width="800">
</p>

<br>

## 安装

**Homebrew**

```bash
brew install zlrkw11/tap/tgo
```

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

# Watch 模式 — 文件变更时自动重跑
tgo --watch ./...
tgo -w ./...
```

### 快捷键

| 按键 | 操作 |
|------|------|
| `↑` `k` | 向上移动 |
| `↓` `j` | 向下移动 |
| `g` `Home` | 跳到顶部 |
| `G` `End` | 跳到底部 |
| `Enter` `空格` | 展开 / 收起包或测试错误详情 |
| `r` | 重跑选中的测试 |
| `q` `Ctrl+C` | 退出 |

### 耗时趋势

tgo 自动将测试耗时记录到 `~/.config/tgo/history.json`。多跑几次后，变化明显的测试旁会显示趋势指标：

- `▲ +45ms` — 测试变慢（红色）
- `▼ -12ms` — 测试变快（绿色）

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
