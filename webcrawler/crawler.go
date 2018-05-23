package webcrawler

import (
	"golang.org/x/net/html"
	"sync"
)

// Crawl should return *sync.Map consisting key - page name and value - higher-level page ("foundOn" page)
type Crawler interface {
	Crawl() *sync.Map
}

type Page struct {
	name    string
	content *html.Node
}

//PageTree is struct created to simplify printing
type PageTree struct {
	name     string
	subPages map[string]PageTree
}
