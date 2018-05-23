package webcrawler

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type IdiomaticCrawler struct {
	domain          string
	visitedPages    *sync.Map
	wg              *sync.WaitGroup
	client          http.Client
	matchSubdomains bool
	toScrap         chan Page
	toFetch         chan string
	reqInterval     time.Duration
	tls             bool
}

func NewIdiomaticCrawler(domain string, matchSubdomains bool, reqInterval time.Duration, tls bool) *IdiomaticCrawler {
	return &IdiomaticCrawler{
		domain:          domain,
		visitedPages:    &sync.Map{},
		wg:              &sync.WaitGroup{},
		client:          clientWithTimeout(tls),
		matchSubdomains: matchSubdomains,
		toScrap:         make(chan Page),
		toFetch:         make(chan string),
		reqInterval:     reqInterval,
		tls:             tls,
	}
}
func (crawler IdiomaticCrawler) Crawl() sync.Map {
	crawler.wg.Add(1)
	go func() {
		for page := range(crawler.toFetch){
			time.Sleep(crawler.reqInterval)
			go crawler.fetch(page)
		}
	}()
	go func() {
		for page := range(crawler.toScrap) {
			go crawler.scrap(page)
		}
	}()
	prot := "https"
	if !crawler.tls {
		prot = "http"
	}
	mainPage := prot + "://" + crawler.domain
	crawler.toFetch <- (mainPage)
	crawler.wg.Wait()
        close(crawler.toFetch)
        close(crawler.toScrap)
	//Avoid infinite loops in printing by deleting main page
	crawler.visitedPages.Delete(mainPage)
	return *crawler.visitedPages
}

func (crawler IdiomaticCrawler) fetch(page string) {
	defer crawler.wg.Done()
	doc, err := fetchAndParse(page, crawler.tls, crawler.client)
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
		if u, inserted := insertURL(link, page.name, crawler.domain, crawler.matchSubdomains, crawler.visitedPages, crawler.tls); inserted {
                        crawler.wg.Add(1)
			crawler.toFetch <- u
		}

	}
	for child := page.content.FirstChild; child != nil; child = child.NextSibling {
                crawler.wg.Add(1)
		crawler.toScrap <- Page{name: page.name, content: child}
	}
}
