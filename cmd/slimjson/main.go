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

// getProfile returns a predefined configuration profile
func getProfile(name string) slimjson.Config {
	profiles := map[string]slimjson.Config{
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

	cfg, ok := profiles[strings.ToLower(name)]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown profile: %s\n", name)
		fmt.Fprintf(os.Stderr, "Available profiles: light, medium, aggressive, ai-optimized\n")
		os.Exit(1)
	}
	return cfg
}

func main() {
	var (
		profile         string
		maxDepth        int
		maxListLength   int
		maxStringLength int
		stripEmpty      bool
		blockList       string
		pretty          bool
	)

	flag.StringVar(&profile, "profile", "", "Use predefined profile: light, medium, aggressive, ai-optimized")
	flag.IntVar(&maxDepth, "depth", 5, "Maximum nesting depth (0 for unlimited)")
	flag.IntVar(&maxListLength, "list-len", 10, "Maximum list length (0 for unlimited)")
	flag.IntVar(&maxStringLength, "string-len", 0, "Maximum string length in characters/runes (0 for unlimited)")
	flag.BoolVar(&stripEmpty, "strip-empty", true, "Remove nulls, empty strings, empty arrays/objects")
	flag.StringVar(&blockList, "block", "", "Comma-separated list of field names to remove")
	flag.BoolVar(&pretty, "pretty", false, "Pretty print output")
	flag.Parse()

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
		cfg = getProfile(profile)
	} else {
		// Use custom parameters
		cfg = slimjson.Config{
			MaxDepth:        maxDepth,
			MaxListLength:   maxListLength,
			MaxStringLength: maxStringLength,
			StripEmpty:      stripEmpty,
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
