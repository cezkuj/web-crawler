package main

import (
	"os"
)

func main() {
	domain := os.Args[1]
	//using keys of map to imitate set
	matchSubdomains := false
	crawler := NewIdiomaticCrawler(domain, matchSubdomains)
	results := crawler.Crawl()
	printResults(domain, results)

}
