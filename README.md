# GitHub Activity CLI

A command-line tool to fetch and display GitHub user activity using the GitHub Events API. This project is part of the [roadmap.sh backend projects](https://roadmap.sh/projects/github-user-activity).

## Features

- Fetch recent GitHub activity for any user
- Filter events by type (push, pull requests, issues, etc.)
- Pagination support for browsing through activity history
- Configurable number of events per page
- Clean, formatted output of user activities

## Installation

### Prerequisites

- Go 1.24.5 or later

### Build from source

1. Clone the repository:
```bash
git clone https://github.com/dmitriy-zverev/github-activity.git
cd github-activity
```

2. Build the application:
```bash
go build -o github-activity
```

3. (Optional) Install globally:
```bash
go install
```

## Usage

### Basic Usage

```bash
./github-activity <username>
```

Example:
```bash
./github-activity octocat
```

### Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-f <event_type>` | Filter events by type | No filter (all events) |
| `-p <page_number>` | Specify page number for pagination | 1 |
| `-n <per_page>` | Number of events per page | 30 |

### Examples

Fetch all recent activity for a user:
```bash
./github-activity dmitriy-zverev
```

Filter only push events:
```bash
./github-activity dmitriy-zverev -f PushEvent
```

Get the second page of results with 10 events per page:
```bash
./github-activity dmitriy-zverev -p 2 -n 10
```

Combine filters and pagination:
```bash
./github-activity dmitriy-zverev -f PullRequestEvent -p 1 -n 5
```

### Supported Event Types

The following GitHub event types can be used with the `-f` filter option:

- `PushEvent` - Code pushes to repositories
- `PullRequestEvent` - Pull request activities
- `CreateEvent` - Repository/branch/tag creation
- `WatchEvent` - Repository starring
- `DeleteEvent` - Repository/branch/tag deletion
- `ForkEvent` - Repository forking
- `IssuesEvent` - Issue activities
- `IssueCommentEvent` - Comments on issues
- `PublicEvent` - Repository made public
- `MemberEvent` - Collaborator activities
- `ReleaseEvent` - Release activities

## Project Structure

```
.
├── main.go              # Main application entry point
├── api_handler.go       # GitHub API interaction logic
├── filter_events.go     # Event filtering functionality
├── models.go            # Data structures for GitHub events
├── printer.go           # Output formatting and display
├── help.go              # Help text and usage information
├── consts.go            # Application constants
├── go.mod               # Go module definition
└── README.md            # This file
```

## API Reference

This tool uses the [GitHub Events API](https://docs.github.com/en/rest/activity/events) to fetch user activity data. No authentication is required for public user data.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is free to use and upgrade.

## Acknowledgments

- This project is part of the [roadmap.sh](https://roadmap.sh) backend development roadmap
- Built using the GitHub REST API
