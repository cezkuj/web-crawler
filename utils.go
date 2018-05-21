package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"github.com/golang-collections/collections/stack"
)

func fetchAndParse(page string) (*html.Node, error) {
	//Default http client does not have timeout
	client := http.Client{Timeout: 150 * time.Second}
	resp, err := client.Get(page)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return doc, nil
}

func getDomain(page string) string {
	u, err := url.Parse(page)
	if err != nil {
		log.Fatal(err)
	}
	hostname := u.Hostname()
	if strings.HasPrefix(hostname, "www.") {
		return hostname[4:]
	}
	return hostname
}

func checkDomain(page, domain string, matchSubdomains bool) bool {
	if matchSubdomains {
		return strings.Contains(getDomain(page), domain)
	}
	return getDomain(page) == domain
}

func buildURL(foundOn, relSuffix string) string {
	if strings.HasPrefix(relSuffix, "/") {
		return "https://" + getDomain(foundOn) + relSuffix
	}
	return foundOn + "/" + relSuffix
}

func removeStringAfterChar(str, ch string) string {
	if i := strings.Index(str, ch); i != -1 {
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
func findLink(page Page) (string, bool) {
	if page.content.Type == html.ElementNode && page.content.Data == "a" {
		for _, attr := range page.content.Attr {
			if attr.Key == "href" {
				return attr.Val, true
			}
			return "", false
		}
	}
	return "", false

}
func insertURL(u, foundOn, domain string, matchSubdomains bool, visitedPages *sync.Map) (string, bool) {
	u = removeGetParams(u)
	u = removeChapterLinks(u)
	//return in case of cases not needed to cover
	if u == "" || u == "/" || strings.HasPrefix(u, "..") || strings.HasPrefix(u, "mailto:") || strings.HasPrefix(u, "tel:") {
		return "", false
	}
	//internal relative links
	if !strings.HasPrefix(u, "http") {
		u = buildURL(foundOn, u)
		//full path links, return if external domain
	} else if !checkDomain(u, domain, matchSubdomains) {
		return "", false
	}
	_, loaded := visitedPages.LoadOrStore(u, foundOn)
	return u, !loaded
}

func stackPages(pageName string, m map[string]string) (*stack.Stack, string) {
	st := stack.New()
	for ; m[pageName] != ""; pageName = m[pageName] {
		st.Push(pageName)
	}
	return st, pageName
}

func printResults(domain string, visitedPages sync.Map) {
	//cast sync.Map to map for simplicity
	pages := make(map[string]string)
	visitedPages.Range(castSyncMapToBultin(pages))
	p_head := PageTree{domain, make(map[string]PageTree)}
	for key := range pages {
		p := p_head
		st, _ := stackPages(key, pages)
		for st.Len() > 0 {
			pageName, ok := st.Pop().(string)
			if !ok {
				log.Println("not ok")
				break
			}
			if _, present := p.subPages[pageName]; !present {
				page := PageTree{pageName, make(map[string]PageTree)}
				p.subPages[page.name] = page
				p = page
			} else {
				p = p.subPages[pageName]
			}
		}
	}
	printPages(p_head, 0)
}

func printPages(page PageTree, depth int) {
	depth++
	log.Println(strings.Repeat(" ", depth*4), depth, page.name)
	for _, subpage := range page.subPages {
		printPages(subpage, depth)
	}
}
func castSyncMapToBultin(m map[string]string) func(key, value interface{}) bool {
	return func(key, value interface{}) bool {
                key_s, ok := key.(string)
                if !ok {
                   log.Println("not ok")
                }
                value_s, ok := value.(string)
                if !ok {
                   log.Println("not ok")
                }
		m[key_s] = value_s
		return true
	}
}
