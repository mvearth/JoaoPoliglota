package config

import "testing"

func TestGetenv(t *testing.T) {
	t.Run("returns fallback when unset", func(t *testing.T) {
		if got := Getenv("JP_TEST_UNSET", "fallback"); got != "fallback" {
			t.Errorf("got %q, want fallback", got)
		}
	})

	t.Run("returns value when set", func(t *testing.T) {
		t.Setenv("JP_TEST_SET", "value")
		if got := Getenv("JP_TEST_SET", "fallback"); got != "value" {
			t.Errorf("got %q, want value", got)
		}
	})

	t.Run("returns fallback when empty", func(t *testing.T) {
		t.Setenv("JP_TEST_EMPTY", "")
		if got := Getenv("JP_TEST_EMPTY", "fallback"); got != "fallback" {
			t.Errorf("got %q, want fallback", got)
		}
	})
}

func TestLoadDBDefaults(t *testing.T) {
	for _, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"} {
		t.Setenv(k, "")
	}
	cfg := LoadDB()
	if cfg.Host != "localhost" || cfg.Port != "5432" || cfg.Name != "joaopoliglota" {
		t.Errorf("unexpected defaults: %+v", cfg)
	}
}
