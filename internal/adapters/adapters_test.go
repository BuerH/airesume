package adapters

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCodexAdapterScan(t *testing.T) {
	root := t.TempDir()
	sessionDir := filepath.Join(root, "sessions", "2026", "07", "07")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	sessionPath := filepath.Join(sessionDir, "rollout.jsonl")
	sessionJSONL := `{"timestamp":"2026-07-07T01:00:00Z","type":"session_meta","payload":{"id":"codex-session","timestamp":"2026-07-07T01:00:00Z","cwd":"/tmp/project"}}
{"timestamp":"2026-07-07T01:01:00Z","type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"hello codex"}]}}
`
	if err := os.WriteFile(sessionPath, []byte(sessionJSONL), 0o644); err != nil {
		t.Fatal(err)
	}
	indexJSONL := `{"id":"codex-session","thread_name":"Indexed title","updated_at":"2026-07-07T01:02:00Z"}`
	if err := os.WriteFile(filepath.Join(root, "session_index.jsonl"), []byte(indexJSONL), 0o644); err != nil {
		t.Fatal(err)
	}

	records, err := NewCodexAdapter(root).Scan(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	rec := records[0]
	if rec.Tool != "codex" || rec.SessionID != "codex-session" || rec.CWD != "/tmp/project" {
		t.Fatalf("unexpected record: %#v", rec)
	}
	if rec.Title != "Indexed title" {
		t.Fatalf("expected indexed title, got %q", rec.Title)
	}
	if rec.FirstUserMessage != "hello codex" {
		t.Fatalf("expected user message, got %q", rec.FirstUserMessage)
	}
	if rec.ResumeCommand != "cd '/tmp/project' && codex 'resume' 'codex-session'" {
		t.Fatalf("unexpected resume command: %q", rec.ResumeCommand)
	}
}

func TestClaudeAdapterScan(t *testing.T) {
	root := t.TempDir()
	projectDir := filepath.Join(root, "projects", "-tmp-project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatal(err)
	}
	sessionPath := filepath.Join(projectDir, "claude-session.jsonl")
	sessionJSONL := `{"type":"user","message":{"role":"user","content":"<ide_opened_file>noise</ide_opened_file>hello claude"},"timestamp":"2026-07-07T01:00:00Z","cwd":"/tmp/project","sessionId":"claude-session","gitBranch":"main","origin":{"kind":"human"}}
`
	if err := os.WriteFile(sessionPath, []byte(sessionJSONL), 0o644); err != nil {
		t.Fatal(err)
	}

	records, err := NewClaudeAdapter(root).Scan(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	rec := records[0]
	if rec.Tool != "claude" || rec.SessionID != "claude-session" || rec.CWD != "/tmp/project" {
		t.Fatalf("unexpected record: %#v", rec)
	}
	if rec.Branch != "main" {
		t.Fatalf("expected branch, got %q", rec.Branch)
	}
	if rec.FirstUserMessage != "hello claude" {
		t.Fatalf("expected cleaned user message, got %q", rec.FirstUserMessage)
	}
	if rec.ResumeCommand != "cd '/tmp/project' && claude '--resume' 'claude-session'" {
		t.Fatalf("unexpected resume command: %q", rec.ResumeCommand)
	}
}

func TestDecodeClaudeProjectDir(t *testing.T) {
	got := decodeClaudeProjectDir("-home-star-project")
	if got != "/home/star/project" {
		t.Fatalf("unexpected decoded dir: %q", got)
	}
}
