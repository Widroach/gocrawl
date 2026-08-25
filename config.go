package main

import (
	"net/url"
	"sync"
	"sync/atomic"
)

type config struct {
	maxPages           atomic.Int64
	pages              map[string]PageData
	externalPages      map[string]any
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	allowSubdomains    bool
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, visited := cfg.pages[normalizedURL]; visited {
		return false
	}

	cfg.pages[normalizedURL] = PageData{URL: normalizedURL}
	return true

}

func (cfg *config) addExternalPage(externalPageURL string) {
	cfg.mu.Lock()
	cfg.externalPages[externalPageURL] = externalPageURL
	cfg.mu.Unlock()

}

func (cfg *config) setPageData(normalizedURL string, data PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normalizedURL] = data
}

func (cfg *config) lenPages() int {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages)
}

func initialize(baseURL string, maxThreads, maxPages int) (*config, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return &config{}, err
	}

	cfg := &config{
		pages:              map[string]PageData{},
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxThreads),
		wg:                 &sync.WaitGroup{},
		baseURL:            parsedBaseURL,
		allowSubdomains:    false,
		externalPages:      map[string]any{},
		maxPages:           atomic.Int64{},
	}
	cfg.maxPages.Store(int64(maxPages))
	return cfg, nil
}
