// Package helper provides common utility functions for parsing map[string]interface{} responses.
package helper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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

// GetFloat retrieves a float64 value from a map.
func GetFloat(m map[string]interface{}, key string) float64 {
	if m == nil {
		return 0
	}
	val, ok := m[key]
	if !ok {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	}
	return 0
}

// GetBool retrieves a bool value from a map.
func GetBool(m map[string]interface{}, key string) bool {
	if m == nil {
		return false
	}
	val, ok := m[key]
	if !ok {
		return false
	}
	b, _ := val.(bool)
	return b
}

// FormatDuration converts seconds to a formatted duration string.
func FormatDuration(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}
	if seconds >= 3600 {
		h := seconds / 3600
		m := (seconds % 3600) / 60
		s := seconds % 60
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	m := seconds / 60
	s := seconds % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

// StripHTML removes HTML tags from a string.
func StripHTML(text string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return strings.TrimSpace(re.ReplaceAllString(text, ""))
}