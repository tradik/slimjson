package slimjson

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ProfileConfig represents a named configuration profile
type ProfileConfig struct {
	Name   string
	Config Config
}

// LoadConfigFile loads configuration from .slimjson file
// Searches in: current directory, user home directory
func LoadConfigFile() (map[string]Config, error) {
	// Try current directory first
	configPath := ".slimjson"
	if _, err := os.Stat(configPath); err != nil {
		// Try home directory
		home, err := os.UserHomeDir()
		if err == nil {
			configPath = filepath.Join(home, ".slimjson")
			if _, err := os.Stat(configPath); err != nil {
				// No config file found - return empty map (not an error)
				return make(map[string]Config), nil
			}
		} else {
			return make(map[string]Config), nil
		}
	}

	return ParseConfigFile(configPath)
}

// ParseConfigFile parses a .slimjson configuration file
func ParseConfigFile(path string) (map[string]Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() { _ = file.Close() }()

	profiles := make(map[string]Config)
	var currentProfile string
	var currentConfig Config

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// Check for profile section [name]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// Save previous profile if exists
			if currentProfile != "" {
				profiles[currentProfile] = currentConfig
			}

			// Start new profile
			currentProfile = strings.TrimSpace(line[1 : len(line)-1])
			currentConfig = Config{
				DecimalPlaces: -1, // Default: no rounding
			}
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid syntax at line %d: %s", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Apply parameter to current config
		if err := applyConfigParameter(&currentConfig, key, value); err != nil {
			return nil, fmt.Errorf("error at line %d: %w", lineNum, err)
		}
	}

	// Save last profile
	if currentProfile != "" {
		profiles[currentProfile] = currentConfig
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return profiles, nil
}

// applyConfigParameter applies a single parameter to config
func applyConfigParameter(cfg *Config, key, value string) error {
	key = strings.ToLower(key)

	// Try basic parameters
	if err := applyBasicParameter(cfg, key, value); err == nil {
		return nil
	} else if err != errUnknownParameter {
		return err
	}

	// Try advanced parameters
	if err := applyAdvancedParameter(cfg, key, value); err == nil {
		return nil
	} else if err != errUnknownParameter {
		return err
	}

	return fmt.Errorf("unknown parameter: %s", key)
}

var errUnknownParameter = fmt.Errorf("unknown parameter")

func applyBasicParameter(cfg *Config, key, value string) error {
	switch key {
	case "depth", "max-depth", "maxdepth":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid depth value: %s", value)
		}
		cfg.MaxDepth = v

	case "list-len", "list-length", "max-list-length", "maxlistlength":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid list-len value: %s", value)
		}
		cfg.MaxListLength = v

	case "string-len", "string-length", "max-string-length", "maxstringlength":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid string-len value: %s", value)
		}
		cfg.MaxStringLength = v

	case "strip-empty", "stripempty":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid strip-empty value: %s", value)
		}
		cfg.StripEmpty = v

	case "block", "block-list", "blocklist":
		if value != "" {
			cfg.BlockList = strings.Split(value, ",")
			for i := range cfg.BlockList {
				cfg.BlockList[i] = strings.TrimSpace(cfg.BlockList[i])
			}
		}

	case "decimal-places", "decimalplaces":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid decimal-places value: %s", value)
		}
		cfg.DecimalPlaces = v

	case "deduplicate", "deduplicate-arrays", "deduplicatearrays":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid deduplicate value: %s", value)
		}
		cfg.DeduplicateArrays = v

	case "sample-strategy", "samplestrategy":
		cfg.SampleStrategy = value

	case "sample-size", "samplesize":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid sample-size value: %s", value)
		}
		cfg.SampleSize = v

	default:
		return errUnknownParameter
	}
	return nil
}

func applyAdvancedParameter(cfg *Config, key, value string) error {
	switch key {
	case "null-compression", "nullcompression":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid null-compression value: %s", value)
		}
		cfg.NullCompression = v

	case "type-inference", "typeinference":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid type-inference value: %s", value)
		}
		cfg.TypeInference = v

	case "bool-compression", "boolcompression":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid bool-compression value: %s", value)
		}
		cfg.BoolCompression = v

	case "timestamp-compression", "timestampcompression":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid timestamp-compression value: %s", value)
		}
		cfg.TimestampCompression = v

	case "string-pooling", "stringpooling":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid string-pooling value: %s", value)
		}
		cfg.StringPooling = v

	case "string-pool-min", "stringpoolmin":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid string-pool-min value: %s", value)
		}
		cfg.StringPoolMinOccurrences = v

	case "number-delta", "numberdelta":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid number-delta value: %s", value)
		}
		cfg.NumberDeltaEncoding = v

	case "number-delta-threshold", "numberdeltathreshold":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid number-delta-threshold value: %s", value)
		}
		cfg.NumberDeltaThreshold = v

	case "enum-detection", "enumdetection":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid enum-detection value: %s", value)
		}
		cfg.EnumDetection = v

	case "enum-max-values", "enummaxvalues":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid enum-max-values value: %s", value)
		}
		cfg.EnumMaxValues = v

	case "strip-emoji", "stripemoji", "strip-utf8-emoji", "striputf8emoji":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid strip-emoji value: %s", value)
		}
		cfg.StripUTF8Emoji = v

	default:
		return errUnknownParameter
	}
	return nil
}

// GetBuiltinProfiles returns the built-in profiles (light, medium, aggressive, ai-optimized)
func GetBuiltinProfiles() map[string]Config {
	return map[string]Config{
		"light": {
			MaxDepth:      10,
			MaxListLength: 20,
			StripEmpty:    true,
		},
		"medium": {
			MaxDepth:      5,
			MaxListLength: 10,
			StripEmpty:    true,
		},
		"aggressive": {
			MaxDepth:      3,
			MaxListLength: 5,
			StripEmpty:    true,
			BlockList:     []string{"description", "summary", "comment", "notes", "bio", "readme"},
		},
		"ai-optimized": {
			MaxDepth:      4,
			MaxListLength: 8,
			StripEmpty:    true,
			BlockList:     []string{"avatar_url", "gravatar_id", "url", "html_url", "followers_url", "following_url", "gists_url", "starred_url", "subscriptions_url", "organizations_url", "repos_url", "events_url", "received_events_url"},
		},
	}
}
