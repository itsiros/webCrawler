package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
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

func (cfg *config) setPageData(normalizedURL string, data PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normalizedURL] = data
}

func configure(rawBaseURL string, maxConcurrency, maxPages int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	return &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}, nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {

	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}

	cfg.mu.Unlock()
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if currentURL.Hostname() != cfg.baseURL.Hostname() {
		return
	}

	normed := NormalizeURL(rawCurrentURL)
	if !cfg.addPageVisit(normed) {
		return
	}

	body, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	p := ExtractPageData(body, rawCurrentURL)
	cfg.setPageData(normed, p)

	for _, url := range p.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}
