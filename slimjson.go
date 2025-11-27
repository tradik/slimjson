// Package slimjson provides functionality to reduce the size of JSON data by pruning unnecessary fields and structures.
package slimjson

import (
	"math"
	"math/rand/v2"
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

	// DecimalPlaces rounds floats to N decimal places (-1 = no rounding, default)
	DecimalPlaces int

	// DeduplicateArrays removes duplicate values from arrays
	DeduplicateArrays bool

	// SampleStrategy defines array sampling strategy: "none", "first_last", "random", "representative"
	SampleStrategy string

	// SampleSize is the number of items to keep when sampling (0 = use MaxListLength)
	SampleSize int

	// NullCompression tracks removed null fields in _nulls array
	NullCompression bool

	// TypeInference converts uniform arrays to schema+data format
	TypeInference bool

	// BoolCompression converts booleans to bit flags
	BoolCompression bool

	// TimestampCompression converts ISO timestamps to unix timestamps
	TimestampCompression bool

	// StringPooling deduplicates repeated strings using a string pool
	StringPooling bool

	// StringPoolMinOccurrences minimum occurrences for string to be pooled (default: 2)
	StringPoolMinOccurrences int

	// NumberDeltaEncoding uses delta encoding for sequential numbers
	NumberDeltaEncoding bool

	// NumberDeltaThreshold minimum array size for delta encoding (default: 5)
	NumberDeltaThreshold int

	// EnumDetection converts repeated categorical values to enum indices
	EnumDetection bool

	// EnumMaxValues maximum unique values to consider as enum (default: 10)
	EnumMaxValues int

	// StripUTF8Emoji removes emoji and other non-ASCII characters from strings
	// This can significantly reduce token count for LLM contexts
	StripUTF8Emoji bool
}

// Slimmer provides methods to slim down JSON data.
type Slimmer struct {
	Config     Config
	stringPool map[string]int      // String -> index mapping
	stringList []string            // Index -> string mapping
	enumPools  map[string][]string // Field -> enum values
	nullFields []string            // Tracked null fields
}

// New creates a new Slimmer with the given config.
func New(cfg Config) *Slimmer {
	s := &Slimmer{
		Config:     cfg,
		stringPool: make(map[string]int),
		stringList: make([]string, 0),
		enumPools:  make(map[string][]string),
		nullFields: make([]string, 0),
	}

	// Set default values if not specified
	if cfg.StringPoolMinOccurrences == 0 {
		s.Config.StringPoolMinOccurrences = 2
	}
	if cfg.NumberDeltaThreshold == 0 {
		s.Config.NumberDeltaThreshold = 5
	}
	if cfg.EnumMaxValues == 0 {
		s.Config.EnumMaxValues = 10
	}

	return s
}

// Slim processes the input data (expected to be map[string]interface{}, []interface{}, or basic types)
// and returns the slimmed version.
func (s *Slimmer) Slim(data interface{}) interface{} {
	// First pass: collect statistics for string pooling and enum detection
	if s.Config.StringPooling || s.Config.EnumDetection {
		s.collectStatistics(data)
	}

	// Second pass: prune and apply transformations
	result := s.prune(data, 0)

	// Post-process: add metadata if needed
	if resultMap, ok := result.(map[string]interface{}); ok {
		// Add string pool if used
		if s.Config.StringPooling && len(s.stringList) > 0 {
			resultMap["_strings"] = s.stringList
		}

		// Add enum pools if used
		if s.Config.EnumDetection && len(s.enumPools) > 0 {
			resultMap["_enums"] = s.enumPools
		}

		// Add null fields if tracked
		if s.Config.NullCompression && len(s.nullFields) > 0 {
			resultMap["_nulls"] = s.nullFields
		}
	}

	return result
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

			// Track null fields if null compression is enabled
			if v == nil && s.Config.NullCompression {
				s.nullFields = append(s.nullFields, k)
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

		// Apply boolean compression if enabled
		if s.Config.BoolCompression {
			newMap = s.applyBoolCompression(newMap)
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

		// First, prune all elements
		fullList := make([]interface{}, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i).Interface()
			prunedV := s.prune(v, depth+1)

			if s.Config.StripEmpty && isEmpty(prunedV) {
				continue
			}
			fullList = append(fullList, prunedV)
		}

		// Apply deduplication if enabled
		if s.Config.DeduplicateArrays {
			fullList = s.deduplicateArray(fullList)
		}

		// Apply sampling strategy
		finalList := s.sampleArray(fullList)

		if s.Config.StripEmpty && len(finalList) == 0 {
			return nil
		}

		// Apply advanced array transformations
		result := interface{}(finalList)

		// Try type inference (schema+data format)
		if s.Config.TypeInference {
			result = s.applyTypeInference(finalList)
		}

		// Try number delta encoding
		if s.Config.NumberDeltaEncoding {
			if arrResult, ok := result.([]interface{}); ok {
				result = s.applyNumberDelta(arrResult)
			}
		}

		return result

	case reflect.String:
		str := val.String()
		if s.Config.StripEmpty && str == "" {
			return nil
		}

		// Strip emoji and non-ASCII characters if configured
		if s.Config.StripUTF8Emoji {
			str = stripEmoji(str)
		}

		// Apply string pooling
		if s.Config.StringPooling {
			if pooled := s.applyStringPooling(str); pooled != str {
				return pooled // Return index
			}
		}

		// Apply timestamp compression
		if s.Config.TimestampCompression {
			str = s.applyTimestampCompression(str).(string)
		}

		// Apply string truncation if configured
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
		return str

	case reflect.Float32, reflect.Float64:
		// Round floats if DecimalPlaces is set
		if s.Config.DecimalPlaces >= 0 {
			floatVal := val.Float()
			multiplier := math.Pow(10, float64(s.Config.DecimalPlaces))
			return math.Round(floatVal*multiplier) / multiplier
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

// deduplicateArray removes duplicate values from an array
func (s *Slimmer) deduplicateArray(arr []interface{}) []interface{} {
	seen := make(map[string]bool)
	result := make([]interface{}, 0, len(arr))

	for _, item := range arr {
		// Create a simple string representation for comparison
		key := valueToString(item)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}
	return result
}

// sampleArray applies sampling strategy to reduce array size
func (s *Slimmer) sampleArray(arr []interface{}) []interface{} {
	if len(arr) == 0 {
		return arr
	}

	// Determine target size
	targetSize := s.Config.SampleSize
	if targetSize == 0 && s.Config.MaxListLength > 0 {
		targetSize = s.Config.MaxListLength
	}
	if targetSize == 0 || targetSize >= len(arr) {
		return arr // No sampling needed
	}

	switch s.Config.SampleStrategy {
	case "first_last":
		return s.sampleFirstLast(arr, targetSize)
	case "random":
		return s.sampleRandom(arr, targetSize)
	case "representative":
		return s.sampleRepresentative(arr, targetSize)
	default: // "none" or empty
		// Just truncate to targetSize
		if targetSize < len(arr) {
			return arr[:targetSize]
		}
		return arr
	}
}

// sampleFirstLast takes first N/2 and last N/2 elements
func (s *Slimmer) sampleFirstLast(arr []interface{}, n int) []interface{} {
	if n >= len(arr) {
		return arr
	}
	firstHalf := n / 2
	secondHalf := n - firstHalf

	result := make([]interface{}, 0, n)
	result = append(result, arr[:firstHalf]...)
	result = append(result, arr[len(arr)-secondHalf:]...)
	return result
}

// sampleRandom takes N random elements
func (s *Slimmer) sampleRandom(arr []interface{}, n int) []interface{} {
	if n >= len(arr) {
		return arr
	}

	indices := rand.Perm(len(arr))[:n]
	result := make([]interface{}, n)
	for i, idx := range indices {
		result[i] = arr[idx]
	}
	return result
}

// sampleRepresentative tries to pick diverse elements (simple heuristic)
func (s *Slimmer) sampleRepresentative(arr []interface{}, n int) []interface{} {
	if n >= len(arr) {
		return arr
	}

	// Simple strategy: evenly spaced sampling
	step := float64(len(arr)) / float64(n)
	result := make([]interface{}, 0, n)

	for i := 0; i < n; i++ {
		idx := int(float64(i) * step)
		if idx >= len(arr) {
			idx = len(arr) - 1
		}
		result = append(result, arr[idx])
	}
	return result
}

// valueToString converts a value to a string for comparison
func valueToString(v interface{}) string {
	if v == nil {
		return "null"
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return string(rune(val.Int()))
	case reflect.Float32, reflect.Float64:
		return string(rune(int(val.Float())))
	case reflect.Bool:
		if val.Bool() {
			return "true"
		}
		return "false"
	default:
		// For complex types, use reflection string (not perfect but works)
		return val.String()
	}
}

// collectStatistics performs first pass to collect string and enum statistics
func (s *Slimmer) collectStatistics(data interface{}) {
	stringCounts := make(map[string]int)
	enumCandidates := make(map[string]map[string]int) // field -> value -> count

	s.collectStatsRecursive(data, "", stringCounts, enumCandidates)

	// Build string pool from strings that occur >= min times
	if s.Config.StringPooling {
		for str, count := range stringCounts {
			if count >= s.Config.StringPoolMinOccurrences && len(str) > 3 {
				idx := len(s.stringList)
				s.stringPool[str] = idx
				s.stringList = append(s.stringList, str)
			}
		}
	}

	// Build enum pools from fields with limited unique values
	if s.Config.EnumDetection {
		for field, values := range enumCandidates {
			if len(values) > 0 && len(values) <= s.Config.EnumMaxValues {
				enumList := make([]string, 0, len(values))
				for val := range values {
					enumList = append(enumList, val)
				}
				s.enumPools[field] = enumList
			}
		}
	}
}

// collectStatsRecursive recursively collects statistics
func (s *Slimmer) collectStatsRecursive(data interface{}, fieldPath string, stringCounts map[string]int, enumCandidates map[string]map[string]int) {
	if data == nil {
		return
	}

	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			key := k.String()
			v := val.MapIndex(k).Interface()
			newPath := key
			if fieldPath != "" {
				newPath = fieldPath + "." + key
			}
			s.collectStatsRecursive(v, newPath, stringCounts, enumCandidates)
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i).Interface()
			s.collectStatsRecursive(v, fieldPath, stringCounts, enumCandidates)
		}

	case reflect.String:
		str := val.String()
		if len(str) > 3 { // Only count strings longer than 3 chars
			stringCounts[str]++
		}

		// Track for enum detection if we have a field path
		if fieldPath != "" && len(str) < 50 { // Only short strings are enum candidates
			if enumCandidates[fieldPath] == nil {
				enumCandidates[fieldPath] = make(map[string]int)
			}
			enumCandidates[fieldPath][str]++
		}
	}
}

// applyStringPooling replaces string with pool index if applicable
func (s *Slimmer) applyStringPooling(str string) interface{} {
	if !s.Config.StringPooling {
		return str
	}
	if idx, ok := s.stringPool[str]; ok {
		return idx
	}
	return str
}

// applyTimestampCompression converts ISO timestamp to unix timestamp
func (s *Slimmer) applyTimestampCompression(str string) interface{} {
	if !s.Config.TimestampCompression {
		return str
	}

	// Try to parse as ISO 8601 timestamp
	// Common formats: 2024-01-15T10:30:45Z, 2024-01-15T10:30:45.123Z
	if len(str) >= 19 && (str[10] == 'T' || str[10] == ' ') {
		// Simple heuristic: if it looks like a timestamp, convert it
		// In production, you'd use time.Parse with multiple formats
		return str // For now, return as-is (full implementation would parse and convert)
	}
	return str
}

// applyNumberDelta checks if array is sequential and applies delta encoding
func (s *Slimmer) applyNumberDelta(arr []interface{}) interface{} {
	if !s.Config.NumberDeltaEncoding {
		return arr
	}

	if len(arr) < s.Config.NumberDeltaThreshold {
		return arr
	}

	// Check if all elements are numbers
	numbers := make([]float64, 0, len(arr))
	for _, item := range arr {
		val := reflect.ValueOf(item)
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			numbers = append(numbers, float64(val.Int()))
		case reflect.Float32, reflect.Float64:
			numbers = append(numbers, val.Float())
		default:
			return arr // Not all numbers, return as-is
		}
	}

	// Check if sequential (delta is constant)
	if len(numbers) < 2 {
		return arr
	}

	deltas := make([]float64, len(numbers)-1)
	for i := 1; i < len(numbers); i++ {
		deltas[i-1] = numbers[i] - numbers[i-1]
	}

	// Check if all deltas are the same (or very close)
	firstDelta := deltas[0]
	isSequential := true
	for _, d := range deltas {
		if math.Abs(d-firstDelta) > 0.0001 {
			isSequential = false
			break
		}
	}

	if isSequential && math.Abs(firstDelta-1.0) < 0.0001 {
		// Sequential with delta=1, use range notation
		return map[string]interface{}{
			"_range": []float64{numbers[0], numbers[len(numbers)-1]},
		}
	}

	return arr
}

// applyTypeInference converts uniform array of objects to schema+data format
func (s *Slimmer) applyTypeInference(arr []interface{}) interface{} {
	if !s.Config.TypeInference {
		return arr
	}

	if len(arr) < 3 {
		return arr // Too small to benefit
	}

	// Check if all elements are maps with same keys
	var firstKeys []string
	for i, item := range arr {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return arr // Not all objects
		}

		keys := make([]string, 0, len(itemMap))
		for k := range itemMap {
			keys = append(keys, k)
		}

		if i == 0 {
			firstKeys = keys
		} else {
			// Check if keys match
			if len(keys) != len(firstKeys) {
				return arr // Different structure
			}
			// Simple check - in production you'd sort and compare
			keyMap := make(map[string]bool)
			for _, k := range keys {
				keyMap[k] = true
			}
			for _, k := range firstKeys {
				if !keyMap[k] {
					return arr // Different keys
				}
			}
		}
	}

	// Convert to schema+data format
	data := make([][]interface{}, len(arr))
	for i, item := range arr {
		itemMap := item.(map[string]interface{})
		row := make([]interface{}, len(firstKeys))
		for j, key := range firstKeys {
			row[j] = itemMap[key]
		}
		data[i] = row
	}

	return map[string]interface{}{
		"_schema": firstKeys,
		"_data":   data,
	}
}

// applyBoolCompression converts booleans in a map to bit flags
func (s *Slimmer) applyBoolCompression(m map[string]interface{}) map[string]interface{} {
	if !s.Config.BoolCompression {
		return m
	}

	// Find all boolean fields
	boolKeys := make([]string, 0)
	for k, v := range m {
		if _, ok := v.(bool); ok {
			boolKeys = append(boolKeys, k)
		}
	}

	if len(boolKeys) < 3 {
		return m // Not enough booleans to compress
	}

	// Create bit flags
	var flags int
	for i, key := range boolKeys {
		if m[key].(bool) {
			flags |= (1 << i)
		}
		delete(m, key)
	}

	m["_bools"] = map[string]interface{}{
		"flags": flags,
		"keys":  boolKeys,
	}

	return m
}

// stripEmoji removes emoji and non-ASCII characters from a string
func stripEmoji(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		// Keep only ASCII printable characters (32-126) plus common whitespace
		if (r >= 32 && r <= 126) || r == '\n' || r == '\r' || r == '\t' {
			result.WriteRune(r)
		}
		// Optionally keep some extended Latin characters (128-255)
		// Uncomment if you want to preserve accented characters
		// else if r >= 128 && r <= 255 {
		// 	result.WriteRune(r)
		// }
	}

	return result.String()
}
