package main

import (
	"os"
        "time"
        "log"
)

func main() {
	domain := os.Args[1]
	//using keys of map to imitate set
	matchSubdomains := false
        reqInterval := time.Millisecond * 0
        start := time.Now()
	crawler := NewIdiomaticCrawler(domain, matchSubdomains, reqInterval)
	results := crawler.Crawl()
        endCrawl := time.Now()
	printResults(domain, results)
        endPrint := time.Now()
        log.Println("Whole program execution: ", endPrint.Sub(start))
        log.Println("Crawling: ", endCrawl.Sub(start))
        log.Println("Printing: ", endPrint.Sub(endCrawl)) 
}
