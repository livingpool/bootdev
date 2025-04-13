package main

import (
	"fmt"
	"log"
	"os"

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

	pages := make(map[string]int)
	err := util.CrawlPage(args[1], args[1], pages)
	if err != nil {
		log.Fatalf("crawler error: %v", err)
	}

	fmt.Println("crawler results:")
	for k, v := range pages {
		fmt.Printf("%s: %d", k, v)
	}
}
