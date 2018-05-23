package webcrawler

import (
	"testing"
)

func TestGetDomain(t *testing.T) {
	testCases := []struct {
		page     string
		expected string
	}{
		{"https://github.com/cezkuj/web-crawler/blob/master/crawl.go", "github.com"},
		{"https://www.reddit.com", "reddit.com"},
	}
	for _, item := range testCases {
		if actual, err := getDomain(item.page); actual != item.expected || err != nil {
			t.Errorf("Test failed actual: %v, expected: %v, err: %v", actual, item.expected, err)

		}
	}
}

func TestCheckDomain(t *testing.T) {
	testCases := []struct {
		page            string
		domain          string
		matchSubdomains bool
		expected        bool
	}{
		{"https://community.monzo.com", "monzo.com", true, true},
		{"https://community.monzo.com", "monzo.com.com", false, false},
		{"https://facebook.com?url=monzo.com", "monzo.com.com", false, false},
		{"https://facebook.com?url=monzo.com", "monzo.com.com", true, false},
	}
	for _, item := range testCases {
		if actual := checkDomain(item.page, item.domain, item.matchSubdomains); actual != item.expected {
			t.Errorf("Test failed actual: %v, expected: %v", actual, item.expected)
		}
	}
}
func TestBuildURL(t *testing.T) {
	testCases := []struct {
		url      string
		path     string
		tls      bool
		expected string
	}{
		{"https://github.com/blog", "/about", true, "https://github.com/about"},
		{"http://github.com/blog", "/about", false, "http://github.com/about"},
		{"http://github.com/blog", "about", false, "http://github.com/blog/about"},
		{"https://github.com/blog", "about", true, "https://github.com/blog/about"},
	}
	for _, item := range testCases {
		if actual := buildURL(item.url, item.path, item.tls); actual != item.expected {
			t.Errorf("Test failed actual: %v, expected: %v", actual, item.expected)
		}
	}
}

func TestRemoveStringAfterChar(t *testing.T) {
	str := "abcdefg@abcd"
	char := "@"
	actual := removeStringAfterChar(str, char)
	expected := "abcdefg"
	if actual != expected {
		t.Errorf("Test failed actual: %v, expected: %v", actual, expected)
	}
}

func TestRemoveGetParams(t *testing.T) {
	url := "manageproducts.do?option=1"
	actual := removeGetParams(url)
	expected := "manageproducts.do"
	if actual != expected {
		t.Errorf("Test failed actual: %v, expected: %v", actual, expected)
	}
}

func TestRemoveChapterLinks(t *testing.T) {
	url := "https://golang.org/pkg/sync/#Map.Delete"
	actual := removeChapterLinks(url)
	expected := "https://golang.org/pkg/sync/"
	if actual != expected {
		t.Errorf("Test failed actual: %v, expected: %v", actual, expected)
	}

}
