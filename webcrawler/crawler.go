package webcrawler

import (
	"golang.org/x/net/html"
	"sync"
)

// Crawl should return sync.Map consisting key - page name and value - page wherepage from key where found
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
