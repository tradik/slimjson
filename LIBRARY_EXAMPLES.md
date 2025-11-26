# SlimJSON Library Examples for Developers

Complete guide for using SlimJSON as a Go library in your applications.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Built-in Profiles](#built-in-profiles)
- [Custom Profiles from File](#custom-profiles-from-file)
- [Real-World Examples](#real-world-examples)
- [Best Practices](#best-practices)

## Installation

```bash
go get github.com/tradik/slimjson
```

## Quick Start

### Basic Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/tradik/slimjson"
)

func main() {
	// Your data
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    123,
			"name":  "John Doe",
			"email": "john@example.com",
			"metadata": map[string]interface{}{
				"created_at": "2024-01-01",
				"updated_at": "2024-01-15",
			},
		},
		"items": []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	// Configure compression
	cfg := slimjson.Config{
		MaxDepth:      2,
		MaxListLength: 5,
		StripEmpty:    true,
	}

	// Process
	slimmer := slimjson.New(cfg)
	result := slimmer.Slim(data)

	// Output
	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
}
```

## Configuration

### Config Structure

```go
type Config struct {
	// Basic Options
	MaxDepth        int      // Maximum nesting depth (0 = unlimited)
	MaxListLength   int      // Maximum array length (0 = unlimited)
	MaxStringLength int      // Maximum string length (0 = unlimited)
	StripEmpty      bool     // Remove nulls, empty strings, empty arrays/objects
	BlockList       []string // Field names to remove (case-insensitive)
	
	// Optimization Options
	DecimalPlaces     int    // Round floats (-1 = no rounding)
	DeduplicateArrays bool   // Remove duplicate array values
	SampleStrategy    string // "none", "first_last", "random", "representative"
	SampleSize        int    // Items to keep when sampling
	
	// Advanced Compression
	NullCompression          bool // Track removed nulls in _nulls
	TypeInference            bool // Convert arrays to schema+data
	BoolCompression          bool // Convert booleans to bit flags
	TimestampCompression     bool // Convert ISO to unix timestamps
	StringPooling            bool // Deduplicate repeated strings
	StringPoolMinOccurrences int  // Min occurrences for pooling (default: 2)
	NumberDeltaEncoding      bool // Delta encoding for sequences
	NumberDeltaThreshold     int  // Min array size for delta (default: 5)
	EnumDetection            bool // Convert categorical values to enums
	EnumMaxValues            int  // Max unique values for enum (default: 10)
}
```

### Configuration Examples

#### Minimal Configuration

```go
cfg := slimjson.Config{
	MaxDepth:   5,
	StripEmpty: true,
}
```

#### Balanced Configuration

```go
cfg := slimjson.Config{
	MaxDepth:          5,
	MaxListLength:     20,
	StripEmpty:        true,
	DecimalPlaces:     2,
	DeduplicateArrays: true,
}
```

#### Maximum Compression

```go
cfg := slimjson.Config{
	MaxDepth:                 3,
	MaxListLength:            10,
	StripEmpty:               true,
	DecimalPlaces:            2,
	DeduplicateArrays:        true,
	SampleStrategy:           "representative",
	SampleSize:               20,
	NullCompression:          true,
	TypeInference:            true,
	BoolCompression:          true,
	TimestampCompression:     true,
	StringPooling:            true,
	StringPoolMinOccurrences: 2,
	NumberDeltaEncoding:      true,
	NumberDeltaThreshold:     5,
	EnumDetection:            true,
	EnumMaxValues:            10,
	BlockList:                []string{"metadata", "debug", "trace"},
}
```

## Built-in Profiles

### Using Built-in Profiles

```go
package main

import (
	"github.com/tradik/slimjson"
)

func main() {
	// Get all built-in profiles
	profiles := slimjson.GetBuiltinProfiles()
	
	// Available profiles: light, medium, aggressive, ai-optimized
	lightCfg := profiles["light"]
	mediumCfg := profiles["medium"]
	aggressiveCfg := profiles["aggressive"]
	aiCfg := profiles["ai-optimized"]
	
	// Use a profile
	slimmer := slimjson.New(mediumCfg)
	result := slimmer.Slim(data)
}
```

### Profile Characteristics

```go
// Light - Preserve most data
profiles["light"] = Config{
	MaxDepth:      10,
	MaxListLength: 20,
	StripEmpty:    true,
}

// Medium - Balanced compression
profiles["medium"] = Config{
	MaxDepth:      5,
	MaxListLength: 10,
	StripEmpty:    true,
}

// Aggressive - Maximum reduction
profiles["aggressive"] = Config{
	MaxDepth:      3,
	MaxListLength: 5,
	StripEmpty:    true,
	BlockList:     []string{"description", "summary", "comment", "notes", "bio", "readme"},
}

// AI-Optimized - For LLM contexts
profiles["ai-optimized"] = Config{
	MaxDepth:      4,
	MaxListLength: 8,
	StripEmpty:    true,
	BlockList:     []string{"avatar_url", "gravatar_id", "url", "html_url", "*_url"},
}
```

## Custom Profiles from File

### Loading .slimjson File

```go
package main

import (
	"fmt"
	"github.com/tradik/slimjson"
)

func main() {
	// Load from .slimjson in current dir or home dir
	customProfiles, err := slimjson.LoadConfigFile()
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
		customProfiles = make(map[string]slimjson.Config)
	}
	
	// Use custom profile
	if cfg, ok := customProfiles["my-api-profile"]; ok {
		slimmer := slimjson.New(cfg)
		result := slimmer.Slim(data)
	}
}
```

### Parsing Specific File

```go
package main

import (
	"github.com/tradik/slimjson"
)

func main() {
	// Parse specific config file
	profiles, err := slimjson.ParseConfigFile("/etc/myapp/.slimjson")
	if err != nil {
		panic(err)
	}
	
	cfg := profiles["production"]
	slimmer := slimjson.New(cfg)
	result := slimmer.Slim(data)
}
```

### Combining Built-in and Custom Profiles

```go
package main

import (
	"github.com/tradik/slimjson"
)

func main() {
	// Start with built-in profiles
	allProfiles := slimjson.GetBuiltinProfiles()
	
	// Load custom profiles
	customProfiles, _ := slimjson.LoadConfigFile()
	
	// Merge (custom profiles override built-in)
	for name, cfg := range customProfiles {
		allProfiles[name] = cfg
	}
	
	// Now use any profile
	cfg := allProfiles["my-custom-profile"]
	slimmer := slimjson.New(cfg)
	result := slimmer.Slim(data)
}
```

## Real-World Examples

### Example 1: REST API Response Compression

```go
package main

import (
	"encoding/json"
	"net/http"
	"github.com/tradik/slimjson"
)

// Middleware for response compression
func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if client wants compressed response
		if r.Header.Get("X-Compress-Response") != "true" {
			next.ServeHTTP(w, r)
			return
		}
		
		// Capture response
		rec := &ResponseRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		
		// Compress if JSON
		if rec.ContentType == "application/json" {
			var data interface{}
			json.Unmarshal(rec.Body, &data)
			
			cfg := slimjson.Config{
				MaxDepth:          5,
				MaxListLength:     50,
				StripEmpty:        true,
				DecimalPlaces:     2,
				DeduplicateArrays: true,
				BlockList:         []string{"internal_id", "debug"},
			}
			
			slimmer := slimjson.New(cfg)
			compressed := slimmer.Slim(data)
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(compressed)
		}
	})
}
```

### Example 2: Database Export Optimization

```go
package main

import (
	"database/sql"
	"encoding/json"
	"github.com/tradik/slimjson"
)

func ExportTableToJSON(db *sql.DB, tableName string) ([]byte, error) {
	// Query database
	rows, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	// Convert to slice of maps
	var records []map[string]interface{}
	// ... populate records from rows ...
	
	// Compress for export
	cfg := slimjson.Config{
		MaxDepth:          10,
		MaxListLength:     1000,
		StripEmpty:        true,
		DecimalPlaces:     4,
		TypeInference:     true,
		NumberDeltaEncoding: true,
	}
	
	data := map[string]interface{}{
		"table":   tableName,
		"records": records,
	}
	
	slimmer := slimjson.New(cfg)
	compressed := slimmer.Slim(data)
	
	return json.Marshal(compressed)
}
```

### Example 3: LLM Context Preparation

```go
package main

import (
	"encoding/json"
	"github.com/tradik/slimjson"
)

type LLMContext struct {
	SystemPrompt string
	UserData     interface{}
	History      []interface{}
}

func PrepareLLMContext(ctx LLMContext) (string, error) {
	// Combine all context data
	data := map[string]interface{}{
		"system":  ctx.SystemPrompt,
		"user":    ctx.UserData,
		"history": ctx.History,
	}
	
	// Maximum compression for token efficiency
	cfg := slimjson.Config{
		MaxDepth:             4,
		MaxListLength:        15,
		StripEmpty:           true,
		DecimalPlaces:        2,
		DeduplicateArrays:    true,
		SampleStrategy:       "representative",
		SampleSize:           10,
		StringPooling:        true,
		TypeInference:        true,
		BoolCompression:      true,
		TimestampCompression: true,
		BlockList:            []string{"avatar_url", "url", "html_url", "metadata"},
	}
	
	slimmer := slimjson.New(cfg)
	compressed := slimmer.Slim(data)
	
	jsonBytes, err := json.Marshal(compressed)
	if err != nil {
		return "", err
	}
	
	return string(jsonBytes), nil
}
```

### Example 4: Analytics Data Processing

```go
package main

import (
	"encoding/json"
	"github.com/tradik/slimjson"
)

type AnalyticsEvent struct {
	Timestamp string
	UserID    int
	Event     string
	Data      map[string]interface{}
}

func ProcessAnalytics(events []AnalyticsEvent) ([]byte, error) {
	// Convert to interface slice
	data := make([]interface{}, len(events))
	for i, e := range events {
		data[i] = map[string]interface{}{
			"timestamp": e.Timestamp,
			"user_id":   e.UserID,
			"event":     e.Event,
			"data":      e.Data,
		}
	}
	
	// Compress with sampling for large datasets
	cfg := slimjson.Config{
		MaxDepth:            3,
		MaxListLength:       100,
		StripEmpty:          true,
		DecimalPlaces:       2,
		SampleStrategy:      "representative",
		SampleSize:          50,
		TypeInference:       true,
		TimestampCompression: true,
		EnumDetection:       true,
	}
	
	wrapper := map[string]interface{}{
		"events": data,
	}
	
	slimmer := slimjson.New(cfg)
	compressed := slimmer.Slim(wrapper)
	
	return json.Marshal(compressed)
}
```

### Example 5: Configuration Management

```go
package main

import (
	"github.com/tradik/slimjson"
)

type AppConfig struct {
	profiles map[string]slimjson.Config
}

func NewAppConfig() *AppConfig {
	ac := &AppConfig{
		profiles: slimjson.GetBuiltinProfiles(),
	}
	
	// Load custom profiles
	customProfiles, _ := slimjson.LoadConfigFile()
	for name, cfg := range customProfiles {
		ac.profiles[name] = cfg
	}
	
	return ac
}

func (ac *AppConfig) GetProfile(name string) (slimjson.Config, bool) {
	cfg, ok := ac.profiles[name]
	return cfg, ok
}

func (ac *AppConfig) Compress(data interface{}, profileName string) interface{} {
	cfg, ok := ac.GetProfile(profileName)
	if !ok {
		// Fallback to medium profile
		cfg = ac.profiles["medium"]
	}
	
	slimmer := slimjson.New(cfg)
	return slimmer.Slim(data)
}
```

## Best Practices

### 1. Choose the Right Profile

```go
// For API responses - preserve structure
cfg := profiles["light"]

// For LLM contexts - maximize compression
cfg := profiles["ai-optimized"]

// For analytics - balance size and detail
cfg := profiles["medium"]
```

### 2. Test Compression Results

```go
func TestCompression(data interface{}, cfg slimjson.Config) {
	original, _ := json.Marshal(data)
	
	slimmer := slimjson.New(cfg)
	compressed := slimmer.Slim(data)
	result, _ := json.Marshal(compressed)
	
	ratio := float64(len(result)) / float64(len(original)) * 100
	fmt.Printf("Compression: %.1f%% of original size\n", ratio)
}
```

### 3. Handle Metadata Fields

```go
result := slimmer.Slim(data)

// Check for metadata
if resultMap, ok := result.(map[string]interface{}); ok {
	if strings, ok := resultMap["_strings"]; ok {
		fmt.Printf("String pool: %v\n", strings)
	}
	if nulls, ok := resultMap["_nulls"]; ok {
		fmt.Printf("Removed nulls: %v\n", nulls)
	}
}
```

### 4. Error Handling

```go
func SafeCompress(data interface{}, profileName string) (interface{}, error) {
	profiles, err := slimjson.LoadConfigFile()
	if err != nil {
		// Fallback to built-in profiles
		profiles = slimjson.GetBuiltinProfiles()
	}
	
	cfg, ok := profiles[profileName]
	if !ok {
		return nil, fmt.Errorf("profile not found: %s", profileName)
	}
	
	slimmer := slimjson.New(cfg)
	return slimmer.Slim(data), nil
}
```

### 5. Performance Considerations

```go
// Reuse slimmer for same configuration
slimmer := slimjson.New(cfg)

for _, item := range items {
	result := slimmer.Slim(item)
	// Process result...
}
```

## Additional Resources

- [Main README](README.md) - CLI usage and installation
- [Examples](EXAMPLES.md) - CLI examples and patterns
- [.slimjson.example](.slimjson.example) - Configuration file examples
- [API Documentation](https://pkg.go.dev/github.com/tradik/slimjson) - Full API reference
