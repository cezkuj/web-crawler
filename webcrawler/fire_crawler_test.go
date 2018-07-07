package webcrawler

import (
	"testing"
)

func testWrapperFire(domain string, tls bool) {
	matchSubdomains := false
	crawler := NewFireAndForgetCrawler(domain, matchSubdomains, tls)
	crawler.Crawl()

}

func TestFireCrawlerLH(*testing.T) {
	testWrapperFire("lhsystems.pl", true)
}

func TestFireCrawlerCampoy(*testing.T) {
	testWrapperFire("campoy.cat", true)
}

func TestFireCrawlerDoug(*testing.T) {
	testWrapperFire("doughellmann.com", true)
}

func TestFireCrawlerBourgon(*testing.T) {
	testWrapperFire("peter.bourgon.org", false)
}

func TestFireCrawlerRakyll(*testing.T) {
	testWrapperFire("rakyll.org", true)
}
