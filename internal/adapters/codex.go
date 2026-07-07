package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BuerH/airesume/internal/session"
)

type CodexAdapter struct {
	root string
}

type codexIndexEntry struct {
	ID        string `json:"id"`
	Thread    string `json:"thread_name"`
	UpdatedAt string `json:"updated_at"`
}

type codexIndexValue struct {
	title     string
	updatedAt time.Time
}

func NewCodexAdapter(root string) CodexAdapter {
	return CodexAdapter{root: root}
}

func (a CodexAdapter) Name() string {
	return "codex"
}

func (a CodexAdapter) Scan(context.Context) ([]session.Record, error) {
	if a.root == "" {
		return nil, nil
	}
	if _, err := os.Stat(a.root); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	index := a.readIndex()
	var records []session.Record
	sessionsRoot := filepath.Join(a.root, "sessions")
	err := filepath.WalkDir(sessionsRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() || !strings.HasSuffix(path, ".jsonl") {
			return nil
		}
		rec, ok := a.parseSession(path, index)
		if ok {
			records = append(records, rec)
		}
		return nil
	})
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return records, err
	}
	return records, nil
}

func (a CodexAdapter) readIndex() map[string]codexIndexValue {
	path := filepath.Join(a.root, "session_index.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	index := map[string]codexIndexValue{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry codexIndexEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if entry.ID == "" {
			continue
		}
		index[entry.ID] = codexIndexValue{
			title:     entry.Thread,
			updatedAt: parseTime(entry.UpdatedAt),
		}
	}
	return index
}

func (a CodexAdapter) parseSession(path string, index map[string]codexIndexValue) (session.Record, bool) {
	rec := session.Record{
		Tool:       "codex",
		SourceFile: path,
		UpdatedAt:  fileModTime(path),
	}

	_ = readJSONL(path, func(obj map[string]any) {
		rec.UpdatedAt = maxTime(rec.UpdatedAt, parseTime(stringField(obj, "timestamp")))
		if stringField(obj, "type") == "session_meta" {
			payload := mapField(obj, "payload")
			if rec.SessionID == "" {
				rec.SessionID = stringField(payload, "id")
			}
			if rec.CWD == "" {
				rec.CWD = cleanAbs(stringField(payload, "cwd"))
			}
			rec.UpdatedAt = maxTime(rec.UpdatedAt, parseTime(stringField(payload, "timestamp")))
			return
		}

		if stringField(obj, "type") != "response_item" {
			return
		}
		payload := mapField(obj, "payload")
		if stringField(payload, "type") != "message" || stringField(payload, "role") != "user" {
			return
		}
		text := cleanUserText(session.CleanText(textFromCodexContent(payload["content"])))
		if shouldSkipSyntheticUserText(text) {
			return
		}
		if rec.FirstUserMessage == "" {
			rec.FirstUserMessage = text
		}
		rec.LastUserMessage = text
	})

	if rec.SessionID == "" || rec.CWD == "" {
		return session.Record{}, false
	}
	if value, ok := index[rec.SessionID]; ok {
		rec.Title = value.title
		rec.UpdatedAt = maxTime(rec.UpdatedAt, value.updatedAt)
	}
	if rec.Title == "" {
		rec.Title = rec.FirstUserMessage
	}
	rec.ResumeCommand = buildResumeCommand(rec.CWD, "codex", "resume", rec.SessionID)
	return rec, true
}
