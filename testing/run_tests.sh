#!/bin/bash

set -e

echo "=== Running SlimJSON Compression Tests ==="
echo ""

# Change to testing directory
cd "$(dirname "$0")"

# Run the compression benchmark
echo "Running compression tests..."
echo ""
go run compression_benchmark.go

echo ""
echo "=== Tests completed ==="
