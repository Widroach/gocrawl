package main

import (
	"net/url"
	"sync"
)

type config struct {
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

func initialize(baseURL string, maxThreads int) (config, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return config{}, err
	}

	return config{
		pages:              map[string]PageData{},
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxThreads),
		wg:                 &sync.WaitGroup{},
		baseURL:            parsedBaseURL,
		allowSubdomains:    false,
		externalPages:      map[string]any{},
	}, nil
}
