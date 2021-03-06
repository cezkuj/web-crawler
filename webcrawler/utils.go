package webcrawler

import (
	"bytes"
	"crypto/tls"
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

func clientWithTimeout(tlsSecure bool) (client http.Client) {
	timeout := 30 * time.Second
	//Default http client does not have timeout
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !tlsSecure},
	}
	return http.Client{Timeout: timeout, Transport: tr}

}

//fetchs provided page, then parses it to *html.Node and returns
func fetchAndParse(page string, tlsSecure bool, client http.Client) (*html.Node, error) {
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
func getProt(tls bool) string {
	prot := "https"
	if !tls {
		prot = "http"
	}
	return prot

}
func getDomain(page string) (string, error) {
	u, err := url.Parse(page)
	if err != nil {
		return "", err
	}
	hostname := u.Hostname()
	if strings.HasPrefix(hostname, "www.") {
		return hostname[4:], nil
	}
	return hostname, nil
}

// checks if page is in expected domain(or one of subdomains in case matchSubdomains is true)
func checkDomain(page, domain string, matchSubdomains bool) bool {
	pageDomain, err := getDomain(page)
	if err != nil {
		log.Println(err)
		return false
	}
	if matchSubdomains {
		return strings.Contains(pageDomain, domain)
	}
	return pageDomain == domain
}

// in case of relative links, appends
func buildURL(foundOn, relSuffix string, tls bool) string {
	if strings.HasPrefix(relSuffix, "/") {
		pageDomain, err := getDomain(foundOn)
		if err != nil {
			return ""
		}
		return getProt(tls) + "://" + pageDomain + relSuffix
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
func insertURL(u, foundOn, domain string, matchSubdomains bool, visitedPages *sync.Map, tls bool) (string, bool) {
	u = removeGetParams(u)
	u = removeChapterLinks(u)
	//return in case of cases not needed to cover
	if u == "" || u == "/" || strings.HasPrefix(u, "..") || strings.HasPrefix(u, "mailto:") || strings.HasPrefix(u, "tel:") {
		return "", false
	}
	//internal relative links
	if !strings.HasPrefix(u, "http") {
		u = buildURL(foundOn, u, tls)
		if u == "" {
			return u, false
		}
		//full path links, return if external domain
	} else if !checkDomain(u, domain, matchSubdomains) {
		return "", false
	}
	_, loaded := visitedPages.LoadOrStore(u, foundOn)
	return u, !loaded
}

func stackPages(pageName interface{}, m *sync.Map) *stack.Stack {
	st := stack.New()
	for {
		if _, found := m.Load(pageName); !found {
			break
		}
		st.Push(pageName)
		pageName, _ = m.Load(pageName)
	}
	return st
}

func PrintResults(domain string, visitedPages *sync.Map) {
	head := PageTree{domain, make(map[string]PageTree)}
	visitedPages.Range(iterateOverKeys(head, visitedPages))
	printPages(head, 0)
}

func printPages(page PageTree, depth int) {
	depth++
	log.Println(strings.Repeat(" ", depth*4), depth, page.name)
	for _, subpage := range page.subPages {
		printPages(subpage, depth)
	}
}

func iterateOverKeys(head PageTree, pages *sync.Map) func(key, value interface{}) bool {
	return func(key, _ interface{}) bool {
		p := head
		st := stackPages(key, pages)
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
		return true
	}
}
