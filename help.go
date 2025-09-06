package main

import "fmt"

func help() {
	fmt.Println("Usage: github-activity <username>")
	fmt.Println("\nAdditional parameters:")
	fmt.Println("  -f (--filter) [event type]")
	fmt.Println("  -p (--page) [page number]")
	fmt.Println("  -n (--number) [per page events]")
}
