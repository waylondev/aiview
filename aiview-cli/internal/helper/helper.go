// Package helper provides common utility functions for parsing map[string]interface{} responses.
package helper

import "strconv"

// GetString retrieves a string value from a nested map.
func GetString(m map[string]interface{}, keys ...string) string {
	if m == nil {
		return ""
	}
	val, ok := m[keys[0]]
	if !ok {
		return ""
	}
	if len(keys) == 1 {
		s, _ := val.(string)
		return s
	}
	sub, ok := val.(map[string]interface{})
	if !ok {
		return ""
	}
	return GetString(sub, keys[1:]...)
}

// GetInt retrieves an int value from a nested map.
func GetInt(m map[string]interface{}, keys ...string) int {
	if m == nil {
		return 0
	}
	val, ok := m[keys[0]]
	if !ok {
		return 0
	}
	if len(keys) == 1 {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case string:
			n, _ := strconv.Atoi(v)
			return n
		}
		return 0
	}
	sub, ok := val.(map[string]interface{})
	if !ok {
		return 0
	}
	return GetInt(sub, keys[1:]...)
}

// GetMap retrieves a map value from a map.
func GetMap(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	val, ok := m[key]
	if !ok {
		return nil
	}
	sub, _ := val.(map[string]interface{})
	return sub
}

// GetSlice retrieves a slice value from a map.
func GetSlice(m map[string]interface{}, key string) []interface{} {
	if m == nil {
		return nil
	}
	val, ok := m[key]
	if !ok {
		return nil
	}
	sub, _ := val.([]interface{})
	return sub
}