package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func CrawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) error {
	parsedRawBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return fmt.Errorf("couldn't parse '%s', err: %v\n", rawBaseURL, err)
	}

	parsedRawCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return fmt.Errorf("couldn't parse '%s', err: %v\n", rawCurrentURL, err)
	}

	// skip other sites
	if parsedRawBaseURL.Host != parsedRawCurrentURL.Host {
		return nil
	}

	normalizedURL, err := NormalizeURL(rawCurrentURL)
	if err != nil {
		return fmt.Errorf("failed to normalize '%s', err: %v", rawCurrentURL, err)
	}

	if _, exists := pages[normalizedURL]; exists {
		pages[normalizedURL]++
		return nil
	} else {
		pages[normalizedURL] = 1
	}

	html, err := GetHTML(rawCurrentURL)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(string(html))

	urls, err := GetURLsFromHTML(string(html), normalizedURL)
	if err != nil {
		return err
	}

	for _, url := range urls {
		err := CrawlPage(rawBaseURL, url, pages)
		if err != nil {
			return err
		}
	}
	return nil
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
