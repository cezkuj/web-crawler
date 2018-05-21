package main

import (
	"golang.org/x/net/html"
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
}

func NewIdiomaticCrawler(domain string, matchSubdomains bool) *IdiomaticCrawler {
	return &IdiomaticCrawler{
		domain:          domain,
		visitedPages:    &sync.Map{},
		wg:              &sync.WaitGroup{},
		matchSubdomains: matchSubdomains,
		toScrap:         make(chan Page),
		toFetch:         make(chan string),
	}
}
func (crawler IdiomaticCrawler) Crawl() sync.Map {
	crawler.wg.Add(1)
	go func() {
		for {
                        time.Sleep(100 * time.Millisecond)
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
	log.Println("finished")
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
	if page.content.Type == html.ElementNode && page.content.Data == "a" {
		for _, attr := range page.content.Attr {
			if attr.Key == "href" {
				if u, inserted := insertURL(attr.Val, page.name, crawler.domain, crawler.matchSubdomains, crawler.visitedPages); inserted {
					crawler.wg.Add(1)
					crawler.toFetch <- u
				}
				break

			}

		}
	}
	for child := page.content.FirstChild; child != nil; child = child.NextSibling {
		crawler.wg.Add(1)
		crawler.toScrap <- Page{name: page.name, content: child}
	}
}
