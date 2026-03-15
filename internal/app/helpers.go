package app

import (
	"context"
	"strings"
)

type closeGroup []func(ctx context.Context) error

func (g *closeGroup) add(fn func(ctx context.Context) error) {
	if fn != nil {
		*g = append(*g, fn)
	}
}

func (g closeGroup) close(ctx context.Context) error {
	for _, fn := range g {
		_ = fn(ctx)
	}
	return nil
}

func normalizeMode(mode string) string {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		return "auto"
	}
	return mode
}

func wantMongo(mode string) bool {
	return mode == "mongo" || mode == "both" || mode == "auto"
}

func wantMySQL(mode string) bool {
	return mode == "mysql" || mode == "both" || mode == "auto"
}
