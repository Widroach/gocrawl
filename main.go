package main

import (
	"fmt"
	"log/slog"
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

	slog.Info("starting to crawl page", "url", baseURL)
	pages := map[string]int{}
	crawlPage(baseURL, baseURL, pages)
	for key, value := range pages {
		fmt.Printf("Domain: %s, count: %d\n", key, value)
	}
	slog.Info("finished crawling URLs", "total", len(pages))

}
