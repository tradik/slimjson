---
layout: default
title: slimjson token optimizer
permalink: /
---

# slimjson 🎯

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](LICENSE)
[![CI/CD](https://github.com/tradik/slimjson/workflows/CI/CD/badge.svg)](https://github.com/tradik/slimjson/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/tradik/slimjson)](https://goreportcard.com/report/github.com/tradik/slimjson)
[![GitHub Pages](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://tradik.github.io/slimjson/)

`slimjson` is a Go tool and library designed to "slim down" JSON data before sending it to AI models or other bandwidth-constrained systems. It removes unnecessary fields, truncates deep nesting, shortens lists and strings, and strips empty values to create an optimized "diet" version of your JSON.

✨ Perfect for reducing token usage when sending large JSON payloads to LLMs and AI APIs.

## Quick Links 📚

- 📖 [CLI Examples](EXAMPLES.md) - Command-line usage examples
- 💻 [Library Guide](LIBRARY_EXAMPLES.md) - Complete guide for developers using SlimJSON as a Go library
- 🌐 [HTTP API Documentation](api/README.md) - REST API reference and integration examples
- 📋 [OpenAPI Specification](api/swagger.yaml) - Swagger/OpenAPI 3.0 spec for API
- ⚙️ [Configuration File](.slimjson.example) - Custom profiles with `.slimjson` file
- 📦 [Go Package Documentation](https://pkg.go.dev/github.com/tradik/slimjson) - Full API reference
- 🧪 [Testing Methodology](testing/METHODOLOGY.md) - Compression testing details

## Features ⚡

- 🗑️ **Prune Unnecessary Fields**: Remove specific fields by name (blocklist).
- 📏 **Truncate Deep Nesting**: Automatically cut off objects/arrays deeper than a specified limit.
- ✂️ **Shorten Lists**: Limit the number of elements in arrays.
- 📝 **Truncate Strings**: Limit string length (UTF-8 aware, counts runes not bytes).
- 🧹 **Strip Empty Values**: Remove `null`, empty strings, empty arrays, and empty objects.
- ⚙️ **Custom Profiles**: Define reusable compression profiles in `.slimjson` config file.
- 🌐 **HTTP Daemon Mode**: Run as a REST API service for JSON compression.
- 🔧 **Go Library**: Use as a library in your Go applications with full programmatic control.
- 🎯 **Config Priority**: `-c` flag for custom config file with highest priority.
- 📋 **Auto Help**: Shows comprehensive usage when run without arguments.
- ⚡ **High Performance**: Process files in 16-47µs with excellent parallel scalability.
- 📊 **Proven Results**: 24-98% size reduction on real-world JSON files.
- 🧪 **Comprehensive Testing**: Full test suite with benchmarks and compression metrics.

## Installation 📦

### Pre-built Binaries

You can download the latest pre-built binaries for Linux, macOS, and FreeBSD from the [Releases](https://github.com/tradik/slimjson/releases) page (once available).

### Installation Guide

#### Linux (amd64/arm64) 🐧

1. Download the binary:
   ```bash
   # For amd64
   wget https://github.com/tradik/slimjson/releases/latest/download/slimjson-linux-amd64 -O slimjson
   
   # For arm64
   wget https://github.com/tradik/slimjson/releases/latest/download/slimjson-linux-arm64 -O slimjson
   ```
2. Make it executable:
   ```bash
   chmod +x slimjson
   ```
3. Move to path:
   ```bash
   sudo mv slimjson /usr/local/bin/
   ```

#### macOS (Intel/Apple Silicon) 🍎

1. Download the binary:
   ```bash
   # For Apple Silicon (M1/M2/etc)
   curl -L https://github.com/tradik/slimjson/releases/latest/download/slimjson-darwin-arm64 -o slimjson
   
   # For Intel
   curl -L https://github.com/tradik/slimjson/releases/latest/download/slimjson-darwin-amd64 -o slimjson
   ```
2. Make it executable:
   ```bash
   chmod +x slimjson
   ```
3. Move to path:
   ```bash
   sudo mv slimjson /usr/local/bin/
   ```
   *Note: You might need to allow the application in System Settings > Privacy & Security if macOS blocks it.*

#### FreeBSD 👹

1. Download the binary:
   ```bash
   # For amd64
   fetch -o slimjson https://github.com/tradik/slimjson/releases/latest/download/slimjson-freebsd-amd64
   
   # For arm64
   fetch -o slimjson https://github.com/tradik/slimjson/releases/latest/download/slimjson-freebsd-arm64
   ```
2. Make it executable:
   ```bash
   chmod +x slimjson
   ```
3. Move to path:
   ```bash
   sudo mv slimjson /usr/local/bin/
   ```

### Build from Source 🔨

If you have Go 1.25+ installed:

```bash
go install github.com/tradik/slimjson/cmd/slimjson@latest
```

Or clone and build using Makefile:

```bash
git clone https://github.com/tradik/slimjson.git
cd slimjson
make build
sudo mv bin/slimjson /usr/local/bin/
```

## Usage 🚀

### Configuration File

SlimJSON supports a `.slimjson` configuration file for defining custom profiles. The file is searched in:
1. Current directory (`./.slimjson`)
2. User home directory (`~/.slimjson`)

**Format:**
```ini
# Comments start with # or //
[profile-name]
parameter=value

[another-profile]
parameter=value
```

**Example `.slimjson`:**
```ini
# Custom profile for API responses
[api-response]
depth=5
list-len=20
strip-empty=true
decimal-places=2
deduplicate=true
block=metadata,debug,trace

# Custom profile for LLM context
[llm-context]
depth=4
list-len=15
strip-empty=true
string-pooling=true
type-inference=true
bool-compression=true
block=avatar_url,url,html_url

# Maximum compression
[maximum]
depth=3
list-len=5
strip-empty=true
decimal-places=2
deduplicate=true
sample-strategy=first_last
sample-size=10
null-compression=true
type-inference=true
bool-compression=true
timestamp-compression=true
string-pooling=true
number-delta=true
enum-detection=true
```

**Using custom profiles:**
```bash
slimjson -profile api-response data.json
slimjson -profile llm-context data.json
```

**Note:** Custom profiles take precedence over built-in profiles. If a parameter is not specified, it defaults to the zero value (disabled).

See [.slimjson.example](.slimjson.example) for a complete configuration file with all available parameters.

### CLI

The `slimjson` CLI reads JSON from stdin or a file and outputs the slimmed JSON to stdout. It can also run as an HTTP daemon for processing JSON via REST API.

#### Quick Start

```bash
# Show help
slimjson

# Process file
slimjson data.json

# Process stdin
cat data.json | slimjson -profile medium

# Run as daemon
slimjson -d -port 8080
```

#### Using Predefined Profiles (Recommended)

```bash
# Light compression - preserve most data
slimjson -profile light input.json > output.json

# Medium compression - balanced reduction
slimjson -profile medium input.json > output.json

# Aggressive compression - maximum reduction
slimjson -profile aggressive input.json > output.json

# AI-Optimized - optimized for LLM token reduction
slimjson -profile ai-optimized input.json > output.json
```

#### Custom Parameters

```bash
# Read from file
slimjson -depth 3 -list-len 5 input.json > output.json

# Read from stdin
cat input.json | slimjson -strip-empty=true -block "password,secret"

# Round numbers and remove duplicates
slimjson -decimal-places 2 -deduplicate data.json

# Sample large arrays - keep first 5 and last 5
slimjson -sample-strategy first_last -sample-size 10 data.json

# Representative sampling - evenly distributed
slimjson -sample-strategy representative -sample-size 20 data.json

# Combine with profile
slimjson -profile medium -decimal-places 2 -deduplicate data.json

# Advanced compression - all optimizations
slimjson -string-pooling -enum-detection -timestamp-compression data.json

# Maximum compression (use all features)
slimjson -profile ai-optimized \
  -decimal-places 2 \
  -deduplicate \
  -sample-strategy representative \
  -sample-size 50 \
  -null-compression \
  -type-inference \
  -bool-compression \
  -timestamp-compression \
  -string-pooling \
  -number-delta \
  -enum-detection \
  data.json
```

**Flags:**

**Basic Options:**
- `-profile string`: Use predefined profile: `light`, `medium`, `aggressive`, `ai-optimized`
- `-depth int`: Maximum nesting depth (default: 5, 0 = unlimited)
- `-list-len int`: Maximum list length (default: 10, 0 = unlimited)
- `-string-len int`: Maximum string length in characters/runes (default: 0 = unlimited)
- `-strip-empty`: Remove nulls, empty strings, empty arrays/objects (default: true)
- `-block string`: Comma-separated list of field names to remove
- `-pretty`: Pretty print output

**Optimization Options:**
- `-decimal-places int`: Round floats to N decimal places (default: -1 = no rounding)
- `-deduplicate`: Remove duplicate values from arrays (default: false)
- `-sample-strategy string`: Array sampling: `none`, `first_last`, `random`, `representative` (default: `none`)
- `-sample-size int`: Number of items when sampling (default: 0 = use list-len)

**Advanced Compression:**
- `-null-compression`: Track removed null fields in _nulls array (default: false)
- `-type-inference`: Convert uniform arrays to schema+data format (default: false)
- `-bool-compression`: Convert booleans to bit flags (default: false)
- `-timestamp-compression`: Convert ISO timestamps to unix timestamps (default: false)
- `-string-pooling`: Deduplicate repeated strings using string pool (default: false)
- `-string-pool-min int`: Minimum occurrences for string pooling (default: 2)
- `-number-delta`: Use delta encoding for sequential numbers (default: false)
- `-number-delta-threshold int`: Minimum array size for delta encoding (default: 5)
- `-enum-detection`: Convert repeated categorical values to enums (default: false)
- `-enum-max-values int`: Maximum unique values to consider as enum (default: 10)

**Profile Details:**

| Profile | MaxDepth | MaxListLength | BlockList | Strategy |
|---------|----------|---------------|-----------|----------|
| **Light** | 10 | 20 | none | Preserve data, limit depth/lists |
| **Medium** | 5 | 10 | none | Balanced reduction |
| **Aggressive** | 3 | 5 | description, summary, comment, notes, bio, readme | Remove verbose text fields |
| **AI-Optimized** | 4 | 8 | avatar_url, gravatar_id, url, html_url, *_url | Remove URLs and metadata |

**Default Values Summary:**

| Parameter | Default | Description |
|-----------|---------|-------------|
| `-depth` | 5 | Maximum nesting depth |
| `-list-len` | 10 | Maximum array length |
| `-string-len` | 0 | No string truncation (unlimited) |
| `-strip-empty` | true | Remove empty values |
| `-decimal-places` | -1 | No rounding |
| `-sample-strategy` | none | No sampling |
| `-sample-size` | 0 | Use list-len value |
| `-string-pool-min` | 2 | Min occurrences for pooling |
| `-number-delta-threshold` | 5 | Min array size for delta |
| `-enum-max-values` | 10 | Max unique values for enum |

**Note**: Profiles do NOT truncate strings to preserve data integrity. Use `-string-len` manually if needed, but be aware this may lose information.

#### Daemon Mode (HTTP Server)

Run SlimJSON as an HTTP daemon to process JSON via REST API:

```bash
# Start daemon on default port 8080
slimjson -d

# Start on custom port
slimjson -d -port 3000

# Use custom config file
slimjson -d -c /path/to/.slimjson
```

**API Endpoints:**

```bash
# Health check
curl http://localhost:8080/health
# Response: {"status":"ok","version":"1.0"}

# List available profiles
curl http://localhost:8080/profiles
# Response: {"builtin":["light","medium","aggressive","ai-optimized"],"custom":["my-profile"]}

# Compress JSON with default settings
curl -X POST http://localhost:8080/slim \
  -H "Content-Type: application/json" \
  -d '{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]}'

# Compress JSON with specific profile
curl -X POST 'http://localhost:8080/slim?profile=medium' \
  -H "Content-Type: application/json" \
  -d '{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]}'

# Use custom profile from config file
curl -X POST 'http://localhost:8080/slim?profile=my-custom-profile' \
  -H "Content-Type: application/json" \
  -d @data.json
```

**Daemon Features:**
- ✅ RESTful API for JSON compression
- ✅ Support for all built-in and custom profiles
- ✅ Health check endpoint for monitoring
- ✅ Profile discovery endpoint
- ✅ Automatic config file loading
- ✅ Production-ready HTTP server

**Use Cases:**
- Microservice for JSON optimization
- API gateway integration
- CI/CD pipeline processing
- Real-time data compression service

#### Custom Config File Priority

```bash
# Priority 1: Specified config file (highest priority)
slimjson -c /path/to/custom.slimjson -profile my-profile data.json

# Priority 2: .slimjson in current directory
slimjson -profile my-profile data.json

# Priority 3: .slimjson in home directory
slimjson -profile my-profile data.json

# Priority 4: Built-in profiles
slimjson -profile medium data.json
```

📚 **See [EXAMPLES.md](EXAMPLES.md) for CLI examples and [LIBRARY_EXAMPLES.md](LIBRARY_EXAMPLES.md) for complete library usage guide.**

### Library (Go Package)

SlimJSON can be used as a Go library in your applications.

#### Installation

```bash
go get github.com/tradik/slimjson
```

#### Basic Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/tradik/slimjson"
)

func main() {
	data := map[string]interface{}{
		"name": "Test",
		"details": map[string]interface{}{
			"deep": map[string]interface{}{
				"too_deep": "value",
			},
		},
		"list": []interface{}{1, 2, 3, 4, 5},
		"empty": "",
	}

	// Create configuration
	cfg := slimjson.Config{
		MaxDepth:      2,
		MaxListLength: 3,
		StripEmpty:    true,
	}

	// Create slimmer and process data
	slimmer := slimjson.New(cfg)
	result := slimmer.Slim(data)

	// Output result
	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
}
```

#### Using Built-in Profiles

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/tradik/slimjson"
)

func main() {
	// Get built-in profile
	profiles := slimjson.GetBuiltinProfiles()
	cfg := profiles["medium"]
	
	// Or use a specific profile directly
	cfg := slimjson.Config{
		MaxDepth:      5,
		MaxListLength: 10,
		StripEmpty:    true,
	}

	slimmer := slimjson.New(cfg)
	result := slimmer.Slim(yourData)
	
	// ... use result
}
```

#### Advanced Configuration

```go
package main

import (
	"github.com/tradik/slimjson"
)

func main() {
	// Full configuration with all options
	cfg := slimjson.Config{
		// Basic options
		MaxDepth:        5,
		MaxListLength:   10,
		MaxStringLength: 0,  // 0 = unlimited
		StripEmpty:      true,
		BlockList:       []string{"password", "secret", "token"},
		
		// Optimization options
		DecimalPlaces:     2,
		DeduplicateArrays: true,
		SampleStrategy:    "first_last",
		SampleSize:        20,
		
		// Advanced compression
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
	}

	slimmer := slimjson.New(cfg)
	result := slimmer.Slim(data)
	
	// Result may contain metadata fields:
	// - _strings: String pool (if StringPooling enabled)
	// - _enums: Enum mappings (if EnumDetection enabled)
	// - _nulls: Tracked null fields (if NullCompression enabled)
}
```

#### Loading Custom Profiles from File

```go
package main

import (
	"fmt"
	"github.com/tradik/slimjson"
)

func main() {
	// Load profiles from .slimjson file
	customProfiles, err := slimjson.LoadConfigFile()
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
		customProfiles = make(map[string]slimjson.Config)
	}
	
	// Use custom profile
	if cfg, ok := customProfiles["my-custom-profile"]; ok {
		slimmer := slimjson.New(cfg)
		result := slimmer.Slim(data)
		// ... use result
	}
	
	// Or combine built-in and custom profiles
	allProfiles := slimjson.GetBuiltinProfiles()
	for name, cfg := range customProfiles {
		allProfiles[name] = cfg
	}
	
	// Use any profile
	cfg := allProfiles["api-response"]
	slimmer := slimjson.New(cfg)
	result := slimmer.Slim(data)
}
```

#### Parsing Config File Manually

```go
package main

import (
	"github.com/tradik/slimjson"
)

func main() {
	// Parse specific config file
	profiles, err := slimjson.ParseConfigFile("/path/to/.slimjson")
	if err != nil {
		panic(err)
	}
	
	cfg := profiles["my-profile"]
	slimmer := slimjson.New(cfg)
	result := slimmer.Slim(data)
}
```

#### Config Structure Reference

```go
type Config struct {
	// Basic options
	MaxDepth        int      // Maximum nesting depth (0 = unlimited)
	MaxListLength   int      // Maximum array length (0 = unlimited)
	MaxStringLength int      // Maximum string length (0 = unlimited)
	StripEmpty      bool     // Remove nulls, empty strings, empty arrays/objects
	BlockList       []string // List of field names to remove (case-insensitive)
	
	// Optimization options
	DecimalPlaces     int    // Round floats to N decimal places (-1 = no rounding)
	DeduplicateArrays bool   // Remove duplicate values from arrays
	SampleStrategy    string // Array sampling: "none", "first_last", "random", "representative"
	SampleSize        int    // Number of items when sampling (0 = use MaxListLength)
	
	// Advanced compression
	NullCompression          bool   // Track removed null fields in _nulls array
	TypeInference            bool   // Convert uniform arrays to schema+data format
	BoolCompression          bool   // Convert booleans to bit flags
	TimestampCompression     bool   // Convert ISO timestamps to unix timestamps
	StringPooling            bool   // Deduplicate repeated strings using string pool
	StringPoolMinOccurrences int    // Minimum occurrences for string pooling (default: 2)
	NumberDeltaEncoding      bool   // Use delta encoding for sequential numbers
	NumberDeltaThreshold     int    // Minimum array size for delta encoding (default: 5)
	EnumDetection            bool   // Convert repeated categorical values to enums
	EnumMaxValues            int    // Maximum unique values to consider as enum (default: 10)
}
```

#### Example: API Response Compression

```go
package main

import (
	"encoding/json"
	"net/http"
	"github.com/tradik/slimjson"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// Get data from database
	data := fetchDataFromDB()
	
	// Configure compression for API response
	cfg := slimjson.Config{
		MaxDepth:          5,
		MaxListLength:     20,
		StripEmpty:        true,
		DecimalPlaces:     2,
		DeduplicateArrays: true,
		BlockList:         []string{"internal_id", "metadata", "debug"},
	}
	
	// Compress data
	slimmer := slimjson.New(cfg)
	compressed := slimmer.Slim(data)
	
	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(compressed)
}
```

#### Example: LLM Context Optimization

```go
package main

import (
	"encoding/json"
	"github.com/tradik/slimjson"
)

func prepareLLMContext(data interface{}) ([]byte, error) {
	// Maximum compression for LLM context
	cfg := slimjson.Config{
		MaxDepth:             4,
		MaxListLength:        15,
		StripEmpty:           true,
		DecimalPlaces:        2,
		DeduplicateArrays:    true,
		StringPooling:        true,
		TypeInference:        true,
		BoolCompression:      true,
		TimestampCompression: true,
		BlockList:            []string{"avatar_url", "url", "html_url"},
	}
	
	slimmer := slimjson.New(cfg)
	compressed := slimmer.Slim(data)
	
	// Convert to JSON for LLM
	return json.Marshal(compressed)
}
```

#### Example: Processing Large Datasets

```go
package main

import (
	"encoding/json"
	"github.com/tradik/slimjson"
)

func processLargeDataset(records []map[string]interface{}) []byte {
	// Wrap in container
	data := map[string]interface{}{
		"records": records,
	}
	
	// Aggressive compression with sampling
	cfg := slimjson.Config{
		MaxDepth:          3,
		MaxListLength:     100,
		StripEmpty:        true,
		DecimalPlaces:     2,
		SampleStrategy:    "representative",
		SampleSize:        50,
		TypeInference:     true,
		NumberDeltaEncoding: true,
	}
	
	slimmer := slimjson.New(cfg)
	compressed := slimmer.Slim(data)
	
	result, _ := json.Marshal(compressed)
	return result
}
```

### Docker / Podman 🐳

Run `slimjson` as a containerized service using Docker or Podman.

#### Pull from GitHub Container Registry

Pre-built multi-arch images are available:

```bash
# Pull latest version
docker pull ghcr.io/tradik/slimjson:latest

# Pull specific version
docker pull ghcr.io/tradik/slimjson:v0.1.X

# Using Podman
podman pull ghcr.io/tradik/slimjson:latest
```

#### Build the image locally

```bash
# Using Docker
make docker-build
# or
docker build -t slimjson:latest .

# Using Podman
make podman-build
# or
podman build -t slimjson:latest .
```

#### Run the container

```bash
# Using Docker (from ghcr.io)
cat input.json | docker run -i --rm ghcr.io/tradik/slimjson:latest -depth 3 -list-len 5

# Using local image
cat input.json | docker run -i --rm slimjson:latest -depth 3 -list-len 5

# Using Podman
cat input.json | podman run -i --rm ghcr.io/tradik/slimjson:latest -depth 3 -list-len 5

# Using docker-compose
docker-compose up
```

#### Example with file mounting

```bash
# Docker
docker run -i --rm -v $(pwd)/data:/data ghcr.io/tradik/slimjson:latest -depth 5 < /data/input.json > /data/output.json

# Podman
podman run -i --rm -v $(pwd)/data:/data:z ghcr.io/tradik/slimjson:latest -depth 5 < /data/input.json > /data/output.json
```

## Compression Results 📊

Real-world compression tests on various JSON files:

| File | Original Size | Config | Compressed Size | Reduction | Reduction % | Original Tokens | Compressed Tokens | Token Reduction % |
|------|---------------|--------|-----------------|-----------|-------------|-----------------|-------------------|-------------------|
| resume.json | 28.2 KB | Light | 21.4 KB | 6.8 KB | 24.2% | 7230 | 5478 | 24.2% |
| resume.json | 28.2 KB | Medium | 18.9 KB | 9.3 KB | 33.1% | 7230 | 4841 | 33.0% |
| resume.json | 28.2 KB | Aggressive | 530 B | 27.7 KB | 98.2% | 7230 | 133 | 98.2% |
| resume.json | 28.2 KB | AI-Optimized | 11.4 KB | 16.8 KB | 59.5% | 7230 | 2928 | 59.5% |
| schema-resume.json | 24.8 KB | Light | 17.9 KB | 7.0 KB | 28.0% | 6359 | 4579 | 28.0% |
| schema-resume.json | 24.8 KB | Medium | 15.3 KB | 9.5 KB | 38.3% | 6359 | 3922 | 38.3% |
| schema-resume.json | 24.8 KB | Aggressive | 530 B | 24.3 KB | 97.9% | 6359 | 133 | 97.9% |
| schema-resume.json | 24.8 KB | AI-Optimized | 9.5 KB | 15.4 KB | 61.8% | 6359 | 2427 | 61.8% |
| users.json | 5.5 KB | Light | 4.0 KB | 1.5 KB | 27.5% | 1412 | 1024 | 27.5% |
| users.json | 5.5 KB | Medium | 4.0 KB | 1.5 KB | 27.5% | 1412 | 1024 | 27.5% |
| users.json | 5.5 KB | Aggressive | 691 B | 4.8 KB | 87.8% | 1412 | 173 | 87.7% |
| users.json | 5.5 KB | AI-Optimized | 2.9 KB | 2.6 KB | 47.9% | 1412 | 736 | 47.9% |

**Token Estimation**: Tokens are estimated using ~4 characters per token, approximating GPT-style tokenization for JSON/English text.

**Statistical Analysis**: Each test is run 10 times. Processing times show mean ± standard deviation (n=10) for statistical reliability. See [testing/METHODOLOGY.md](testing/METHODOLOGY.md) for detailed methodology.

### Configuration Profiles

- **Light**: `MaxDepth: 10, MaxListLength: 20, StripEmpty: true` - Preserves most data structure
- **Medium**: `MaxDepth: 5, MaxListLength: 10, StripEmpty: true` - Balanced compression
- **Aggressive**: `MaxDepth: 3, MaxListLength: 5, StripEmpty: true, BlockList: [description, summary, comment, notes, bio, readme]` - Removes verbose text fields
- **AI-Optimized**: `MaxDepth: 4, MaxListLength: 8, StripEmpty: true, BlockList: [*_url fields]` - Removes URLs and metadata for LLM optimization

### Performance Benchmarks

Benchmarks run on Apple M2 (arm64):

```
BenchmarkSlim_Small-8              196236    16327 ns/op    19384 B/op    442 allocs/op
BenchmarkSlim_Medium-8              90555    38951 ns/op    42472 B/op    955 allocs/op
BenchmarkSlim_Large-8               77553    46666 ns/op    50032 B/op   1143 allocs/op
BenchmarkSlim_Aggressive-8         181198    20002 ns/op    12720 B/op    598 allocs/op
BenchmarkSlim_Parallel-8           177736    20209 ns/op    42473 B/op    955 allocs/op
```

**Key Performance Metrics:**
- **Small files (5KB)**: ~16µs per operation
- **Medium files (25KB)**: ~39µs per operation
- **Large files (28KB)**: ~47µs per operation
- **Aggressive compression**: ~20µs per operation (faster due to early truncation)
- **Parallel processing**: Excellent scalability with minimal overhead

Run your own compression tests:
```bash
cd testing
./run_tests.sh
```

Run performance benchmarks:
```bash
go test -bench=. -benchmem
```

## Development 🛠️

### Requirements

- Go 1.25+

### Testing

```bash
go test ./...
```

### Benchmarking

```bash
go test -bench=. -benchmem -benchtime=3s
```

### Compression Testing

```bash
cd testing
go run compression_benchmark.go
```

### Linting

```bash
golangci-lint run
```

## Documentation 📚

- **[QUICK_START.md](QUICK_START.md)** - Get started with testing in 5 minutes
- **[EXAMPLES.md](EXAMPLES.md)** - Comprehensive usage examples and patterns
- **[CHANGELOG.md](CHANGELOG.md)** - Version history and changes
- **[testing/README.md](testing/README.md)** - Testing suite documentation
- **[TESTING_SUMMARY.md](TESTING_SUMMARY.md)** - Complete testing overview

## Contributing 🤝

Contributions are welcome! Please feel free to submit a Pull Request.

## License 📄

BSD-3-Clause License - see [LICENSE](LICENSE) file for details.
