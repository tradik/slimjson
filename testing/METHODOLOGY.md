# Testing Methodology

## Statistical Analysis

### Standard Deviation Calculation

All compression tests are run **10 times** for each configuration to ensure statistical reliability. We calculate:

- **Mean (Average)**: Average processing time across all iterations
- **Standard Deviation (σ)**: Measure of variability in processing times
- **Iterations (n)**: Number of test runs (n=10)

### Formula

**Standard Deviation:**
```
σ = √(Σ(xi - μ)² / n)

Where:
- xi = individual measurement
- μ = mean (average)
- n = number of measurements
```

**Coefficient of Variation (CV):**
```
CV = (σ / μ) × 100%
```

Lower CV indicates more consistent performance.

## Token Counting

Tokens are estimated using a simplified approximation:

```
tokens ≈ characters / 4
```

This approximates GPT-style tokenization for JSON/English text:
- **1 token** ≈ **4 characters** (average)
- Actual tokenization varies by model (GPT-3.5, GPT-4, Claude, etc.)
- This is a conservative estimate for planning purposes

### Why This Approximation?

1. **Simplicity**: No external dependencies on tokenizer libraries
2. **Consistency**: Same calculation across all tests
3. **Accuracy**: Within 10-15% of actual GPT tokenization for JSON
4. **Speed**: Instant calculation without API calls

### Real Tokenization

For production use, consider using actual tokenizers:
- **tiktoken** (OpenAI): `pip install tiktoken`
- **transformers** (HuggingFace): `pip install transformers`
- **anthropic** (Claude): Use their API

## Test Execution

### Iteration Count

We use **10 iterations** because:
- Sufficient for statistical significance (n≥10)
- Fast enough for quick feedback
- Captures performance variability
- Standard in microbenchmarking

### Warm-up

The first iteration includes:
- Memory allocation
- JIT compilation (if applicable)
- Cache warming

Subsequent iterations measure steady-state performance.

### Environment

Tests should be run on:
- **Idle system**: Minimize background processes
- **Consistent hardware**: Same CPU/memory
- **Stable temperature**: Avoid thermal throttling

## Metrics Explained

### Processing Time

**Mean ± StdDev (n=10)**

Example: `45.2µs ± 3.1µs (n=10)`

- **Mean**: 45.2µs average processing time
- **StdDev**: ±3.1µs variation (68% of runs within this range)
- **n=10**: Based on 10 test runs

**Interpretation:**
- Low StdDev (< 10% of mean): Consistent performance
- High StdDev (> 20% of mean): Variable performance, investigate

### Size Reduction

**Percentage Reduction:**
```
reduction% = ((original - compressed) / original) × 100
```

Example: 60.5% reduction means:
- Original: 28.2 KB
- Compressed: 11.2 KB
- Saved: 17.0 KB (60.5%)

### Token Reduction

Same formula as size reduction, but for estimated tokens:

```
token_reduction% = ((original_tokens - compressed_tokens) / original_tokens) × 100
```

Example: 60.5% token reduction means:
- Original: 7230 tokens
- Compressed: 2859 tokens
- Saved: 4371 tokens (60.5%)

**Cost Impact:**
If API costs $0.01 per 1000 tokens:
- Original cost: $0.0723
- Compressed cost: $0.0286
- Savings: $0.0437 (60.5%)

## Benchmark Profiles

### Light Compression
- **Target**: Minimal data loss
- **Use Case**: Preserve structure, remove empties
- **Expected Reduction**: 20-30%

### Medium Compression
- **Target**: Balanced reduction
- **Use Case**: General purpose, API responses
- **Expected Reduction**: 30-40%

### Aggressive Compression
- **Target**: Maximum reduction
- **Use Case**: Previews, summaries, extreme optimization
- **Expected Reduction**: 85-98%

### AI-Optimized
- **Target**: Token reduction for LLMs
- **Use Case**: Sending to GPT/Claude/etc.
- **Expected Reduction**: 50-65%

## Statistical Confidence

### Confidence Intervals

For 10 iterations, 95% confidence interval:

```
CI = mean ± (1.96 × σ / √n)
CI = mean ± (1.96 × σ / √10)
CI = mean ± (0.62 × σ)
```

### Sample Size Justification

**Why n=10?**

| n | Confidence | Speed | Accuracy |
|---|------------|-------|----------|
| 3 | Low | Fast | ±30% |
| 5 | Medium | Fast | ±20% |
| 10 | Good | Medium | ±10% |
| 30 | High | Slow | ±5% |
| 100 | Very High | Very Slow | ±2% |

We chose n=10 as optimal balance between accuracy and speed.

## Reproducibility

To reproduce results:

```bash
# Run tests
cd testing
go run compression_benchmark.go

# Results will vary slightly due to:
# - System load
# - CPU frequency scaling
# - Memory pressure
# - Cache state
```

### Expected Variance

Typical standard deviation ranges:
- **Small files (5KB)**: ±1-2µs
- **Medium files (25KB)**: ±3-5µs
- **Large files (50KB+)**: ±5-10µs

Higher variance indicates:
- System under load
- Thermal throttling
- Background processes
- Memory pressure

## Validation

### Sanity Checks

1. **Compression never increases size** (except edge cases)
2. **Token count proportional to size**
3. **StdDev < 20% of mean** (consistent performance)
4. **Aggressive > Medium > Light** (reduction order)

### Known Limitations

1. **Token estimation**: ±10-15% accuracy vs real tokenizers
2. **Small sample size**: n=10 may miss rare outliers
3. **Single-threaded**: Doesn't test parallel performance
4. **Cold start**: First iteration may be slower

## References

- **Standard Deviation**: https://en.wikipedia.org/wiki/Standard_deviation
- **Coefficient of Variation**: https://en.wikipedia.org/wiki/Coefficient_of_variation
- **Microbenchmarking**: https://go.dev/blog/benchmarks
- **GPT Tokenization**: https://platform.openai.com/tokenizer

## Future Improvements

1. **Real tokenizer integration**: Use tiktoken for accurate counts
2. **Larger sample sizes**: Configurable n (10, 30, 100)
3. **Percentile reporting**: P50, P95, P99 latencies
4. **Outlier detection**: Identify and report anomalies
5. **Parallel testing**: Multi-threaded performance
6. **Memory profiling**: Track allocation patterns
