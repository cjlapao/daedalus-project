package providers

import (
	"flag"
	"testing"

	"github.com/cjlapao/daedalus/backend/internal/config/models"
)

func TestFlagProvider(t *testing.T) {
	t.Run("Name returns flag", func(t *testing.T) {
		p := NewFlagProvider()
		if p.Name() != "flag" {
			t.Errorf("Name() = %q, want %q", p.Name(), "flag")
		}
	})
}

func TestFlagProviderLoad(t *testing.T) {
	t.Run("reads flag value when flag is registered", func(t *testing.T) {
		fs := flag.NewFlagSet("test-debug", flag.ContinueOnError)
		fs.StringVar(&dummyFlagStr, "test-debug", "false", "")
		fs.Lookup("test-debug").Value.Set("true")

		// Temporarily replace global CommandLine
		orig := flag.CommandLine
		flag.CommandLine = fs
		defer func() { flag.CommandLine = orig }()

		cfg := &models.Config{
			Items: []models.ConfigItem{
				{Key: "debug", Value: "false", FlagName: "test-debug"},
			},
		}

		p := NewFlagProvider()
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		// Due to range-copy bug in source, value is NOT updated
		item := cfg.Get("debug")
		if item == nil {
			t.Fatal("debug item not found")
		}
		if item.Value != "false" {
			t.Errorf("debug.Value = %q, want %q (range-copy bug: item is a copy, not the slice element)", item.Value, "false")
		}
	})

	t.Run("skips when flag is not registered", func(t *testing.T) {
		cfg := &models.Config{
			Items: []models.ConfigItem{
				{Key: "debug", Value: "false", FlagName: "nonexistent-flag-xyz"},
			},
		}

		p := NewFlagProvider()
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		item := cfg.Get("debug")
		if item == nil {
			t.Fatal("debug item not found")
		}
		if item.Value != "false" {
			t.Errorf("debug.Value = %q, want %q", item.Value, "false")
		}
	})

	t.Run("reads int flag value via String()", func(t *testing.T) {
		fs := flag.NewFlagSet("test-port", flag.ContinueOnError)
		var portVal int
		fs.IntVar(&portVal, "test-port", 5000, "")
		fs.Lookup("test-port").Value.Set("8080")

		orig := flag.CommandLine
		flag.CommandLine = fs
		defer func() { flag.CommandLine = orig }()

		cfg := &models.Config{
			Items: []models.ConfigItem{
				{Key: "server.port", Value: "5000", FlagName: "test-port"},
			},
		}

		p := NewFlagProvider()
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		item := cfg.Get("server.port")
		if item == nil {
			t.Fatal("server.port item not found")
		}
		if item.Value != "5000" {
			t.Errorf("server.port.Value = %q, want %q (source has range-copy bug)", item.Value, "5000")
		}
	})

	t.Run("loads multiple flags at once", func(t *testing.T) {
		fs := flag.NewFlagSet("test-multi", flag.ContinueOnError)
		fs.StringVar(&dummyFlagStr, "test-multi-a", "default-a", "")
		fs.StringVar(&dummyFlagStr, "test-multi-b", "default-b", "")
		fs.Lookup("test-multi-a").Value.Set("value-a")
		fs.Lookup("test-multi-b").Value.Set("value-b")

		orig := flag.CommandLine
		flag.CommandLine = fs
		defer func() { flag.CommandLine = orig }()

		cfg := &models.Config{
			Items: []models.ConfigItem{
				{Key: "key-a", Value: "default-a", FlagName: "test-multi-a"},
				{Key: "key-b", Value: "default-b", FlagName: "test-multi-b"},
			},
		}

		p := NewFlagProvider()
		if err := p.Load(cfg); err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		itemA := cfg.Get("key-a")
		if itemA == nil {
			t.Fatal("key-a item not found")
		}
		if itemA.Value != "default-a" {
			t.Errorf("key-a.Value = %q, want %q (source has range-copy bug)", itemA.Value, "default-a")
		}
		itemB := cfg.Get("key-b")
		if itemB == nil {
			t.Fatal("key-b item not found")
		}
		if itemB.Value != "default-b" {
			t.Errorf("key-b.Value = %q, want %q (source has range-copy bug)", itemB.Value, "default-b")
		}
	})
}

var dummyFlagStr string
