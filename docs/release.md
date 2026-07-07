# Release Guide

This document is for maintainers.

## Module Path

The Go module path is:

```go
module github.com/BuerH/airesume
```

## Local Checks

```bash
make check
make release-snapshot VERSION=v0.1.0
```

`make release-snapshot` writes release assets to `dist/`:

- `airesume-linux-amd64`
- `airesume-linux-arm64`
- `airesume-darwin-amd64`
- `airesume-darwin-arm64`
- `airesume-windows-amd64.exe`
- `checksums.txt`

## Publish

Commit all changes:

```bash
git add .
git commit -m "Prepare release v0.1.0"
git push origin master
```

Create and push the tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow builds Linux, macOS, and Windows binaries, then uploads all artifacts and `checksums.txt` to the GitHub Release.

If a tag was pushed before the release workflow existed, delete and recreate the tag:

```bash
git tag -d v0.1.0
git push origin :refs/tags/v0.1.0
git tag v0.1.0
git push origin v0.1.0
```
