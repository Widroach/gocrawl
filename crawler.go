package main

import (
	"log/slog"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	if !isSameDomain(rawBaseURL, rawCurrentURL) {
		slog.Info("skipping page, it is the same domain as its parent", "current", rawCurrentURL, "parent", rawBaseURL)
		return
	}

	currentURL, _ := normalizeURL(rawCurrentURL)
	if _, exist := pages[currentURL]; exist {
		pages[currentURL]++
		return
	} else {
		pages[currentURL] = 1
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		slog.Warn("failed to get html", "warn", err)
		return
	}
	slog.Info("retrieved html content", "url", rawBaseURL, "length", len([]byte(html)))

	pageData := extractPageData(html, rawBaseURL)
	for _, url := range pageData.OutgoingLinks {
		slog.Info("crawling new page", "url", url)
		crawlPage(rawBaseURL, url, pages)
		slog.Info("finished crawling page", "url", url)
	}
}

func isSameDomain(rawBaseURL, rawCurrentURL string) bool {
	parsedBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		slog.Warn("failed to parse url", "url", rawBaseURL)
		return false
	}
	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		slog.Warn("failed to parse url", "url", rawCurrentURL)
		return false
	}

	if parsedBaseURL.Hostname() != parsedCurrentURL.Hostname() {
		return false
	}
	return true
}
