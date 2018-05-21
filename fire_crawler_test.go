package main

import (
	"testing"
)

func TTestFireAndForgetCrawler(*testing.T) {
	domain := "lhsystems.pl"
	matchSubdomains := false
	crawler := NewFireAndForgetCrawler(domain, matchSubdomains)
	crawler.Crawl()

}
