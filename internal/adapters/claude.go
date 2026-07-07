package adapters

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/BuerH/airesume/internal/session"
)

type ClaudeAdapter struct {
	root string
}

func NewClaudeAdapter(root string) ClaudeAdapter {
	return ClaudeAdapter{root: root}
}

func (a ClaudeAdapter) Name() string {
	return "claude"
}

func (a ClaudeAdapter) Scan(context.Context) ([]session.Record, error) {
	if a.root == "" {
		return nil, nil
	}
	if _, err := os.Stat(a.root); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var records []session.Record
	projectsRoot := filepath.Join(a.root, "projects")
	err := filepath.WalkDir(projectsRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() || !strings.HasSuffix(path, ".jsonl") {
			return nil
		}
		rec, ok := a.parseSession(path)
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

func (a ClaudeAdapter) parseSession(path string) (session.Record, bool) {
	rec := session.Record{
		Tool:       "claude",
		SessionID:  strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		SourceFile: path,
		UpdatedAt:  fileModTime(path),
	}

	_ = readJSONL(path, func(obj map[string]any) {
		rec.UpdatedAt = maxTime(rec.UpdatedAt, parseTime(stringField(obj, "timestamp")))
		rec.UpdatedAt = maxTime(rec.UpdatedAt, timestampMillis(obj["timestamp"]))
		if rec.CWD == "" {
			rec.CWD = cleanAbs(stringField(obj, "cwd"))
		}
		if rec.SessionID == "" {
			rec.SessionID = stringField(obj, "sessionId")
		}
		if rec.Branch == "" {
			rec.Branch = stringField(obj, "gitBranch")
		}

		if stringField(obj, "type") != "user" {
			return
		}
		origin := mapField(obj, "origin")
		if kind := stringField(origin, "kind"); kind != "" && kind != "human" {
			return
		}
		message := mapField(obj, "message")
		if stringField(message, "role") != "user" {
			return
		}
		text := cleanUserText(session.CleanText(textFromClaudeContent(message["content"])))
		if shouldSkipSyntheticUserText(text) {
			return
		}
		if rec.FirstUserMessage == "" {
			rec.FirstUserMessage = text
		}
		rec.LastUserMessage = text
	})

	if rec.CWD == "" {
		projectDir := filepath.Base(filepath.Dir(path))
		rec.CWD = decodeClaudeProjectDir(projectDir)
	}
	if rec.SessionID == "" || rec.CWD == "" {
		return session.Record{}, false
	}
	rec.Title = rec.FirstUserMessage
	rec.ResumeCommand = buildResumeCommand(rec.CWD, "claude", "--resume", rec.SessionID)
	return rec, true
}
