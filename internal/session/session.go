package session

import (
	"encoding/json"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Record struct {
	Tool             string    `json:"tool"`
	SessionID        string    `json:"sessionId"`
	CWD              string    `json:"cwd"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Title            string    `json:"title,omitempty"`
	FirstUserMessage string    `json:"firstUserMessage,omitempty"`
	LastUserMessage  string    `json:"lastUserMessage,omitempty"`
	Branch           string    `json:"branch,omitempty"`
	ResumeCommand    string    `json:"resumeCommand"`
	SourceFile       string    `json:"sourceFile"`
	Summary          string    `json:"summary,omitempty"`
	SummaryProvider  string    `json:"summaryProvider,omitempty"`
}

type DirectoryGroup struct {
	CWD         string    `json:"cwd"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Count       int       `json:"count"`
	Tools       []string  `json:"tools"`
	LatestTitle string    `json:"latestTitle,omitempty"`
}

func SortByUpdated(records []Record) {
	sort.SliceStable(records, func(i, j int) bool {
		return records[i].UpdatedAt.After(records[j].UpdatedAt)
	})
}

func FilterByCWD(records []Record, cwd string) []Record {
	if cwd == "" {
		return records
	}
	cwd = filepath.Clean(cwd)
	filtered := records[:0]
	for _, rec := range records {
		if filepath.Clean(rec.CWD) == cwd {
			filtered = append(filtered, rec)
		}
	}
	return filtered
}

func GroupByDirectory(records []Record) []DirectoryGroup {
	byDir := map[string]*DirectoryGroup{}
	toolSets := map[string]map[string]bool{}
	for _, rec := range records {
		group := byDir[rec.CWD]
		if group == nil {
			group = &DirectoryGroup{CWD: rec.CWD}
			byDir[rec.CWD] = group
			toolSets[rec.CWD] = map[string]bool{}
		}
		group.Count++
		toolSets[rec.CWD][rec.Tool] = true
		if rec.UpdatedAt.After(group.UpdatedAt) {
			group.UpdatedAt = rec.UpdatedAt
			group.LatestTitle = firstNonEmpty(rec.Title, rec.LastUserMessage, rec.FirstUserMessage)
		}
	}

	dirs := make([]DirectoryGroup, 0, len(byDir))
	for cwd, group := range byDir {
		for tool := range toolSets[cwd] {
			group.Tools = append(group.Tools, tool)
		}
		sort.Strings(group.Tools)
		dirs = append(dirs, *group)
	}
	sort.SliceStable(dirs, func(i, j int) bool {
		return dirs[i].UpdatedAt.After(dirs[j].UpdatedAt)
	})
	return dirs
}

func WriteJSON(w io.Writer, value any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}

func CleanText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
