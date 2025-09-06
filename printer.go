package main

import (
	"fmt"
)

func printer(userActivities []githubUserData) error {
	for _, activity := range userActivities {
		userActivityString, err := activityString(activity)
		if err != nil {
			return err
		}
		fmt.Printf("  - %s\n", userActivityString)
	}
	return nil
}

func activityString(userActivity githubUserData) (string, error) {
	switch userActivity.Type {
	case PUSH_EVENT:
		return fmt.Sprintf(
			"Pushed %d commits to %s",
			len(userActivity.Payload.Commits),
			userActivity.Repo.Name,
		), nil
	case CREATE_EVENT:
		if userActivity.Payload.RefType == "repository" {
			return fmt.Sprintf(
				"Created %s at %s",
				userActivity.Payload.RefType,
				userActivity.Repo.Name,
			), nil
		}

		return fmt.Sprintf(
			"Created %s '%s' at %s",
			userActivity.Payload.RefType,
			userActivity.Payload.Ref,
			userActivity.Repo.Name,
		), nil
	case WATCH_EVENT:
		if userActivity.Payload.Action == "started" {
			return fmt.Sprintf(
				"Started watching %s",
				userActivity.Repo.Name,
			), nil
		}

		return fmt.Sprintf(
			"Ended watching %s",
			userActivity.Repo.Name,
		), nil
	case DELETE_EVENT:
		return fmt.Sprintf(
			"Deleted %s '%s' at %s",
			userActivity.Payload.RefType,
			userActivity.Payload.Ref,
			userActivity.Repo.Name,
		), nil
	case FORK_EVENT:
		return fmt.Sprintf(
			"Forked %s to %s",
			userActivity.Repo.Name,
			userActivity.Payload.Forkee.FullName,
		), nil
	case ISSUES_EVENT:
		return fmt.Sprintf(
			"Issue '%s' %s at %s",
			userActivity.Payload.Issue.Title,
			userActivity.Payload.Action,
			userActivity.Repo.Name,
		), nil
	case ISSUES_COMMENT_EVENT:
		return fmt.Sprintf(
			"Commented at issue '%s' at %s",
			userActivity.Payload.Issue.Title,
			userActivity.Repo.Name,
		), nil
	case PULL_REQUEST_EVENT:
		return fmt.Sprintf(
			"Pull request '%s' %s at %s",
			userActivity.Payload.PullReq.Title,
			userActivity.Payload.Action,
			userActivity.Repo.Name,
		), nil
	case PUBLIC_EVENT:
		return fmt.Sprintf(
			"Repo %s is now public",
			userActivity.Repo.Name,
		), nil
	case MEMBER_EVENT:
		return fmt.Sprintf(
			"Added %s to %s",
			userActivity.Payload.Member.Login,
			userActivity.Repo.Name,
		), nil
	case RELEASE_EVENT:
		return fmt.Sprintf(
			"Released '%s' at %s",
			userActivity.Payload.Release.Name,
			userActivity.Repo.Name,
		), nil
	default:
		return fmt.Sprintf(
			"%s to %s",
			userActivity.Type,
			userActivity.Repo.Name,
		), nil
	}
}
