package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[1]

	fmt.Printf("starting crawl of: %s", baseURL)
	fmt.Print(getHTML(baseURL))
}
