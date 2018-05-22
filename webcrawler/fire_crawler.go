package webcrawler

import (
	"log"
	"sync"
)

type FireAndForgetCrawler struct {
	domain          string
	visitedPages    *sync.Map
	wg              *sync.WaitGroup
	matchSubdomains bool
        tls bool
}

func NewFireAndForgetCrawler(domain string, matchSubdomains bool, tls bool) *FireAndForgetCrawler {
	return &FireAndForgetCrawler{
		domain:          domain,
		visitedPages:    &sync.Map{},
		wg:              &sync.WaitGroup{},
		matchSubdomains: matchSubdomains,
                tls: tls,
	}
}
func (crawler FireAndForgetCrawler) Crawl() sync.Map {
	crawler.wg.Add(1)
        prot := "https"
        if crawler.tls {
           prot = "http"
        }
        mainPage := prot + "://" + crawler.domain
	go crawler.fetch(mainPage)
	crawler.wg.Wait()
	//Avoid infinite loops in printing by deleting main page
	crawler.visitedPages.Delete(mainPage)
	log.Println("finished")
	return *crawler.visitedPages

}

func (crawler FireAndForgetCrawler) fetch(page string) {
	defer crawler.wg.Done()
	doc, err := fetchAndParse(page, crawler.tls)
	if err != nil {
		log.Println(err)
		return
	}
	crawler.wg.Add(1)
	go crawler.parse(Page{name: page, content: doc})
}

func (crawler FireAndForgetCrawler) parse(page Page) {
	defer crawler.wg.Done()
	if link, found := findLink(page); found {
		if u, inserted := insertURL(link, page.name, crawler.domain, crawler.matchSubdomains, crawler.visitedPages); inserted {
			crawler.wg.Add(1)
			go crawler.fetch(u)
		}

	}
	for child := page.content.FirstChild; child != nil; child = child.NextSibling {
		crawler.wg.Add(1)
		go crawler.parse(Page{name: page.name, content: child})
	}
}
