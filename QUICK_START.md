# Quick Start Guide

## Testing Your Compression

### 1. Run Compression Tests
```bash
make compression-test
```

This will show you real compression results on 3 test files with 4 different profiles.

### 2. Run Performance Benchmarks
```bash
make bench
```

See how fast slimjson processes different file sizes.

### 3. Run Unit Tests
```bash
make test
```

Verify all functionality works correctly.

## Test Files Included

1. **resume.json** (28.2 KB) - Complex nested resume
2. **schema-resume.json** (24.8 KB) - Schema.org resume format
3. **users.json** (5.5 KB) - Simple user list

## Compression Profiles

### Light (24-28% reduction)
```bash
slimjson -depth 10 -list-len 20 -strip-empty input.json
```

### Medium (28-39% reduction)
```bash
slimjson -depth 5 -list-len 10 -string-len 200 -strip-empty input.json
```

### Aggressive (88-98% reduction)
```bash
slimjson -depth 3 -list-len 5 -string-len 100 -strip-empty -block "description,summary" input.json
```

### AI-Optimized (48-63% reduction)
```bash
slimjson -depth 4 -list-len 8 -string-len 150 -strip-empty -block "avatar_url,url" input.json
```

## Real Results

| File | Profile | Original | Compressed | Reduction |
|------|---------|----------|------------|-----------|
| resume.json | Medium | 28.2 KB | 18.8 KB | 33.5% |
| schema-resume.json | AI-Optimized | 24.8 KB | 9.2 KB | 62.9% |
| users.json | Aggressive | 5.5 KB | 691 B | 87.8% |

## Performance

- Small files (5KB): ~16µs
- Medium files (25KB): ~39µs
- Large files (28KB): ~47µs

## More Information

- **Examples**: See [EXAMPLES.md](EXAMPLES.md) for detailed usage
- **Testing**: See [testing/README.md](testing/README.md) for test documentation
- **Changes**: See [CHANGELOG.md](CHANGELOG.md) for version history
- **Summary**: See [TESTING_SUMMARY.md](TESTING_SUMMARY.md) for complete overview
