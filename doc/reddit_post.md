**Title:** I built an interactive TUI for `go test` — rerun failing tests, watch mode, and duration trends

I got tired of squinting at `go test` output, so I built **tgo** — an interactive terminal UI that makes running Go tests actually enjoyable.

![demo](https://raw.githubusercontent.com/zlrkw11/tgo/main/doc/demo.gif)

**What it does:**
- Real-time streaming of test results with pass/fail highlighting
- Press `r` on any failing test to rerun just that one — no restart needed
- `tgo --watch` auto-reruns tests when you save a `.go` file (like Jest for Go)
- Tracks test durations across runs and flags tests that got slower
- Expandable packages, inline error previews, progress bar

**Install:**
```
go install github.com/zlrkw11/tgo@latest
```

Then just run `tgo ./...` in any Go project.

Built with Bubble Tea + Lip Gloss. MIT licensed.

GitHub: https://github.com/zlrkw11/tgo

Would love feedback — what features would make this useful for your workflow?
