package main

import (
	"./GTFS"
	"./News"
	"github.com/robfig/cron"
	"log"
	"net/http"
)

func main() {
	driver := GTFS.OpenDB()
	router := GTFS.CreateRouter(driver)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	newsDb := News.OpenDatabase()

	c := cron.New()
	c.AddFunc("@every 15m", func() { News.UpdateNews(newsDb) })
	c.Start()

	News.UpdateNews(newsDb)

	select {} // "sleep" forever
}
