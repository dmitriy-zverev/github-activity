package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You did not specify username")
		return
	}

	fmt.Printf("Fetching activity for '%s'...\n", os.Args[1])

	_, err := fetchGithubUserData(os.Args[1])
	if err != nil {
		fmt.Printf("Error fetching user activity: %v\n", err)
		return
	}
}
