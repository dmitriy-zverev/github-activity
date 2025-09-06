package main

import "strings"

func filterEvents(events []githubUserData, eventType string) []githubUserData {
	if eventType == DEFAULT_FILTER_TYPE {
		return events
	}

	var newEvents []githubUserData

	for _, event := range events {
		if strings.Contains(strings.ToLower(event.Type), strings.ToLower(eventType)) {
			newEvents = append(newEvents, event)
		}
	}

	return newEvents
}
