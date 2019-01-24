package News

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

func getSleepTime() time.Duration {
	r := rand.Intn(1000)
	return 500 + time.Duration(r)*time.Millisecond
}

func crawlNews(client http.Client, newsStub *NewsItem, chFinished chan bool) {
	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	url := newsStub.Url

	// Request the HTML page.
	log.Printf("Getting url %s", url)
	time.Sleep(getSleepTime())
	res, err := client.Get(url)
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

	pageTitle := doc.Find(".page-title")
	publishedOn := pageTitle.Find(".timestamp").Text()
	title := pageTitle.Find("h1").Text()

	infoParagraphs := doc.Find(".paragraph-info")
	if infoParagraphs.Length() != 1 {
		log.Panicf("Expected 1 paragraph-info for url %s", url)
	}
	affectsLines, affectsDays := parseInfoParagraph(url, infoParagraphs)

	body, bodyErr := doc.Find(".field-item").First().Html()
	if bodyErr != nil {
		log.Printf("Body is nil for url %s", url)
	}

	newsStub.Title = cleanUpTitle(title)
	newsStub.PublishedOn = parsePublishedDateTime(publishedOn)
	newsStub.AffectsDay = affectsDays
	newsStub.AffectsLines = affectsLines
	newsStub.Body = fixImageUrls(body)

	log.Printf("Found news article. Title: '%s', publishedOn: '%s', affectsLines: '%s', affectsDays: '%s'", title, publishedOn, affectsLines, affectsDays)
}

func parseInfoParagraph(url string, paragraph *goquery.Selection) (string, string) {
	var affectsLines string
	var affectsDays string
	text := paragraph.Text()
	text = strings.TrimSpace(text)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if strings.Contains(line, "Dotyczy linii") {
			parts := strings.Split(line, ":")
			if len(parts) == 0 {
				log.Panicf("affectsLine -> len(parts) == 0 for url %s", url)
			}

			affectsLines = parts[len(parts)-1]
		} else if strings.Contains(line, "Obowiązuje w dniach") {
			parts := strings.Split(line, ":")
			if len(parts) == 0 {
				log.Panicf("affectsDays -> len(parts) == 0 for url %s", url)
			}

			affectsDays = parts[len(parts)-1]
		}
	}

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

func parsePublishedDateTime(publishedOn string) time.Time {
	layout := "02.01.2006 15:04"
	t, err := time.Parse(layout, publishedOn)
	if err != nil {
		log.Println(err)
	}

	return t
}

// https://stackoverflow.com/a/24894202
func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func fillOutNewsStubs(client http.Client, newsStubs []NewsItem) {
	// Channels
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	for idx, _ := range newsStubs {
		go crawlNews(client, &newsStubs[idx], chFinished)
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
