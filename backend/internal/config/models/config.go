package models

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValueResolver defines an interface for resolving configuration values
type ValueResolver interface {
	Prefix() string
	Resolve(ctx context.Context, value string) (string, error)
}

type Config struct {
	Debug     bool         `json:"debug"`
	Items     []ConfigItem `json:"items"`
	resolvers map[string]ValueResolver
}

func (c *Config) RegisterResolver(resolver ValueResolver) {
	if c.resolvers == nil {
		c.resolvers = make(map[string]ValueResolver)
	}
	c.resolvers[resolver.Prefix()] = resolver
}

var variableRegex = regexp.MustCompile(`^\$\{\{\s*([a-zA-Z0-9_]+)::(.+?)\s*\}\}$`)

func (c *Config) ResolveValue(ctx context.Context, value string) string {
	if c.resolvers == nil {
		return value
	}

	match := variableRegex.FindStringSubmatch(value)
	if len(match) == 3 {
		prefix := match[1]
		content := match[2]

		if resolver, ok := c.resolvers[prefix]; ok {
			resolved, err := resolver.Resolve(ctx, content)
			if err != nil {
				// Log error? For now just return original or empty?
				// Maybe returning original allows debugging.
				return value
			}
			return resolved
		}
	}
	return value
}

func (c *Config) Get(key string) *ConfigItem {
	for _, item := range c.Items {
		if item.Key == key {
			return &item
		}
	}

	return nil
}

func (c *Config) IsDebug(ctx context.Context) bool {
	return c.GetBool(ctx, DebugKey, false)
}

func (c *Config) GetValue(ctx context.Context, key string, defaultValue interface{}) interface{} {
	item := c.Get(key)
	if item == nil || !item.IsSet() {
		return defaultValue
	}

	valStr := item.Value
	// If it's a string type, try to resolve it?
	// The problem is defaultValue interface{} doesn't tell us type easily.
	// But GetValue returns interface{}.
	// We'll resolve if it looks like a variable.
	return c.ResolveValue(ctx, valStr)
}

func (c *Config) GetString(ctx context.Context, key string, defaultValue string) string {
	item := c.Get(key)
	if item == nil || !item.IsSet() {
		return defaultValue
	}

	return c.ResolveValue(ctx, item.GetString())
}

func (c *Config) GetBool(ctx context.Context, key string, defaultValue bool) bool {
	item := c.Get(key)
	if item == nil || !item.IsSet() {
		return defaultValue
	}

	val := c.ResolveValue(ctx, item.Value)

	return val == "true" || val == "1" || val == "yes" || val == "on"
}

func (c *Config) GetInt(ctx context.Context, key string, defaultValue int) int {
	item := c.Get(key)
	if item == nil || !item.IsSet() {
		return defaultValue
	}

	val := c.ResolveValue(ctx, item.Value)

	v, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return v
}

func (c *Config) GetDuration(ctx context.Context, key string, defaultValue time.Duration) time.Duration {
	item := c.Get(key)
	if item == nil || !item.IsSet() {
		return defaultValue
	}

	val := c.ResolveValue(ctx, item.Value)

	dur, err := time.ParseDuration(val)
	if err != nil {
		return defaultValue
	}

	return dur
}

func (c *Config) Set(key string, value string) {
	for i, item := range c.Items {
		if item.Key == key {
			c.Items[i].Value = value
			return
		}
	}

	c.Items = append(c.Items, ConfigItem{Key: key, Value: value})
}

func (c *Config) StoragePath() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		userHome = "."
	}

	return filepath.Join(userHome, DefaultStoragePath)
}

// Bind maps a section of configuration to a struct.
func (c *Config) Bind(ctx context.Context, section string, target interface{}) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	elem := val.Elem()
	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typ.Field(i)

		configKey := fieldType.Tag.Get("config")
		if configKey == "" {
			configKey = toSnakeCase(fieldType.Name)
		}

		// If section is provided, prepend it
		fullKey := configKey
		if section != "" {
			fullKey = section + "." + configKey
		}

		if err := c.bindValue(ctx, fullKey, field); err != nil {
			return fmt.Errorf("failed to bind key %s: %w", fullKey, err)
		}
	}

	return nil
}

func (c *Config) bindValue(ctx context.Context, key string, field reflect.Value) error {
	item := c.Get(key)
	var resolvedValue string
	if item != nil && item.IsSet() {
		resolvedValue = c.ResolveValue(ctx, item.Value)
	}

	switch field.Kind() {
	case reflect.Struct:
		// Handle nested struct
		typ := field.Type()
		for i := 0; i < field.NumField(); i++ {
			structField := typ.Field(i)
			configTag := structField.Tag.Get("config")
			if configTag == "" {
				configTag = toSnakeCase(structField.Name)
			}

			fieldKey := key + "." + configTag
			if err := c.bindValue(ctx, fieldKey, field.Field(i)); err != nil {
				return err
			}
		}
	case reflect.String:
		if resolvedValue != "" {
			field.SetString(resolvedValue)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if resolvedValue != "" {
			if field.Type() == reflect.TypeOf(time.Duration(0)) {
				dur, _ := time.ParseDuration(resolvedValue)
				field.SetInt(int64(dur))
			} else {
				val, _ := strconv.Atoi(resolvedValue)
				field.SetInt(int64(val))
			}
		}
	case reflect.Bool:
		if resolvedValue != "" {
			isTrue := resolvedValue == "true" || resolvedValue == "1" || resolvedValue == "yes" || resolvedValue == "on"
			field.SetBool(isTrue)
		}
	case reflect.Float32, reflect.Float64:
		if resolvedValue != "" {
			val, err := strconv.ParseFloat(resolvedValue, 64)
			if err == nil {
				field.SetFloat(val)
			}
		}
	case reflect.Slice:
		// Handle slice of structs or primitives
		sliceType := field.Type()
		elemType := sliceType.Elem()
		isPtr := elemType.Kind() == reflect.Ptr
		if isPtr {
			elemType = elemType.Elem()
		}

		// Create a new slice
		newSlice := reflect.MakeSlice(sliceType, 0, 0)

		// Iterate until we find no more items
		index := 0
		for {
			indexKey := fmt.Sprintf("%s.%d", key, index)

			// If element is a struct, Check if any key starts with indexKey
			found := false
			if elemType.Kind() == reflect.Struct {
				// For structs, we check if there are any config items that start with "indexKey."
				prefix := indexKey + "."
				for _, cfgItem := range c.Items {
					if strings.HasPrefix(cfgItem.Key, prefix) {
						found = true
						break
					}
				}
			} else {
				// For primitives, check specific key
				idxItem := c.Get(indexKey)
				if idxItem != nil && idxItem.IsSet() {
					found = true
				}
			}

			if !found {
				break
			}

			// Create new element
			newElem := reflect.New(elemType)

			if elemType.Kind() == reflect.Struct {
				// Bind struct fields
				for i := 0; i < elemType.NumField(); i++ {
					structField := elemType.Field(i)
					configTag := structField.Tag.Get("config")
					if configTag == "" {
						configTag = toSnakeCase(structField.Name)
					}

					fieldKey := indexKey + "." + configTag
					if err := c.bindValue(ctx, fieldKey, newElem.Elem().Field(i)); err != nil {
						return err
					}
				}
			} else {
				// Bind primitive
				if err := c.bindValue(ctx, indexKey, newElem.Elem()); err != nil {
					return err
				}
			}

			if isPtr {
				newSlice = reflect.Append(newSlice, newElem)
			} else {
				newSlice = reflect.Append(newSlice, newElem.Elem())
			}
			index++
		}

		if newSlice.Len() > 0 {
			field.Set(newSlice)
		} else if field.Type().Elem().Kind() == reflect.String {
			// Fallback for simple string slice
			if resolvedValue != "" {
				vals := strings.Split(resolvedValue, ",")
				field.Set(reflect.ValueOf(vals))
			}
		}
	case reflect.Map:
		if field.Type().Key().Kind() != reflect.String {
			return nil // Only support string keys for now
		}

		// Create new map
		if field.IsNil() {
			field.Set(reflect.MakeMap(field.Type()))
		}

		// Find all items with prefix
		prefix := key + "."
		for _, item := range c.Items {
			if strings.HasPrefix(item.Key, prefix) {
				mapKey := strings.TrimPrefix(item.Key, prefix)
				// Resolve value here too
				resolvedMapVal := c.ResolveValue(ctx, item.Value)
				field.SetMapIndex(reflect.ValueOf(mapKey), reflect.ValueOf(resolvedMapVal))
			}
		}
	}
	return nil
}

func toSnakeCase(str string) string {
	matchFirstCap := true
	var sb strings.Builder
	for i, c := range str {
		if c >= 'A' && c <= 'Z' {
			if !matchFirstCap && i > 0 {
				sb.WriteRune('_')
			}
			sb.WriteRune(c + 32)
			matchFirstCap = true
		} else {
			sb.WriteRune(c)
			matchFirstCap = false
		}
	}
	return sb.String()
}

// GetSection returns a structured representation (map or slice) of a configuration section.
func (c *Config) GetSection(key string) interface{} {
	prefix := key + "."
	sectionItems := make(map[string]string)

	// Collect relevant items
	for _, item := range c.Items {
		if strings.HasPrefix(item.Key, prefix) {
			relativePath := strings.TrimPrefix(item.Key, prefix)
			sectionItems[relativePath] = item.Value
		} else if item.Key == key {
			// Exact match (leaf node requested as section?)
			return item.Value
		}
	}

	if len(sectionItems) == 0 {
		return nil
	}

	return c.unflatten(sectionItems)
}

func (c *Config) unflatten(flat map[string]string) interface{} {
	result := make(map[string]interface{})

	for path, value := range flat {
		parts := strings.Split(path, ".")
		current := result

		for i, part := range parts {
			if i == len(parts)-1 {
				// Last part, set value
				current[part] = value
			} else {
				// Intermediate part, create map if not exists
				if _, ok := current[part]; !ok {
					current[part] = make(map[string]interface{})
				}

				if nextMap, ok := current[part].(map[string]interface{}); ok {
					current = nextMap
				} else {
					// Conflict: current[part] is not a map (maybe it was a leaf in another key?)
					// Priority to latest or deeper structure?
					// For now, we overwrite if it's not a map, or ignore?
					// Let's assume consistent config structure.
					// If we encounter a primitive where we need a map, we replace it?
					// Or if we encounter a map where we need a primitive?
					// Let's just create a new map.
					newMap := make(map[string]interface{})
					current[part] = newMap
					current = newMap
				}
			}
		}
	}

	return c.detectSlices(result)
}

// detectSlices recursively converts maps with integer-like keys (0, 1, 2...) into slices.
func (c *Config) detectSlices(data interface{}) interface{} {
	m, isMap := data.(map[string]interface{})
	if !isMap {
		return data
	}

	// Check if keys are all integers and sequential starting from 0 (or close enough valid array)
	// Or simply if keys look like integers.
	// The flatten logic produced "0", "1", ...

	allKeysInt := true
	maxInt := -1
	count := 0

	for k := range m {
		i, err := strconv.Atoi(k)
		if err != nil {
			allKeysInt = false
			break
		}
		if i > maxInt {
			maxInt = i
		}
		count++
	}

	// Recursively process children first
	for k, v := range m {
		m[k] = c.detectSlices(v)
	}

	// If all keys are integers, convert to slice
	// We might have sparse arrays if some config is missing, but config usually is sequential.
	if allKeysInt && count > 0 {
		// Create slice
		slice := make([]interface{}, maxInt+1)
		for k, v := range m {
			i, _ := strconv.Atoi(k)
			slice[i] = v
		}
		return slice
	}

	return m
}
