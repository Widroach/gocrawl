package main

import (
	"log/slog"
	"net/url"
	"strings"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	maxPages := int(cfg.maxPages.Load())
	if maxPages > 0 && cfg.lenPages() >= maxPages {
		return
	}

	if !isSameDomain(cfg.baseURL.String(), rawCurrentURL, cfg.allowSubdomains) {
		slog.Debug("skipping page crawling, current URL is a different domain than its parent", "current", rawCurrentURL, "parent", cfg.baseURL)
		cfg.addExternalPage(rawCurrentURL)
		return
	}

	normalizedURL, _ := normalizeURL(rawCurrentURL)
	if isFirst := cfg.addPageVisit(normalizedURL); !isFirst {
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		slog.Warn("failed to get html", "warn", err)
		return
	}
	pageData := extractPageData(html, rawCurrentURL)
	cfg.setPageData(normalizedURL, pageData)
	for _, url := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}

	slog.Debug("finished crawling the page", "url", rawCurrentURL)
}

func isSameDomain(rawBaseURL, rawCurrentURL string, allowSubdomains bool) bool {
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

	baseHost := parsedBaseURL.Hostname()
	currentHost := parsedCurrentURL.Hostname()

	if allowSubdomains {
		return getRootDomain(baseHost) == getRootDomain(currentHost)
	}
	return baseHost == currentHost
}

func getRootDomain(hostname string) string {
	parts := strings.Split(hostname, ".")
	if len(parts) <= 2 {
		return hostname
	}
	return strings.Join(parts[len(parts)-2:], ".")
}
