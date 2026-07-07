# airesume

Unified local session launcher for Codex and Claude Code.

`airesume` scans local session files, shows conversations by update time, groups them by directory, and prints or runs a copyable resume command.

## Features

- Read-only scan of `~/.codex` and `~/.claude`.
- Sessions sorted by latest update time.
- Directory view sorted by latest update time.
- Optional `fzf` picker for interactive selection.
- Colored table fallback selector when `fzf` is unavailable.
- Copyable resume commands:
  - `cd '<cwd>' && codex resume '<session_id>'`
  - `cd '<cwd>' && claude --resume '<session_id>'`
- Clipboard copy through `pbcopy`, `wl-copy`, `xclip`, `xsel`, `clip.exe`, or OSC52-capable terminals.
- Direct resume by ID, prefix, latest session, or numbered prompt.
- Adapter interface for future AI coding tools.
- Placeholder summary provider interface for future AI API summaries.

## Installation

Install the latest version with Go:

```bash
go install github.com/BuerH/airesume/cmd/airesume@latest
```

Or install a specific release:

```bash
go install github.com/BuerH/airesume/cmd/airesume@v0.1.0
```

Prebuilt binaries are attached to GitHub Releases for Linux, macOS, and Windows.

## Build From Source

```bash
go build -o bin/airesume ./cmd/airesume
```

If the directory is not a normal git checkout, use:

```bash
go build -buildvcs=false -o bin/airesume ./cmd/airesume
```

For release builds:

```bash
make release-snapshot VERSION=v0.1.0
```

## Usage

```bash
airesume
airesume --all
airesume dirs --all
airesume cmd --last
airesume cmd --last --copy
airesume resume --last
airesume resume <session_id>
airesume version
```

By default, `airesume` only shows sessions whose `cwd` matches the current directory. Use `--all` to show every known session, or `--cwd <path>` to filter another directory.

## Commands

```text
airesume [list] [--all] [--cwd PATH] [--limit N] [--json]
airesume dirs [--all] [--cwd PATH] [--json]
airesume cmd [SESSION_ID] [--all] [--cwd PATH] [--last] [--copy]
airesume resume [SESSION_ID] [--all] [--cwd PATH] [--last]
airesume summarize [SESSION_ID]
airesume version
```

`summarize` is intentionally a placeholder in the first version. The interface is in `internal/summary` and can later be backed by OpenAI, Anthropic, or a local model.

When `cmd` or `resume` needs you to choose a session, `airesume` uses `fzf` if it is installed and the current process is attached to a terminal. Without `fzf`, it falls back to a numbered prompt.

The numbered prompt is a fixed-width table with tool, update time, directory, branch, short session ID, and title. Colors are enabled automatically in terminals. Set `NO_COLOR=1` to disable colors or `FORCE_COLOR=1` to force ANSI color output.

Use `cmd --copy` to copy the selected resume command to the system clipboard and print it:

```bash
airesume cmd --all --copy
```

## Architecture

The core record is `session.Record`. Tools implement the adapter interface:

```go
type Adapter interface {
    Name() string
    Scan(ctx context.Context) ([]session.Record, error)
}
```

Current adapters:

- `internal/adapters/codex`
- `internal/adapters/claude`

Future adapters can be added to `DefaultRegistry()`.

## Release

This section is for maintainers.

The Go module path is:

```go
module github.com/BuerH/airesume
```

Create the GitHub repository and push the default branch first:

```bash
git remote add origin git@github.com:BuerH/airesume.git
git branch -M main
git push -u origin main
```

Run checks locally:

```bash
make check
make release-snapshot VERSION=v0.1.0
```

Publish a release by pushing a version tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow builds Linux, macOS, and Windows binaries, then uploads all artifacts and `checksums.txt` to the GitHub Release.
