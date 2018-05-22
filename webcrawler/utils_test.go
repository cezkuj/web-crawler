package webcrawler

import (
	"testing"
)

func TestGetDomain(t *testing.T) {
	page := "https://github.com/cezkuj/web-crawler/blob/master/crawl.go"
	got, err := getDomain(page)
	want := "github.com"
	if got != want || err != nil {
		t.Errorf("TestGetDomain failed got: %v, want: %v, err: %v", got, want, err)
	}
	page = "https://www.reddit.com"
	got, err = getDomain(page)
	want = "reddit.com"
	if got != want || err != nil {
		t.Errorf("TestGetDomain failed got: %v, want: %v, err: %v", got, want, err)
	}
}

func TestCheckDomain(t *testing.T) {
	matchSubdomains := true
	page := "https://community.monzo.com"
	domain := "monzo.com"
	got := checkDomain(page, domain, matchSubdomains)
	want := true
	if got != want {
		t.Errorf("TestCheckDomain failed got: %v, want: %v", got, want)
	}
	matchSubdomains = false
	got = checkDomain(page, domain, matchSubdomains)
	want = false
	if got != want {
		t.Errorf("TestCheckDomain failed got: %v, want: %v", got, want)
	}
	page = "https://facebook.com?url=monzo.com"
	got = checkDomain(page, domain, matchSubdomains)
	want = false
	if got != want {
		t.Errorf("TestCheckDomain failed got: %v, want: %v", got, want)
	}
	matchSubdomains = true
	got = checkDomain(page, domain, matchSubdomains)
	if got != want {
		t.Errorf("TestCheckDomain failed got: %v, want: %v", got, want)
	}
}

func TestBuildURL(t *testing.T) {
	path := "/about"
	url := "https://github.com/blog"
	got := buildURL(url, path)
	want := "https://github.com/about"
	if got != want {
		t.Errorf("TestBuildUrl failed got: %v, want: %v", got, want)
	}
	path = "about"
	got = buildURL(url, path)
	want = "https://github.com/blog/about"
	if got != want {
		t.Errorf("TestBuildUrl failed got: %v, want: %v", got, want)
	}
}

func TestRemoveStringAfterChar(t *testing.T) {
	str := "abcdefg@abcd"
	char := "@"
	got := removeStringAfterChar(str, char)
	want := "abcdefg"
	if got != want {
		t.Errorf("TestRemoveStringAfterChar failed got: %v, want: %v", got, want)
	}
}

func TestRemoveGetParams(t *testing.T) {
	url := "manageproducts.do?option=1"
	got := removeGetParams(url)
	want := "manageproducts.do"
	if got != want {
		t.Errorf("TestRemoveGetParams failed got: %v, want: %v", got, want)
	}
}

func TestRemoveChapterLinks(t *testing.T) {
	url := "https://golang.org/pkg/sync/#Map.Delete"
	got := removeChapterLinks(url)
	want := "https://golang.org/pkg/sync/"
	if got != want {
		t.Errorf("TestRemoveChapterLinks failed got: %v, want: %v", got, want)
	}
}
