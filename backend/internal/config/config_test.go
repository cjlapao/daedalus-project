package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackageFacade(t *testing.T) {
	t.Run("Reset clears singleton", func(t *testing.T) {
		Reset()
		if GetInstance() != nil {
			t.Error("GetInstance() should return nil after Reset()")
		}
	})
}

func TestInitialize(t *testing.T) {
	defer Reset()

	t.Run("delegates to service.Initialize", func(t *testing.T) {
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
}

func TestGetInstance(t *testing.T) {
	defer Reset()

	t.Run("returns nil before Initialize", func(t *testing.T) {
		Reset()
		if GetInstance() != nil {
			t.Error("GetInstance() should return nil before Initialize")
		}
	})

	t.Run("returns instance after Initialize", func(t *testing.T) {
		Reset()
		Initialize()
		if GetInstance() == nil {
			t.Fatal("GetInstance() should return non-nil after Initialize")
		}
	})
}

func TestReset(t *testing.T) {
	defer Reset()

	t.Run("clears singleton", func(t *testing.T) {
		Initialize()
		Reset()
		if GetInstance() != nil {
			t.Error("GetInstance() should return nil after Reset()")
		}
	})
}

func TestGetSection(t *testing.T) {
	defer Reset()

	t.Run("returns error when not initialized", func(t *testing.T) {
		Reset()
		_, err := GetSection[string]("nonexistent")
		if err == nil {
			t.Error("GetSection() should error when not initialized")
		}
	})

	t.Run("delegates to service.GetSection", func(t *testing.T) {
		Reset()
		svc, err := Initialize()
		if err != nil {
			t.Fatalf("Initialize() error: %v", err)
		}
		// Set a value via the service
		svc.Get().Set("test.key", "test.value")

		type TestConf struct {
			Key string `config:"test.key"`
		}
		result, err := GetSection[TestConf]("")
		if err != nil {
			t.Fatalf("GetSection() error: %v", err)
		}
		if result.Key != "test.value" {
			t.Errorf("Key = %q, want %q", result.Key, "test.value")
		}
	})
}

func TestLoadFromFile(t *testing.T) {
	defer Reset()

	t.Run("creates new service when no instance exists", func(t *testing.T) {
		Reset()
		tmpDir := t.TempDir()
		yamlPath := filepath.Join(tmpDir, "test.yaml")
		err := os.WriteFile(yamlPath, []byte(`key: value
`), 0o644)
		if err != nil {
			t.Fatalf("failed to write test YAML: %v", err)
		}

		cfg, err := LoadFromFile(yamlPath)
		if err != nil {
			t.Fatalf("LoadFromFile() error: %v", err)
		}
		if cfg == nil {
			t.Fatal("LoadFromFile() returned nil")
		}

		item := cfg.Get("key")
		if item == nil || item.Value != "value" {
			t.Errorf("key = %v, want value", item)
		}
	})

	t.Run("delegates to existing instance", func(t *testing.T) {
		Reset()
		Initialize()
		tmpDir := t.TempDir()
		yamlPath := filepath.Join(tmpDir, "test2.yaml")
		err := os.WriteFile(yamlPath, []byte(`other: data
`), 0o644)
		if err != nil {
			t.Fatalf("failed to write test YAML: %v", err)
		}

		cfg, err := LoadFromFile(yamlPath)
		if err != nil {
			t.Fatalf("LoadFromFile() error: %v", err)
		}
		if cfg == nil {
			t.Fatal("LoadFromFile() returned nil")
		}

		item := cfg.Get("other")
		if item == nil || item.Value != "data" {
			t.Errorf("other = %v, want data", item)
		}
	})
}
