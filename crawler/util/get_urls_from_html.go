package util

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func GetURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	reader := strings.NewReader(htmlBody)
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse HTML: %v", err)
	}

	var urls []string
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					href, err := url.Parse(a.Val)
					if err != nil {
						fmt.Printf("couldn't parse href '%v': %v\n", a.Val, err)
						continue
					}
					resolvedURL := baseURL.ResolveReference(href)
					urls = append(urls, resolvedURL.String())
				}
			}
		}
	}

	return urls, nil
}
