package main

import (
	"testing"
)

// Test validation with JSONPlaceholder API response format
// Example response from https://jsonplaceholder.typicode.com/posts/1:
// {
//   "userId": 1,
//   "id": 1,
//   "title": "sunt aut facere...",
//   "body": "quia et suscipit..."
// }

func TestValidateResponse_ExistingFields(t *testing.T) {
	// Sample JSONPlaceholder response
	jsonResponse := `{
		"userId": 1,
		"id": 1,
		"title": "Sample title",
		"body": "Sample body text"
	}`

	config := &Config{
		Benchmark: BenchmarkConfig{
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{
						Path: "userId",
						Type: "number",
					},
					{
						Path:  "id",
						Type:  "number",
						Value: float64(1),
					},
					{
						Path: "title",
						Type: "string",
					},
					{
						Path: "body",
						Type: "string",
					},
				},
			},
		},
	}

	engine := &BenchmarkEngine{config: config}
	err := engine.validateResponse([]byte(jsonResponse))
	if err != nil {
		t.Errorf("Validation should pass for existing fields: %v", err)
	}
}

func TestValidateResponse_NonExistentField(t *testing.T) {
	// Sample JSONPlaceholder response
	jsonResponse := `{
		"userId": 1,
		"id": 1,
		"title": "Sample title",
		"body": "Sample body text"
	}`

	config := &Config{
		Benchmark: BenchmarkConfig{
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{
						Path: "nonExistentField",
						Type: "string",
					},
				},
			},
		},
	}

	engine := &BenchmarkEngine{config: config}
	err := engine.validateResponse([]byte(jsonResponse))
	if err == nil {
		t.Error("Validation should fail for non-existent field")
	}
	expectedError := "validation failed for path 'nonExistentField'"
	if err != nil && !contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got: %v", expectedError, err)
	}
}

func TestValidateResponse_WrongType(t *testing.T) {
	jsonResponse := `{
		"userId": "string_instead_of_number",
		"id": 1,
		"title": "Sample title",
		"body": "Sample body text"
	}`

	config := &Config{
		Benchmark: BenchmarkConfig{
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{
						Path: "userId",
						Type: "number",
					},
				},
			},
		},
	}

	engine := &BenchmarkEngine{config: config}
	err := engine.validateResponse([]byte(jsonResponse))
	if err == nil {
		t.Error("Validation should fail for wrong type")
	}
	if err != nil && !contains(err.Error(), "expected number") {
		t.Errorf("Expected type error, got: %v", err)
	}
}

func TestValidateResponse_WrongValue(t *testing.T) {
	jsonResponse := `{
		"userId": 1,
		"id": 99,
		"title": "Sample title",
		"body": "Sample body text"
	}`

	config := &Config{
		Benchmark: BenchmarkConfig{
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{
						Path:  "id",
						Type:  "number",
						Value: float64(1),
					},
				},
			},
		},
	}

	engine := &BenchmarkEngine{config: config}
	err := engine.validateResponse([]byte(jsonResponse))
	if err == nil {
		t.Error("Validation should fail for wrong value")
	}
	if err != nil && !contains(err.Error(), "expected value 1") {
		t.Errorf("Expected value mismatch error, got: %v", err)
	}
}

func TestValidateResponse_NestedFields(t *testing.T) {
	// Test with nested JSON structure (like GitHub API)
	jsonResponse := `{
		"user": {
			"id": 123,
			"name": "John Doe",
			"profile": {
				"verified": true,
				"followers": 42
			}
		},
		"status": "success"
	}`

	config := &Config{
		Benchmark: BenchmarkConfig{
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{
						Path: "user.id",
						Type: "number",
					},
					{
						Path: "user.name",
						Type: "string",
					},
					{
						Path:  "user.profile.verified",
						Type:  "boolean",
						Value: true,
					},
					{
						Path: "user.profile.followers",
						Type: "number",
					},
					{
						Path:  "status",
						Type:  "string",
						Value: "success",
					},
				},
			},
		},
	}

	engine := &BenchmarkEngine{config: config}
	err := engine.validateResponse([]byte(jsonResponse))
	if err != nil {
		t.Errorf("Validation should pass for nested fields: %v", err)
	}
}

func TestValidateResponse_ArrayType(t *testing.T) {
	// Test array response (like JSONPlaceholder /posts)
	jsonResponse := `{
		"data": [
			{"id": 1, "title": "Post 1"},
			{"id": 2, "title": "Post 2"}
		],
		"total": 2
	}`

	config := &Config{
		Benchmark: BenchmarkConfig{
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{
						Path: "data",
						Type: "array",
					},
					{
						Path: "total",
						Type: "number",
					},
				},
			},
		},
	}

	engine := &BenchmarkEngine{config: config}
	err := engine.validateResponse([]byte(jsonResponse))
	if err != nil {
		t.Errorf("Validation should pass for array type: %v", err)
	}
}

func TestValidateResponse_DisabledValidation(t *testing.T) {
	jsonResponse := `{"invalid": "json"}`

	config := &Config{
		Benchmark: BenchmarkConfig{
			ResponseValidation: &ResponseValidation{
				Enabled: false,
				Checks: []ValidationCheck{
					{
						Path: "nonExistent",
						Type: "string",
					},
				},
			},
		},
	}

	engine := &BenchmarkEngine{config: config}
	err := engine.validateResponse([]byte(jsonResponse))
	if err != nil {
		t.Error("Validation should not run when disabled")
	}
}

func TestValidateResponse_InvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json}`

	config := &Config{
		Benchmark: BenchmarkConfig{
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{
						Path: "any",
						Type: "string",
					},
				},
			},
		},
	}

	engine := &BenchmarkEngine{config: config}
	err := engine.validateResponse([]byte(invalidJSON))
	if err == nil {
		t.Error("Validation should fail for invalid JSON")
	}
	if err != nil && !contains(err.Error(), "failed to parse JSON") {
		t.Errorf("Expected JSON parse error, got: %v", err)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

// Test with actual API configurations
func TestValidateResponse_GitHubAPI(t *testing.T) {
	// Simulating GitHub API user response
	jsonResponse := `{
		"login": "octocat",
		"id": 1,
		"avatar_url": "https://github.com/images/error/octocat_happy.gif",
		"type": "User",
		"site_admin": false,
		"name": "monalisa octocat",
		"company": "GitHub",
		"public_repos": 2,
		"public_gists": 1,
		"followers": 20,
		"following": 0
	}`

	config := &Config{
		Benchmark: BenchmarkConfig{
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{
						Path: "login",
						Type: "string",
					},
					{
						Path: "id",
						Type: "number",
					},
					{
						Path:  "type",
						Type:  "string",
						Value: "User",
					},
					{
						Path:  "site_admin",
						Type:  "boolean",
						Value: false,
					},
					{
						Path: "public_repos",
						Type: "number",
					},
				},
			},
		},
	}

	engine := &BenchmarkEngine{config: config}
	err := engine.validateResponse([]byte(jsonResponse))
	if err != nil {
		t.Errorf("GitHub API validation should pass: %v", err)
	}
}

// Benchmark tests
func BenchmarkValidateResponse_SimpleFields(b *testing.B) {
	jsonResponse := []byte(`{
		"userId": 1,
		"id": 1,
		"title": "Sample title",
		"body": "Sample body text"
	}`)

	config := &Config{
		Benchmark: BenchmarkConfig{
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{Path: "userId", Type: "number"},
					{Path: "id", Type: "number"},
					{Path: "title", Type: "string"},
					{Path: "body", Type: "string"},
				},
			},
		},
	}

	engine := &BenchmarkEngine{config: config}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.validateResponse(jsonResponse)
	}
}

func BenchmarkValidateResponse_NestedFields(b *testing.B) {
	jsonResponse := []byte(`{
		"user": {
			"profile": {
				"details": {
					"verified": true,
					"score": 100
				}
			}
		}
	}`)

	config := &Config{
		Benchmark: BenchmarkConfig{
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{Path: "user.profile.details.verified", Type: "boolean"},
					{Path: "user.profile.details.score", Type: "number"},
				},
			},
		},
	}

	engine := &BenchmarkEngine{config: config}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.validateResponse(jsonResponse)
	}
}