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

## Features

- **Real-time streaming** - watch tests pass and fail as they run
- **Interactive navigation** - browse packages and tests with keyboard controls
- **Expandable packages** - drill into any package to see individual test results
- **Inline errors** - failed test output shown right where you need it
- **Progress bar** - visual pass/fail ratio at a glance
- **Beautiful UI** - rounded borders, color-coded status icons, highlighted cursor

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

## Installation

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
```

### Keyboard Controls

| Key | Action |
|-----|--------|
| `↑` `k` | Move up |
| `↓` `j` | Move down |
| `Enter` `Space` | Expand / collapse package |
| `q` `Ctrl+C` | Quit |

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
