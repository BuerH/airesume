# airesume

`airesume` is a small command-line session picker for AI coding tools.

If you use both Codex and Claude Code, it can be hard to remember which conversation to resume. `airesume` scans your local session history, shows recent conversations in one place, and prints or runs the correct resume command.

Currently supported:

- Codex CLI
- Claude Code

`airesume` is read-only. It reads local session files from `~/.codex` and `~/.claude` and does not upload your history anywhere.

## Install

With Go:

```bash
go install github.com/BuerH/airesume/cmd/airesume@latest
```

Or download a prebuilt binary from the GitHub Releases page.

Optional tools:

- `fzf` for interactive fuzzy selection.
- `pbcopy`, `wl-copy`, `xclip`, `xsel`, `clip.exe`, or an OSC52-capable terminal for clipboard copy.

## Quick Start

Show sessions for the current directory:

```bash
airesume
```

Show sessions from all directories:

```bash
airesume --all
```

Pick a session and print the command to resume it:

```bash
airesume cmd --all
```

Resume a session directly:

```bash
airesume resume --all
```

Resume the latest matching session:

```bash
airesume resume --last
```

Copy the resume command:

```bash
airesume cmd --all --copy
```

## What It Looks Like

When `fzf` is installed, `airesume cmd` and `airesume resume` open a fuzzy picker.

Without `fzf`, `airesume` falls back to a fixed-width table:

```text
  # TOOL    UPDATED          DIR                                    BRANCH           SESSION  TITLE
--- ------- ---------------- -------------------------------------- ---------------- -------- --------------------
 1. codex   2026-07-07 11:04 ~/AI_resume_tree                       -                019f227e Fix session picker display
 2. claude  2026-07-06 18:41 ~/project                              main             92851d15 Write a short commit message
```

Colors are enabled automatically in terminals.

Disable colors:

```bash
NO_COLOR=1 airesume cmd --all
```

Force colors:

```bash
FORCE_COLOR=1 airesume cmd --all
```

## Common Workflows

List directories by recent activity:

```bash
airesume dirs --all
```

Filter another project directory:

```bash
airesume --cwd ~/my-project
airesume resume --cwd ~/my-project
```

Print JSON for scripting:

```bash
airesume --all --json
airesume dirs --all --json
```

Resume by session ID or prefix:

```bash
airesume resume 019f227e
airesume cmd 92851d15
```

## Commands

```text
airesume [list] [--all] [--cwd PATH] [--limit N] [--json]
airesume dirs [--all] [--cwd PATH] [--json]
airesume cmd [SESSION_ID] [--all] [--cwd PATH] [--last] [--copy]
airesume resume [SESSION_ID] [--all] [--cwd PATH] [--last]
airesume summarize [SESSION_ID]
airesume version
```

Default behavior:

- `airesume` only shows sessions whose recorded working directory matches your current directory.
- `--all` shows sessions from every known directory.
- `--cwd PATH` filters sessions for a specific directory.

## Privacy

`airesume` only reads local files:

- Codex: `~/.codex`
- Claude Code: `~/.claude`

It does not modify those files. It does not send session content to any network service.

The `summarize` command is a placeholder for future AI summaries and is not implemented yet. Any future summary provider should be opt-in.

## Development

Build from source:

```bash
go build -o bin/airesume ./cmd/airesume
```

Run tests:

```bash
go test ./...
```

Run the local checks used by CI:

```bash
make check
```

## Contributing

Issues and pull requests are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT. See [LICENSE](LICENSE).
