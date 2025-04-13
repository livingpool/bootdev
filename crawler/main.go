package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/livingpool/bootdev-crawler/util"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) >= 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	fmt.Println("starting crawler of:", args[1])

	baseURL, err := url.Parse(args[1])
	if err != nil {
		log.Fatalf("error parsing '%s': err: %v", args[1], err)
	}

	config := util.Config{
		Pages:              make(map[string]int),
		BaseURL:            baseURL,
		Mu:                 &sync.RWMutex{},
		ConcurrencyControl: make(chan struct{}, 1),
		Wg:                 &sync.WaitGroup{},
	}

	config.Wg.Add(1)
	go config.CrawlPage(baseURL.String())
	config.Wg.Wait()

	fmt.Println("crawler results:")
	for k, v := range config.Pages {
		fmt.Printf("%s: %d\n", k, v)
	}
}
