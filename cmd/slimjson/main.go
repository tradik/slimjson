// Package main provides the CLI for slimjson.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tradik/slimjson"
)

// getProfile returns a configuration profile (built-in or from config file)
func getProfile(name string, customProfiles map[string]slimjson.Config) slimjson.Config {
	// First check custom profiles from config file
	if cfg, ok := customProfiles[strings.ToLower(name)]; ok {
		return cfg
	}

	// Then check built-in profiles
	builtinProfiles := slimjson.GetBuiltinProfiles()
	if cfg, ok := builtinProfiles[strings.ToLower(name)]; ok {
		return cfg
	}

	// Profile not found
	fmt.Fprintf(os.Stderr, "Unknown profile: %s\n", name)
	fmt.Fprintf(os.Stderr, "\nBuilt-in profiles: light, medium, aggressive, ai-optimized\n")

	if len(customProfiles) > 0 {
		fmt.Fprintf(os.Stderr, "\nCustom profiles from .slimjson:\n")
		for profileName := range customProfiles {
			fmt.Fprintf(os.Stderr, "  - %s\n", profileName)
		}
	}

	os.Exit(1)
	return slimjson.Config{}
}

func main() {
	var (
		profile                  string
		maxDepth                 int
		maxListLength            int
		maxStringLength          int
		stripEmpty               bool
		blockList                string
		pretty                   bool
		decimalPlaces            int
		deduplicateArrays        bool
		sampleStrategy           string
		sampleSize               int
		nullCompression          bool
		typeInference            bool
		boolCompression          bool
		timestampCompression     bool
		stringPooling            bool
		stringPoolMinOccurrences int
		numberDeltaEncoding      bool
		numberDeltaThreshold     int
		enumDetection            bool
		enumMaxValues            int
	)

	flag.StringVar(&profile, "profile", "", "Use predefined profile: light, medium, aggressive, ai-optimized")
	flag.IntVar(&maxDepth, "depth", 5, "Maximum nesting depth (0 for unlimited)")
	flag.IntVar(&maxListLength, "list-len", 10, "Maximum list length (0 for unlimited)")
	flag.IntVar(&maxStringLength, "string-len", 0, "Maximum string length in characters/runes (0 for unlimited)")
	flag.BoolVar(&stripEmpty, "strip-empty", true, "Remove nulls, empty strings, empty arrays/objects")
	flag.StringVar(&blockList, "block", "", "Comma-separated list of field names to remove")
	flag.BoolVar(&pretty, "pretty", false, "Pretty print output")
	flag.IntVar(&decimalPlaces, "decimal-places", -1, "Round floats to N decimal places (-1 for no rounding)")
	flag.BoolVar(&deduplicateArrays, "deduplicate", false, "Remove duplicate values from arrays")
	flag.StringVar(&sampleStrategy, "sample-strategy", "none", "Array sampling: none, first_last, random, representative")
	flag.IntVar(&sampleSize, "sample-size", 0, "Number of items when sampling (0 = use list-len)")
	flag.BoolVar(&nullCompression, "null-compression", false, "Track removed null fields in _nulls array")
	flag.BoolVar(&typeInference, "type-inference", false, "Convert uniform arrays to schema+data format")
	flag.BoolVar(&boolCompression, "bool-compression", false, "Convert booleans to bit flags")
	flag.BoolVar(&timestampCompression, "timestamp-compression", false, "Convert ISO timestamps to unix timestamps")
	flag.BoolVar(&stringPooling, "string-pooling", false, "Deduplicate repeated strings using string pool")
	flag.IntVar(&stringPoolMinOccurrences, "string-pool-min", 2, "Minimum occurrences for string pooling")
	flag.BoolVar(&numberDeltaEncoding, "number-delta", false, "Use delta encoding for sequential numbers")
	flag.IntVar(&numberDeltaThreshold, "number-delta-threshold", 5, "Minimum array size for delta encoding")
	flag.BoolVar(&enumDetection, "enum-detection", false, "Convert repeated categorical values to enums")
	flag.IntVar(&enumMaxValues, "enum-max-values", 10, "Maximum unique values to consider as enum")
	flag.Parse()

	// Load custom profiles from .slimjson config file
	customProfiles, err := slimjson.LoadConfigFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load .slimjson config file: %v\n", err)
		customProfiles = make(map[string]slimjson.Config)
	}

	var input io.Reader
	args := flag.Args()
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer func() { _ = f.Close() }()
		input = f
	} else {
		input = os.Stdin
	}

	decoder := json.NewDecoder(input)
	var data interface{}
	if err := decoder.Decode(&data); err != nil {
		if err == io.EOF {
			return
		}
		fmt.Fprintf(os.Stderr, "Error decoding JSON: %v\n", err)
		os.Exit(1)
	}

	// Apply profile if specified
	var cfg slimjson.Config
	if profile != "" {
		cfg = getProfile(profile, customProfiles)
		// Allow overriding profile settings with explicit flags
		if decimalPlaces >= 0 {
			cfg.DecimalPlaces = decimalPlaces
		}
		if deduplicateArrays {
			cfg.DeduplicateArrays = deduplicateArrays
		}
		if sampleStrategy != "none" {
			cfg.SampleStrategy = sampleStrategy
			cfg.SampleSize = sampleSize
		}
		// Apply advanced optimizations if specified
		if nullCompression {
			cfg.NullCompression = nullCompression
		}
		if typeInference {
			cfg.TypeInference = typeInference
		}
		if boolCompression {
			cfg.BoolCompression = boolCompression
		}
		if timestampCompression {
			cfg.TimestampCompression = timestampCompression
		}
		if stringPooling {
			cfg.StringPooling = stringPooling
			cfg.StringPoolMinOccurrences = stringPoolMinOccurrences
		}
		if numberDeltaEncoding {
			cfg.NumberDeltaEncoding = numberDeltaEncoding
			cfg.NumberDeltaThreshold = numberDeltaThreshold
		}
		if enumDetection {
			cfg.EnumDetection = enumDetection
			cfg.EnumMaxValues = enumMaxValues
		}
	} else {
		// Use custom parameters
		cfg = slimjson.Config{
			MaxDepth:                 maxDepth,
			MaxListLength:            maxListLength,
			MaxStringLength:          maxStringLength,
			StripEmpty:               stripEmpty,
			DecimalPlaces:            decimalPlaces,
			DeduplicateArrays:        deduplicateArrays,
			SampleStrategy:           sampleStrategy,
			SampleSize:               sampleSize,
			NullCompression:          nullCompression,
			TypeInference:            typeInference,
			BoolCompression:          boolCompression,
			TimestampCompression:     timestampCompression,
			StringPooling:            stringPooling,
			StringPoolMinOccurrences: stringPoolMinOccurrences,
			NumberDeltaEncoding:      numberDeltaEncoding,
			NumberDeltaThreshold:     numberDeltaThreshold,
			EnumDetection:            enumDetection,
			EnumMaxValues:            enumMaxValues,
		}
		if blockList != "" {
			cfg.BlockList = strings.Split(blockList, ",")
		}
	}

	slimmer := slimjson.New(cfg)
	result := slimmer.Slim(data)

	encoder := json.NewEncoder(os.Stdout)
	if pretty {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}
