package adapters

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var syntheticBlocks = []*regexp.Regexp{
	regexp.MustCompile(`(?s)<ide_opened_file>.*?</ide_opened_file>`),
	regexp.MustCompile(`(?s)<local-command-caveat>.*?</local-command-caveat>`),
	regexp.MustCompile(`(?s)<turn_aborted>.*?</turn_aborted>`),
	regexp.MustCompile(`(?s)<environment_context>.*?</environment_context>`),
	regexp.MustCompile(`(?s)<command-[^>]+>.*?</command-[^>]+>`),
	regexp.MustCompile(`(?s)<local-command-[^>]+>.*?</local-command-[^>]+>`),
}

func readJSONL(path string, handle func(map[string]any)) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var obj map[string]any
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			continue
		}
		handle(obj)
	}
	return nil
}

func parseTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15-04-05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t
		}
	}
	return time.Time{}
}

func maxTime(current time.Time, next time.Time) time.Time {
	if next.After(current) {
		return next
	}
	return current
}

func stringField(obj map[string]any, key string) string {
	value, _ := obj[key].(string)
	return value
}

func mapField(obj map[string]any, key string) map[string]any {
	value, _ := obj[key].(map[string]any)
	return value
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func buildResumeCommand(cwd string, tool string, args ...string) string {
	parts := []string{"cd", shellQuote(cwd), "&&", tool}
	for _, arg := range args {
		parts = append(parts, shellQuote(arg))
	}
	return strings.Join(parts, " ")
}

func fileModTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

func cleanAbs(path string) string {
	if path == "" {
		return ""
	}
	if abs, err := filepath.Abs(path); err == nil {
		return filepath.Clean(abs)
	}
	return filepath.Clean(path)
}

func textFromClaudeContent(value any) string {
	switch content := value.(type) {
	case string:
		return content
	case []any:
		var parts []string
		for _, item := range content {
			obj, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if stringField(obj, "type") == "text" {
				parts = append(parts, stringField(obj, "text"))
			}
		}
		return strings.Join(parts, " ")
	default:
		return ""
	}
}

func textFromCodexContent(value any) string {
	items, ok := value.([]any)
	if !ok {
		return ""
	}
	var parts []string
	for _, item := range items {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		switch stringField(obj, "type") {
		case "input_text", "text":
			parts = append(parts, stringField(obj, "text"))
		}
	}
	return strings.Join(parts, " ")
}

func shouldSkipSyntheticUserText(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return true
	}
	prefixes := []string{
		"<environment_context>",
		"<permissions instructions>",
		"<collaboration_mode>",
		"<skills_instructions>",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(text, prefix) {
			return true
		}
	}
	return false
}

func cleanUserText(text string) string {
	text = strings.TrimSpace(text)
	for _, block := range syntheticBlocks {
		text = block.ReplaceAllString(text, " ")
	}
	markers := []string{
		"## My request for Codex:",
		"My request for Codex:",
	}
	for _, marker := range markers {
		if idx := strings.LastIndex(text, marker); idx >= 0 {
			text = text[idx+len(marker):]
			break
		}
	}
	return strings.Join(strings.Fields(text), " ")
}

func decodeClaudeProjectDir(name string) string {
	if name == "" || name[0] != '-' {
		return ""
	}
	return "/" + strings.ReplaceAll(strings.TrimPrefix(name, "-"), "-", "/")
}

func timestampMillis(value any) time.Time {
	switch v := value.(type) {
	case float64:
		return time.UnixMilli(int64(v))
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return time.Time{}
		}
		return time.UnixMilli(n)
	default:
		return time.Time{}
	}
}
