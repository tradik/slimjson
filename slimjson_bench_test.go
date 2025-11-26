package slimjson

import (
	"encoding/json"
	"os"
	"testing"
)

// BenchmarkSlim_Small tests performance on small JSON (5KB)
func BenchmarkSlim_Small(b *testing.B) {
	data := loadTestData(b, "testing/fixtures/users.json")
	cfg := Config{
		MaxDepth:      5,
		MaxListLength: 10,
		StripEmpty:    true,
	}
	slimmer := New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slimmer.Slim(data)
	}
}

// BenchmarkSlim_Medium tests performance on medium JSON (25KB)
func BenchmarkSlim_Medium(b *testing.B) {
	data := loadTestData(b, "testing/fixtures/schema-resume.json")
	cfg := Config{
		MaxDepth:      5,
		MaxListLength: 10,
		StripEmpty:    true,
	}
	slimmer := New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slimmer.Slim(data)
	}
}

// BenchmarkSlim_Large tests performance on large JSON (28KB)
func BenchmarkSlim_Large(b *testing.B) {
	data := loadTestData(b, "testing/fixtures/resume.json")
	cfg := Config{
		MaxDepth:      5,
		MaxListLength: 10,
		StripEmpty:    true,
	}
	slimmer := New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slimmer.Slim(data)
	}
}

// BenchmarkSlim_DeepNesting tests performance with deep nesting limits
func BenchmarkSlim_DeepNesting(b *testing.B) {
	data := loadTestData(b, "testing/fixtures/schema-resume.json")
	cfg := Config{
		MaxDepth:      10,
		MaxListLength: 20,
		StripEmpty:    true,
	}
	slimmer := New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slimmer.Slim(data)
	}
}

// BenchmarkSlim_Aggressive tests performance with aggressive compression
func BenchmarkSlim_Aggressive(b *testing.B) {
	data := loadTestData(b, "testing/fixtures/schema-resume.json")
	cfg := Config{
		MaxDepth:        3,
		MaxListLength:   5,
		MaxStringLength: 100,
		StripEmpty:      true,
		BlockList:       []string{"description", "summary", "comment"},
	}
	slimmer := New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slimmer.Slim(data)
	}
}

// BenchmarkSlim_NoLimits tests performance with no compression limits
func BenchmarkSlim_NoLimits(b *testing.B) {
	data := loadTestData(b, "testing/fixtures/schema-resume.json")
	cfg := Config{
		MaxDepth:      0,
		MaxListLength: 0,
		StripEmpty:    false,
	}
	slimmer := New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slimmer.Slim(data)
	}
}

// BenchmarkSlim_StringTruncation tests performance with string truncation
func BenchmarkSlim_StringTruncation(b *testing.B) {
	data := loadTestData(b, "testing/fixtures/schema-resume.json")
	cfg := Config{
		MaxDepth:        5,
		MaxListLength:   10,
		MaxStringLength: 50,
		StripEmpty:      true,
	}
	slimmer := New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slimmer.Slim(data)
	}
}

// BenchmarkSlim_BlockList tests performance with blocklist filtering
func BenchmarkSlim_BlockList(b *testing.B) {
	data := loadTestData(b, "testing/fixtures/schema-resume.json")
	cfg := Config{
		MaxDepth:      5,
		MaxListLength: 10,
		StripEmpty:    true,
		BlockList:     []string{"url", "avatar_url", "html_url", "gravatar_id", "description"},
	}
	slimmer := New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slimmer.Slim(data)
	}
}

// BenchmarkSlim_Parallel tests parallel processing performance
func BenchmarkSlim_Parallel(b *testing.B) {
	data := loadTestData(b, "testing/fixtures/schema-resume.json")
	cfg := Config{
		MaxDepth:      5,
		MaxListLength: 10,
		StripEmpty:    true,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		slimmer := New(cfg)
		for pb.Next() {
			_ = slimmer.Slim(data)
		}
	})
}

// Helper function to load test data
func loadTestData(b *testing.B, filepath string) interface{} {
	b.Helper()

	fileData, err := os.ReadFile(filepath)
	if err != nil {
		b.Fatalf("Failed to read test file %s: %v", filepath, err)
	}

	var data interface{}
	if err := json.Unmarshal(fileData, &data); err != nil {
		b.Fatalf("Failed to unmarshal test data: %v", err)
	}

	return data
}
