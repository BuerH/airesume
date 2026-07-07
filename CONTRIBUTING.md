# Contributing

Thanks for considering a contribution to `airesume`.

## Development

Requirements:

- Go 1.22 or newer

Common commands:

```bash
make test
make build
make check
```

If you are not using `make`, the equivalent commands are:

```bash
gofmt -w ./cmd ./internal
go test ./...
go build -o bin/airesume ./cmd/airesume
```

## Pull Requests

- Keep changes focused.
- Add or update tests for parser, filtering, sorting, and display behavior.
- Do not write to `~/.codex` or `~/.claude`; adapters must remain read-only.
- Do not commit generated binaries, local caches, or private session data.

## Adding An Adapter

Adapters implement:

```go
type Adapter interface {
    Name() string
    Scan(ctx context.Context) ([]session.Record, error)
}
```

Add the adapter under `internal/adapters`, then register it in `DefaultRegistry`.
