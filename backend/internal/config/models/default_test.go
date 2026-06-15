package models

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	t.Run("returns non-nil config", func(t *testing.T) {
		cfg := DefaultConfig()
		if cfg == nil {
			t.Fatal("DefaultConfig() returned nil")
		}
	})

	t.Run("has Items slice", func(t *testing.T) {
		cfg := DefaultConfig()
		if cfg.Items == nil {
			t.Fatal("DefaultConfig().Items is nil")
		}
	})

	t.Run("contains expected number of items", func(t *testing.T) {
		cfg := DefaultConfig()
		want := 50 // Count of default config items in default.go
		if len(cfg.Items) != want {
			t.Errorf("DefaultConfig().Items len = %d, want %d", len(cfg.Items), want)
		}
	})

	t.Run("contains DebugKey", func(t *testing.T) {
		cfg := DefaultConfig()
		item := cfg.Get(DebugKey)
		if item == nil {
			t.Fatal("DebugKey not found in DefaultConfig")
		}
		if item.Value != "false" {
			t.Errorf("DebugKey.Value = %q, want %q", item.Value, "false")
		}
		if item.EnvName != DebugEnvKey {
			t.Errorf("DebugKey.EnvName = %q, want %q", item.EnvName, DebugEnvKey)
		}
	})

	t.Run("contains EnvironmentKey with default value", func(t *testing.T) {
		cfg := DefaultConfig()
		item := cfg.Get(EnvironmentKey)
		if item == nil {
			t.Fatal("EnvironmentKey not found")
		}
		if item.Value != "development" {
			t.Errorf("EnvironmentKey.Value = %q, want %q", item.Value, "development")
		}
	})

	t.Run("contains Logger keys", func(t *testing.T) {
		cfg := DefaultConfig()
		for _, key := range []string{LogLevelKey, LogFormatKey, LogEnableCallerKey, LogUseStdoutKey, LogFilePathKey} {
			item := cfg.Get(key)
			if item == nil {
				t.Errorf("key %q not found in DefaultConfig", key)
			}
		}
	})

	t.Run("contains Server keys", func(t *testing.T) {
		cfg := DefaultConfig()
		for _, key := range []string{ServerAPIPortKey, ServerBindAddressKey, ServerBaseURLKey, ServerAPIPrefixKey} {
			item := cfg.Get(key)
			if item == nil {
				t.Errorf("key %q not found in DefaultConfig", key)
			}
		}
		port := cfg.Get(ServerAPIPortKey)
		if port != nil && port.Value != "5000" {
			t.Errorf("ServerAPIPortKey.Value = %q, want %q", port.Value, "5000")
		}
	})

	t.Run("contains Database keys", func(t *testing.T) {
		cfg := DefaultConfig()
		for _, key := range []string{DatabaseTypeKey, DatabaseHostKey, DatabasePortKey, DatabaseDatabaseKey, DatabaseUsernameKey, DatabasePasswordKey} {
			item := cfg.Get(key)
			if item == nil {
				t.Errorf("key %q not found in DefaultConfig", key)
			}
		}
		typ := cfg.Get(DatabaseTypeKey)
		if typ != nil && typ.Value != "sqlite" {
			t.Errorf("DatabaseTypeKey.Value = %q, want %q", typ.Value, "sqlite")
		}
	})

	t.Run("contains Auth keys", func(t *testing.T) {
		cfg := DefaultConfig()
		for _, key := range []string{AuthRootPasswordKey, JwtAuthSecretKey, JwtIssuerKey} {
			item := cfg.Get(key)
			if item == nil {
				t.Errorf("key %q not found in DefaultConfig", key)
			}
		}
	})

	t.Run("contains Security keys", func(t *testing.T) {
		cfg := DefaultConfig()
		for _, key := range []string{SecurityPasswordMinLengthKey, SecurityPasswordRequireNumberKey, SecurityPasswordRequireSpecialKey, SecurityPasswordRequireUppercaseKey} {
			item := cfg.Get(key)
			if item == nil {
				t.Errorf("key %q not found in DefaultConfig", key)
			}
		}
	})

	t.Run("contains CORS keys", func(t *testing.T) {
		cfg := DefaultConfig()
		for _, key := range []string{CorsAllowOriginsKey, CorsAllowMethodsKey, CorsAllowHeadersKey, CorsExposeHeadersKey} {
			item := cfg.Get(key)
			if item == nil {
				t.Errorf("key %q not found in DefaultConfig", key)
			}
		}
	})

	t.Run("most default config items have EnvName or FlagName set", func(t *testing.T) {
		cfg := DefaultConfig()
		// api_key and vaults are intentionally without EnvName/FlagName
		noProvider := map[string]bool{"api_key": true, "vaults": true}
		for _, item := range cfg.Items {
			if noProvider[item.Key] {
				continue
			}
			if item.EnvName == "" && item.FlagName == "" {
				t.Errorf("ConfigItem key=%q has neither EnvName nor FlagName set", item.Key)
			}
		}
	})
}
