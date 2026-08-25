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
	cfg, err := initialize(baseURL, 10)
	if err != nil {
		slog.Error("failed to initialize", "err", err)
		return
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(baseURL)
	cfg.wg.Wait()
	for key, value := range cfg.pages {
		fmt.Printf("Domain: %s, count: %d\n", key, len(value.OutgoingLinks))
	}
	for key, _ := range cfg.externalPages {
		fmt.Printf("External domain: %s\n", key)
	}
	slog.Info("finished crawling URLs", "total_url", len(cfg.pages) + len(cfg.externalPages))

}
