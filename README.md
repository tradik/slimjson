# slimjson

`slimjson` is a Go tool and library designed to "slim down" JSON data before sending it to AI models or other bandwidth-constrained systems. It removes unnecessary fields, truncates deep nesting, shortens lists and strings, and strips empty values to create an optimized "diet" version of your JSON.

Perfect for reducing token usage when sending large JSON payloads to LLMs and AI APIs.

## Features

- **Prune Unnecessary Fields**: Remove specific fields by name (blocklist).
- **Truncate Deep Nesting**: Automatically cut off objects/arrays deeper than a specified limit.
- **Shorten Lists**: Limit the number of elements in arrays.
- **Truncate Strings**: Limit string length (UTF-8 aware, counts runes not bytes).
- **Strip Empty Values**: Remove `null`, empty strings, empty arrays, and empty objects.

## Installation

### Pre-built Binaries

You can download the latest pre-built binaries for Linux, macOS, and FreeBSD from the [Releases](https://github.com/tradik/slimjson/releases) page (once available).

### Installation Guide

#### Linux (amd64/arm64)

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

#### macOS (Intel/Apple Silicon)

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

#### FreeBSD

1. Download the binary:
   ```bash
   fetch -o slimjson https://github.com/tradik/slimjson/releases/latest/download/slimjson-freebsd-amd64
   ```
2. Make it executable:
   ```bash
   chmod +x slimjson
   ```
3. Move to path:
   ```bash
   sudo mv slimjson /usr/local/bin/
   ```

### Build from Source

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

## Usage

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

## Development

### Requirements

- Go 1.25+

### Testing

```bash
go test ./...
```

### Linting

```bash
golangci-lint run
```
