package main

import (
	"github.com/robfig/cron"
	"log"
	"net/http"
)

func main() {
	var localStatus Status
	var repoStatus RepositoryStatus

	err := loadStatus(&localStatus)
	if err != nil {
		log.Fatal(err)
	}

	checkForUpdatesWrapped := func() { checkForUpdates(&localStatus, &repoStatus) }

	c := cron.New()
	c.AddFunc("@hourly", checkForUpdatesWrapped)
	c.Start()
	checkForUpdatesWrapped()

	driver := openDB()
	router := createRouter(driver)
	log.Fatal(http.ListenAndServe(":8080", router))
}
