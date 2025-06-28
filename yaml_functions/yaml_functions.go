package yaml_functions

import (
	"strings"
)

// GetCaseInsensitiveMap searches for a key in the map (case-insensitively) and returns its value as a map[string]interface{}.
// Returns nil if the key is not found or the value is not a map.
func GetCaseInsensitiveMap(m map[string]interface{}, key string) map[string]interface{} {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if result, ok := v.(map[string]interface{}); ok {
				return result
			}
		}
	}
	return nil
}

// GetCaseInsensitiveList searches for a key in the map (case-insensitively) and returns its value as a []string.
// Only string elements are included in the returned slice.
// Returns nil if the key is not found or the value is not a list.
func GetCaseInsensitiveList(m map[string]interface{}, key string) []string {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			raw, ok := v.([]interface{})
			if !ok {
				return nil
			}
			var result []string
			for _, val := range raw {
				if s, ok := val.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return nil
}

// GetCaseInsensitiveString searches for a key in the map (case-insensitively) and returns its value as a string.
// Returns an empty string if the key is not found or the value is not a string.
func GetCaseInsensitiveString(m map[string]interface{}, key string) string {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if str, ok := v.(string); ok {
				return str
			}
		}
	}
	return ""
}

// GetNestedString attempts to retrieve a string value for the given key from the map.
// First, it tries to get the string directly using GetCaseInsensitiveString.
// If not found, it then checks if the key maps to a nested map and returns the first string value from that map.
// Returns an empty string if no suitable value is found.
func GetNestedString(m map[string]interface{}, key string) string {
	if val := GetCaseInsensitiveString(m, key); val != "" {
		return val
	}
	if sub := GetCaseInsensitiveMap(m, key); sub != nil {
		for _, v := range sub {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

// GetNestedMap is a convenience wrapper that retrieves a nested map for a given key using case-insensitive matching.
// Returns nil if the key is not found or the value is not a map.
func GetNestedMap(m map[string]interface{}, key string) map[string]interface{} {
	return GetCaseInsensitiveMap(m, key)
}