package main

import (
	"testing"
)

func TestIdiomaticCrawler(*testing.T) {
	domain := "lhsystems.pl"
	matchSubdomains := true
	crawler := NewIdiomaticCrawler(domain, matchSubdomains)
	crawler.Crawl()

}
