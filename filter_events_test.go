package main

import (
	"reflect"
	"testing"
)

func TestFilterEvents(t *testing.T) {
	// Sample test data
	testEvents := []githubUserData{
		{Type: "PushEvent", Repo: struct {
			Name string `json:"name"`
		}{Name: "test-repo"}},
		{Type: "PullRequestEvent", Repo: struct {
			Name string `json:"name"`
		}{Name: "test-repo"}},
		{Type: "CreateEvent", Repo: struct {
			Name string `json:"name"`
		}{Name: "test-repo"}},
		{Type: "WatchEvent", Repo: struct {
			Name string `json:"name"`
		}{Name: "test-repo"}},
		{Type: "IssuesEvent", Repo: struct {
			Name string `json:"name"`
		}{Name: "test-repo"}},
	}

	tests := []struct {
		name      string
		events    []githubUserData
		eventType string
		expected  []githubUserData
	}{
		{
			name:      "No filter - returns all events",
			events:    testEvents,
			eventType: DEFAULT_FILTER_TYPE,
			expected:  testEvents,
		},
		{
			name:      "Filter by PushEvent",
			events:    testEvents,
			eventType: "pushevent",
			expected:  []githubUserData{testEvents[0]},
		},
		{
			name:      "Filter by PullRequestEvent",
			events:    testEvents,
			eventType: "pullrequestevent",
			expected:  []githubUserData{testEvents[1]},
		},
		{
			name:      "Filter by partial match (case insensitive)",
			events:    testEvents,
			eventType: "push",
			expected:  []githubUserData{testEvents[0]},
		},
		{
			name:      "Filter by partial match - request",
			events:    testEvents,
			eventType: "request",
			expected:  []githubUserData{testEvents[1]},
		},
		{
			name:      "Filter with no matches",
			events:    testEvents,
			eventType: "nonexistentevent",
			expected:  []githubUserData{},
		},
		{
			name:      "Empty events list",
			events:    []githubUserData{},
			eventType: "pushevent",
			expected:  []githubUserData{},
		},
		{
			name:      "Empty filter on empty events",
			events:    []githubUserData{},
			eventType: DEFAULT_FILTER_TYPE,
			expected:  []githubUserData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterEvents(tt.events, tt.eventType)

			// Handle nil vs empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("filterEvents() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFilterEventsMultipleMatches(t *testing.T) {
	// Test data with multiple events of the same type
	testEvents := []githubUserData{
		{Type: "PushEvent", Repo: struct {
			Name string `json:"name"`
		}{Name: "repo1"}},
		{Type: "PullRequestEvent", Repo: struct {
			Name string `json:"name"`
		}{Name: "repo2"}},
		{Type: "PushEvent", Repo: struct {
			Name string `json:"name"`
		}{Name: "repo3"}},
		{Type: "CreateEvent", Repo: struct {
			Name string `json:"name"`
		}{Name: "repo4"}},
		{Type: "PushEvent", Repo: struct {
			Name string `json:"name"`
		}{Name: "repo5"}},
	}

	result := filterEvents(testEvents, "pushevent")
	expected := []githubUserData{testEvents[0], testEvents[2], testEvents[4]}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("filterEvents() with multiple matches = %v, want %v", result, expected)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 PushEvents, got %d", len(result))
	}
}
