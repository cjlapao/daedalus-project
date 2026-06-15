package providers

import (
	"os"
	"testing"

	"github.com/cjlapao/daedalus/backend/internal/config/models"
)

func TestEnvProvider(t *testing.T) {
	t.Run("Name returns env", func(t *testing.T) {
		p := NewEnvProvider()
		if p.Name() != "env" {
			t.Errorf("Name() = %q, want %q", p.Name(), "env")
		}
	})
}

func TestEnvProviderLoad(t *testing.T) {
	// Clean up env vars after each test
	cleanup := func(keys ...string) {
		for _, k := range keys {
			os.Unsetenv(k)
		}
	}

	t.Run("respects explicit EnvName", func(t *testing.T) {
		defer cleanup("MY_DEBUG")
		os.Setenv("MY_DEBUG", "true")

		cfg := &models.Config{
			Items: []models.ConfigItem{
				{Key: "debug", Value: "false", EnvName: "MY_DEBUG"},
			},
		}

		p := NewEnvProvider()
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		item := cfg.Get("debug")
		if item == nil {
			t.Fatal("debug item not found")
		}
		if item.Value != "true" {
			t.Errorf("debug.Value = %q, want %q", item.Value, "true")
		}
	})

	t.Run("auto-maps dot.key to UPPER_KEY", func(t *testing.T) {
		defer cleanup("SERVER_PORT")
		os.Setenv("SERVER_PORT", "8080")

		cfg := &models.Config{
			Items: []models.ConfigItem{
				{Key: "server.port", Value: "5000"},
			},
		}

		p := NewEnvProvider()
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		item := cfg.Get("server.port")
		if item == nil {
			t.Fatal("server.port item not found")
		}
		if item.Value != "8080" {
			t.Errorf("server.port.Value = %q, want %q", item.Value, "8080")
		}
	})

	t.Run("explicit EnvName takes priority over auto-mapped", func(t *testing.T) {
		defer cleanup("MY_KEY", "AUTO_KEY")
		os.Setenv("MY_KEY", "from-env")
		os.Setenv("AUTO_KEY", "from-auto")

		cfg := &models.Config{
			Items: []models.ConfigItem{
				{Key: "my.key", Value: "default", EnvName: "MY_KEY"},
			},
		}

		p := NewEnvProvider()
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		item := cfg.Get("my.key")
		if item == nil {
			t.Fatal("my.key item not found")
		}
		// Explicit EnvName takes priority
		if item.Value != "from-env" {
			t.Errorf("my.key.Value = %q, want %q", item.Value, "from-env")
		}
	})

	t.Run("no env set leaves value unchanged", func(t *testing.T) {
		defer cleanup("NONEXISTENT_KEY_XYZ")

		cfg := &models.Config{
			Items: []models.ConfigItem{
				{Key: "test.key", Value: "default", EnvName: "NONEXISTENT_KEY_XYZ"},
			},
		}

		p := NewEnvProvider()
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		item := cfg.Get("test.key")
		if item == nil {
			t.Fatal("test.key item not found")
		}
		if item.Value != "default" {
			t.Errorf("test.key.Value = %q, want %q", item.Value, "default")
		}
	})

	t.Run("auto-mapped key with no env set leaves value unchanged", func(t *testing.T) {
		defer cleanup("NONEXISTENT_AUTO_KEY_XYZ")

		cfg := &models.Config{
			Items: []models.ConfigItem{
				{Key: "nonexistent.auto.key", Value: "default"},
			},
		}

		p := NewEnvProvider()
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		item := cfg.Get("nonexistent.auto.key")
		if item == nil {
			t.Fatal("item not found")
		}
		if item.Value != "default" {
			t.Errorf("Value = %q, want %q", item.Value, "default")
		}
	})
}
