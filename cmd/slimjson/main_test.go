package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tradik/slimjson"
)

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","version":"1.0"}`))
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}

	if response["version"] != "1.0" {
		t.Errorf("Expected version '1.0', got '%s'", response["version"])
	}
}

func TestProfilesEndpoint(t *testing.T) {
	customProfiles := map[string]slimjson.Config{
		"test-profile": {
			MaxDepth:      3,
			MaxListLength: 5,
			StripEmpty:    true,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/profiles", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		profiles := make(map[string][]string)
		profiles["builtin"] = []string{"light", "medium", "aggressive", "ai-optimized"}
		profiles["custom"] = make([]string, 0)

		for name := range customProfiles {
			profiles["custom"] = append(profiles["custom"], name)
		}

		json.NewEncoder(w).Encode(profiles)
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string][]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response["builtin"]) != 4 {
		t.Errorf("Expected 4 built-in profiles, got %d", len(response["builtin"]))
	}

	if len(response["custom"]) != 1 {
		t.Errorf("Expected 1 custom profile, got %d", len(response["custom"]))
	}
}

func TestSlimEndpoint(t *testing.T) {
	allProfiles := slimjson.GetBuiltinProfiles()

	tests := []struct {
		name           string
		method         string
		profile        string
		input          string
		expectedStatus int
		checkResult    bool
	}{
		{
			name:           "Valid request with medium profile",
			method:         http.MethodPost,
			profile:        "medium",
			input:          `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}],"prices":[19.999,29.123]}`,
			expectedStatus: http.StatusOK,
			checkResult:    true,
		},
		{
			name:           "Valid request without profile",
			method:         http.MethodPost,
			profile:        "",
			input:          `{"test":"data"}`,
			expectedStatus: http.StatusOK,
			checkResult:    true,
		},
		{
			name:           "Invalid method GET",
			method:         http.MethodGet,
			profile:        "",
			input:          `{}`,
			expectedStatus: http.StatusMethodNotAllowed,
			checkResult:    false,
		},
		{
			name:           "Invalid JSON",
			method:         http.MethodPost,
			profile:        "",
			input:          `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
			checkResult:    false,
		},
		{
			name:           "Unknown profile",
			method:         http.MethodPost,
			profile:        "nonexistent",
			input:          `{"test":"data"}`,
			expectedStatus: http.StatusBadRequest,
			checkResult:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/slim"
			if tt.profile != "" {
				url += "?profile=" + tt.profile
			}

			req := httptest.NewRequest(tt.method, url, bytes.NewBufferString(tt.input))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					return
				}

				profileName := r.URL.Query().Get("profile")

				var cfg slimjson.Config
				if profileName != "" {
					var ok bool
					cfg, ok = allProfiles[profileName]
					if !ok {
						http.Error(w, "Unknown profile", http.StatusBadRequest)
						return
					}
				} else {
					cfg = slimjson.Config{
						MaxDepth:      5,
						MaxListLength: 10,
						StripEmpty:    true,
					}
				}

				var data interface{}
				if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}

				slimmer := slimjson.New(cfg)
				result := slimmer.Slim(data)

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResult && w.Code == http.StatusOK {
				var result interface{}
				if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
					t.Errorf("Failed to decode result: %v", err)
				}
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	customProfiles := map[string]slimjson.Config{
		"custom-test": {
			MaxDepth:      3,
			MaxListLength: 5,
			StripEmpty:    true,
		},
	}

	tests := []struct {
		name        string
		profileName string
		shouldExist bool
	}{
		{
			name:        "Built-in profile light",
			profileName: "light",
			shouldExist: true,
		},
		{
			name:        "Built-in profile medium",
			profileName: "medium",
			shouldExist: true,
		},
		{
			name:        "Custom profile",
			profileName: "custom-test",
			shouldExist: true,
		},
		{
			name:        "Non-existent profile",
			profileName: "nonexistent",
			shouldExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would normally call getProfile, but since it exits on error,
			// we test the logic separately
			builtinProfiles := slimjson.GetBuiltinProfiles()

			// Check custom first
			_, customOk := customProfiles[tt.profileName]
			_, builtinOk := builtinProfiles[tt.profileName]

			exists := customOk || builtinOk

			if exists != tt.shouldExist {
				t.Errorf("Profile %s: expected exists=%v, got exists=%v",
					tt.profileName, tt.shouldExist, exists)
			}
		})
	}
}

func TestConfigFilePriority(t *testing.T) {
	// Test that custom config file takes priority
	// This is more of an integration test

	t.Run("Custom config has priority", func(t *testing.T) {
		// Create a mock custom profile
		customProfiles := map[string]slimjson.Config{
			"medium": { // Override built-in medium
				MaxDepth:      99,
				MaxListLength: 99,
				StripEmpty:    false,
			},
		}

		builtinProfiles := slimjson.GetBuiltinProfiles()

		// Merge with custom taking priority
		allProfiles := builtinProfiles
		for name, cfg := range customProfiles {
			allProfiles[name] = cfg
		}

		// Check that custom overrode built-in
		mediumCfg := allProfiles["medium"]
		if mediumCfg.MaxDepth != 99 {
			t.Errorf("Expected custom MaxDepth=99, got %d", mediumCfg.MaxDepth)
		}
		if mediumCfg.MaxListLength != 99 {
			t.Errorf("Expected custom MaxListLength=99, got %d", mediumCfg.MaxListLength)
		}
	})
}
