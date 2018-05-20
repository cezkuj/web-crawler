package main

import (
	"bytes"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)
type Crawler interface {
   Crawl(domain string) sync.Map
}
type FireAndForgetCrawler struct{
     domain string
     visitedPages *sync.Map
     wg *sync.WaitGroup
     matchSubdomains bool
}
func NewFireAndForgetCrawler(domain string, matchSubdomains bool) *FireAndForgetCrawler{
  return &FireAndForgetCrawler{domain: domain, visitedPages: &sync.Map{}, wg: &sync.WaitGroup{}, matchSubdomains: matchSubdomains}
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
	//Default http client does not have timeout
	client := http.Client{Timeout: 150 * time.Second}
	resp, err := client.Get(page)
	if err != nil {
		log.Println(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Println(err)
		return
	}
	crawler.wg.Add(1)
	go crawler.parse(doc, page)
}

func (crawler FireAndForgetCrawler) parse(node *html.Node, foundOn string) {
	defer crawler.wg.Done()
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				crawler.validateUrl(attr.Val, foundOn)
				break
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		crawler.wg.Add(1)
		go crawler.parse(child, foundOn)
	}
}
func (crawler FireAndForgetCrawler) validateUrl(u, foundOn string) {
	u = removeGetParams(u)
	u = removeChapterLinks(u)
	//return in case of cases not needed to cover
	if u == "" || u == "/" || strings.HasPrefix(u, "..") || strings.HasPrefix(u, "mailto:") || strings.HasPrefix(u, "tel:") {
		return
	}
	log.Println("Hi " + u + ", found on: " + foundOn)
	//internal relative links
	if !strings.HasPrefix(u, "http") {
		u = buildUrl(foundOn, u)
		//full path links, return if external domain
	} else if !checkDomain(u, crawler.domain, crawler.matchSubdomains) {
		return
	}
	if !keyInMap(u, *crawler.visitedPages) {
		crawler.visitedPages.Store(u, true)
		crawler.wg.Add(1)
		go crawler.fetch(u)
	}
}
