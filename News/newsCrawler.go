package News

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

func crawlNews(newsStub *NewsItem, chFinished chan bool) {
	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	url := newsStub.Url

	// Request the HTML page.
	log.Printf("Getting url %s", url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	pageTitle := doc.Find(".page-title")
	publishedOn := pageTitle.Find(".timestamp").Text()
	title := pageTitle.Find("h1").Text()

	infoParagraphs := doc.Find(".paragraph-info p")
	if infoParagraphs.Length() == 0 {
		log.Fatalf("No info paragraphs found for url %s", url)
	}
	affectsLines, affectsDays := parseInfoParagraph(url, infoParagraphs)

	body, bodyErr := doc.Find(".field-item").First().Html()
	if bodyErr != nil {
		log.Fatalf("Body is nil for url %s", url)
	}

	newsStub.Title = cleanUpTitle(title)
	newsStub.PublishedOn = publishedOn
	newsStub.AffectsDay = affectsDays
	newsStub.AffectsLines = affectsLines
	newsStub.Body = fixImageUrls(body)

	log.Printf("Found news article. Title: '%s', publishedOn: '%s', affectsLines: '%s', affectsDays: '%s'", title, publishedOn, affectsLines, affectsDays)
}

func parseInfoParagraph(url string, paragraphs *goquery.Selection) (string, string) {
	var affectsLines string
	var affectsDays string
	paragraphs.Each(func(i int, s *goquery.Selection) {
		text, err := s.Html()
		if err != nil {
			log.Fatalf("paragraph -> Html is nil for url %s", url)
		}
		if strings.Contains(text, "Dotyczy linii") {
			parts := strings.Split(text, ":")
			if len(parts) == 0 {
				log.Fatalf("affectsLine -> len(parts) == 0 for url %s", url)
			}

			affectsLines = parts[len(parts)-1]
		} else if strings.Contains(text, "Obowiązuje w dniach") {
			parts := strings.Split(text, ":")
			if len(parts) == 0 {
				log.Fatalf("affectsDays -> len(parts) == 0 for url %s", url)
			}

			affectsDays = parts[len(parts)-1]
		}
	})

	return strings.TrimSpace(affectsLines), strings.TrimSpace(affectsDays)
}

func fixImageUrls(body string) string {
	pattern := `(?P<tag><img src=\")(?P<url>\/[\w\/\-\_.]+)\"`
	re := regexp.MustCompile(pattern)
	replacement := fmt.Sprintf(`${tag}%s${url}"`, baseUrl)
	return re.ReplaceAllString(body, replacement)
}

func cleanUpTitle(title string) string {
	title = strings.TrimRight(title, ".")
	parts := strings.Split(title, " - ")
	if len(parts) == 1 {
		return title
	}

	title = parts[1]
	upperCaseLetter := unicode.ToUpper(rune(title[0]))
	return replaceAtIndex(title, upperCaseLetter, 0)
}

// https://stackoverflow.com/a/24894202
func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func fillOutNewsStubs(newsStubs []NewsItem) {
	// Channels
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	for idx, _ := range newsStubs {
		go crawlNews(&newsStubs[idx], chFinished)
	}

	// Subscribe to both channels
	for c := 0; c < len(newsStubs); {
		select {
		case <-chFinished:
			c++
		}
	}

	log.Printf("Finished populating news stubs")
}
