package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}

	pageNum := DEFAULT_PAGE_NUM
	perPageNum := DEFAULT_PER_PAGE_EVENTS

	if slices.Contains(os.Args, "-p") {
		idx := slices.Index(os.Args, "-p")
		if len(os.Args) > idx+1 {
			if _, err := strconv.Atoi(os.Args[idx+1]); err != nil {
				fmt.Printf("Error while parsing page number: %v is not a number\n", os.Args[idx+1])
				return
			}
			pageNum = os.Args[idx+1]
		}
	}

	if slices.Contains(os.Args, "-n") {
		idx := slices.Index(os.Args, "-n")
		if len(os.Args) > idx+1 {
			if _, err := strconv.Atoi(os.Args[idx+1]); err != nil {
				fmt.Printf("Error while parsing per page events number: %v is not a number\n", os.Args[idx+1])
				return
			}
			perPageNum = os.Args[idx+1]
		}
	}

	fmt.Printf(
		"Fetching activity for '%s' at page %s with %s per page events...\n",
		os.Args[1],
		pageNum,
		perPageNum,
	)

	activities, err := fetchGithubUserData(os.Args[1], pageNum, perPageNum)
	if err != nil {
		fmt.Printf("Error fetching user activity: %v\n", err)
		return
	}

	eventTypeFiler := DEFAULT_FILTER_TYPE
	if slices.Contains(os.Args, "-f") {
		idx := slices.Index(os.Args, "-f")
		if len(os.Args) > idx+1 {
			eventTypeFiler = os.Args[idx+1]
		}
	}

	activities = filterEvents(activities, eventTypeFiler)

	if len(activities) < 1 && eventTypeFiler != DEFAULT_FILTER_TYPE {
		fmt.Printf("  No result for '%s' filter.\n", eventTypeFiler)
		return
	}

	if err := printer(activities); err != nil {
		fmt.Printf("Couldn't print user activity: %v\n", err)
		return
	}
}
