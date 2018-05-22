package webcrawler

import (
	"testing"
)

func TestFireAndForgetCrawler(*testing.T) {
	domain := "lhsystems.pl"
	matchSubdomains := false
	crawler := NewFireAndForgetCrawler(domain, matchSubdomains)
	crawler.Crawl()

}
