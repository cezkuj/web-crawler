package main

import (
	"os"
)

func main() {
	domain := os.Args[1]
	//using keys of map to imitate set
	matchSubdomains := false
	crawler := NewFireAndForgetCrawler(domain, matchSubdomains)
	results := crawler.Crawl()
	results.Range(printMap)

}
