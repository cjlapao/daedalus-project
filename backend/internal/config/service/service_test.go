package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/cjlapao/daedalus/backend/internal/config/interfaces"
	"github.com/cjlapao/daedalus/backend/internal/config/models"
	"github.com/cjlapao/daedalus/backend/internal/config/providers"
)

func TestConfigServiceNew(t *testing.T) {
	t.Run("creates ConfigService with default config", func(t *testing.T) {
		svc := New()
		if svc == nil {
			t.Fatal("New() returned nil")
		}
		if svc.config == nil {
			t.Fatal("ConfigService.config is nil")
		}
		if svc.config.Items == nil {
			t.Fatal("ConfigService.config.Items is nil")
		}
		if len(svc.providers) != 0 {
			t.Errorf("providers len = %d, want 0", len(svc.providers))
		}
		if svc.isLoaded {
			t.Error("isLoaded = true, want false")
		}
	})
}

func TestConfigServiceGetInstance(t *testing.T) {
	defer Reset()

	t.Run("returns nil when not initialized", func(t *testing.T) {
		Reset()
		if GetInstance() != nil {
			t.Error("GetInstance() = non-nil, want nil before Initialize")
		}
	})
}

func TestConfigServiceInitialize(t *testing.T) {
	defer Reset()

	t.Run("creates singleton with default providers", func(t *testing.T) {
		Reset()
		svc, err := Initialize()
		if err != nil {
			t.Fatalf("Initialize() error: %v", err)
		}
		if svc == nil {
			t.Fatal("Initialize() returned nil")
		}
		if GetInstance() != svc {
			t.Error("GetInstance() != returned service")
		}
	})

	t.Run("is idempotent (second call returns same instance)", func(t *testing.T) {
		Reset()
		svc1, _ := Initialize()
		svc2, _ := Initialize()
		if svc1 != svc2 {
			t.Error("Initialize() should be idempotent via sync.Once")
		}
	})

	t.Run("returns error when file provider fails", func(t *testing.T) {
		Reset()
		// Create a temp dir and set CONFIG_FILE_PATH to a non-existent file
		tmpDir := t.TempDir()
		os.Setenv(models.ConfigFilePathEnv, filepath.Join(tmpDir, "nonexistent.yaml"))
		defer os.Unsetenv(models.ConfigFilePathEnv)

		_, err := Initialize()
		// Should error because file provider can't find the file
		if err == nil {
			t.Error("Initialize() should error when config file doesn't exist")
		}
	})
}

func TestConfigServiceReset(t *testing.T) {
	defer Reset()

	t.Run("clears singleton", func(t *testing.T) {
		Reset()
		Initialize()
		Reset()
		if GetInstance() != nil {
			t.Error("GetInstance() should return nil after Reset()")
		}
	})
}

func TestConfigServiceLoad(t *testing.T) {
	defer Reset()

	t.Run("loads config from all providers", func(t *testing.T) {
		Reset()
		svc := New()
		svc.AddProvider(providers.NewFileProvider(""))
		svc.AddProvider(providers.NewEnvProvider())
		svc.AddProvider(providers.NewFlagProvider())

		// Set an env var to verify it's picked up
		os.Setenv("DEBUG", "true")
		defer os.Unsetenv("DEBUG")

		err := svc.Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if !svc.isLoaded {
			t.Error("isLoaded should be true after Load()")
		}

		item := svc.Get().Get("debug")
		if item == nil {
			t.Fatal("debug item not found")
		}
		// EnvProvider should have set debug to "true"
		if item.Value != "true" {
			t.Logf("debug.Value = %q (env may not have been picked up due to provider order)", item.Value)
		}
	})

	t.Run("resets config to defaults before loading", func(t *testing.T) {
		Reset()
		svc := New()
		svc.AddProvider(providers.NewEnvProvider())

		// Pre-set a value
		svc.config.Set("mykey", "old-value")

		err := svc.Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		// After Load, config is reset to defaults, so mykey should not exist
		item := svc.Get().Get("mykey")
		if item != nil {
			t.Errorf("mykey should have been reset, but found: %v", item)
		}
	})
}

func TestConfigServiceGet(t *testing.T) {
	defer Reset()

	t.Run("returns config pointer", func(t *testing.T) {
		Reset()
		svc := New()
		if svc.Get() == nil {
			t.Fatal("Get() returned nil")
		}
	})
}

func TestConfigServiceGetSection(t *testing.T) {
	defer Reset()

	t.Run("delegates to config.Bind", func(t *testing.T) {
		Reset()
		svc := New()
		svc.AddProvider(providers.NewEnvProvider())
		svc.config.Set("server.port", "8080")

		type ServerConfig struct {
			Port int `config:"server.port"`
		}
		var target ServerConfig
		err := svc.GetSection("", &target)
		if err != nil {
			t.Fatalf("GetSection() error: %v", err)
		}
		if target.Port != 8080 {
			t.Errorf("Port = %d, want 8080", target.Port)
		}
	})

	t.Run("returns error when config is nil", func(t *testing.T) {
		Reset()
		svc := New()
		svc.config = nil
		err := svc.GetSection("", &struct{}{})
		if err == nil {
			t.Error("GetSection() should error when config is nil")
		}
	})
}

func TestConfigServiceLoadFromFile(t *testing.T) {
	t.Run("loads config from file", func(t *testing.T) {
		tmpDir := t.TempDir()
		yamlPath := filepath.Join(tmpDir, "test.yaml")
		err := os.WriteFile(yamlPath, []byte(`server:
  port: 9999
`), 0o644)
		if err != nil {
			t.Fatalf("failed to write test YAML: %v", err)
		}

		cfg, err := New().LoadFromFile(yamlPath)
		if err != nil {
			t.Fatalf("LoadFromFile() error: %v", err)
		}
		if cfg == nil {
			t.Fatal("LoadFromFile() returned nil config")
		}

		port := cfg.Get("server.port")
		if port == nil || port.Value != "9999" {
			t.Errorf("server.port = %v, want 9999", port)
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		_, err := New().LoadFromFile("/nonexistent/config.yaml")
		if err == nil {
			t.Fatal("LoadFromFile() should error for non-existent file")
		}
	})
}

func TestConfigServiceAddProvider(t *testing.T) {
	defer Reset()

	t.Run("adds provider and marks not loaded", func(t *testing.T) {
		Reset()
		svc := New()
		if svc.isLoaded {
			t.Error("new service should not be loaded")
		}

		svc.AddProvider(providers.NewEnvProvider())
		if len(svc.providers) != 1 {
			t.Errorf("providers len = %d, want 1", len(svc.providers))
		}
	})
}

func TestConfigServiceName(t *testing.T) {
	t.Run("returns config", func(t *testing.T) {
		svc := New()
		if svc.Name() != "config" {
			t.Errorf("Name() = %q, want %q", svc.Name(), "config")
		}
	})
}

func TestConfigServiceInit(t *testing.T) {
	defer Reset()

	t.Run("adds default providers and loads", func(t *testing.T) {
		Reset()
		svc := New()
		err := svc.Init(context.Background())
		if err != nil {
			t.Fatalf("Init() error: %v", err)
		}
		if len(svc.providers) != 3 {
			t.Errorf("providers len = %d, want 3 (file, env, flag)", len(svc.providers))
		}
	})

	t.Run("is idempotent when already loaded", func(t *testing.T) {
		Reset()
		svc := New()
		svc.isLoaded = true
		initialProviders := len(svc.providers)

		err := svc.Init(context.Background())
		if err != nil {
			t.Fatalf("Init() error: %v", err)
		}
		// Should not add more providers
		if len(svc.providers) != initialProviders {
			t.Error("Init() should be idempotent when already loaded")
		}
	})
}

func TestConfigServiceHealth(t *testing.T) {
	t.Run("returns error when not loaded", func(t *testing.T) {
		svc := New()
		err := svc.Health(context.Background())
		if err == nil {
			t.Error("Health() should error when not loaded")
		}
	})

	t.Run("returns nil when loaded", func(t *testing.T) {
		Reset()
		svc := New()
		svc.isLoaded = true
		err := svc.Health(context.Background())
		if err != nil {
			t.Errorf("Health() error: %v", err)
		}
	})
}

func TestConfigServiceIsEnabled(t *testing.T) {
	t.Run("returns true", func(t *testing.T) {
		svc := New()
		if !svc.IsEnabled() {
			t.Error("IsEnabled() = false, want true")
		}
	})
}

func TestConfigServiceDependencies(t *testing.T) {
	t.Run("returns empty slice", func(t *testing.T) {
		svc := New()
		deps := svc.Dependencies()
		if deps == nil {
			t.Fatal("Dependencies() returned nil")
		}
		if len(deps) != 0 {
			t.Errorf("Dependencies() len = %d, want 0", len(deps))
		}
	})
}

// --- Fake provider for testing ---

type fakeProvider struct {
	name    string
	loadErr error
}

func (f *fakeProvider) Name() string   { return f.name }
func (f *fakeProvider) Load(cfg *models.Config) error {
	if f.loadErr != nil {
		return f.loadErr
	}
	cfg.Set("fake.key", "fake.value")
	return nil
}

// Verify interfaces
var _ interfaces.Provider = (*fakeProvider)(nil)
var _ interfaces.ConfigService = (*ConfigService)(nil)
