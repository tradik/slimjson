package slimjson

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".slimjson")

	configContent := `# Test config
[test-profile]
depth=5
list-len=10
strip-empty=true
decimal-places=2
deduplicate=true
block=field1,field2

[another-profile]
depth=3
string-pooling=true
string-pool-min=3
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Parse config file
	profiles, err := ParseConfigFile(configPath)
	if err != nil {
		t.Fatalf("Failed to parse config file: %v", err)
	}

	// Check test-profile
	testProfile, ok := profiles["test-profile"]
	if !ok {
		t.Fatal("Expected test-profile to exist")
	}

	if testProfile.MaxDepth != 5 {
		t.Errorf("Expected MaxDepth=5, got %d", testProfile.MaxDepth)
	}

	if testProfile.MaxListLength != 10 {
		t.Errorf("Expected MaxListLength=10, got %d", testProfile.MaxListLength)
	}

	if !testProfile.StripEmpty {
		t.Error("Expected StripEmpty=true")
	}

	if testProfile.DecimalPlaces != 2 {
		t.Errorf("Expected DecimalPlaces=2, got %d", testProfile.DecimalPlaces)
	}

	if !testProfile.DeduplicateArrays {
		t.Error("Expected DeduplicateArrays=true")
	}

	if len(testProfile.BlockList) != 2 {
		t.Errorf("Expected 2 blocked fields, got %d", len(testProfile.BlockList))
	}

	// Check another-profile
	anotherProfile, ok := profiles["another-profile"]
	if !ok {
		t.Fatal("Expected another-profile to exist")
	}

	if anotherProfile.MaxDepth != 3 {
		t.Errorf("Expected MaxDepth=3, got %d", anotherProfile.MaxDepth)
	}

	if !anotherProfile.StringPooling {
		t.Error("Expected StringPooling=true")
	}

	if anotherProfile.StringPoolMinOccurrences != 3 {
		t.Errorf("Expected StringPoolMinOccurrences=3, got %d", anotherProfile.StringPoolMinOccurrences)
	}

	t.Logf("Successfully parsed %d profiles", len(profiles))
}

func TestParseConfigFileWithComments(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".slimjson")

	configContent := `# This is a comment
// This is also a comment

[profile1]
# Comment before parameter
depth=5
// Another comment
list-len=10

# Empty lines should be ignored

[profile2]
depth=3
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	profiles, err := ParseConfigFile(configPath)
	if err != nil {
		t.Fatalf("Failed to parse config file: %v", err)
	}

	if len(profiles) != 2 {
		t.Errorf("Expected 2 profiles, got %d", len(profiles))
	}
}

func TestParseConfigFileInvalidSyntax(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".slimjson")

	configContent := `[profile]
invalid line without equals sign
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	_, err = ParseConfigFile(configPath)
	if err == nil {
		t.Error("Expected error for invalid syntax, got nil")
	}
}

func TestParseConfigFileInvalidValue(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".slimjson")

	configContent := `[profile]
depth=not-a-number
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	_, err = ParseConfigFile(configPath)
	if err == nil {
		t.Error("Expected error for invalid value, got nil")
	}
}

func TestApplyConfigParameter(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		checkFunc func(*Config) bool
	}{
		{
			name:  "depth",
			key:   "depth",
			value: "5",
			checkFunc: func(c *Config) bool {
				return c.MaxDepth == 5
			},
		},
		{
			name:  "list-len",
			key:   "list-len",
			value: "10",
			checkFunc: func(c *Config) bool {
				return c.MaxListLength == 10
			},
		},
		{
			name:  "strip-empty",
			key:   "strip-empty",
			value: "true",
			checkFunc: func(c *Config) bool {
				return c.StripEmpty == true
			},
		},
		{
			name:  "decimal-places",
			key:   "decimal-places",
			value: "2",
			checkFunc: func(c *Config) bool {
				return c.DecimalPlaces == 2
			},
		},
		{
			name:  "string-pooling",
			key:   "string-pooling",
			value: "true",
			checkFunc: func(c *Config) bool {
				return c.StringPooling == true
			},
		},
		{
			name:  "block-list",
			key:   "block",
			value: "field1,field2,field3",
			checkFunc: func(c *Config) bool {
				return len(c.BlockList) == 3
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			err := applyConfigParameter(cfg, tt.key, tt.value)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.checkFunc(cfg) {
				t.Errorf("Config parameter not applied correctly")
			}
		})
	}
}

func TestGetBuiltinProfiles(t *testing.T) {
	profiles := GetBuiltinProfiles()

	expectedProfiles := []string{"light", "medium", "aggressive", "ai-optimized"}
	for _, name := range expectedProfiles {
		if _, ok := profiles[name]; !ok {
			t.Errorf("Expected built-in profile '%s' not found", name)
		}
	}

	if len(profiles) != len(expectedProfiles) {
		t.Errorf("Expected %d built-in profiles, got %d", len(expectedProfiles), len(profiles))
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	// Change to a directory where .slimjson doesn't exist
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() { _ = os.Chdir(originalDir) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	profiles, err := LoadConfigFile()
	if err != nil {
		t.Errorf("Expected no error when config file not found, got: %v", err)
	}

	if len(profiles) != 0 {
		t.Errorf("Expected empty profiles map, got %d profiles", len(profiles))
	}
}

func TestConfigFileAllParameters(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".slimjson")

	configContent := `[full-config]
depth=5
list-len=10
string-len=100
strip-empty=true
block=field1,field2
decimal-places=2
deduplicate=true
sample-strategy=first_last
sample-size=20
null-compression=true
type-inference=true
bool-compression=true
timestamp-compression=true
string-pooling=true
string-pool-min=3
number-delta=true
number-delta-threshold=10
enum-detection=true
enum-max-values=5
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	profiles, err := ParseConfigFile(configPath)
	if err != nil {
		t.Fatalf("Failed to parse config file: %v", err)
	}

	cfg := profiles["full-config"]

	// Verify all parameters
	if cfg.MaxDepth != 5 {
		t.Errorf("MaxDepth: expected 5, got %d", cfg.MaxDepth)
	}
	if cfg.MaxListLength != 10 {
		t.Errorf("MaxListLength: expected 10, got %d", cfg.MaxListLength)
	}
	if cfg.MaxStringLength != 100 {
		t.Errorf("MaxStringLength: expected 100, got %d", cfg.MaxStringLength)
	}
	if !cfg.StripEmpty {
		t.Error("StripEmpty: expected true")
	}
	if len(cfg.BlockList) != 2 {
		t.Errorf("BlockList: expected 2 items, got %d", len(cfg.BlockList))
	}
	if cfg.DecimalPlaces != 2 {
		t.Errorf("DecimalPlaces: expected 2, got %d", cfg.DecimalPlaces)
	}
	if !cfg.DeduplicateArrays {
		t.Error("DeduplicateArrays: expected true")
	}
	if cfg.SampleStrategy != "first_last" {
		t.Errorf("SampleStrategy: expected 'first_last', got '%s'", cfg.SampleStrategy)
	}
	if cfg.SampleSize != 20 {
		t.Errorf("SampleSize: expected 20, got %d", cfg.SampleSize)
	}
	if !cfg.NullCompression {
		t.Error("NullCompression: expected true")
	}
	if !cfg.TypeInference {
		t.Error("TypeInference: expected true")
	}
	if !cfg.BoolCompression {
		t.Error("BoolCompression: expected true")
	}
	if !cfg.TimestampCompression {
		t.Error("TimestampCompression: expected true")
	}
	if !cfg.StringPooling {
		t.Error("StringPooling: expected true")
	}
	if cfg.StringPoolMinOccurrences != 3 {
		t.Errorf("StringPoolMinOccurrences: expected 3, got %d", cfg.StringPoolMinOccurrences)
	}
	if !cfg.NumberDeltaEncoding {
		t.Error("NumberDeltaEncoding: expected true")
	}
	if cfg.NumberDeltaThreshold != 10 {
		t.Errorf("NumberDeltaThreshold: expected 10, got %d", cfg.NumberDeltaThreshold)
	}
	if !cfg.EnumDetection {
		t.Error("EnumDetection: expected true")
	}
	if cfg.EnumMaxValues != 5 {
		t.Errorf("EnumMaxValues: expected 5, got %d", cfg.EnumMaxValues)
	}

	t.Log("All parameters parsed correctly")
}
