# SlimJSON Testing Suite

This directory contains comprehensive compression tests and benchmarks for slimjson.

## Structure

```
testing/
├── fixtures/              # Test JSON files
│   ├── resume.json       # 28KB - Complex resume with nested data
│   ├── schema-resume.json # 25KB - Schema.org resume format
│   └── users.json        # 5KB - Simple user list
├── compression_benchmark.go # Compression testing tool
├── run_tests.sh          # Test runner script
└── README.md            # This file
```

## Running Tests

### Quick Start

```bash
./run_tests.sh
```

### Manual Execution

```bash
go run compression_benchmark.go
```

## Test Files

### resume.json (28.2 KB)
Complex resume with deep nesting, arrays, and various data types. Tests real-world JSON structure handling.

### schema-resume.json (24.8 KB)
Schema.org compatible resume format with semantic markup. Tests structured data compression.

### users.json (5.5 KB)
Simple user list from JSONPlaceholder API. Tests basic array and object compression.

## Compression Profiles

The test suite evaluates four compression profiles:

### Light Compression
- **MaxDepth**: 10
- **MaxListLength**: 20
- **StripEmpty**: true
- **Use Case**: Preserve most data while removing empty values
- **Typical Reduction**: 24-28%

### Medium Compression
- **MaxDepth**: 5
- **MaxListLength**: 10
- **MaxStringLength**: 200
- **StripEmpty**: true
- **Use Case**: Balanced compression for general use
- **Typical Reduction**: 28-39%

### Aggressive Compression
- **MaxDepth**: 3
- **MaxListLength**: 5
- **MaxStringLength**: 100
- **StripEmpty**: true
- **BlockList**: description, summary, comment, notes
- **Use Case**: Maximum size reduction
- **Typical Reduction**: 88-98%

### AI-Optimized Compression
- **MaxDepth**: 4
- **MaxListLength**: 8
- **MaxStringLength**: 150
- **StripEmpty**: true
- **BlockList**: gravatar_id, avatar_url, url, html_url
- **Use Case**: Optimized for LLM token reduction
- **Typical Reduction**: 48-63%

## Output Format

The test tool provides:

1. **Detailed Results**: Per-file, per-config compression metrics
2. **Markdown Table**: Ready-to-use table for documentation
3. **Processing Time**: Performance metrics for each operation

## Adding New Test Files

1. Place JSON file in `fixtures/` directory
2. Run `./run_tests.sh`
3. Results will automatically include the new file

## Performance Considerations

- Small files (< 10KB): ~16µs processing time
- Medium files (10-30KB): ~39µs processing time
- Large files (> 30KB): ~47µs processing time
- Aggressive compression is faster due to early truncation

## Integration with CI/CD

Add to your CI pipeline:

```yaml
- name: Run compression tests
  run: |
    cd testing
    ./run_tests.sh
```

## Benchmarking

For detailed performance benchmarks, run from the project root:

```bash
go test -bench=. -benchmem -benchtime=3s
```

See `slimjson_bench_test.go` for benchmark implementations.
