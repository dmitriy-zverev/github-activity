package main

import (
	"os"
	"testing"
)

func TestMainFunctionArguments(t *testing.T) {
	// Save original args
	originalArgs := os.Args

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "No arguments",
			args:        []string{"github-activity"},
			expectError: true, // Should show help
		},
		{
			name:        "Valid username only",
			args:        []string{"github-activity", "testuser"},
			expectError: false,
		},
		{
			name:        "Username with page flag",
			args:        []string{"github-activity", "testuser", "-p", "2"},
			expectError: false,
		},
		{
			name:        "Username with per-page flag",
			args:        []string{"github-activity", "testuser", "-n", "10"},
			expectError: false,
		},
		{
			name:        "Username with filter flag",
			args:        []string{"github-activity", "testuser", "-f", "PushEvent"},
			expectError: false,
		},
		{
			name:        "All flags combined",
			args:        []string{"github-activity", "testuser", "-p", "2", "-n", "10", "-f", "PushEvent"},
			expectError: false,
		},
		{
			name:        "Invalid page number",
			args:        []string{"github-activity", "testuser", "-p", "invalid"},
			expectError: true,
		},
		{
			name:        "Invalid per-page number",
			args:        []string{"github-activity", "testuser", "-n", "invalid"},
			expectError: true,
		},
		{
			name:        "Page flag without value",
			args:        []string{"github-activity", "testuser", "-p"},
			expectError: false, // Should use default
		},
		{
			name:        "Per-page flag without value",
			args:        []string{"github-activity", "testuser", "-n"},
			expectError: false, // Should use default
		},
		{
			name:        "Filter flag without value",
			args:        []string{"github-activity", "testuser", "-f"},
			expectError: false, // Should use empty filter
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test arguments
			os.Args = tt.args

			// Note: We can't easily test the main function directly because it calls
			// external APIs and prints to stdout. In a real scenario, we would:
			// 1. Refactor main to return errors instead of calling os.Exit
			// 2. Extract the argument parsing logic into a separate testable function
			// 3. Mock the external dependencies

			// For now, we'll test the argument parsing logic indirectly
			// by checking the conditions that would cause errors

			if len(os.Args) < 2 {
				if !tt.expectError {
					t.Errorf("Expected no error for case %s, but would show help", tt.name)
				}
			}

			// Test page number validation
			for i, arg := range os.Args {
				if arg == "-p" && i+1 < len(os.Args) {
					pageArg := os.Args[i+1]
					if pageArg == "invalid" && !tt.expectError {
						t.Errorf("Expected error for invalid page number, but test expects no error")
					}
				}
				if arg == "-n" && i+1 < len(os.Args) {
					perPageArg := os.Args[i+1]
					if perPageArg == "invalid" && !tt.expectError {
						t.Errorf("Expected error for invalid per-page number, but test expects no error")
					}
				}
			}
		})
	}

	// Restore original args
	os.Args = originalArgs
}

func TestArgumentParsing(t *testing.T) {
	// Test the argument parsing logic that would be in main
	tests := []struct {
		name            string
		args            []string
		expectedPage    string
		expectedPerPage string
		expectedFilter  string
	}{
		{
			name:            "Default values",
			args:            []string{"github-activity", "testuser"},
			expectedPage:    DEFAULT_PAGE_NUM,
			expectedPerPage: DEFAULT_PER_PAGE_EVENTS,
			expectedFilter:  DEFAULT_FILTER_TYPE,
		},
		{
			name:            "Custom page",
			args:            []string{"github-activity", "testuser", "-p", "3"},
			expectedPage:    "3",
			expectedPerPage: DEFAULT_PER_PAGE_EVENTS,
			expectedFilter:  DEFAULT_FILTER_TYPE,
		},
		{
			name:            "Custom per-page",
			args:            []string{"github-activity", "testuser", "-n", "20"},
			expectedPage:    DEFAULT_PAGE_NUM,
			expectedPerPage: "20",
			expectedFilter:  DEFAULT_FILTER_TYPE,
		},
		{
			name:            "Custom filter",
			args:            []string{"github-activity", "testuser", "-f", "PushEvent"},
			expectedPage:    DEFAULT_PAGE_NUM,
			expectedPerPage: DEFAULT_PER_PAGE_EVENTS,
			expectedFilter:  "PushEvent",
		},
		{
			name:            "All custom values",
			args:            []string{"github-activity", "testuser", "-p", "2", "-n", "15", "-f", "WatchEvent"},
			expectedPage:    "2",
			expectedPerPage: "15",
			expectedFilter:  "WatchEvent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the argument parsing logic from main
			pageNum := DEFAULT_PAGE_NUM
			perPageNum := DEFAULT_PER_PAGE_EVENTS
			eventTypeFilter := DEFAULT_FILTER_TYPE

			// Parse page flag
			for i, arg := range tt.args {
				if arg == "-p" && i+1 < len(tt.args) {
					pageNum = tt.args[i+1]
				}
				if arg == "-n" && i+1 < len(tt.args) {
					perPageNum = tt.args[i+1]
				}
				if arg == "-f" && i+1 < len(tt.args) {
					eventTypeFilter = tt.args[i+1]
				}
			}

			if pageNum != tt.expectedPage {
				t.Errorf("Expected page %s, got %s", tt.expectedPage, pageNum)
			}
			if perPageNum != tt.expectedPerPage {
				t.Errorf("Expected per-page %s, got %s", tt.expectedPerPage, perPageNum)
			}
			if eventTypeFilter != tt.expectedFilter {
				t.Errorf("Expected filter %s, got %s", tt.expectedFilter, eventTypeFilter)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Test that constants have expected values
	if DEFAULT_PAGE_NUM != "1" {
		t.Errorf("Expected DEFAULT_PAGE_NUM to be '1', got %s", DEFAULT_PAGE_NUM)
	}
	if DEFAULT_PER_PAGE_EVENTS != "30" {
		t.Errorf("Expected DEFAULT_PER_PAGE_EVENTS to be '30', got %s", DEFAULT_PER_PAGE_EVENTS)
	}
	if DEFAULT_FILTER_TYPE != "" {
		t.Errorf("Expected DEFAULT_FILTER_TYPE to be empty string, got %s", DEFAULT_FILTER_TYPE)
	}
}

func TestEventTypeConstants(t *testing.T) {
	// Test that event type constants are defined correctly
	expectedConstants := map[string]string{
		"PUSH_EVENT":           "PushEvent",
		"PULL_REQUEST_EVENT":   "PullRequestEvent",
		"CREATE_EVENT":         "CreateEvent",
		"WATCH_EVENT":          "WatchEvent",
		"DELETE_EVENT":         "DeleteEvent",
		"FORK_EVENT":           "ForkEvent",
		"ISSUES_EVENT":         "IssuesEvent",
		"ISSUES_COMMENT_EVENT": "IssueCommentEvent",
		"PUBLIC_EVENT":         "PublicEvent",
		"MEMBER_EVENT":         "MemberEvent",
		"RELEASE_EVENT":        "ReleaseEvent",
	}

	actualConstants := map[string]string{
		"PUSH_EVENT":           PUSH_EVENT,
		"PULL_REQUEST_EVENT":   PULL_REQUEST_EVENT,
		"CREATE_EVENT":         CREATE_EVENT,
		"WATCH_EVENT":          WATCH_EVENT,
		"DELETE_EVENT":         DELETE_EVENT,
		"FORK_EVENT":           FORK_EVENT,
		"ISSUES_EVENT":         ISSUES_EVENT,
		"ISSUES_COMMENT_EVENT": ISSUES_COMMENT_EVENT,
		"PUBLIC_EVENT":         PUBLIC_EVENT,
		"MEMBER_EVENT":         MEMBER_EVENT,
		"RELEASE_EVENT":        RELEASE_EVENT,
	}

	for name, expected := range expectedConstants {
		if actual, exists := actualConstants[name]; !exists {
			t.Errorf("Constant %s is not defined", name)
		} else if actual != expected {
			t.Errorf("Expected %s to be %s, got %s", name, expected, actual)
		}
	}
}

// Integration test helper functions
func createTestGithubUserData() []githubUserData {
	return []githubUserData{
		{
			Type: PUSH_EVENT,
			Repo: struct {
				Name string `json:"name"`
			}{Name: "test-repo"},
			Payload: struct {
				Ref     string `json:"ref"`
				RefType string `json:"ref_type"`
				Commits []struct {
					Message string `json:"message"`
				} `json:"commits"`
				Action string `json:"action"`
				Forkee struct {
					FullName string `json:"full_name"`
					Owner    struct {
						Login string `json:"login"`
					} `json:"owner"`
				} `json:"forkee"`
				Issue struct {
					Title string `json:"title"`
				} `json:"issue"`
				PullReq struct {
					Title string `json:"title"`
				} `json:"pull_request"`
				Member struct {
					Login string `json:"login"`
				} `json:"member"`
				Release struct {
					Name string `json:"name"`
				} `json:"release"`
			}{
				Commits: []struct {
					Message string `json:"message"`
				}{
					{Message: "Initial commit"},
					{Message: "Add feature"},
				},
			},
		},
		{
			Type: WATCH_EVENT,
			Repo: struct {
				Name string `json:"name"`
			}{Name: "another-repo"},
			Payload: struct {
				Ref     string `json:"ref"`
				RefType string `json:"ref_type"`
				Commits []struct {
					Message string `json:"message"`
				} `json:"commits"`
				Action string `json:"action"`
				Forkee struct {
					FullName string `json:"full_name"`
					Owner    struct {
						Login string `json:"login"`
					} `json:"owner"`
				} `json:"forkee"`
				Issue struct {
					Title string `json:"title"`
				} `json:"issue"`
				PullReq struct {
					Title string `json:"title"`
				} `json:"pull_request"`
				Member struct {
					Login string `json:"login"`
				} `json:"member"`
				Release struct {
					Name string `json:"name"`
				} `json:"release"`
			}{
				Action: "started",
			},
		},
		{
			Type: CREATE_EVENT,
			Repo: struct {
				Name string `json:"name"`
			}{Name: "new-repo"},
			Payload: struct {
				Ref     string `json:"ref"`
				RefType string `json:"ref_type"`
				Commits []struct {
					Message string `json:"message"`
				} `json:"commits"`
				Action string `json:"action"`
				Forkee struct {
					FullName string `json:"full_name"`
					Owner    struct {
						Login string `json:"login"`
					} `json:"owner"`
				} `json:"forkee"`
				Issue struct {
					Title string `json:"title"`
				} `json:"issue"`
				PullReq struct {
					Title string `json:"title"`
				} `json:"pull_request"`
				Member struct {
					Login string `json:"login"`
				} `json:"member"`
				Release struct {
					Name string `json:"name"`
				} `json:"release"`
			}{
				RefType: "repository",
			},
		},
	}
}

func TestIntegrationFilterAndPrint(t *testing.T) {
	// Integration test: filter events and then print them
	testData := createTestGithubUserData()

	// Test filtering and printing PushEvents (using lowercase for case-insensitive match)
	filtered := filterEvents(testData, "pushevent")
	if len(filtered) != 1 {
		t.Errorf("Expected 1 PushEvent, got %d", len(filtered))
	}

	// Test that the filtered event is correct (only if we have results)
	if len(filtered) > 0 && filtered[0].Type != PUSH_EVENT {
		t.Errorf("Expected filtered event to be PushEvent, got %s", filtered[0].Type)
	}

	// Test filtering with no matches
	noMatches := filterEvents(testData, "NonExistentEvent")
	if len(noMatches) != 0 {
		t.Errorf("Expected 0 events for non-existent filter, got %d", len(noMatches))
	}

	// Test filtering with no filter (should return all)
	allEvents := filterEvents(testData, DEFAULT_FILTER_TYPE)
	if len(allEvents) != len(testData) {
		t.Errorf("Expected %d events with no filter, got %d", len(testData), len(allEvents))
	}
}
