package main

import (
	"sync"
	"testing"
)

func TestKeyInMap(t *testing.T) {
	m := &sync.Map{}
	m.Store("test", true)
	got := keyInMap("test", *m)
	want := true
	if got != want {
		t.Errorf("TestKeyInMap failed got: %v, want: %v", got, want)
	}
	got = keyInMap("test2", *m)
	want = false
	if got != want {
		t.Errorf("TestKeyInMap failed got: %v, want: %v", got, want)
	}
}

func TestGetDomain(t *testing.T) {
	page := "https://github.com/cezkuj/web-crawler/blob/master/crawl.go"
	got := getDomain(page)
	want := "github.com"
	if got != want {
		t.Errorf("TestGetDomain failed got: %v, want: %v", got, want)
	}
}

func TestCheckDomain(t *testing.T) {
	page := "https://community.monzo.com"
	domain := "monzo.com"
	got := checkDomain(page, domain)
	want := true
	if got != want {
		t.Errorf("TestCheckDomain failed got: %v, want: %v", got, want)
	}
	page = "https://facebook.com?url=monzo.com"
	got = checkDomain(page, domain)
	want = false
	if got != want {
		t.Errorf("TestCheckDomain failed got: %v, want: %v", got, want)
	}
}

func TestBuildUrl(t *testing.T) {
	path := "/about"
	url := "https://github.com/blog"
	got := buildUrl(url, path)
	want := "https://github.com/about"
	if got != want {
		t.Errorf("TestBuildUrl failed got: %v, want: %v", got, want)
	}
	path = "about"
	got = buildUrl(url, path)
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
