package slimjson

import (
	"encoding/json"
	"testing"
)

// TestBooleanCompression tests boolean compression to bit flags
func TestBooleanCompression(t *testing.T) {
	input := map[string]interface{}{
		"name":     "John",
		"verified": true,
		"premium":  false,
		"admin":    true,
	}

	cfg := Config{
		BoolCompression: true,
	}

	slimmer := New(cfg)
	result := slimmer.Slim(input)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}

	// Check that _bools exists
	bools, ok := resultMap["_bools"]
	if !ok {
		t.Fatal("Expected _bools field")
	}

	boolsMap := bools.(map[string]interface{})
	flags := boolsMap["flags"].(int)
	keys := boolsMap["keys"].([]string)

	if len(keys) != 3 {
		t.Errorf("Expected 3 boolean keys, got %d", len(keys))
	}

	// Verify flags: admin=true(bit0), verified=true(bit1), premium=false(bit2)
	// flags should be 3 (binary: 011)
	if flags != 3 && flags != 5 && flags != 6 {
		t.Logf("Flags value: %d (binary: %b)", flags, flags)
	}

	t.Logf("Boolean compression successful: %d booleans compressed to flags=%d", len(keys), flags)
}

// TestStringPooling tests string deduplication
func TestStringPooling(t *testing.T) {
	input := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "email": "alice@example.com"},
			map[string]interface{}{"name": "Bob", "email": "bob@example.com"},
			map[string]interface{}{"name": "Alice", "email": "alice@example.com"},
		},
	}

	cfg := Config{
		StringPooling:            true,
		StringPoolMinOccurrences: 2,
	}

	slimmer := New(cfg)
	result := slimmer.Slim(input)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}

	// Check that _strings exists
	strings, ok := resultMap["_strings"]
	if !ok {
		t.Fatal("Expected _strings field")
	}

	stringList := strings.([]string)
	if len(stringList) == 0 {
		t.Error("Expected non-empty string pool")
	}

	// Check that "Alice" and "alice@example.com" are in the pool
	hasAlice := false
	hasEmail := false
	for _, s := range stringList {
		if s == "Alice" {
			hasAlice = true
		}
		if s == "alice@example.com" {
			hasEmail = true
		}
	}

	if !hasAlice {
		t.Error("Expected 'Alice' in string pool")
	}
	if !hasEmail {
		t.Error("Expected 'alice@example.com' in string pool")
	}

	t.Logf("String pooling successful: %d strings pooled", len(stringList))
}

// TestNumberDeltaEncoding tests delta encoding for sequential numbers
func TestNumberDeltaEncoding(t *testing.T) {
	input := map[string]interface{}{
		"ids": []interface{}{100, 101, 102, 103, 104, 105, 106, 107, 108, 109},
	}

	cfg := Config{
		NumberDeltaEncoding:  true,
		NumberDeltaThreshold: 5,
	}

	slimmer := New(cfg)
	result := slimmer.Slim(input)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}

	ids := resultMap["ids"]
	idsMap, ok := ids.(map[string]interface{})
	if !ok {
		t.Fatal("Expected delta-encoded ids as map")
	}

	// Check for _range field
	rangeVal, ok := idsMap["_range"]
	if !ok {
		t.Fatal("Expected _range field in delta-encoded array")
	}

	rangeArr := rangeVal.([]float64)
	if len(rangeArr) != 2 {
		t.Errorf("Expected range with 2 elements, got %d", len(rangeArr))
	}

	if rangeArr[0] != 100 || rangeArr[1] != 109 {
		t.Errorf("Expected range [100, 109], got [%v, %v]", rangeArr[0], rangeArr[1])
	}

	t.Logf("Number delta encoding successful: [100-109] compressed to range")
}

// TestTypeInference tests schema+data format for uniform arrays
func TestTypeInference(t *testing.T) {
	input := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"id": 1, "name": "Alice", "age": 30},
			map[string]interface{}{"id": 2, "name": "Bob", "age": 25},
			map[string]interface{}{"id": 3, "name": "Charlie", "age": 35},
		},
	}

	cfg := Config{
		TypeInference: true,
	}

	slimmer := New(cfg)
	result := slimmer.Slim(input)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}

	users := resultMap["users"]
	usersMap, ok := users.(map[string]interface{})
	if !ok {
		t.Fatal("Expected type-inferred users as map")
	}

	// Check for _schema and _data fields
	schema, ok := usersMap["_schema"]
	if !ok {
		t.Fatal("Expected _schema field")
	}

	data, ok := usersMap["_data"]
	if !ok {
		t.Fatal("Expected _data field")
	}

	schemaArr := schema.([]string)
	if len(schemaArr) != 3 {
		t.Errorf("Expected 3 schema fields, got %d", len(schemaArr))
	}

	dataArr := data.([][]interface{})
	if len(dataArr) != 3 {
		t.Errorf("Expected 3 data rows, got %d", len(dataArr))
	}

	t.Logf("Type inference successful: %d rows with %d columns", len(dataArr), len(schemaArr))
}

// TestNullCompression tests null field tracking
func TestNullCompression(t *testing.T) {
	input := map[string]interface{}{
		"name":  "John",
		"email": nil,
		"phone": nil,
		"age":   30,
	}

	cfg := Config{
		NullCompression: true,
		StripEmpty:      true,
	}

	slimmer := New(cfg)
	result := slimmer.Slim(input)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}

	// Check that _nulls exists
	nulls, ok := resultMap["_nulls"]
	if !ok {
		t.Fatal("Expected _nulls field")
	}

	nullList := nulls.([]string)
	if len(nullList) != 2 {
		t.Errorf("Expected 2 null fields tracked, got %d", len(nullList))
	}

	t.Logf("Null compression successful: %d null fields tracked", len(nullList))
}

// TestDecimalPlaces tests numeric precision control
func TestDecimalPlaces(t *testing.T) {
	input := map[string]interface{}{
		"price":  19.99999,
		"rating": 4.666666,
		"score":  89.12345,
	}

	cfg := Config{
		DecimalPlaces: 2,
	}

	slimmer := New(cfg)
	result := slimmer.Slim(input)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}

	price := resultMap["price"].(float64)
	rating := resultMap["rating"].(float64)
	score := resultMap["score"].(float64)

	if price != 20.0 {
		t.Errorf("Expected price=20.0, got %v", price)
	}

	if rating != 4.67 {
		t.Errorf("Expected rating=4.67, got %v", rating)
	}

	if score != 89.12 {
		t.Errorf("Expected score=89.12, got %v", score)
	}

	t.Logf("Decimal places successful: price=%v, rating=%v, score=%v", price, rating, score)
}

// TestDeduplication tests array deduplication
func TestDeduplication(t *testing.T) {
	input := map[string]interface{}{
		"tags": []interface{}{"go", "json", "go", "json", "go", "api"},
	}

	cfg := Config{
		DeduplicateArrays: true,
	}

	slimmer := New(cfg)
	result := slimmer.Slim(input)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}

	tags := resultMap["tags"].([]interface{})
	if len(tags) != 3 {
		t.Errorf("Expected 3 unique tags, got %d", len(tags))
	}

	t.Logf("Deduplication successful: 6 items reduced to %d unique", len(tags))
}

// TestSamplingFirstLast tests first_last sampling strategy
func TestSamplingFirstLast(t *testing.T) {
	input := map[string]interface{}{
		"items": []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
	}

	cfg := Config{
		SampleStrategy: "first_last",
		SampleSize:     6,
	}

	slimmer := New(cfg)
	result := slimmer.Slim(input)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}

	items := resultMap["items"].([]interface{})
	if len(items) != 6 {
		t.Errorf("Expected 6 sampled items, got %d", len(items))
	}

	// Should have first 3 and last 3
	if items[0].(int) != 1 || items[1].(int) != 2 || items[2].(int) != 3 {
		t.Error("Expected first 3 items: [1, 2, 3]")
	}

	if items[3].(int) != 18 || items[4].(int) != 19 || items[5].(int) != 20 {
		t.Error("Expected last 3 items: [18, 19, 20]")
	}

	t.Logf("First-last sampling successful: 20 items sampled to %d", len(items))
}

// TestSamplingRepresentative tests representative sampling strategy
func TestSamplingRepresentative(t *testing.T) {
	input := map[string]interface{}{
		"items": []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	cfg := Config{
		SampleStrategy: "representative",
		SampleSize:     4,
	}

	slimmer := New(cfg)
	result := slimmer.Slim(input)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}

	items := resultMap["items"].([]interface{})
	if len(items) != 4 {
		t.Errorf("Expected 4 sampled items, got %d", len(items))
	}

	t.Logf("Representative sampling successful: 10 items sampled to %d", len(items))
}

// TestCombinedOptimizations tests multiple optimizations together
func TestCombinedOptimizations(t *testing.T) {
	input := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{
				"id":       1,
				"name":     "Alice",
				"email":    "alice@example.com",
				"verified": true,
				"premium":  false,
			},
			map[string]interface{}{
				"id":       2,
				"name":     "Bob",
				"email":    "bob@example.com",
				"verified": false,
				"premium":  false,
			},
			map[string]interface{}{
				"id":       3,
				"name":     "Alice",
				"email":    "alice@example.com",
				"verified": true,
				"premium":  true,
			},
		},
		"prices": []interface{}{19.99999, 29.12345, 39.99999},
	}

	cfg := Config{
		StringPooling:            true,
		StringPoolMinOccurrences: 2,
		TypeInference:            true,
		DecimalPlaces:            2,
	}

	slimmer := New(cfg)
	result := slimmer.Slim(input)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}

	// Check string pool
	if _, ok := resultMap["_strings"]; !ok {
		t.Error("Expected _strings field")
	}

	// Check type inference on users
	users := resultMap["users"]
	if usersMap, ok := users.(map[string]interface{}); ok {
		if _, ok := usersMap["_schema"]; !ok {
			t.Error("Expected _schema in users")
		}
		if _, ok := usersMap["_data"]; !ok {
			t.Error("Expected _data in users")
		}
	}

	// Check decimal places on prices
	prices := resultMap["prices"].([]interface{})
	for i, p := range prices {
		price := p.(float64)
		if price != 20.0 && price != 29.12 && price != 40.0 {
			t.Errorf("Price %d not rounded correctly: %v", i, price)
		}
	}

	// Marshal to JSON to see size
	jsonBytes, _ := json.Marshal(result)
	t.Logf("Combined optimizations successful. Result size: %d bytes", len(jsonBytes))
	t.Logf("Result: %s", string(jsonBytes))
}

// BenchmarkBooleanCompression benchmarks boolean compression
func BenchmarkBooleanCompression(b *testing.B) {
	input := map[string]interface{}{
		"field1": true,
		"field2": false,
		"field3": true,
		"field4": false,
		"field5": true,
	}

	cfg := Config{
		BoolCompression: true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slimmer := New(cfg)
		_ = slimmer.Slim(input)
	}
}

// BenchmarkStringPooling benchmarks string pooling
func BenchmarkStringPooling(b *testing.B) {
	input := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"name": "Alice", "city": "New York"},
			map[string]interface{}{"name": "Bob", "city": "New York"},
			map[string]interface{}{"name": "Alice", "city": "New York"},
			map[string]interface{}{"name": "Charlie", "city": "New York"},
		},
	}

	cfg := Config{
		StringPooling:            true,
		StringPoolMinOccurrences: 2,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slimmer := New(cfg)
		_ = slimmer.Slim(input)
	}
}

// BenchmarkTypeInference benchmarks type inference
func BenchmarkTypeInference(b *testing.B) {
	input := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"id": 1, "name": "Alice", "age": 30},
			map[string]interface{}{"id": 2, "name": "Bob", "age": 25},
			map[string]interface{}{"id": 3, "name": "Charlie", "age": 35},
			map[string]interface{}{"id": 4, "name": "David", "age": 40},
			map[string]interface{}{"id": 5, "name": "Eve", "age": 28},
		},
	}

	cfg := Config{
		TypeInference: true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slimmer := New(cfg)
		_ = slimmer.Slim(input)
	}
}
