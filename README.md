# slimjson 🎯

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](LICENSE)
[![CI/CD](https://github.com/tradik/slimjson/workflows/CI/CD/badge.svg)](https://github.com/tradik/slimjson/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/tradik/slimjson)](https://goreportcard.com/report/github.com/tradik/slimjson)
[![GitHub Pages](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://tradik.github.io/slimjson/)

`slimjson` is a Go tool and library designed to "slim down" JSON data before sending it to AI models or other bandwidth-constrained systems. It removes unnecessary fields, truncates deep nesting, shortens lists and strings, and strips empty values to create an optimized "diet" version of your JSON.

✨ Perfect for reducing token usage when sending large JSON payloads to LLMs and AI APIs.

## Features ⚡

- 🗑️ **Prune Unnecessary Fields**: Remove specific fields by name (blocklist).
- 📏 **Truncate Deep Nesting**: Automatically cut off objects/arrays deeper than a specified limit.
- ✂️ **Shorten Lists**: Limit the number of elements in arrays.
- 📝 **Truncate Strings**: Limit string length (UTF-8 aware, counts runes not bytes).
- 🧹 **Strip Empty Values**: Remove `null`, empty strings, empty arrays, and empty objects.
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

### CLI

The `slimjson` CLI reads JSON from stdin or a file and outputs the slimmed JSON to stdout.

```bash
# Read from file
slimjson -depth 3 -list-len 5 input.json > output.json

# Read from stdin
cat input.json | slimjson -strip-empty=true -block "password,secret"
```

**Flags:**

- `-depth int`: Maximum nesting depth (default 5). 0 for unlimited.
- `-list-len int`: Maximum list length (default 10). 0 for unlimited.
- `-string-len int`: Maximum string length in characters/runes (default 0 = unlimited). UTF-8 aware.
- `-strip-empty`: Remove nulls, empty strings, empty arrays/objects (default true).
- `-block string`: Comma-separated list of field names to remove.
- `-pretty`: Pretty print output.

📚 **See [EXAMPLES.md](EXAMPLES.md) for detailed usage examples and common patterns.**

### Library

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

	cfg := slimjson.Config{
		MaxDepth:      2,
		MaxListLength: 3,
		StripEmpty:    true,
	}

	slimmer := slimjson.New(cfg)
	result := slimmer.Slim(data)

	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
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
| resume.json | 28.2 KB | Medium | 18.8 KB | 9.5 KB | 33.5% | 7230 | 4805 | 33.5% |
| resume.json | 28.2 KB | Aggressive | 530 B | 27.7 KB | 98.2% | 7230 | 133 | 98.2% |
| resume.json | 28.2 KB | AI-Optimized | 11.2 KB | 17.1 KB | 60.5% | 7230 | 2859 | 60.5% |
| schema-resume.json | 24.8 KB | Light | 17.9 KB | 7.0 KB | 28.0% | 6359 | 4579 | 28.0% |
| schema-resume.json | 24.8 KB | Medium | 15.2 KB | 9.7 KB | 38.9% | 6359 | 3887 | 38.9% |
| schema-resume.json | 24.8 KB | Aggressive | 530 B | 24.3 KB | 97.9% | 6359 | 133 | 97.9% |
| schema-resume.json | 24.8 KB | AI-Optimized | 9.2 KB | 15.6 KB | 62.9% | 6359 | 2358 | 62.9% |
| users.json | 5.5 KB | Light | 4.0 KB | 1.5 KB | 27.5% | 1412 | 1024 | 27.5% |
| users.json | 5.5 KB | Medium | 4.0 KB | 1.5 KB | 27.5% | 1412 | 1024 | 27.5% |
| users.json | 5.5 KB | Aggressive | 691 B | 4.8 KB | 87.8% | 1412 | 173 | 87.7% |
| users.json | 5.5 KB | AI-Optimized | 2.9 KB | 2.6 KB | 47.9% | 1412 | 736 | 47.9% |

**Token Estimation**: Tokens are estimated using ~4 characters per token, approximating GPT-style tokenization for JSON/English text.

**Statistical Analysis**: Each test is run 10 times. Processing times show mean ± standard deviation (n=10) for statistical reliability. See [testing/METHODOLOGY.md](testing/METHODOLOGY.md) for detailed methodology.

### Configuration Profiles

- **Light**: `MaxDepth: 10, MaxListLength: 20, StripEmpty: true` - Preserves most data structure
- **Medium**: `MaxDepth: 5, MaxListLength: 10, MaxStringLength: 200, StripEmpty: true` - Balanced compression
- **Aggressive**: `MaxDepth: 3, MaxListLength: 5, MaxStringLength: 100, StripEmpty: true` - Maximum size reduction
- **AI-Optimized**: `MaxDepth: 4, MaxListLength: 8, MaxStringLength: 150, StripEmpty: true` - Optimized for LLM token reduction

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
