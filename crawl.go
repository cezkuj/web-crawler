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
	"time"
)

func main() {
	domain := os.Args[1]
	//using keys of map to imitate set
	visitedPages := make(map[string]bool)
	crawl(domain, "https://"+domain, visitedPages)
	log.Println(visitedPages)

}

func crawl(domain string, page string, visitedPages map[string]bool) {

	//Default http client does not have timeout
	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(page)
	if err != nil {
		log.Println(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	parse(doc, domain, visitedPages, page)
}

func parse(n *html.Node, domain string, visitedPages map[string]bool, foundOn string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				log.Println(a.Val)
				validateUrl(a.Val, domain, foundOn, visitedPages)
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse(c, domain, visitedPages, foundOn)
	}
}
func validateUrl(value, domain, foundOn string, visitedPages map[string]bool) {
	u := removeGetParams(value)
	u = removeChapterLinks(u)
        //return in case of cases not needed to cover
        if u == "" || u == "/" || strings.HasPrefix(u, "..") || strings.HasPrefix(u, "mailto:") ||  strings.HasPrefix(u, "tel:"){
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
	if !keyInMap(u, visitedPages) {
		log.Println("Added")
		visitedPages[u] = true
		crawl(domain, u, visitedPages)
	}
}

func keyInMap(key string, m map[string]bool) bool {
	_, ok := m[key]
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
        return getDomain(page) == domain
	//return strings.Contains(getDomain(page), domain)
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
