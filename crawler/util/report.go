package util

import (
	"fmt"
	"slices"
	"strings"
)

type Page struct {
	Count int
	URL   string
}

func PrintReport(pages map[string]int, baseURL string) {
	startLine := fmt.Sprintf(`
=============================
  REPORT for %s
=============================
`, baseURL)

	fmt.Println(startLine)

	res := []Page{}
	for k, v := range pages {
		res = append(res, Page{v, k})
	}

	slices.SortFunc(res, SortPages)

	for _, line := range res {
		fmt.Printf("Found %d internal links to %s\n", line.Count, line.URL)
	}
}

func SortPages(a, b Page) int {
	if a.Count > b.Count {
		return -1
	} else if a.Count < b.Count {
		return 1
	} else {
		return strings.Compare(a.URL, b.URL)
	}
}
