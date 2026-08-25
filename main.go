package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	maxThreads := flag.Int("threads", 3, "max concurrent threads")
	maxPages := flag.Int("limit", 0, "max pages to crawl (0 = unlimited)")
	flag.Parse()
	baseURL := flag.Arg(0)

	if flag.NArg() < 1 {
		fmt.Println("usage: gocrawl  [--threads N] [--limit N] <url>")
		os.Exit(1)
	}

	cfg, err := initialize(baseURL, *maxThreads, *maxPages)
	if err != nil {
		slog.Error("failed to initialize", "err", err)
		return
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(baseURL)
	cfg.wg.Wait()
	slog.Info("finished crawling URLs", "total_url", len(cfg.pages)+len(cfg.externalPages))
	if err := writeJSONReport(cfg.pages, "report.json"); err != nil {
		slog.Error("error saving the result", "err", err)
	}

}
