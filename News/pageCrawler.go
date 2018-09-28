package News

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

const baseUrl = "http://mpk.wroc.pl"

func crawlPage(url string, ch chan NewsItem, chFinished chan bool) {
	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	// Request the HTML page.
	log.Printf("Getting url %s", url)
	res, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Print(err)
	}

	doc.Find(".info").Each(func(i int, s *goquery.Selection) {
		aTag := s.Find(".title a")
		newsUrl, linkExists := aTag.Attr("href")
		if linkExists == false {
			log.Printf("news url doesnt exist at block %d at page %s", i, url)
		}

		teaser := s.Find(".teaser")

		if linkExists && teaser.Length() == 1 {
			var newsItem NewsItem
			newsUrl = fmt.Sprintf("%s%s", baseUrl, newsUrl)
			newsItem.Url = newsUrl
			newsItem.Synopsis = teaser.First().Text()
			ch <- newsItem
		}
	})
}

const nPages = 1

func getNewsStubs(seedUrl string) []NewsItem {
	newsStubs := make([]NewsItem, 0)
	seedUrls := make([]string, nPages)
	for idx := 0; idx < nPages; idx++ {
		seedUrls[idx] = fmt.Sprintf(seedUrl, idx)
	}

	// Channels
	chNewsStubs := make(chan NewsItem)
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	for _, url := range seedUrls {
		go crawlPage(url, chNewsStubs, chFinished)
	}

	// Subscribe to both channels
	for c := 0; c < len(seedUrls); {
		select {
		case newsStub := <-chNewsStubs:
			newsStubs = append(newsStubs, newsStub)
		case <-chFinished:
			c++
		}
	}

	log.Printf("Finished shallow crawling, got %d news stubs", len(newsStubs))
	return newsStubs
}
