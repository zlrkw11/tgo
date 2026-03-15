<p align="center">
  <h1 align="center">⚡ TGO</h1>
  <p align="center">
    <strong>A beautiful, interactive TUI for <code>go test</code></strong>
  </p>
  <p align="center">
    <a href="https://github.com/zlrkw11/tgo/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
    <a href="https://goreportcard.com/report/github.com/zlrkw11/tgo"><img src="https://goreportcard.com/badge/github.com/zlrkw11/tgo" alt="Go Report Card"></a>
    <a href="https://pkg.go.dev/github.com/zlrkw11/tgo"><img src="https://pkg.go.dev/badge/github.com/zlrkw11/tgo.svg" alt="Go Reference"></a>
  </p>
  <p align="center">
    English | <a href="README_CN.md">中文</a>
  </p>
</p>

<br>

> Stop squinting at raw `go test` output. **tgo** gives you a real-time, interactive terminal UI with pass/fail highlighting, expandable packages, progress bars, and inline error messages.

<br>

## Why tgo?

<table>
<tr>
<td><h3>Rerun failed tests</h3>Press <code>r</code> to rerun a single failing test — no restart, no re-typing commands. Fix, rerun, repeat.</td>
<td><h3>Watch mode</h3><code>tgo --watch</code> monitors your <code>.go</code> files and reruns tests automatically on save. Like Jest, but for Go.</td>
<td><h3>Duration trends</h3>tgo tracks test speed across runs. If a test gets 1.5x slower, you'll see it flagged with <code>▲ +45ms</code>.</td>
</tr>
</table>

<br>

## Features

- **Real-time streaming** - watch tests pass and fail as they run
- **Interactive navigation** - browse packages and tests with keyboard controls
- **Expandable packages** - drill into any package to see individual test results
- **Error summaries** - failed tests show a one-line error preview, press Enter for full details
- **Rerun failed tests** - press `r` on any failed test to rerun it instantly without restarting
- **Watch mode** - auto-reruns tests when `.go` files change
- **Duration trends** - tracks test speed across runs, flags tests that got slower
- **Progress bar** - visual pass/fail ratio at a glance
- **Beautiful UI** - rounded borders, color-coded status icons, highlighted cursor

<p align="center">
  <img src="doc/demo.gif" alt="tgo demo" width="800">
</p>
<br>

## Installation

**Homebrew**

```bash
brew install zlrkw11/tap/tgo
```

**Go 1.22+**

```bash
go install github.com/zlrkw11/tgo@latest
```

**From source**

```bash
git clone https://github.com/zlrkw11/tgo.git
cd tgo
go build -o tgo .
```

<br>

## Usage

```bash
# Run all tests in the current project
tgo ./...

# Run tests in a specific package
tgo ./pkg/auth/...

# Run with default (./...) if no args given
tgo

# Watch mode — auto-rerun on file changes
tgo --watch ./...
tgo -w ./...
```

### Keyboard Controls

| Key | Action |
|-----|--------|
| `↑` `k` | Move up |
| `↓` `j` | Move down |
| `g` `Home` | Jump to top |
| `G` `End` | Jump to bottom |
| `Enter` `Space` | Expand / collapse package or test errors |
| `r` | Rerun the selected test |
| `q` `Ctrl+C` | Quit |

### Duration Trends

tgo automatically records test durations to `~/.config/tgo/history.json`. After a few runs, you'll see trend indicators next to tests that got significantly slower or faster:

- `▲ +45ms` — test got slower (red)
- `▼ -12ms` — test got faster (green)

<br>

## How It Works

tgo runs `go test -json` under the hood and parses the structured JSON output in real time. Each test event is streamed into a [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI that renders an interactive, navigable view of your test results.

```
go test -json ./...  →  parse events  →  Bubble Tea TUI
```

<br>

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions

<br>

## License

[MIT](LICENSE)
