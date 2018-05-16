package main

import (
	"log"
	"sync"
	"time"
)

type IdiomaticCrawler struct {
	domain          string
	visitedPages    *sync.Map
	wg              *sync.WaitGroup
	matchSubdomains bool
	toScrap         chan Page
	toFetch         chan string
        reqInterval     time.Duration
}

func NewIdiomaticCrawler(domain string, matchSubdomains bool, reqInterval time.Duration) *IdiomaticCrawler {
	return &IdiomaticCrawler{
		domain:          domain,
		visitedPages:    &sync.Map{},
		wg:              &sync.WaitGroup{},
		matchSubdomains: matchSubdomains,
		toScrap:         make(chan Page),
		toFetch:         make(chan string),
                reqInterval: reqInterval,
	}
}
func (crawler IdiomaticCrawler) Crawl() sync.Map {
	crawler.wg.Add(1)
	go func() {
		for {
			time.Sleep(crawler.reqInterval)
			go crawler.fetch(<-crawler.toFetch)
		}
	}()
	go func() {
		for {
			go crawler.scrap(<-crawler.toScrap)
		}
	}()
	crawler.toFetch <- ("https://" + crawler.domain)
	crawler.wg.Wait()
        //Avoid infinite loops in printing by deleting main page
        crawler.visitedPages.Delete("https://" + crawler.domain)
	return *crawler.visitedPages
}

func (crawler IdiomaticCrawler) fetch(page string) {
	defer crawler.wg.Done()
	doc, err := fetchAndParse(page)
	if err != nil {
		log.Println(err)
		return
	}
	crawler.wg.Add(1)
	crawler.toScrap <- Page{name: page, content: doc}
}

func (crawler IdiomaticCrawler) scrap(page Page) {
	defer crawler.wg.Done()
	if link, found := findLink(page); found {
                log.Println("Found link: " + link)
		if u, inserted := insertURL(link, page.name, crawler.domain, crawler.matchSubdomains, crawler.visitedPages); inserted {
			crawler.wg.Add(1)
			crawler.toFetch <- u
		}

	}
	for child := page.content.FirstChild; child != nil; child = child.NextSibling {
		crawler.wg.Add(1)
		crawler.toScrap <- Page{name: page.name, content: child}
	}
}
