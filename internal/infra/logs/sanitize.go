package logs

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var sensitiveKeys = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"authorization",
	"cookie",
	"session",
}

func sanitizeArgs(args ...any) []any {
	if len(args) == 0 {
		return nil
	}
	safe := make([]any, 0, len(args))
	for i := 0; i < len(args); i += 2 {
		key := args[i]
		safe = append(safe, key)
		if i+1 >= len(args) {
			continue
		}
		name, _ := key.(string)
		if isSensitiveName(name) {
			safe = append(safe, "[REDACTED]")
			continue
		}
		safe = append(safe, sanitize(args[i+1]))
	}
	return safe
}

func sanitize(value any) any {
	return sanitizeValue(reflect.ValueOf(value), "")
}

func sanitizeValue(v reflect.Value, name string) any {
	if !v.IsValid() {
		return nil
	}
	if isSensitiveName(name) {
		return "[REDACTED]"
	}
	if shouldPreserveValue(v) {
		return v.Interface()
	}

	switch v.Kind() {
	case reflect.Interface, reflect.Pointer:
		if v.IsNil() {
			return nil
		}
		return sanitizeValue(v.Elem(), name)
	case reflect.Struct:
		m := make(map[string]any)
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue
			}
			fieldName := jsonFieldName(field)
			if fieldName == "-" {
				continue
			}
			if field.Anonymous && fieldName == "" {
				embedded := sanitizeValue(v.Field(i), field.Name)
				if sub, ok := embedded.(map[string]any); ok {
					for k, val := range sub {
						m[k] = val
					}
					continue
				}
			}
			if fieldName == "" {
				fieldName = field.Name
			}
			m[fieldName] = sanitizeValue(v.Field(i), fieldName)
		}
		return m
	case reflect.Map:
		if v.IsNil() {
			return nil
		}
		m := make(map[string]any, v.Len())
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			keyName := fmt.Sprint(key.Interface())
			m[keyName] = sanitizeValue(iter.Value(), keyName)
		}
		return m
	case reflect.Slice, reflect.Array:
		length := v.Len()
		items := make([]any, 0, length)
		for i := 0; i < length; i++ {
			items = append(items, sanitizeValue(v.Index(i), name))
		}
		return items
	default:
		return v.Interface()
	}
}

func shouldPreserveValue(v reflect.Value) bool {
	if !v.CanInterface() {
		return false
	}
	value := v.Interface()
	if _, ok := value.(json.Marshaler); ok {
		return true
	}
	if _, ok := value.(encoding.TextMarshaler); ok {
		return true
	}
	if _, ok := value.(fmt.Stringer); ok && v.Kind() != reflect.Struct {
		return true
	}
	return false
}

func jsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return ""
	}
	name := strings.Split(tag, ",")[0]
	return strings.TrimSpace(name)
}

func isSensitiveName(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return false
	}
	for _, token := range sensitiveKeys {
		if strings.Contains(name, token) {
			return true
		}
	}
	return false
}
