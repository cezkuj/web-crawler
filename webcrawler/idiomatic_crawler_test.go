package webcrawler

import (
	"testing"
        "time"
)

func TestIdiomaticCrawler(*testing.T) {
	domain := "lhsystems.pl"
	matchSubdomains := true
        reqInterval := time.Millisecond * 0
	crawler := NewIdiomaticCrawler(domain, matchSubdomains, reqInterval)
	crawler.Crawl()

}
