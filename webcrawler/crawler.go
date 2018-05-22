package webcrawler

import (
	"golang.org/x/net/html"
	"sync"
)

type Crawler interface {
	Crawl() sync.Map
}

type Page struct {
	name    string
	content *html.Node
}

type PageTree struct {
	name     string
	subPages map[string]PageTree
}
