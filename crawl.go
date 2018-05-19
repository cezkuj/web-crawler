package main

import (
	"io/ioutil"
	"log"
	"net/http"
        "net/url"
        "time"
        "bytes"
        "strings"
        "golang.org/x/net/html"
        "os"
)

func main() {
        domain := os.Args[1]
        //using keys of map to imitate set
        visitedPages := make(map[string]bool)
        crawl(domain, "https://" + domain, visitedPages)
        log.Println(visitedPages)

}

func crawl(domain string, page string, visitedPages map[string]bool){

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
                log.Println("Hi " + a.Val + ", found on: " + foundOn)
                if strings.HasPrefix(a.Val, "/") || strings.HasPrefix(a.Val, "../"){
                    u := buildUrl(foundOn, a.Val)
                    log.Println("Rel: " + u)
                    if !keyInMap(u, visitedPages) {
                        log.Println("Added")
                        visitedPages[u] = true
                        crawl(domain, u, visitedPages)
                        break
                    }
                } else if checkDomain(a.Val, domain) && !keyInMap(a.Val, visitedPages){
                  log.Println("Full: " + a.Val)
                  visitedPages[a.Val] = true
                  crawl(domain, a.Val, visitedPages)
                  break
                }
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        parse(c, domain, visitedPages, foundOn)
    }
}

func keyInMap(key string, m map[string]bool) bool{
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
func checkDomain(page, domain string) bool{
   return strings.Contains(getDomain(page), domain)
}

func buildUrl(foundOn, relSuffix string) string{
  //covering case with / needed before ../ suffix
  if strings.HasPrefix(relSuffix, "../") {
     relSuffix = "/" + relSuffix
  } 
  return "https://" + getDomain(foundOn) + relSuffix 
} 

