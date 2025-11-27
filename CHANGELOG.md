# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Daemon Mode (HTTP Server)**: Run SlimJSON as an HTTP service
  - `-d` or `-daemon` flag to run as HTTP server
  - `-port` flag to specify port (default: 8080)
  - RESTful API with 3 endpoints:
    - `POST /slim?profile=<name>` - Compress JSON
    - `GET /health` - Health check
    - `GET /profiles` - List available profiles
  - Supports all built-in and custom profiles
  - Production-ready HTTP server
- **Custom Config File Priority**: `-c` or `-config` flag
  - Specify custom config file path
  - Takes priority over `.slimjson` search
  - Enables multiple config files for different environments
- **Default Help Message**: Shows usage when run without arguments
  - Comprehensive help with all options
  - Examples for common use cases
  - API documentation for daemon mode
- **API Documentation**: Complete REST API documentation
  - OpenAPI 3.0 specification (Swagger)
  - API reference with examples
  - Integration guides for Python, JavaScript, Go
  - Docker and Kubernetes deployment examples
- **Package Documentation**: Comprehensive Go package documentation
  - Full API reference for pkg.go.dev
  - Usage examples and best practices
  - Real-world integration patterns
  - Thread safety guarantees
- **HTTP API Tests**: Complete test suite for daemon mode
  - Health endpoint tests
  - Profiles endpoint tests
  - Compression endpoint tests with multiple scenarios
  - Error handling tests
  - Config priority tests
- **Automated Release Pipeline**: GitHub Actions for publishing
  - Automatic binary builds for multiple platforms (linux, darwin, freebsd × amd64, arm64)
  - GitHub Releases with changelog and checksums
  - Automatic publication to pkg.go.dev
  - Docker image builds and push to ghcr.io
  - Multi-architecture Docker support (amd64, arm64)
  - Documentation validation workflow
  - Release process documentation (RELEASING.md)
- **Configuration File Support**: `.slimjson` configuration file for custom profiles
  - Searches in current directory and user home directory
  - Simple INI-style format: `[profile-name]` followed by `key=value` pairs
  - Custom profiles take precedence over built-in profiles
  - Supports all compression parameters
  - See `.slimjson.example` for complete reference
- **CLI Predefined Profiles**: New `-profile` flag with 4 predefined compression profiles:
  - `light`: Light compression, preserves most data (MaxDepth: 10, MaxListLength: 20)
  - `medium`: Balanced compression (MaxDepth: 5, MaxListLength: 10)
  - `aggressive`: Removes verbose text fields (MaxDepth: 3, MaxListLength: 5, BlockList: description, summary, etc.)
  - `ai-optimized`: Removes URLs and metadata (MaxDepth: 4, MaxListLength: 8, BlockList: *_url fields)

**Basic Optimizations:**
- **Numeric Precision Control**: `-decimal-places N` rounds floats to N decimal places (e.g., 19.999 → 20.0)
- **Array Deduplication**: `-deduplicate` removes duplicate values from arrays
- **Array Sampling Strategies**: `-sample-strategy` with 4 modes:
  - `none`: No sampling (default)
  - `first_last`: Keep first N/2 and last N/2 elements
  - `random`: Random N elements
  - `representative`: Evenly distributed sampling
- **Sample Size Control**: `-sample-size N` specifies number of items to keep when sampling
- **String Truncation with Ellipsis**: Truncated strings now end with `...` to indicate content was cut

**Advanced Compression Features:**
- **Null Compression**: `-null-compression` tracks removed null fields in `_nulls` array
- **Type Inference**: `-type-inference` converts uniform arrays to schema+data format (40% savings)
- **Boolean Compression**: `-bool-compression` converts booleans to bit flags (60% savings)
- **Timestamp Compression**: `-timestamp-compression` converts ISO timestamps to unix timestamps (50% savings)
- **String Pooling**: `-string-pooling` deduplicates repeated strings (30-50% savings)
  - `-string-pool-min N` sets minimum occurrences (default: 2)
- **Number Delta Encoding**: `-number-delta` uses delta encoding for sequential numbers (50% savings)
  - `-number-delta-threshold N` sets minimum array size (default: 5)
- **Enum Detection**: `-enum-detection` converts repeated categorical values to enums (20-40% savings)
  - `-enum-max-values N` sets maximum unique values (default: 10)

### Changed
- **Profiles no longer truncate strings** to preserve data integrity - use BlockList instead to remove entire unnecessary fields
- **Profile flags can be overridden**: Use `-profile medium -decimal-places 2` to combine profile with custom settings
- Comprehensive compression testing suite in `testing/` directory
- Three real-world JSON test files (resume.json, schema-resume.json, users.json)
- Compression benchmark tool (`compression_benchmark.go`) with detailed metrics
- **Token counting** in compression tests - shows estimated token reduction for AI/LLM use cases
- **Statistical analysis** with standard deviation calculation (10 iterations per test)
- Detailed methodology documentation in `testing/METHODOLOGY.md`
- Performance benchmarks (`slimjson_bench_test.go`) with 9 different scenarios
- Compression results table in README showing real-world reduction percentages
- Four compression profiles: Light, Medium, Aggressive, AI-Optimized
- Performance metrics documentation (16-47µs processing time)
- Testing documentation in `testing/README.md`
- Makefile targets: `make bench` and `make compression-test`
- Shell script `run_tests.sh` for easy test execution
- Documentation files: EXAMPLES.md, QUICK_START.md, TESTING_SUMMARY.md
- **GitHub Pages** with Jekyll and Cayman theme for documentation
- Automatic deployment via GitHub Actions

### Changed
- Updated README with comprehensive compression results and benchmarks
- Enhanced .gitignore to exclude testing artifacts
- Improved Makefile with new testing and benchmarking targets

### Performance
- Small files (5KB): ~16µs per operation
- Medium files (25KB): ~39µs per operation
- Large files (28KB): ~47µs per operation
- Excellent parallel processing scalability

## [0.1.6] - Previous Release

### Features
- Core JSON slimming functionality
- CLI tool with multiple configuration options
- Docker/Podman support
- Multi-platform binaries (Linux, macOS, FreeBSD)
- Basic unit tests
