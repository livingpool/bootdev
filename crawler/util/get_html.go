package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Config struct {
	Pages              map[string]int
	BaseURL            *url.URL
	Mu                 *sync.Mutex
	ConcurrencyControl chan struct{}
	Wg                 *sync.WaitGroup
	MaxPages           int
}

func (cfg *Config) CrawlPage(rawCurrentURL string) {
	defer func() {
		cfg.Wg.Done()
		<-cfg.ConcurrencyControl
	}()

	cfg.ConcurrencyControl <- struct{}{}

	cfg.Mu.Lock()
	if len(cfg.Pages) >= cfg.MaxPages {
		cfg.Mu.Unlock()
		return
	}
	cfg.Mu.Unlock()

	fmt.Println("crawling", rawCurrentURL)

	parsedRawCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("couldn't parse '%s', err: %v\n", rawCurrentURL, err)
		return
	}

	// skip other sites
	if cfg.BaseURL.Host != parsedRawCurrentURL.Host {
		return
	}

	normalizedURL, err := NormalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("failed to normalize '%s', err: %v\n", rawCurrentURL, err)
		return
	}

	// dont parse an already visited site
	if !cfg.addPageVisit(normalizedURL) {
		return
	}

	html, err := GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("failed to get html: %v\n", err)
		return
	}

	urls, err := GetURLsFromHTML(string(html), rawCurrentURL)
	if err != nil {
		fmt.Printf("failed to get urls: %v\n", err)
		return
	}

	for _, url := range urls {
		cfg.Wg.Add(1)
		go cfg.CrawlPage(url)
	}
}

func (cfg *Config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.Mu.Lock()
	defer cfg.Mu.Unlock()

	if _, exists := cfg.Pages[normalizedURL]; exists {
		cfg.Pages[normalizedURL]++
		return false
	}

	cfg.Pages[normalizedURL] = 1
	return true
}
func GetHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("network error: %v", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("http error: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("got non-HTML response: %s", contentType)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(data), nil
}
