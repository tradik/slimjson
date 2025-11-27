// Package slimjson provides JSON compression and optimization for AI/LLM contexts.
//
// SlimJSON is designed to "slim down" JSON data before sending it to AI models
// or other bandwidth-constrained systems. It removes unnecessary fields, truncates
// deep nesting, shortens lists and strings, and strips empty values to create an
// optimized "diet" version of your JSON.
//
// # Features
//
//   - Prune unnecessary fields by name (blocklist)
//   - Truncate deep nesting automatically
//   - Shorten arrays to specified length
//   - Truncate strings (UTF-8 aware)
//   - Strip empty values (null, "", [], {})
//   - Custom compression profiles
//   - Advanced compression (string pooling, type inference, etc.)
//   - High performance (16-47µs per operation)
//   - Proven results (24-98% size reduction)
//
// # Quick Start
//
// Basic usage with default configuration:
//
//	package main
//
//	import (
//	    "encoding/json"
//	    "fmt"
//	    "github.com/tradik/slimjson"
//	)
//
//	func main() {
//	    data := map[string]interface{}{
//	        "name": "Test",
//	        "details": map[string]interface{}{
//	            "deep": map[string]interface{}{
//	                "too_deep": "value",
//	            },
//	        },
//	        "list": []interface{}{1, 2, 3, 4, 5},
//	        "empty": "",
//	    }
//
//	    cfg := slimjson.Config{
//	        MaxDepth:      2,
//	        MaxListLength: 3,
//	        StripEmpty:    true,
//	    }
//
//	    slimmer := slimjson.New(cfg)
//	    result := slimmer.Slim(data)
//
//	    out, _ := json.MarshalIndent(result, "", "  ")
//	    fmt.Println(string(out))
//	}
//
// # Using Built-in Profiles
//
// SlimJSON provides four built-in compression profiles:
//
//	profiles := slimjson.GetBuiltinProfiles()
//
//	// Light - preserves most data
//	lightCfg := profiles["light"]
//
//	// Medium - balanced compression
//	mediumCfg := profiles["medium"]
//
//	// Aggressive - maximum reduction
//	aggressiveCfg := profiles["aggressive"]
//
//	// AI-Optimized - for LLM contexts
//	aiCfg := profiles["ai-optimized"]
//
//	slimmer := slimjson.New(mediumCfg)
//	result := slimmer.Slim(data)
//
// # Configuration Options
//
// The Config struct provides comprehensive control over compression:
//
//	type Config struct {
//	    // Basic options
//	    MaxDepth        int      // Maximum nesting depth (0 = unlimited)
//	    MaxListLength   int      // Maximum array length (0 = unlimited)
//	    MaxStringLength int      // Maximum string length (0 = unlimited)
//	    StripEmpty      bool     // Remove nulls, empty strings, arrays, objects
//	    BlockList       []string // Field names to remove
//
//	    // Optimization options
//	    DecimalPlaces     int    // Round floats to N decimal places
//	    DeduplicateArrays bool   // Remove duplicate array values
//	    SampleStrategy    string // Array sampling strategy
//	    SampleSize        int    // Number of items when sampling
//
//	    // Advanced compression
//	    NullCompression          bool // Track removed nulls
//	    TypeInference            bool // Convert arrays to schema+data
//	    BoolCompression          bool // Convert booleans to bit flags
//	    TimestampCompression     bool // Convert ISO to unix timestamps
//	    StringPooling            bool // Deduplicate repeated strings
//	    StringPoolMinOccurrences int  // Min occurrences for pooling
//	    NumberDeltaEncoding      bool // Delta encoding for sequences
//	    NumberDeltaThreshold     int  // Min array size for delta
//	    EnumDetection            bool // Convert categorical values to enums
//	    EnumMaxValues            int  // Max unique values for enum
//	    StripUTF8Emoji           bool // Remove emoji and non-ASCII characters
//	}
//
// # Advanced Compression
//
// Enable advanced compression features for maximum size reduction:
//
//	cfg := slimjson.Config{
//	    MaxDepth:             4,
//	    MaxListLength:        15,
//	    StripEmpty:           true,
//	    DecimalPlaces:        2,
//	    StringPooling:        true,
//	    TypeInference:        true,
//	    BoolCompression:      true,
//	    TimestampCompression: true,
//	    StripUTF8Emoji:       true, // Remove emoji for LLM contexts
//	}
//
//	slimmer := slimjson.New(cfg)
//	result := slimmer.Slim(data)
//
//	// Result may contain metadata fields:
//	// - _strings: String pool (if StringPooling enabled)
//	// - _enums: Enum mappings (if EnumDetection enabled)
//	// - _nulls: Tracked null fields (if NullCompression enabled)
//
// # Emoji and Non-ASCII Character Removal
//
// Remove emoji and non-ASCII characters to reduce token count for LLMs:
//
//	cfg := slimjson.Config{
//	    StripUTF8Emoji: true,
//	}
//
//	data := map[string]interface{}{
//	    "message": "Hello 👋 World 🌍!",
//	    "status":  "✅ Completed",
//	}
//
//	slimmer := slimjson.New(cfg)
//	result := slimmer.Slim(data)
//	// Result: {"message": "Hello  World !", "status": " Completed"}
//
// # Custom Profiles from File
//
// Load custom profiles from a .slimjson configuration file:
//
//	// Load from current dir or home dir
//	customProfiles, err := slimjson.LoadConfigFile()
//	if err != nil {
//	    // Handle error
//	}
//
//	// Use custom profile
//	cfg := customProfiles["my-custom-profile"]
//	slimmer := slimjson.New(cfg)
//	result := slimmer.Slim(data)
//
// Or parse a specific config file:
//
//	profiles, err := slimjson.ParseConfigFile("/path/to/.slimjson")
//	if err != nil {
//	    // Handle error
//	}
//
//	cfg := profiles["production"]
//	slimmer := slimjson.New(cfg)
//	result := slimmer.Slim(data)
//
// # Real-World Examples
//
// API Response Compression:
//
//	func compressAPIResponse(data interface{}) ([]byte, error) {
//	    cfg := slimjson.Config{
//	        MaxDepth:          5,
//	        MaxListLength:     20,
//	        StripEmpty:        true,
//	        DecimalPlaces:     2,
//	        DeduplicateArrays: true,
//	        BlockList:         []string{"internal_id", "metadata", "debug"},
//	    }
//
//	    slimmer := slimjson.New(cfg)
//	    compressed := slimmer.Slim(data)
//
//	    return json.Marshal(compressed)
//	}
//
// LLM Context Optimization:
//
//	func prepareLLMContext(data interface{}) ([]byte, error) {
//	    cfg := slimjson.Config{
//	        MaxDepth:             4,
//	        MaxListLength:        15,
//	        StripEmpty:           true,
//	        StringPooling:        true,
//	        TypeInference:        true,
//	        BoolCompression:      true,
//	        TimestampCompression: true,
//	        BlockList:            []string{"avatar_url", "url", "html_url"},
//	    }
//
//	    slimmer := slimjson.New(cfg)
//	    compressed := slimmer.Slim(data)
//
//	    return json.Marshal(compressed)
//	}
//
// # HTTP Daemon Mode
//
// The slimjson CLI can run as an HTTP daemon:
//
//	# Start daemon on port 8080
//	slimjson -d
//
//	# Start on custom port
//	slimjson -d -port 3000
//
// API Endpoints:
//   - GET  /health - Health check
//   - GET  /profiles - List available profiles
//   - POST /slim?profile=<name> - Compress JSON
//
// Example API usage:
//
//	curl -X POST 'http://localhost:8080/slim?profile=medium' \
//	  -H "Content-Type: application/json" \
//	  -d '{"users":[{"id":1,"name":"Alice"}]}'
//
// # Performance
//
// SlimJSON is highly optimized for performance:
//   - Process files in 16-47µs
//   - Excellent parallel scalability
//   - Minimal memory allocations
//   - Zero-copy string operations where possible
//
// Benchmark results on real-world data:
//   - GitHub Users API: 24% size reduction
//   - JSON Resume: 98% size reduction
//   - Complex nested objects: 50-70% reduction
//
// # Configuration File Format
//
// The .slimjson configuration file uses a simple INI-style format:
//
//	# Custom profile for API responses
//	[api-response]
//	depth=5
//	list-len=20
//	strip-empty=true
//	decimal-places=2
//	deduplicate=true
//	block=metadata,debug,trace
//
//	# Custom profile for LLM context
//	[llm-context]
//	depth=4
//	list-len=15
//	strip-empty=true
//	string-pooling=true
//	type-inference=true
//	bool-compression=true
//
// The file is searched in:
//  1. Path specified by -c/--config flag (highest priority)
//  2. Current directory (./.slimjson)
//  3. User home directory (~/.slimjson)
//
// # Thread Safety
//
// The Slimmer type is safe for concurrent use. You can create a single
// Slimmer instance and use it from multiple goroutines:
//
//	slimmer := slimjson.New(cfg)
//
//	// Safe to use concurrently
//	go func() {
//	    result1 := slimmer.Slim(data1)
//	    // ...
//	}()
//
//	go func() {
//	    result2 := slimmer.Slim(data2)
//	    // ...
//	}()
//
// # Error Handling
//
// The Slim method does not return errors. Instead, it gracefully handles
// edge cases:
//   - Invalid data types are passed through unchanged
//   - Nil values are handled appropriately
//   - Deep recursion is prevented by MaxDepth
//
// For config file operations, errors are returned:
//
//	profiles, err := slimjson.LoadConfigFile()
//	if err != nil {
//	    // File not found or parse error
//	    log.Printf("Warning: %v", err)
//	    // Use built-in profiles as fallback
//	    profiles = slimjson.GetBuiltinProfiles()
//	}
//
// # Links
//
//   - GitHub: https://github.com/tradik/slimjson
//   - Documentation: https://pkg.go.dev/github.com/tradik/slimjson
//   - Examples: https://github.com/tradik/slimjson/blob/main/LIBRARY_EXAMPLES.md
//   - API Spec: https://github.com/tradik/slimjson/blob/main/api/swagger.yaml
package slimjson
