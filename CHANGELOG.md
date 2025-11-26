# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive compression testing suite in `testing/` directory
- Three real-world JSON test files (resume.json, schema-resume.json, users.json)
- Compression benchmark tool (`compression_benchmark.go`) with detailed metrics
- **Token counting** in compression tests - shows estimated token reduction for AI/LLM use cases
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
