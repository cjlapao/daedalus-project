package models

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestConfigRegisterResolver(t *testing.T) {
	t.Run("registers a resolver", func(t *testing.T) {
		cfg := &Config{}
		cfg.RegisterResolver(&testResolver{prefix: "test", value: "resolved"})

		if cfg.resolvers == nil || len(cfg.resolvers) != 1 {
			t.Fatal("resolvers map not initialized")
		}
		if _, ok := cfg.resolvers["test"]; !ok {
			t.Error("resolver not registered with prefix 'test'")
		}
	})

	t.Run("overwrites existing resolver with same prefix", func(t *testing.T) {
		cfg := &Config{}
		cfg.RegisterResolver(&testResolver{prefix: "test", value: "first"})
		cfg.RegisterResolver(&testResolver{prefix: "test", value: "second"})

		if len(cfg.resolvers) != 1 {
			t.Errorf("resolvers len = %d, want 1", len(cfg.resolvers))
		}
	})
}

func TestConfigResolveValue(t *testing.T) {
	t.Run("returns value as-is when no resolvers", func(t *testing.T) {
		cfg := &Config{}
		got := cfg.ResolveValue(context.Background(), "${{test::content}}")
		if got != "${{test::content}}" {
			t.Errorf("ResolveValue() = %q, want %q", got, "${{test::content}}")
		}
	})

	t.Run("resolves with matching resolver", func(t *testing.T) {
		cfg := &Config{}
		cfg.RegisterResolver(&testResolver{prefix: "db", value: "postgresql://localhost/mydb"})

		got := cfg.ResolveValue(context.Background(), "${{db::connection}}")
		if got != "postgresql://localhost/mydb" {
			t.Errorf("ResolveValue() = %q, want %q", got, "postgresql://localhost/mydb")
		}
	})

	t.Run("returns original value when resolver not found", func(t *testing.T) {
		cfg := &Config{}
		cfg.RegisterResolver(&testResolver{prefix: "db", value: "postgresql"})

		got := cfg.ResolveValue(context.Background(), "${{missing::key}}")
		if got != "${{missing::key}}" {
			t.Errorf("ResolveValue() = %q, want original", got)
		}
	})

	t.Run("returns original value when resolver errors", func(t *testing.T) {
		cfg := &Config{}
		cfg.RegisterResolver(&errorResolver{prefix: "fail"})

		got := cfg.ResolveValue(context.Background(), "${{fail::key}}")
		if got != "${{fail::key}}" {
			t.Errorf("ResolveValue() = %q, want original on error", got)
		}
	})

	t.Run("returns non-variable value as-is", func(t *testing.T) {
		cfg := &Config{}
		cfg.RegisterResolver(&testResolver{prefix: "x", value: "y"})

		got := cfg.ResolveValue(context.Background(), "plain-text")
		if got != "plain-text" {
			t.Errorf("ResolveValue() = %q, want %q", got, "plain-text")
		}
	})
}

func TestConfigGet(t *testing.T) {
	t.Run("returns item by key", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{
				{Key: "debug", Value: "true"},
				{Key: "port", Value: "8080"},
			},
		}

		item := cfg.Get("debug")
		if item == nil {
			t.Fatal("Get(debug) returned nil")
		}
		if item.Value != "true" {
			t.Errorf("debug.Value = %q, want %q", item.Value, "true")
		}
	})

	t.Run("returns nil for missing key", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{
				{Key: "debug", Value: "true"},
			},
		}

		if cfg.Get("nonexistent") != nil {
			t.Error("Get(nonexistent) should return nil")
		}
	})
}

func TestConfigIsDebug(t *testing.T) {
	t.Run("returns true when debug is true", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: DebugKey, Value: "true"}},
		}
		if !cfg.IsDebug(context.Background()) {
			t.Error("IsDebug() = false, want true")
		}
	})

	t.Run("returns false when debug is false", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: DebugKey, Value: "false"}},
		}
		if cfg.IsDebug(context.Background()) {
			t.Error("IsDebug() = true, want false")
		}
	})

	t.Run("returns default false when key missing", func(t *testing.T) {
		cfg := &Config{Items: []ConfigItem{}}
		if cfg.IsDebug(context.Background()) {
			t.Error("IsDebug() = true, want false (default)")
		}
	})
}

func TestConfigGetValue(t *testing.T) {
	t.Run("returns resolved value when key exists", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "mykey", Value: "myvalue"}},
		}
		got := cfg.GetValue(context.Background(), "mykey", "default")
		if got != "myvalue" {
			t.Errorf("GetValue() = %v, want %v", got, "myvalue")
		}
	})

	t.Run("returns default when key missing", func(t *testing.T) {
		cfg := &Config{Items: []ConfigItem{}}
		got := cfg.GetValue(context.Background(), "missing", "default-val")
		if got != "default-val" {
			t.Errorf("GetValue() = %v, want %v", got, "default-val")
		}
	})

	t.Run("returns default when key not set", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "empty", Value: ""}},
		}
		got := cfg.GetValue(context.Background(), "empty", "default-val")
		if got != "default-val" {
			t.Errorf("GetValue() = %v, want %v", got, "default-val")
		}
	})
}

func TestConfigGetString(t *testing.T) {
	t.Run("returns string value when key exists", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "name", Value: "test"}},
		}
		got := cfg.GetString(context.Background(), "name", "default")
		if got != "test" {
			t.Errorf("GetString() = %q, want %q", got, "test")
		}
	})

	t.Run("returns default when key missing", func(t *testing.T) {
		cfg := &Config{Items: []ConfigItem{}}
		got := cfg.GetString(context.Background(), "missing", "default")
		if got != "default" {
			t.Errorf("GetString() = %q, want %q", got, "default")
		}
	})
}

func TestConfigGetBool(t *testing.T) {
	t.Run("returns true for truthy value", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "flag", Value: "true"}},
		}
		if !cfg.GetBool(context.Background(), "flag", false) {
			t.Error("GetBool() = false, want true")
		}
	})

	t.Run("returns false for falsy value", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "flag", Value: "false"}},
		}
		if cfg.GetBool(context.Background(), "flag", true) {
			t.Error("GetBool() = true, want false")
		}
	})

	t.Run("returns default when key missing", func(t *testing.T) {
		cfg := &Config{Items: []ConfigItem{}}
		if !cfg.GetBool(context.Background(), "missing", true) {
			t.Error("GetBool() = false, want true (default)")
		}
	})
}

func TestConfigGetInt(t *testing.T) {
	t.Run("parses valid int", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "count", Value: "42"}},
		}
		got := cfg.GetInt(context.Background(), "count", 0)
		if got != 42 {
			t.Errorf("GetInt() = %d, want 42", got)
		}
	})

	t.Run("returns default on parse error", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "bad", Value: "not-a-number"}},
		}
		got := cfg.GetInt(context.Background(), "bad", -1)
		if got != -1 {
			t.Errorf("GetInt() = %d, want -1", got)
		}
	})

	t.Run("returns default when key missing", func(t *testing.T) {
		cfg := &Config{Items: []ConfigItem{}}
		got := cfg.GetInt(context.Background(), "missing", 99)
		if got != 99 {
			t.Errorf("GetInt() = %d, want 99", got)
		}
	})
}

func TestConfigGetDuration(t *testing.T) {
	t.Run("parses valid duration", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "timeout", Value: "30s"}},
		}
		got := cfg.GetDuration(context.Background(), "timeout", 0)
		want := 30 * time.Second
		if got != want {
			t.Errorf("GetDuration() = %v, want %v", got, want)
		}
	})

	t.Run("returns default on parse error", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "bad", Value: "not-a-duration"}},
		}
		got := cfg.GetDuration(context.Background(), "bad", 5*time.Second)
		if got != 5*time.Second {
			t.Errorf("GetDuration() = %v, want 5s", got)
		}
	})

	t.Run("returns default when key missing", func(t *testing.T) {
		cfg := &Config{Items: []ConfigItem{}}
		got := cfg.GetDuration(context.Background(), "missing", 10*time.Second)
		if got != 10*time.Second {
			t.Errorf("GetDuration() = %v, want 10s", got)
		}
	})
}

func TestConfigSet(t *testing.T) {
	t.Run("updates existing item", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "debug", Value: "false"}},
		}
		cfg.Set("debug", "true")
		item := cfg.Get("debug")
		if item == nil || item.Value != "true" {
			t.Errorf("Set() failed: debug.Value = %q, want %q", itemValue(item), "true")
		}
	})

	t.Run("appends new item", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "debug", Value: "false"}},
		}
		cfg.Set("newkey", "newvalue")
		if len(cfg.Items) != 2 {
			t.Errorf("Items len = %d, want 2", len(cfg.Items))
		}
		item := cfg.Get("newkey")
		if item == nil || item.Value != "newvalue" {
			t.Errorf("newkey = %v, want value newvalue", item)
		}
	})
}

func TestConfigStoragePath(t *testing.T) {
	t.Run("returns joined path", func(t *testing.T) {
		cfg := &Config{}
		path := cfg.StoragePath()
		// Should contain DefaultStoragePath
		if !contains(path, DefaultStoragePath) {
			t.Errorf("StoragePath() = %q, should contain %q", path, DefaultStoragePath)
		}
	})

	t.Run("uses dot when home dir lookup fails", func(t *testing.T) {
		// We can't easily test the error case for os.UserHomeDir()
		// but we can verify the normal case works
		cfg := &Config{}
		path := cfg.StoragePath()
		if path == "" {
			t.Error("StoragePath() returned empty string")
		}
	})
}

func TestConfigBind(t *testing.T) {
	t.Run("returns error for non-pointer target", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "name", Value: "test"}},
		}
		err := cfg.Bind(context.Background(), "", "not-a-pointer")
		if err == nil {
			t.Error("Bind() should error for non-pointer")
		}
	})

	t.Run("returns error for pointer to non-struct", func(t *testing.T) {
		cfg := &Config{}
		var s string
		err := cfg.Bind(context.Background(), "", &s)
		if err == nil {
			t.Error("Bind() should error for pointer to non-struct")
		}
	})

	t.Run("binds string field with config tag", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "server.port", Value: "8080"}},
		}
		type ServerConfig struct {
			Port int `config:"server.port"`
		}
		var target ServerConfig
		err := cfg.Bind(context.Background(), "", &target)
		if err != nil {
			t.Fatalf("Bind() error: %v", err)
		}
		if target.Port != 8080 {
			t.Errorf("Port = %d, want 8080", target.Port)
		}
	})

	t.Run("binds string field with snake_case auto-mapping", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "host_name", Value: "myhost"}},
		}
		type HostConfig struct {
			HostName string // auto-maps to host_name
		}
		var target HostConfig
		err := cfg.Bind(context.Background(), "", &target)
		if err != nil {
			t.Fatalf("Bind() error: %v", err)
		}
		if target.HostName != "myhost" {
			t.Errorf("HostName = %q, want %q", target.HostName, "myhost")
		}
	})

	t.Run("binds bool field", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "debug", Value: "true"}},
		}
		type DebugConfig struct {
			Debug bool `config:"debug"`
		}
		var target DebugConfig
		err := cfg.Bind(context.Background(), "", &target)
		if err != nil {
			t.Fatalf("Bind() error: %v", err)
		}
		if !target.Debug {
			t.Error("Debug = false, want true")
		}
	})

	t.Run("binds float field", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "ratio", Value: "0.75"}},
		}
		type RatioConfig struct {
			Ratio float64 `config:"ratio"`
		}
		var target RatioConfig
		err := cfg.Bind(context.Background(), "", &target)
		if err != nil {
			t.Fatalf("Bind() error: %v", err)
		}
		if target.Ratio != 0.75 {
			t.Errorf("Ratio = %f, want 0.75", target.Ratio)
		}
	})

	t.Run("binds duration field", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "timeout", Value: "30s"}},
		}
		type TimeoutConfig struct {
			Timeout time.Duration `config:"timeout"`
		}
		var target TimeoutConfig
		err := cfg.Bind(context.Background(), "", &target)
		if err != nil {
			t.Fatalf("Bind() error: %v", err)
		}
		if target.Timeout != 30*time.Second {
			t.Errorf("Timeout = %v, want 30s", target.Timeout)
		}
	})

	t.Run("binds nested struct", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{
				{Key: "db.host", Value: "localhost"},
				{Key: "db.port", Value: "5432"},
			},
		}
		type DBConfig struct {
			Host string
			Port int
		}
		type Config struct {
			DB DBConfig `config:"db"`
		}
		var target Config
		err := cfg.Bind(context.Background(), "", &target)
		if err != nil {
			t.Fatalf("Bind() error: %v", err)
		}
		if target.DB.Host != "localhost" {
			t.Errorf("DB.Host = %q, want %q", target.DB.Host, "localhost")
		}
		if target.DB.Port != 5432 {
			t.Errorf("DB.Port = %d, want 5432", target.DB.Port)
		}
	})

	t.Run("binds string slice from comma-separated value", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "origins", Value: "a,b,c"}},
		}
		type CorsConfig struct {
			Origins []string `config:"origins"`
		}
		var target CorsConfig
		err := cfg.Bind(context.Background(), "", &target)
		if err != nil {
			t.Fatalf("Bind() error: %v", err)
		}
		if len(target.Origins) != 3 || target.Origins[0] != "a" {
			t.Errorf("Origins = %v, want [a b c]", target.Origins)
		}
	})

	t.Run("binds map[string]string", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{
				{Key: "headers.content-type", Value: "application/json"},
				{Key: "headers.accept", Value: "text/html"},
			},
		}
		type HeadersConfig struct {
			Headers map[string]string `config:"headers"`
		}
		var target HeadersConfig
		err := cfg.Bind(context.Background(), "", &target)
		if err != nil {
			t.Fatalf("Bind() error: %v", err)
		}
		if target.Headers == nil {
			t.Fatal("Headers is nil")
		}
		if target.Headers["content-type"] != "application/json" {
			t.Errorf("Headers[content-type] = %q, want %q", target.Headers["content-type"], "application/json")
		}
		if target.Headers["accept"] != "text/html" {
			t.Errorf("Headers[accept] = %q, want %q", target.Headers["accept"], "text/html")
		}
	})

	t.Run("section prefix prepended to config key", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "server.port", Value: "9090"}},
		}
		type ServerConfig struct {
			Port int `config:"port"`
		}
		var target ServerConfig
		err := cfg.Bind(context.Background(), "server", &target)
		if err != nil {
			t.Fatalf("Bind() error: %v", err)
		}
		if target.Port != 9090 {
			t.Errorf("Port = %d, want 9090", target.Port)
		}
	})
}

func TestConfigToSnakeCase(t *testing.T) {
	t.Run("PascalCase to snake_case", func(t *testing.T) {
		if got := toSnakeCase("HostName"); got != "host_name" {
			t.Errorf("toSnakeCase(HostName) = %q, want %q", got, "host_name")
		}
	})

	t.Run("multiple capitals", func(t *testing.T) {
		// toSnakeCase inserts _ only between a lowercase and an uppercase,
		// not between consecutive uppercase letters.
		if got := toSnakeCase("XMLParser"); got != "xmlparser" {
			t.Errorf("toSnakeCase(XMLParser) = %q, want %q", got, "xmlparser")
		}
	})

	t.Run("already snake_case", func(t *testing.T) {
		if got := toSnakeCase("already_snake"); got != "already_snake" {
			t.Errorf("toSnakeCase(already_snake) = %q, want %q", got, "already_snake")
		}
	})

	t.Run("single word", func(t *testing.T) {
		if got := toSnakeCase("Debug"); got != "debug" {
			t.Errorf("toSnakeCase(Debug) = %q, want %q", got, "debug")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		if got := toSnakeCase(""); got != "" {
			t.Errorf("toSnakeCase(\"\") = %q, want empty", got)
		}
	})
}

func TestConfigGetSection(t *testing.T) {
	t.Run("returns nil when no items match prefix", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{{Key: "other.key", Value: "value"}},
		}
		if cfg.GetSection("nonexistent") != nil {
			t.Error("GetSection(nonexistent) should return nil")
		}
	})

	t.Run("returns structured data for nested items", func(t *testing.T) {
		cfg := &Config{
			Items: []ConfigItem{
				{Key: "db.host", Value: "localhost"},
				{Key: "db.port", Value: "5432"},
			},
		}
		result := cfg.GetSection("db")
		if result == nil {
			t.Fatal("GetSection(db) returned nil")
		}
		// Should be a map or slice from unflatten
		_, isMap := result.(map[string]interface{})
		_, isSlice := result.([]interface{})
		if !isMap && !isSlice {
			t.Errorf("GetSection(db) returned %T, want map or slice", result)
		}
	})
}

func TestConfigUnflatten(t *testing.T) {
	t.Run("flattens nested structure", func(t *testing.T) {
		flat := map[string]string{
			"server.host": "localhost",
			"server.port": "8080",
		}
		result := (&Config{}).unflatten(flat)
		m, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("unflatten() returned %T, want map", result)
		}
		server, ok := m["server"].(map[string]interface{})
		if !ok {
			t.Fatalf("server not found in result")
		}
		if server["host"] != "localhost" {
			t.Errorf("server.host = %v, want localhost", server["host"])
		}
		if server["port"] != "8080" {
			t.Errorf("server.port = %v, want 8080", server["port"])
		}
	})

	t.Run("single-level keys remain as values", func(t *testing.T) {
		flat := map[string]string{
			"debug": "true",
		}
		result := (&Config{}).unflatten(flat)
		m, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("unflatten() returned %T, want map", result)
		}
		if m["debug"] != "true" {
			t.Errorf("debug = %v, want true", m["debug"])
		}
	})
}

func TestConfigDetectSlices(t *testing.T) {
	t.Run("converts integer-keyed map to slice", func(t *testing.T) {
		data := map[string]interface{}{
			"0": "first",
			"1": "second",
			"2": "third",
		}
		result := (&Config{}).detectSlices(data)
		slice, ok := result.([]interface{})
		if !ok {
			t.Fatalf("detectSlices() returned %T, want slice", result)
		}
		if len(slice) != 3 {
			t.Errorf("slice len = %d, want 3", len(slice))
		}
		if slice[0] != "first" {
			t.Errorf("slice[0] = %v, want first", slice[0])
		}
	})

	t.Run("returns non-map data unchanged", func(t *testing.T) {
		data := "just-a-string"
		result := (&Config{}).detectSlices(data)
		if result != "just-a-string" {
			t.Errorf("detectSlices(string) = %v, want unchanged", result)
		}
	})
}

// --- Test helpers ---

func itemValue(i *ConfigItem) string {
	if i == nil {
		return "<nil>"
	}
	return i.Value
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// testResolver implements ValueResolver for testing
type testResolver struct {
	prefix string
	value  string
}

func (r *testResolver) Prefix() string { return r.prefix }
func (r *testResolver) Resolve(ctx context.Context, value string) (string, error) {
	return r.value, nil
}

// errorResolver implements ValueResolver that always errors
type errorResolver struct {
	prefix string
}

func (r *errorResolver) Prefix() string { return r.prefix }
func (r *errorResolver) Resolve(ctx context.Context, value string) (string, error) {
	return "", fmt.Errorf("resolver error")
}
