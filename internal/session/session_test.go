package session

import (
	"testing"
	"time"
)

func TestSortByUpdated(t *testing.T) {
	oldTime := time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC)
	newTime := time.Date(2026, 1, 2, 1, 0, 0, 0, time.UTC)
	records := []Record{
		{SessionID: "old", UpdatedAt: oldTime},
		{SessionID: "new", UpdatedAt: newTime},
	}

	SortByUpdated(records)

	if records[0].SessionID != "new" {
		t.Fatalf("expected newest session first, got %q", records[0].SessionID)
	}
}

func TestFilterByCWD(t *testing.T) {
	records := []Record{
		{SessionID: "keep", CWD: "/tmp/project"},
		{SessionID: "drop", CWD: "/tmp/other"},
	}

	filtered := FilterByCWD(records, "/tmp/project")

	if len(filtered) != 1 || filtered[0].SessionID != "keep" {
		t.Fatalf("unexpected filtered records: %#v", filtered)
	}
}

func TestGroupByDirectory(t *testing.T) {
	oldTime := time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC)
	newTime := time.Date(2026, 1, 2, 1, 0, 0, 0, time.UTC)
	records := []Record{
		{Tool: "codex", CWD: "/tmp/project", UpdatedAt: oldTime, Title: "old"},
		{Tool: "claude", CWD: "/tmp/project", UpdatedAt: newTime, Title: "new"},
		{Tool: "codex", CWD: "/tmp/other", UpdatedAt: oldTime, Title: "other"},
	}

	groups := GroupByDirectory(records)

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].CWD != "/tmp/project" {
		t.Fatalf("expected newest directory first, got %q", groups[0].CWD)
	}
	if groups[0].Count != 2 {
		t.Fatalf("expected project count 2, got %d", groups[0].Count)
	}
	if groups[0].LatestTitle != "new" {
		t.Fatalf("expected latest title, got %q", groups[0].LatestTitle)
	}
	if len(groups[0].Tools) != 2 || groups[0].Tools[0] != "claude" || groups[0].Tools[1] != "codex" {
		t.Fatalf("expected sorted tools, got %#v", groups[0].Tools)
	}
}
