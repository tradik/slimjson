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

func main() {
	var (
		maxDepth        int
		maxListLength   int
		maxStringLength int
		stripEmpty      bool
		blockList       string
		pretty          bool
	)

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

	cfg := slimjson.Config{
		MaxDepth:        maxDepth,
		MaxListLength:   maxListLength,
		MaxStringLength: maxStringLength,
		StripEmpty:      stripEmpty,
	}
	if blockList != "" {
		cfg.BlockList = strings.Split(blockList, ",")
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
