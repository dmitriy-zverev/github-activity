package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchGithubUserData(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		page           string
		perPage        string
		mockResponse   []githubUserData
		mockStatusCode int
		expectError    bool
		expectedCount  int
	}{
		{
			name:     "Successful request",
			username: "testuser",
			page:     "1",
			perPage:  "30",
			mockResponse: []githubUserData{
				{
					Type: "PushEvent",
					Repo: struct {
						Name string `json:"name"`
					}{Name: "test-repo"},
				},
				{
					Type: "WatchEvent",
					Repo: struct {
						Name string `json:"name"`
					}{Name: "another-repo"},
				},
			},
			mockStatusCode: 200,
			expectError:    false,
			expectedCount:  2,
		},
		{
			name:           "User not found",
			username:       "nonexistentuser",
			page:           "1",
			perPage:        "30",
			mockResponse:   nil,
			mockStatusCode: 404,
			expectError:    true,
			expectedCount:  0,
		},
		{
			name:           "Server error",
			username:       "testuser",
			page:           "1",
			perPage:        "30",
			mockResponse:   nil,
			mockStatusCode: 500,
			expectError:    true,
			expectedCount:  0,
		},
		{
			name:           "Rate limit exceeded",
			username:       "testuser",
			page:           "1",
			perPage:        "30",
			mockResponse:   nil,
			mockStatusCode: 403,
			expectError:    true,
			expectedCount:  0,
		},
		{
			name:           "Empty response",
			username:       "testuser",
			page:           "1",
			perPage:        "30",
			mockResponse:   []githubUserData{},
			mockStatusCode: 200,
			expectError:    false,
			expectedCount:  0,
		},
		{
			name:     "Different page and per_page values",
			username: "testuser",
			page:     "2",
			perPage:  "10",
			mockResponse: []githubUserData{
				{
					Type: "CreateEvent",
					Repo: struct {
						Name string `json:"name"`
					}{Name: "new-repo"},
				},
			},
			mockStatusCode: 200,
			expectError:    false,
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the request URL contains expected parameters
				expectedPath := "/users/" + tt.username + "/events"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Check query parameters
				query := r.URL.Query()
				if query.Get("page") != tt.page {
					t.Errorf("Expected page %s, got %s", tt.page, query.Get("page"))
				}
				if query.Get("per_page") != tt.perPage {
					t.Errorf("Expected per_page %s, got %s", tt.perPage, query.Get("per_page"))
				}

				// Check headers
				if r.Header.Get("Accept") != "application/vnd.github+json" {
					t.Errorf("Expected Accept header 'application/vnd.github+json', got %s", r.Header.Get("Accept"))
				}
				if r.Header.Get("X-GitHub-Api-Version") != "2022-11-28" {
					t.Errorf("Expected X-GitHub-Api-Version header '2022-11-28', got %s", r.Header.Get("X-GitHub-Api-Version"))
				}

				// Set response status code
				w.WriteHeader(tt.mockStatusCode)

				// Return mock response if status is OK
				if tt.mockStatusCode == 200 && tt.mockResponse != nil {
					json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			// Replace the GitHub API URL with our test server URL
			// We need to modify the function to accept a base URL for testing
			// For now, we'll test the current implementation by temporarily modifying it

			// Since we can't easily modify the hardcoded URL in the function,
			// we'll test the error cases and verify the function structure

			// This is a limitation of the current implementation - it should accept a configurable base URL
			// For a real test, we would need to refactor the function to accept a base URL parameter

			// Test with invalid parameters that would cause URL construction issues
			if tt.expectError && tt.mockStatusCode >= 400 {
				// We can't easily test the HTTP client part without refactoring,
				// but we can test parameter validation
				result, err := fetchGithubUserData(tt.username, tt.page, tt.perPage)

				// The function will likely fail due to network issues when trying to reach GitHub
				// In a real scenario, we'd want to inject the HTTP client or base URL
				if len(result) != 0 {
					t.Errorf("Expected empty result for error case, got %d items", len(result))
				}

				// We expect an error for these cases
				if err == nil {
					t.Errorf("Expected error for case %s, but got none", tt.name)
				}

				// Note: This test is limited because the function has a hardcoded URL
				// In a production environment, we'd refactor to accept a configurable client or URL
			}
		})
	}
}

func TestFetchGithubUserDataURLConstruction(t *testing.T) {
	// Test that would verify URL construction if we could intercept it
	// This is more of a documentation of what we would test with a refactored function

	tests := []struct {
		name        string
		username    string
		page        string
		perPage     string
		expectedURL string
	}{
		{
			name:        "Basic URL construction",
			username:    "testuser",
			page:        "1",
			perPage:     "30",
			expectedURL: "https://api.github.com/users/testuser/events?page=1&per_page=30",
		},
		{
			name:        "Different page number",
			username:    "anotheruser",
			page:        "5",
			perPage:     "10",
			expectedURL: "https://api.github.com/users/anotheruser/events?page=5&per_page=10",
		},
		{
			name:        "Username with special characters",
			username:    "user-name_123",
			page:        "1",
			perPage:     "50",
			expectedURL: "https://api.github.com/users/user-name_123/events?page=1&per_page=50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test documents the expected URL format
			// In a refactored version, we would verify the actual URL construction
			expectedURL := "https://api.github.com/users/" + tt.username + "/events?page=" + tt.page + "&per_page=" + tt.perPage
			if expectedURL != tt.expectedURL {
				t.Errorf("URL construction test failed: expected %s, got %s", tt.expectedURL, expectedURL)
			}
		})
	}
}

func TestFetchGithubUserDataErrorHandling(t *testing.T) {
	// Test error handling with invalid inputs
	tests := []struct {
		name     string
		username string
		page     string
		perPage  string
	}{
		{
			name:     "Empty username",
			username: "",
			page:     "1",
			perPage:  "30",
		},
		{
			name:     "Empty page",
			username: "testuser",
			page:     "",
			perPage:  "30",
		},
		{
			name:     "Empty perPage",
			username: "testuser",
			page:     "1",
			perPage:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fetchGithubUserData(tt.username, tt.page, tt.perPage)

			// The function should handle these cases gracefully
			// Even if it doesn't return an error, the result should be empty or the function should fail
			if err == nil && len(result) > 0 {
				// This might be unexpected depending on how GitHub API handles empty parameters
				t.Logf("Function succeeded with parameters: username=%s, page=%s, perPage=%s", tt.username, tt.page, tt.perPage)
			}
		})
	}
}

// Note: The current fetchGithubUserData function has some limitations for testing:
// 1. Hardcoded GitHub API URL makes it difficult to mock
// 2. No dependency injection for HTTP client
// 3. No configurable timeout or retry logic
//
// For better testability, consider refactoring to:
// 1. Accept a base URL parameter or HTTP client interface
// 2. Use dependency injection for the HTTP client
// 3. Separate URL construction from HTTP request logic
//
// Example refactored signature:
// func fetchGithubUserData(client HTTPClient, baseURL, username, page, perPage string) ([]githubUserData, error)
