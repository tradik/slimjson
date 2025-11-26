package slimjson

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSlimmer_Slim(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		input    string
		expected string
	}{
		{
			name: "Strip empty fields",
			config: Config{
				StripEmpty: true,
			},
			input:    `{"a": 1, "b": "", "c": null, "d": [], "e": {}}`,
			expected: `{"a": 1}`,
		},
		{
			name: "Max depth",
			config: Config{
				MaxDepth: 2,
			},
			input:    `{"a": {"b": {"c": 1}}}`,
			expected: `{"a": {"b": null}}`,
		},
		{
			name: "Max list length",
			config: Config{
				MaxListLength: 2,
			},
			input:    `{"list": [1, 2, 3, 4]}`,
			expected: `{"list": [1, 2]}`,
		},
		{
			name: "Block list",
			config: Config{
				BlockList: []string{"secret", "password"},
			},
			input:    `{"user": "me", "password": "123", "secret": "shh"}`,
			expected: `{"user": "me"}`,
		},
		{
			name: "Max string length with UTF-8",
			config: Config{
				MaxStringLength: 5,
			},
			input:    `{"text": "Hello World", "emoji": "🎉🎊🎈🎁🎀🎂"}`,
			expected: `{"text": "He...", "emoji": "🎉🎊..."}`,
		},
		{
			name: "Complex combination",
			config: Config{
				MaxDepth:      10,
				MaxListLength: 2,
				StripEmpty:    true,
				BlockList:     []string{"ignore"},
			},
			input: `
			{
				"keep": "yes",
				"ignore": "no",
				"empty": "",
				"deep": {
					"level1": {
						"level2": {
							"level3": "deep enough"
						}
					}
				},
				"list": [
					{"id": 1},
					{"id": 2},
					{"id": 3}
				],
				"emptyList": [null, "", {}] 
			}`,
			expected: `{"keep": "yes", "deep": {"level1": {"level2": {"level3": "deep enough"}}}, "list": [{"id": 1}, {"id": 2}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputData interface{}
			if err := json.Unmarshal([]byte(tt.input), &inputData); err != nil {
				t.Fatalf("Failed to unmarshal input: %v", err)
			}

			slimmer := New(tt.config)
			got := slimmer.Slim(inputData)

			// Marshal got to compare with expected JSON string
			gotBytes, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("Failed to marshal result: %v", err)
			}

			// Normalize expected JSON
			var expectedData interface{}
			if err := json.Unmarshal([]byte(tt.expected), &expectedData); err != nil {
				t.Fatalf("Failed to unmarshal expected: %v", err)
			}
			expectedBytes, _ := json.Marshal(expectedData)

			// Compare as interface{} to avoid ordering issues, or just compare normalized strings
			if !reflect.DeepEqual(got, expectedData) {
				t.Errorf("Slim() = %s, want %s", string(gotBytes), string(expectedBytes))
			}
		})
	}
}
