package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsDotEnvWithoutOverridingExistingEnv(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	dotenv := []byte("SERVER_PORT=19001\nCORS_ORIGINS=http://localhost:3000\n")
	if err := os.WriteFile(filepath.Join(tmp, ".env"), dotenv, 0o600); err != nil {
		t.Fatal(err)
	}
	oldCORS, hadCORS := os.LookupEnv("CORS_ORIGINS")
	_ = os.Unsetenv("CORS_ORIGINS")
	t.Cleanup(func() {
		if hadCORS {
			_ = os.Setenv("CORS_ORIGINS", oldCORS)
		} else {
			_ = os.Unsetenv("CORS_ORIGINS")
		}
	})
	t.Setenv("SERVER_PORT", "19000")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != 19000 {
		t.Fatalf("SERVER_PORT from environment should win, got %d", cfg.Server.Port)
	}
	if cfg.CORS.Origins != "http://localhost:3000" {
		t.Fatalf("CORS_ORIGINS should be loaded from .env, got %q", cfg.CORS.Origins)
	}
}
