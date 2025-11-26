package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tradik/slimjson"
)

// CompressionResult holds metrics for a single test
type CompressionResult struct {
	Filename             string
	OriginalSize         int
	CompressedSize       int
	Reduction            float64
	ReductionPct         float64
	OriginalTokens       int
	CompressedTokens     int
	TokenReduction       float64
	TokenReductionPct    float64
	ProcessingTime       time.Duration
	ProcessingTimeStdDev time.Duration
	Iterations           int
	ConfigUsed           string
}

// Statistics holds statistical metrics
type Statistics struct {
	Mean   float64
	StdDev float64
	Min    float64
	Max    float64
}

// TestConfig defines a compression configuration to test
type TestConfig struct {
	Name        string
	Config      slimjson.Config
	Description string
}

func main() {
	fixturesDir := "fixtures"

	// Define test configurations
	configs := []TestConfig{
		{
			Name: "Light",
			Config: slimjson.Config{
				MaxDepth:      10,
				MaxListLength: 20,
				StripEmpty:    true,
			},
			Description: "Light compression - preserve most data",
		},
		{
			Name: "Medium",
			Config: slimjson.Config{
				MaxDepth:      5,
				MaxListLength: 10,
				StripEmpty:    true,
			},
			Description: "Medium compression - balanced reduction",
		},
		{
			Name: "Aggressive",
			Config: slimjson.Config{
				MaxDepth:      3,
				MaxListLength: 5,
				StripEmpty:    true,
				BlockList:     []string{"description", "summary", "comment", "notes", "bio", "readme"},
			},
			Description: "Aggressive compression - removes verbose fields",
		},
		{
			Name: "AI-Optimized",
			Config: slimjson.Config{
				MaxDepth:      4,
				MaxListLength: 8,
				StripEmpty:    true,
				BlockList:     []string{"avatar_url", "gravatar_id", "url", "html_url", "followers_url", "following_url", "gists_url", "starred_url", "subscriptions_url", "organizations_url", "repos_url", "events_url", "received_events_url"},
			},
			Description: "Optimized for AI/LLM - removes URLs and metadata",
		},
	}

	// Get all JSON files in fixtures directory
	files, err := filepath.Glob(filepath.Join(fixturesDir, "*.json"))
	if err != nil {
		fmt.Printf("Error reading fixtures directory: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No JSON files found in fixtures directory")
		os.Exit(1)
	}

	fmt.Println("=== SlimJSON Compression Test Results ===")
	fmt.Println()

	// Test each file with each configuration
	for _, file := range files {
		filename := filepath.Base(file)
		fmt.Printf("Testing: %s\n", filename)
		fmt.Println(strings.Repeat("-", 80))

		// Read original file
		originalData, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			continue
		}

		// Parse JSON
		var data interface{}
		if err := json.Unmarshal(originalData, &data); err != nil {
			fmt.Printf("Error parsing JSON from %s: %v\n", file, err)
			continue
		}

		originalSize := len(originalData)

		// Test each configuration
		for _, testCfg := range configs {
			result := testCompression(filename, data, originalSize, originalData, testCfg)
			printResult(result)
		}

		fmt.Println()
	}

	// Generate summary table
	fmt.Println()
	fmt.Println("=== Summary Table (for README) ===")
	fmt.Println()
	generateMarkdownTable(files, configs)
}

func testCompression(filename string, data interface{}, originalSize int, originalData []byte, testCfg TestConfig) CompressionResult {
	slimmer := slimjson.New(testCfg.Config)

	// Run multiple iterations to calculate statistics
	const iterations = 10
	times := make([]float64, iterations)

	var compressedData []byte
	var err error

	for i := 0; i < iterations; i++ {
		start := time.Now()
		compressed := slimmer.Slim(data)
		elapsed := time.Since(start)
		times[i] = float64(elapsed.Nanoseconds())

		// Marshal on last iteration
		if i == iterations-1 {
			compressedData, err = json.Marshal(compressed)
			if err != nil {
				fmt.Printf("Error marshaling compressed data: %v\n", err)
				return CompressionResult{}
			}
		}
	}

	// Calculate statistics
	stats := calculateStatistics(times)
	avgTime := time.Duration(stats.Mean)
	stdDevTime := time.Duration(stats.StdDev)

	compressedSize := len(compressedData)
	reduction := float64(originalSize - compressedSize)
	reductionPct := (reduction / float64(originalSize)) * 100

	// Count tokens
	originalTokens := countTokens(string(originalData))
	compressedTokens := countTokens(string(compressedData))
	tokenReduction := float64(originalTokens - compressedTokens)
	tokenReductionPct := (tokenReduction / float64(originalTokens)) * 100

	return CompressionResult{
		Filename:             filename,
		OriginalSize:         originalSize,
		CompressedSize:       compressedSize,
		Reduction:            reduction,
		ReductionPct:         reductionPct,
		OriginalTokens:       originalTokens,
		CompressedTokens:     compressedTokens,
		TokenReduction:       tokenReduction,
		TokenReductionPct:    tokenReductionPct,
		ProcessingTime:       avgTime,
		ProcessingTimeStdDev: stdDevTime,
		Iterations:           iterations,
		ConfigUsed:           testCfg.Name,
	}
}

// calculateStatistics computes mean, standard deviation, min, and max
func calculateStatistics(values []float64) Statistics {
	if len(values) == 0 {
		return Statistics{}
	}

	// Calculate mean
	var sum float64
	min := values[0]
	max := values[0]

	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	mean := sum / float64(len(values))

	// Calculate standard deviation
	var varianceSum float64
	for _, v := range values {
		diff := v - mean
		varianceSum += diff * diff
	}
	variance := varianceSum / float64(len(values))
	stdDev := math.Sqrt(variance)

	return Statistics{
		Mean:   mean,
		StdDev: stdDev,
		Min:    min,
		Max:    max,
	}
}

func printResult(result CompressionResult) {
	fmt.Printf("  Config: %s\n", result.ConfigUsed)
	fmt.Printf("    Original:    %s (%d tokens)\n", formatBytes(result.OriginalSize), result.OriginalTokens)
	fmt.Printf("    Compressed:  %s (%d tokens)\n", formatBytes(result.CompressedSize), result.CompressedTokens)
	fmt.Printf("    Reduction:   %s (%.2f%%) | Tokens: %d (%.2f%%)\n",
		formatBytes(int(result.Reduction)), result.ReductionPct,
		int(result.TokenReduction), result.TokenReductionPct)
	fmt.Printf("    Time:        %v ± %v (n=%d)\n",
		result.ProcessingTime, result.ProcessingTimeStdDev, result.Iterations)
	fmt.Println()
}

// countTokens estimates token count using a simple word-based approach
// This approximates GPT-style tokenization (roughly 1 token per 4 characters for English)
func countTokens(text string) int {
	// Remove whitespace and count characters
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}

	// Approximate: 1 token ≈ 4 characters (common for JSON/English in GPT models)
	// This is a rough estimate; real tokenization varies by model
	charCount := len(text)
	return (charCount + 3) / 4 // Round up
}

func formatBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func generateMarkdownTable(files []string, configs []TestConfig) {
	fmt.Println("| File | Original Size | Config | Compressed Size | Reduction | Reduction % | Original Tokens | Compressed Tokens | Token Reduction % |")
	fmt.Println("|------|---------------|--------|-----------------|-----------|-------------|-----------------|-------------------|-------------------|")

	for _, file := range files {
		filename := filepath.Base(file)

		// Read original file
		originalData, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Parse JSON
		var data interface{}
		if err := json.Unmarshal(originalData, &data); err != nil {
			continue
		}

		originalSize := len(originalData)

		// Test each configuration
		for _, testCfg := range configs {
			result := testCompression(filename, data, originalSize, originalData, testCfg)
			fmt.Printf("| %s | %s | %s | %s | %s | %.1f%% | %d | %d | %.1f%% |\n",
				result.Filename,
				formatBytes(result.OriginalSize),
				result.ConfigUsed,
				formatBytes(result.CompressedSize),
				formatBytes(int(result.Reduction)),
				result.ReductionPct,
				result.OriginalTokens,
				result.CompressedTokens,
				result.TokenReductionPct,
			)
		}
	}
}
