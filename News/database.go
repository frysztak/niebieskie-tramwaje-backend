package News

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os/user"
	"path"
)

const schema = `
CREATE TABLE IF NOT EXISTS news (
	url TEXT PRIMARY KEY UNIQUE,
	title TEXT,
	published_on DATETIME,
	synopsis TEXT,
	affects_lines TEXT,
	affects_days TEXT,
	body TEXT
);
`

func OpenDatabase() *sqlx.DB {
	log.Println("Opening database...")

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dbPath := path.Join(usr.HomeDir, "MPK-API", "storage.db")
	log.Printf("DB path: %s", dbPath)

	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	return db
}

func insertNewsIntoDB(db *sqlx.DB, news []NewsItem) {
	tx := db.MustBegin()
	for _, newsItem := range news {
		tx.NamedExec(`
			INSERT OR IGNORE INTO news (url, title, published_on, synopsis, affects_lines, affects_days, body)
			VALUES (:url, :title, :published_on, :synopsis, :affects_lines, :affects_days, :body)`,
			&newsItem)
	}
	err := tx.Commit()

	if err != nil {
		log.Panic("DB commit failed: %s", err)
	}
	log.Println("Commited to DB")
}

const itemsPerPage = 10

func getNews(db *sqlx.DB, limit int, page int) []NewsItem {
	log.Printf("Retrieving news, page %d, limit %d", page, limit)
	news := []NewsItem{}

	if db == nil {
		panic("DB is nil")
		return news
	}

	offset := page * itemsPerPage
	db.Select(&news, `
		SELECT * FROM news
		ORDER BY published_on DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	return news
}
