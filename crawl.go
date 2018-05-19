package main

import (
	"bytes"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)
// docu
func Doc(){
}
func main() {
	domain := os.Args[1]
	//using keys of map to imitate set
	visitedPages := &sync.Map{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go crawl(domain, "https://"+domain, visitedPages, wg)
	wg.Wait()
        log.Println("finished")
	visitedPages.Range(printMap)

}

func crawl(domain string, page string, visitedPages *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
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
	wg.Add(1)
	go parse(doc, domain, visitedPages, page, wg)
}

func parse(n *html.Node, domain string, visitedPages *sync.Map, foundOn string, wg *sync.WaitGroup) {
	defer wg.Done()
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				log.Println(a.Val)
				validateUrl(a.Val, domain, foundOn, visitedPages, wg)
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		wg.Add(1)
		go parse(c, domain, visitedPages, foundOn, wg)
	}
}
func validateUrl(value, domain, foundOn string, visitedPages *sync.Map, wg *sync.WaitGroup) {
	u := removeGetParams(value)
	u = removeChapterLinks(u)
	//return in case of cases not needed to cover
	if u == "" || u == "/" || strings.HasPrefix(u, "..") || strings.HasPrefix(u, "mailto:") || strings.HasPrefix(u, "tel:") {
		return
	}
	log.Println("Hi " + u + ", found on: " + foundOn)
	//internal relative links
	if !strings.HasPrefix(u, "http") {
		u = buildUrl(foundOn, u)
		log.Println("Rel: " + u)
		//full path links, return if external domain
	} else if !checkDomain(u, domain) {
		return
	}
	if !keyInMap(u, *visitedPages) {
		log.Println("Added")
		visitedPages.Store(u, true)
		wg.Add(1)
		go crawl(domain, u, visitedPages, wg)
	}
}

func keyInMap(key string, m sync.Map) bool {
	_, ok := m.Load(key)
	return ok
}
func getDomain(page string) string {
	u, err := url.Parse(page)
	if err != nil {
		log.Fatal(err)
	}
	return u.Hostname()
}
func checkDomain(page, domain string) bool {
        //log.Println(getDomain(page), domain)
	//return getDomain(page) == domain
	return strings.Contains(getDomain(page), domain)
}

func buildUrl(foundOn, relSuffix string) string {
	log.Println(foundOn, relSuffix)
	if strings.HasPrefix(relSuffix, "/") {
		return "https://" + getDomain(foundOn) + relSuffix
	}
	return foundOn + "/" + relSuffix
}
func removeStringAfterChar(str, ch string) string {
	if i := strings.Index(str, ch); i != -1 {
		log.Println(str + " found")
		log.Println(str, str[:i])
		return str[:i]
	}
	return str

}
func removeGetParams(u string) string {
	return removeStringAfterChar(u, "?")

}
func removeChapterLinks(u string) string {
	return removeStringAfterChar(u, "#")
}

func printMap(key, value interface{}) bool {
	log.Println(key, value)
	return true
}
