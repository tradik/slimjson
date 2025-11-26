// Package main provides the CLI for slimjson.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
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

// printUsage prints the usage information
func printUsage() {
	fmt.Fprintf(os.Stderr, `slimjson - JSON optimizer for AI/LLM contexts

Usage:
  slimjson [options] [file]              Process JSON file or stdin
  slimjson -d [options]                  Run as HTTP daemon
  slimjson -h                            Show this help

Daemon Mode:
  -d, -daemon                Run as HTTP daemon listening on specified port
  -port int                  Port for daemon mode (default: 8080)

Configuration:
  -c, -config string         Path to custom config file (takes priority over .slimjson)
  -profile string            Use predefined profile: light, medium, aggressive, ai-optimized

Basic Options:
  -depth int                 Maximum nesting depth (default: 5, 0 = unlimited)
  -list-len int              Maximum list length (default: 10, 0 = unlimited)
  -string-len int            Maximum string length (default: 0 = unlimited)
  -strip-empty               Remove nulls, empty strings, empty arrays/objects (default: true)
  -block string              Comma-separated list of field names to remove
  -pretty                    Pretty print output

Optimization Options:
  -decimal-places int        Round floats to N decimal places (default: -1 = no rounding)
  -deduplicate               Remove duplicate values from arrays
  -sample-strategy string    Array sampling: none, first_last, random, representative (default: none)
  -sample-size int           Number of items when sampling (default: 0 = use list-len)

Advanced Compression:
  -null-compression          Track removed null fields in _nulls array
  -type-inference            Convert uniform arrays to schema+data format
  -bool-compression          Convert booleans to bit flags
  -timestamp-compression     Convert ISO timestamps to unix timestamps
  -string-pooling            Deduplicate repeated strings using string pool
  -string-pool-min int       Minimum occurrences for string pooling (default: 2)
  -number-delta              Use delta encoding for sequential numbers
  -number-delta-threshold int Minimum array size for delta encoding (default: 5)
  -enum-detection            Convert repeated categorical values to enums
  -enum-max-values int       Maximum unique values to consider as enum (default: 10)

Examples:
  # Process file with medium profile
  slimjson -profile medium data.json

  # Run as daemon on port 3000
  slimjson -d -port 3000

  # Use custom config file
  slimjson -c /path/to/config.slimjson -profile my-profile data.json

  # Process stdin with custom settings
  cat data.json | slimjson -depth 3 -list-len 5 -pretty

Daemon API:
  POST /slim                 Compress JSON (use ?profile=name for profiles)
  GET  /health               Health check
  GET  /profiles             List available profiles

For more information: https://github.com/tradik/slimjson
`)
}

// runDaemon starts the HTTP server
func runDaemon(port int, customProfiles map[string]slimjson.Config) {
	// Combine built-in and custom profiles
	allProfiles := slimjson.GetBuiltinProfiles()
	for name, cfg := range customProfiles {
		allProfiles[name] = cfg
	}

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","version":"1.0"}`)
	})

	// List profiles endpoint
	http.HandleFunc("/profiles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		profiles := make(map[string][]string)
		profiles["builtin"] = []string{"light", "medium", "aggressive", "ai-optimized"}
		profiles["custom"] = make([]string, 0)

		for name := range customProfiles {
			profiles["custom"] = append(profiles["custom"], name)
		}

		json.NewEncoder(w).Encode(profiles)
	})

	// Slim endpoint
	http.HandleFunc("/slim", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get profile from query parameter
		profileName := r.URL.Query().Get("profile")

		var cfg slimjson.Config
		if profileName != "" {
			var ok bool
			cfg, ok = allProfiles[strings.ToLower(profileName)]
			if !ok {
				http.Error(w, fmt.Sprintf("Unknown profile: %s", profileName), http.StatusBadRequest)
				return
			}
		} else {
			// Default config
			cfg = slimjson.Config{
				MaxDepth:      5,
				MaxListLength: 10,
				StripEmpty:    true,
			}
		}

		// Parse JSON from request body
		var data interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Process
		slimmer := slimjson.New(cfg)
		result := slimmer.Slim(data)

		// Return result
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode result: %v", err), http.StatusInternalServerError)
			return
		}
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("SlimJSON daemon starting on http://localhost%s", addr)
	log.Printf("Endpoints:")
	log.Printf("  POST /slim?profile=<name>  - Compress JSON")
	log.Printf("  GET  /health               - Health check")
	log.Printf("  GET  /profiles             - List profiles")
	log.Printf("Available profiles: %d built-in, %d custom", len(slimjson.GetBuiltinProfiles()), len(customProfiles))

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	var (
		daemon                   bool
		configFile               string
		port                     int
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

	flag.BoolVar(&daemon, "d", false, "Run as HTTP daemon")
	flag.BoolVar(&daemon, "daemon", false, "Run as HTTP daemon")
	flag.StringVar(&configFile, "c", "", "Path to custom config file")
	flag.StringVar(&configFile, "config", "", "Path to custom config file")
	flag.IntVar(&port, "port", 8080, "Port for daemon mode")
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

	// Custom usage message
	flag.Usage = printUsage

	flag.Parse()

	// Show help if no arguments and not daemon mode
	if !daemon && len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}

	// Load custom profiles from config file
	var customProfiles map[string]slimjson.Config
	var err error

	if configFile != "" {
		// Priority: use specified config file
		customProfiles, err = slimjson.ParseConfigFile(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load config file %s: %v\n", configFile, err)
			os.Exit(1)
		}
	} else {
		// Fallback: search for .slimjson in current dir and home dir
		customProfiles, err = slimjson.LoadConfigFile()
		if err != nil {
			// Not an error if file doesn't exist
			customProfiles = make(map[string]slimjson.Config)
		}
	}

	// Run daemon mode if requested
	if daemon {
		runDaemon(port, customProfiles)
		return
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
