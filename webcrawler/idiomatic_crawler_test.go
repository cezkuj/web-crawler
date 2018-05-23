package webcrawler

import (
	"log"
	"strconv"
	"testing"
	"time"
)

func testWrapperIdiomatic(domain string, intervalMillis int, tls bool) {
	matchSubdomains := false
	reqInterval, err := time.ParseDuration(strconv.Itoa(intervalMillis) + "ms")
	if err != nil {
		log.Fatal(err)
	}
	crawler := NewIdiomaticCrawler(domain, matchSubdomains, reqInterval, tls)
	crawler.Crawl()

}
func TestIdiomaticCrawlerLH(*testing.T) {
	testWrapperIdiomatic("lhsystems.pl", 1, true)
}

func TestIdiomaticCrawlerMonzo(*testing.T) {
	testWrapperIdiomatic("monzo.com", 0, true)
}

func TestIdiomaticCrawlerCampoy(*testing.T) {
	testWrapperIdiomatic("campoy.cat", 0, true)
}

func TestIdiomaticCrawlerDoug(*testing.T) {
	testWrapperIdiomatic("doughellmann.com", 10, true)
}

func TestIdiomaticCrawlerPB(*testing.T) {
	testWrapperIdiomatic("pbpython.com", 0, false)
}

func TestIdiomaticCrawlerBourgon(*testing.T) {
	testWrapperIdiomatic("peter.bourgon.org", 0, false)
}

func TestIdiomaticCrawlerRakyll(*testing.T) {
	testWrapperIdiomatic("rakyll.org", 0, true)
}
