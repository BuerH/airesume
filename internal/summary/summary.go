package summary

import (
	"context"
	"errors"
)

type Provider interface {
	Summarize(ctx context.Context, sessionID string) (string, error)
}

type NoopProvider struct{}

func (NoopProvider) Summarize(context.Context, string) (string, error) {
	return "", errors.New("AI summary provider is not configured yet")
}
