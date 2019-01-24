package News

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type NewsItem struct {
	Url          string    `db:"url"`
	Title        string    `db:"title"`
	PublishedOn  time.Time `db:"published_on"`
	Synopsis     string    `db:"synopsis"`
	AffectsLines string    `db:"affects_lines"`
	AffectsDay   string    `db:"affects_days"`
	Body         string    `db:"body"`
}

func UpdateNews(db *sqlx.DB) {
	seedUrl := "http://mpk.wroc.pl/informacje/zmiany-w-komunikacji?page=%d"

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	news := getNewsStubs(client, seedUrl)
	fillOutNewsStubs(client, news)
	insertNewsIntoDB(db, news)
}
