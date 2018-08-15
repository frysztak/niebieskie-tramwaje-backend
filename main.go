package main

import (
	"github.com/robfig/cron"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var localStatus Status
	var repoStatus RepositoryStatus

	err := loadStatus(&localStatus)
	if err != nil {
		log.Fatal(err)
	}

	client, err := createDockerClient(localStatus.Directory)
	if err != nil {
		log.Fatal(err)
	}
	client.up()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		client.down()
		os.Exit(1)
	}()

	checkForUpdatesWrapped := func() { checkForUpdates(&localStatus, &repoStatus, &client) }

	c := cron.New()
	c.AddFunc("@hourly", checkForUpdatesWrapped)
	c.Start()
	checkForUpdatesWrapped()

	driver := openDB()
	router := createRouter(driver)
	log.Fatal(http.ListenAndServe(":8080", router))
}
