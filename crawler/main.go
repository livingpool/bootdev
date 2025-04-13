package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/livingpool/bootdev-crawler/util"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) >= 5 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	// parse command line arguments
	baseURL, err := url.Parse(args[1])
	if err != nil {
		log.Fatalf("error parsing '%s': err: %v", args[1], err)
	}

	var maxConcurrency, maxPages int = 1, 1
	if len(args) >= 3 {
		maxConcurrency, err = strconv.Atoi(args[2])
		if err != nil {
			log.Fatalf("maxConcurrency should be integer, got=%s", args[2])
		}
	}

	if len(args) >= 4 {
		maxPages, err = strconv.Atoi(args[3])
		if err != nil {
			log.Fatalf("maxPages should be integer, got=%s", args[2])
		}
	}

	fmt.Println("starting crawler of:", args[1])

	config := util.Config{
		Pages:              make(map[string]int),
		BaseURL:            baseURL,
		Mu:                 &sync.Mutex{},
		ConcurrencyControl: make(chan struct{}, maxConcurrency),
		Wg:                 &sync.WaitGroup{},
		MaxPages:           maxPages,
	}

	config.Wg.Add(1)
	go config.CrawlPage(baseURL.String())
	config.Wg.Wait()

	util.PrintReport(config.Pages, baseURL.String())
}
