package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOverloadsDotEnvInDev(t *testing.T) {
	t.Setenv("APP_ENV", "dev")
	t.Setenv("LOG_LEVEL", "error")

	withTempEnvFile(t, "LOG_LEVEL=warn\n", func() {
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if cfg.Log.Level != "warn" {
			t.Fatalf("expected LOG_LEVEL from .env in dev, got %q", cfg.Log.Level)
		}
	})
}

func TestLoadKeepsSystemEnvInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("LOG_LEVEL", "error")

	withTempEnvFile(t, "LOG_LEVEL=warn\n", func() {
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if cfg.Log.Level != "error" {
			t.Fatalf("expected system LOG_LEVEL in production, got %q", cfg.Log.Level)
		}
	})
}

func withTempEnvFile(t *testing.T, content string, fn func()) {
	t.Helper()

	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWD)
	}()

	fn()
}
