package main

import (
	"bytes"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
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
func keyInMap(key string, m sync.Map) bool {
	_, ok := m.Load(key)
	return ok
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

func buildUrl(foundOn, relSuffix string) string {
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
func insertURL(u, foundOn, domain string, matchSubdomains bool, visitedPages *sync.Map) (string, bool) {
        u = removeGetParams(u)
        u = removeChapterLinks(u)
        //return in case of cases not needed to cover
        if u == "" || u == "/" || strings.HasPrefix(u, "..") || strings.HasPrefix(u, "mailto:") || strings.HasPrefix(u, "tel:") {
                return "", false
        }
        log.Println("Hi " + u + ", found on: " + foundOn)
        //internal relative links
        if !strings.HasPrefix(u, "http") {
                u = buildUrl(foundOn, u)
                //full path links, return if external domain
        } else if !checkDomain(u, domain, matchSubdomains) {
                return "", false
        }
        _, loaded := visitedPages.LoadOrStore(u, foundOn)
        return u, !loaded
}
