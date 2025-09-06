package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrinter(t *testing.T) {
	tests := []struct {
		name           string
		activities     []githubUserData
		expectedError  bool
		expectedOutput string
	}{
		{
			name:          "Empty activities list",
			activities:    []githubUserData{},
			expectedError: true,
		},
		{
			name: "Single activity",
			activities: []githubUserData{
				{
					Type: "PushEvent",
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
			},
			expectedError:  false,
			expectedOutput: "  - Pushed 2 commits to test-repo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := printer(tt.activities)

			// Restore stdout
			w.Close()
			os.Stdout = old

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectedError && output != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

func TestActivityString(t *testing.T) {
	tests := []struct {
		name     string
		activity githubUserData
		expected string
		hasError bool
	}{
		{
			name: "PushEvent",
			activity: githubUserData{
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
						{Message: "commit1"},
						{Message: "commit2"},
					},
				},
			},
			expected: "Pushed 2 commits to test-repo",
			hasError: false,
		},
		{
			name: "CreateEvent - repository",
			activity: githubUserData{
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
			expected: "Created repository at new-repo",
			hasError: false,
		},
		{
			name: "CreateEvent - branch",
			activity: githubUserData{
				Type: CREATE_EVENT,
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
					Ref:     "feature-branch",
					RefType: "branch",
				},
			},
			expected: "Created branch 'feature-branch' at test-repo",
			hasError: false,
		},
		{
			name: "WatchEvent - started",
			activity: githubUserData{
				Type: WATCH_EVENT,
				Repo: struct {
					Name string `json:"name"`
				}{Name: "watched-repo"},
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
			expected: "Started watching watched-repo",
			hasError: false,
		},
		{
			name: "WatchEvent - ended",
			activity: githubUserData{
				Type: WATCH_EVENT,
				Repo: struct {
					Name string `json:"name"`
				}{Name: "watched-repo"},
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
					Action: "stopped",
				},
			},
			expected: "Ended watching watched-repo",
			hasError: false,
		},
		{
			name: "DeleteEvent",
			activity: githubUserData{
				Type: DELETE_EVENT,
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
					Ref:     "old-branch",
					RefType: "branch",
				},
			},
			expected: "Deleted branch 'old-branch' at test-repo",
			hasError: false,
		},
		{
			name: "ForkEvent",
			activity: githubUserData{
				Type: FORK_EVENT,
				Repo: struct {
					Name string `json:"name"`
				}{Name: "original-repo"},
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
					Forkee: struct {
						FullName string `json:"full_name"`
						Owner    struct {
							Login string `json:"login"`
						} `json:"owner"`
					}{
						FullName: "user/forked-repo",
					},
				},
			},
			expected: "Forked original-repo to user/forked-repo",
			hasError: false,
		},
		{
			name: "IssuesEvent",
			activity: githubUserData{
				Type: ISSUES_EVENT,
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
					Action: "opened",
					Issue: struct {
						Title string `json:"title"`
					}{
						Title: "Bug report",
					},
				},
			},
			expected: "Issue 'Bug report' opened at test-repo",
			hasError: false,
		},
		{
			name: "IssueCommentEvent",
			activity: githubUserData{
				Type: ISSUES_COMMENT_EVENT,
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
					Issue: struct {
						Title string `json:"title"`
					}{
						Title: "Bug report",
					},
				},
			},
			expected: "Commented at issue 'Bug report' at test-repo",
			hasError: false,
		},
		{
			name: "PullRequestEvent",
			activity: githubUserData{
				Type: PULL_REQUEST_EVENT,
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
					Action: "opened",
					PullReq: struct {
						Title string `json:"title"`
					}{
						Title: "Add new feature",
					},
				},
			},
			expected: "Pull request 'Add new feature' opened at test-repo",
			hasError: false,
		},
		{
			name: "PublicEvent",
			activity: githubUserData{
				Type: PUBLIC_EVENT,
				Repo: struct {
					Name string `json:"name"`
				}{Name: "test-repo"},
			},
			expected: "Repo test-repo is now public",
			hasError: false,
		},
		{
			name: "MemberEvent",
			activity: githubUserData{
				Type: MEMBER_EVENT,
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
					Member: struct {
						Login string `json:"login"`
					}{
						Login: "newuser",
					},
				},
			},
			expected: "Added newuser to test-repo",
			hasError: false,
		},
		{
			name: "ReleaseEvent",
			activity: githubUserData{
				Type: RELEASE_EVENT,
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
					Release: struct {
						Name string `json:"name"`
					}{
						Name: "v1.0.0",
					},
				},
			},
			expected: "Released 'v1.0.0' at test-repo",
			hasError: false,
		},
		{
			name: "Unknown event type",
			activity: githubUserData{
				Type: "UnknownEvent",
				Repo: struct {
					Name string `json:"name"`
				}{Name: "test-repo"},
			},
			expected: "UnknownEvent to test-repo",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := activityString(tt.activity)

			if tt.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestPrinterMultipleActivities(t *testing.T) {
	activities := []githubUserData{
		{
			Type: PUSH_EVENT,
			Repo: struct {
				Name string `json:"name"`
			}{Name: "repo1"},
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
				}{{Message: "commit1"}},
			},
		},
		{
			Type: WATCH_EVENT,
			Repo: struct {
				Name string `json:"name"`
			}{Name: "repo2"},
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
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := printer(activities)

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedLines := []string{
		"  - Pushed 1 commits to repo1",
		"  - Started watching repo2",
	}

	for _, expectedLine := range expectedLines {
		if !strings.Contains(output, expectedLine) {
			t.Errorf("Expected output to contain %q, got %q", expectedLine, output)
		}
	}
}
