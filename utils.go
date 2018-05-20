package main

import (
	"net/url"
        "log"
        "sync"
        "strings"
)

func keyInMap(key string, m sync.Map) bool {
	_, ok := m.Load(key)
	return ok
}
func getDomain(page string) string {
	u, err := url.Parse(page)
	if err != nil {
		log.Fatal(err)
	}
	hostname := u.Hostname()
	if strings.HasPrefix(hostname, "www.") {
		return hostname[4:]
	}
	return hostname
}
func checkDomain(page, domain string, matchSubdomains bool) bool {
        if matchSubdomains {
	 return strings.Contains(getDomain(page), domain)
        }
        return getDomain(page) == domain
}

func buildUrl(foundOn, relSuffix string) string {
	if strings.HasPrefix(relSuffix, "/") {
		return "https://" + getDomain(foundOn) + relSuffix
	}
        return foundOn + "/" + relSuffix
}
func removeStringAfterChar(str, ch string) string {
	if i := strings.Index(str, ch); i != -1 {
		return str[:i]
	}
	return str

}
func removeGetParams(u string) string {
	return removeStringAfterChar(u, "?")

}
func removeChapterLinks(u string) string {
	return removeStringAfterChar(u, "#")
}

func printMap(key, value interface{}) bool {
	log.Println(key, value)
        return true
}
