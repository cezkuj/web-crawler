package main

import (
	"golang.org/x/net/html"
	"log"
	"sync"
)

type FireAndForgetCrawler struct {
	domain          string
	visitedPages    *sync.Map
	wg              *sync.WaitGroup
	matchSubdomains bool
}

func NewFireAndForgetCrawler(domain string, matchSubdomains bool) *FireAndForgetCrawler {
	return &FireAndForgetCrawler{
		domain:          domain,
		visitedPages:    &sync.Map{},
		wg:              &sync.WaitGroup{},
		matchSubdomains: matchSubdomains,
	}
}
func (crawler FireAndForgetCrawler) Crawl() sync.Map {
	crawler.wg.Add(1)
	go crawler.fetch("https://" + crawler.domain)
	crawler.wg.Wait()
	log.Println("finished")
	return *crawler.visitedPages

}

func (crawler FireAndForgetCrawler) fetch(page string) {
	defer crawler.wg.Done()
	doc, err := fetchAndParse(page)
	if err != nil {
		log.Println(err)
		return
	}
	crawler.wg.Add(1)
	go crawler.parse(Page{name: page, content: doc})
}

func (crawler FireAndForgetCrawler) parse(page Page) {
	defer crawler.wg.Done()
	if page.content.Type == html.ElementNode && page.content.Data == "a" {
		for _, attr := range page.content.Attr {
                          if u, inserted := insertURL(attr.Val, page.name, crawler.domain, crawler.matchSubdomains, crawler.visitedPages); inserted {
                                        crawler.wg.Add(1)
                                        go crawler.fetch(u)
                                }
                                break
		}
	}
	for child := page.content.FirstChild; child != nil; child = child.NextSibling {
		crawler.wg.Add(1)
		go crawler.parse(Page{name: page.name, content: child})
	}
}
