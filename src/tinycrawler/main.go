package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// fetchSite fetches a single URL and returns its content as a strong
func fetchSite(site string) (string, error) {
	// GET the URL
	resp, err := http.Get(site)
	if err != nil {
		return "", err
	}

	// read the Body
	pageText, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	// return Content
	return string(pageText), nil
}

// crawl begins crawling a URL (ad infinitum)
func crawl(site string) {
	// hash of URLs we've already crawled
	sitesCrawled := make(map[string]struct{})

	// slice of sites we want to crawl
	sitesToCrawl := []string{site}

	// loop as long as we have sites left to crawl
	for len(sitesToCrawl) > 0 {
		// pop the next site to crawl
		var curSite string
		curSite, sitesToCrawl = sitesToCrawl[len(sitesToCrawl)-1], sitesToCrawl[:len(sitesToCrawl)-1]

		// add it to our crawled sites hash
		sitesCrawled[curSite] = struct{}{}

		// fetch page
		fmt.Println("CRAWLING: ", curSite)
		pageText, err := fetchSite(curSite)
		if err != nil {
			fmt.Println("got error", err)
			continue
		}

		// parse page, add all "new" links to the sitesToCrawl slice
		_, links := parseLinks(pageText)
		for _, l := range links {
			if _, ok := sitesCrawled[l]; !ok {
				// new site found, add it to sitesToCrawl
				sitesToCrawl = append(sitesToCrawl, l)
				fmt.Println("FOUND NEW LINK TO CRAWL: ", l)
			}
		}
	}
}

func main() {
	crawl("http://reddit.com/")
}
