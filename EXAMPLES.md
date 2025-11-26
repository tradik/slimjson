# SlimJSON Examples

This document provides practical examples of using slimjson for various use cases.

## Table of Contents

- [Basic Usage](#basic-usage)
- [AI/LLM Optimization](#aillm-optimization)
- [API Response Compression](#api-response-compression)
- [Configuration Profiles](#configuration-profiles)
- [Real-World Examples](#real-world-examples)

## Basic Usage

### Strip Empty Values

Remove all null, empty strings, empty arrays, and empty objects:

```bash
echo '{"name":"John","email":"","tags":[],"meta":null}' | slimjson -strip-empty
```

Output:
```json
{"name":"John"}
```

### Limit Array Length

Keep only the first N elements of arrays:

```bash
echo '{"users":[1,2,3,4,5,6,7,8,9,10]}' | slimjson -list-len 3
```

Output:
```json
{"users":[1,2,3]}
```

### Truncate Deep Nesting

Limit object/array nesting depth:

```bash
echo '{"a":{"b":{"c":{"d":"too deep"}}}}' | slimjson -depth 2
```

Output:
```json
{"a":{"b":null}}
```

## AI/LLM Optimization

### Reduce Token Count for GPT

Optimize JSON for sending to language models:

```bash
cat large-data.json | slimjson \
  -depth 4 \
  -list-len 8 \
  -string-len 150 \
  -strip-empty \
  -block "avatar_url,gravatar_id,html_url"
```

**Result**: ~60% size reduction, perfect for token-limited APIs.

### Resume/CV Compression

Compress resume data for AI processing:

```bash
cat resume.json | slimjson \
  -depth 5 \
  -list-len 10 \
  -string-len 200 \
  -strip-empty \
  -block "description,summary" \
  -pretty
```

**Result**: ~33% size reduction while preserving key information.

## API Response Compression

### GitHub API Response

Compress GitHub API responses before processing:

```bash
curl -s https://api.github.com/users/octocat/repos | slimjson \
  -depth 3 \
  -list-len 5 \
  -block "owner,license,permissions" \
  -strip-empty
```

### JSONPlaceholder Users

Simplify user data from APIs:

```bash
curl -s https://jsonplaceholder.typicode.com/users | slimjson \
  -depth 2 \
  -list-len 10 \
  -strip-empty
```

## Configuration Profiles

### Light Compression (24-28% reduction)

Preserve most data structure:

```bash
slimjson -depth 10 -list-len 20 -strip-empty input.json
```

**Use Case**: General cleanup, remove empty values only.

### Medium Compression (28-39% reduction)

Balanced compression:

```bash
slimjson -depth 5 -list-len 10 -string-len 200 -strip-empty input.json
```

**Use Case**: API responses, data pipelines, general purpose.

### Aggressive Compression (88-98% reduction)

Maximum size reduction:

```bash
slimjson \
  -depth 3 \
  -list-len 5 \
  -string-len 100 \
  -strip-empty \
  -block "description,summary,comment,notes" \
  input.json
```

**Use Case**: Extreme token reduction, data sampling, previews.

### AI-Optimized (48-63% reduction)

Optimized for LLM processing:

```bash
slimjson \
  -depth 4 \
  -list-len 8 \
  -string-len 150 \
  -strip-empty \
  -block "avatar_url,gravatar_id,url,html_url" \
  input.json
```

**Use Case**: Sending to GPT, Claude, or other LLMs.

## Real-World Examples

### Example 1: E-commerce Product Catalog

Compress product data for AI recommendations:

```bash
cat products.json | slimjson \
  -depth 4 \
  -list-len 5 \
  -string-len 100 \
  -block "internal_id,warehouse_location,supplier_code" \
  -strip-empty > products-slim.json
```

### Example 2: Log Analysis

Reduce log file size before processing:

```bash
cat application.log.json | slimjson \
  -depth 3 \
  -list-len 10 \
  -string-len 500 \
  -block "stack_trace,debug_info" \
  -strip-empty > logs-slim.json
```

### Example 3: Configuration Files

Clean up configuration with unused fields:

```bash
cat config.json | slimjson \
  -strip-empty \
  -block "deprecated,legacy,old_format" \
  -pretty > config-clean.json
```

### Example 4: Database Export

Compress database exports for backup:

```bash
pg_dump --format=plain mydb | \
  jq -s '.' | \
  slimjson -depth 5 -list-len 100 -strip-empty > backup-slim.json
```

### Example 5: CI/CD Pipeline

Optimize test results for storage:

```bash
cat test-results.json | slimjson \
  -depth 4 \
  -list-len 20 \
  -string-len 200 \
  -block "stdout,stderr,logs" \
  -strip-empty > test-results-slim.json
```

## Library Usage Examples

### Go Library - Basic

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/tradik/slimjson"
)

func main() {
    data := map[string]interface{}{
        "users": []interface{}{
            map[string]interface{}{"id": 1, "name": "Alice"},
            map[string]interface{}{"id": 2, "name": "Bob"},
            map[string]interface{}{"id": 3, "name": "Charlie"},
        },
        "empty": "",
        "null": nil,
    }

    cfg := slimjson.Config{
        MaxListLength: 2,
        StripEmpty:    true,
    }

    slimmer := slimjson.New(cfg)
    result := slimmer.Slim(data)

    output, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(output))
}
```

### Go Library - Advanced

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "github.com/tradik/slimjson"
)

func main() {
    // Read JSON file
    fileData, _ := os.ReadFile("input.json")
    
    var data interface{}
    json.Unmarshal(fileData, &data)

    // Configure slimmer
    cfg := slimjson.Config{
        MaxDepth:        4,
        MaxListLength:   8,
        MaxStringLength: 150,
        StripEmpty:      true,
        BlockList:       []string{"password", "secret", "token"},
    }

    slimmer := slimjson.New(cfg)
    result := slimmer.Slim(data)

    // Write result
    output, _ := json.MarshalIndent(result, "", "  ")
    os.WriteFile("output.json", output, 0644)
    
    fmt.Printf("Original: %d bytes\n", len(fileData))
    fmt.Printf("Compressed: %d bytes\n", len(output))
    fmt.Printf("Reduction: %.1f%%\n", 
        float64(len(fileData)-len(output))/float64(len(fileData))*100)
}
```

## Performance Tips

1. **Use appropriate depth limits**: Deeper structures take longer to process
2. **Enable StripEmpty**: Fastest way to reduce size with minimal overhead
3. **BlockList for known fields**: More efficient than string truncation
4. **Parallel processing**: Use goroutines for multiple files
5. **Benchmark your use case**: Run `go test -bench=.` to measure performance

## Common Patterns

### Pipeline Processing

```bash
# Multi-stage compression
cat data.json | \
  slimjson -strip-empty | \
  slimjson -depth 5 -list-len 10 | \
  slimjson -string-len 100 -pretty
```

### Conditional Compression

```bash
# Compress only if file is large
if [ $(stat -f%z input.json) -gt 100000 ]; then
    slimjson -depth 4 -list-len 8 input.json > output.json
else
    cp input.json output.json
fi
```

### Batch Processing

```bash
# Process all JSON files in directory
for file in *.json; do
    slimjson -depth 5 -list-len 10 -strip-empty "$file" > "slim-$file"
done
```

## Troubleshooting

### Output is too aggressive

Increase limits:
```bash
slimjson -depth 10 -list-len 50 -string-len 500 input.json
```

### Output is too large

Decrease limits or add blocklist:
```bash
slimjson -depth 2 -list-len 3 -string-len 50 -block "field1,field2" input.json
```

### Need to preserve specific fields

Don't use blocklist for important fields, only for removable ones.

### Performance issues

- Use `-strip-empty` only (fastest)
- Reduce depth limit
- Process files in parallel
- Use aggressive compression profile
