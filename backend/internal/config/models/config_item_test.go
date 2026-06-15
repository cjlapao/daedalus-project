package models

import (
	"testing"
	"time"
)

func TestConfigItem(t *testing.T) {
	t.Run("Get returns Value", func(t *testing.T) {
		item := &ConfigItem{Key: "test", Value: "hello"}
		if got := item.Get(); got != "hello" {
			t.Errorf("Get() = %q, want %q", got, "hello")
		}
	})

	t.Run("Set updates Value", func(t *testing.T) {
		item := &ConfigItem{Key: "test"}
		item.Set("world")
		if item.Value != "world" {
			t.Errorf("Set() failed: Value = %q, want %q", item.Value, "world")
		}
	})

	t.Run("IsSet returns true when Value is non-empty", func(t *testing.T) {
		item := &ConfigItem{Key: "test", Value: "data"}
		if !item.IsSet() {
			t.Error("IsSet() = false, want true for non-empty Value")
		}
	})

	t.Run("IsSet returns false when Value is empty", func(t *testing.T) {
		item := &ConfigItem{Key: "test", Value: ""}
		if item.IsSet() {
			t.Error("IsSet() = true, want false for empty Value")
		}
	})

	t.Run("IsFlagSet returns true when FlagName is non-empty", func(t *testing.T) {
		item := &ConfigItem{Key: "test", FlagName: "debug"}
		if !item.IsFlagSet() {
			t.Error("IsFlagSet() = false, want true for non-empty FlagName")
		}
	})

	t.Run("IsFlagSet returns false when FlagName is empty", func(t *testing.T) {
		item := &ConfigItem{Key: "test", FlagName: ""}
		if item.IsFlagSet() {
			t.Error("IsFlagSet() = true, want false for empty FlagName")
		}
	})

	t.Run("IsEnvSet returns true when EnvName is non-empty", func(t *testing.T) {
		item := &ConfigItem{Key: "test", EnvName: "DEBUG"}
		if !item.IsEnvSet() {
			t.Error("IsEnvSet() = false, want true for non-empty EnvName")
		}
	})

	t.Run("IsEnvSet returns false when EnvName is empty", func(t *testing.T) {
		item := &ConfigItem{Key: "test", EnvName: ""}
		if item.IsEnvSet() {
			t.Error("IsEnvSet() = true, want false for empty EnvName")
		}
	})
}

func TestConfigItemGetBool(t *testing.T) {
	t.Run("returns true for 'true'", func(t *testing.T) {
		item := &ConfigItem{Value: "true"}
		if !item.GetBool() {
			t.Error("GetBool() = false, want true")
		}
	})

	t.Run("returns true for '1'", func(t *testing.T) {
		item := &ConfigItem{Value: "1"}
		if !item.GetBool() {
			t.Error("GetBool() = false, want true")
		}
	})

	t.Run("returns true for 'yes'", func(t *testing.T) {
		item := &ConfigItem{Value: "yes"}
		if !item.GetBool() {
			t.Error("GetBool() = false, want true")
		}
	})

	t.Run("returns true for 'on'", func(t *testing.T) {
		item := &ConfigItem{Value: "on"}
		if !item.GetBool() {
			t.Error("GetBool() = false, want true")
		}
	})

	t.Run("returns false for 'false'", func(t *testing.T) {
		item := &ConfigItem{Value: "false"}
		if item.GetBool() {
			t.Error("GetBool() = true, want false")
		}
	})

	t.Run("returns false for '0'", func(t *testing.T) {
		item := &ConfigItem{Value: "0"}
		if item.GetBool() {
			t.Error("GetBool() = true, want false")
		}
	})

	t.Run("returns false for nil receiver", func(t *testing.T) {
		var item *ConfigItem
		if item.GetBool() {
			t.Error("GetBool() on nil receiver = true, want false")
		}
	})
}

func TestConfigItemGetInt(t *testing.T) {
	t.Run("parses valid int", func(t *testing.T) {
		item := &ConfigItem{Value: "42"}
		if item.GetInt() != 42 {
			t.Errorf("GetInt() = %d, want 42", item.GetInt())
		}
	})

	t.Run("parses negative int", func(t *testing.T) {
		item := &ConfigItem{Value: "-7"}
		if item.GetInt() != -7 {
			t.Errorf("GetInt() = %d, want -7", item.GetInt())
		}
	})

	t.Run("returns 0 on parse error", func(t *testing.T) {
		item := &ConfigItem{Value: "not-a-number"}
		if item.GetInt() != 0 {
			t.Errorf("GetInt() = %d, want 0", item.GetInt())
		}
	})

	t.Run("returns 0 on empty value", func(t *testing.T) {
		item := &ConfigItem{Value: ""}
		if item.GetInt() != 0 {
			t.Errorf("GetInt() = %d, want 0", item.GetInt())
		}
	})

	t.Run("returns 0 on nil receiver", func(t *testing.T) {
		var item *ConfigItem
		if item.GetInt() != 0 {
			t.Errorf("GetInt() on nil = %d, want 0", item.GetInt())
		}
	})
}

func TestConfigItemGetDuration(t *testing.T) {
	t.Run("parses valid duration", func(t *testing.T) {
		item := &ConfigItem{Value: "5s"}
		want := 5 * time.Second
		if item.GetDuration() != want {
			t.Errorf("GetDuration() = %v, want %v", item.GetDuration(), want)
		}
	})

	t.Run("parses duration with minutes", func(t *testing.T) {
		item := &ConfigItem{Value: "2m30s"}
		want := 2*time.Minute + 30*time.Second
		if item.GetDuration() != want {
			t.Errorf("GetDuration() = %v, want %v", item.GetDuration(), want)
		}
	})

	t.Run("returns 0 on parse error", func(t *testing.T) {
		item := &ConfigItem{Value: "not-a-duration"}
		if item.GetDuration() != 0 {
			t.Errorf("GetDuration() = %v, want 0", item.GetDuration())
		}
	})

	t.Run("returns 0 on nil receiver", func(t *testing.T) {
		var item *ConfigItem
		if item.GetDuration() != 0 {
			t.Errorf("GetDuration() on nil = %v, want 0", item.GetDuration())
		}
	})
}

func TestConfigItemGetString(t *testing.T) {
	t.Run("returns Value for non-empty", func(t *testing.T) {
		item := &ConfigItem{Value: "hello"}
		if item.GetString() != "hello" {
			t.Errorf("GetString() = %q, want %q", item.GetString(), "hello")
		}
	})

	t.Run("returns empty string for empty Value", func(t *testing.T) {
		item := &ConfigItem{Value: ""}
		if item.GetString() != "" {
			t.Errorf("GetString() = %q, want empty", item.GetString())
		}
	})

	t.Run("returns empty string for nil receiver", func(t *testing.T) {
		var item *ConfigItem
		if item.GetString() != "" {
			t.Errorf("GetString() on nil = %q, want empty", item.GetString())
		}
	})
}

func TestConfigItemGetStringSlice(t *testing.T) {
	t.Run("splits comma-separated values", func(t *testing.T) {
		item := &ConfigItem{Value: "a,b,c"}
		want := []string{"a", "b", "c"}
		got := item.GetStringSlice()
		if len(got) != len(want) {
			t.Fatalf("GetStringSlice() len = %d, want %d", len(got), len(want))
		}
		for i, v := range want {
			if got[i] != v {
				t.Errorf("GetStringSlice()[%d] = %q, want %q", i, got[i], v)
			}
		}
	})

	t.Run("single value returns single-element slice", func(t *testing.T) {
		item := &ConfigItem{Value: "only"}
		got := item.GetStringSlice()
		if len(got) != 1 || got[0] != "only" {
			t.Errorf("GetStringSlice() = %v, want [only]", got)
		}
	})

	t.Run("empty value returns single empty string slice", func(t *testing.T) {
		item := &ConfigItem{Value: ""}
		got := item.GetStringSlice()
		// strings.Split("", ",") returns []string{""}, not []
		if len(got) != 1 || got[0] != "" {
			t.Errorf("GetStringSlice() = %v, want [\"\"]", got)
		}
	})

	t.Run("nil receiver returns empty slice", func(t *testing.T) {
		var item *ConfigItem
		got := item.GetStringSlice()
		if got == nil {
			t.Error("GetStringSlice() on nil = nil, want empty slice")
		}
		if len(got) != 0 {
			t.Errorf("GetStringSlice() on nil = %v, want empty slice", got)
		}
	})
}
