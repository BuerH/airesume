package adapters

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/BuerH/airesume/internal/session"
)

type Adapter interface {
	Name() string
	Scan(ctx context.Context) ([]session.Record, error)
}

type Registry struct {
	adapters []Adapter
}

func DefaultRegistry() Registry {
	home, _ := os.UserHomeDir()
	return Registry{
		adapters: []Adapter{
			NewCodexAdapter(filepath.Join(home, ".codex")),
			NewClaudeAdapter(filepath.Join(home, ".claude")),
		},
	}
}

func (r Registry) Scan(ctx context.Context) ([]session.Record, error) {
	var records []session.Record
	var joined error
	for _, adapter := range r.adapters {
		scanned, err := adapter.Scan(ctx)
		if err != nil {
			joined = errors.Join(joined, err)
			continue
		}
		records = append(records, scanned...)
	}
	return records, joined
}
