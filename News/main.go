package News

import (
	"github.com/jmoiron/sqlx"
	"time"
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
	news := getNewsStubs(seedUrl)
	fillOutNewsStubs(news)
	insertNewsIntoDB(db, news)
}
