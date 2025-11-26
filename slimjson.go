// Package slimjson provides functionality to reduce the size of JSON data by pruning unnecessary fields and structures.
package slimjson

import (
	"reflect"
	"strings"
)

// Config holds the configuration for the slimming process.
type Config struct {
	// MaxDepth is the maximum nesting depth allowed.
	// Objects/Arrays deeper than this will be truncated (removed or replaced).
	// 0 means no limit (or use a very high default if preferred, but let's say 0 is unlimited).
	// However, to "cut too deep nesting", we should probably default to something reasonable if 0.
	// Let's make 0 mean "unlimited" and user must set it, or we handle it in logic.
	MaxDepth int

	// MaxListLength is the maximum number of elements allowed in a list.
	// Elements beyond this count are removed.
	MaxListLength int

	// MaxStringLength is the maximum number of characters (runes) allowed in a string.
	// Strings longer than this will be truncated.
	MaxStringLength int

	// StripEmpty removes fields with null values, empty strings, empty arrays, or empty objects.
	StripEmpty bool

	// BlockList is a list of field names to remove.
	BlockList []string
}

// Slimmer provides methods to slim down JSON data.
type Slimmer struct {
	Config Config
}

// New creates a new Slimmer with the given config.
func New(cfg Config) *Slimmer {
	return &Slimmer{Config: cfg}
}

// Slim processes the input data (expected to be map[string]interface{}, []interface{}, or basic types)
// and returns the slimmed version.
func (s *Slimmer) Slim(data interface{}) interface{} {
	return s.prune(data, 0)
}

func (s *Slimmer) prune(data interface{}, depth int) interface{} {
	if data == nil {
		if s.Config.StripEmpty {
			return nil // Caller should handle nil removal if in object/array
		}
		return nil
	}

	// Check depth
	if s.Config.MaxDepth > 0 && depth >= s.Config.MaxDepth {
		return nil
	}

	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.Map:
		// Handle Object
		if val.Len() == 0 {
			if s.Config.StripEmpty {
				return nil
			}
			return data
		}

		newMap := make(map[string]interface{})
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key().String()
			v := iter.Value().Interface()

			// Check BlockList
			if s.isBlocked(k) {
				continue
			}

			prunedV := s.prune(v, depth+1)

			if s.Config.StripEmpty && isEmpty(prunedV) {
				continue
			}

			newMap[k] = prunedV
		}

		if s.Config.StripEmpty && len(newMap) == 0 {
			return nil
		}
		return newMap

	case reflect.Slice, reflect.Array:
		// Handle Array
		if val.Len() == 0 {
			if s.Config.StripEmpty {
				return nil
			}
			return data
		}

		length := val.Len()
		if s.Config.MaxListLength > 0 && length > s.Config.MaxListLength {
			length = s.Config.MaxListLength
		}

		newList := make([]interface{}, 0, length)
		for i := 0; i < length; i++ {
			v := val.Index(i).Interface()
			prunedV := s.prune(v, depth+1)

			if s.Config.StripEmpty && isEmpty(prunedV) {
				continue
			}
			newList = append(newList, prunedV)
		}

		if s.Config.StripEmpty && len(newList) == 0 {
			return nil
		}
		return newList

	case reflect.String:
		str := val.String()
		if s.Config.StripEmpty && str == "" {
			return nil
		}
		if s.Config.MaxStringLength > 0 {
			runes := []rune(str)
			if len(runes) > s.Config.MaxStringLength {
				// Truncate and add ellipsis to indicate truncation
				if s.Config.MaxStringLength > 3 {
					return string(runes[:s.Config.MaxStringLength-3]) + "..."
				}
				return string(runes[:s.Config.MaxStringLength])
			}
		}
		return data

	default:
		return data
	}
}

func (s *Slimmer) isBlocked(key string) bool {
	for _, blocked := range s.Config.BlockList {
		if strings.EqualFold(blocked, key) {
			return true
		}
	}
	return false
}

func isEmpty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Map, reflect.Slice, reflect.Array:
		return v.Len() == 0
	}
	return false
}
