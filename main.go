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

	dockerClient, err := createDockerClient(localStatus.Directory)
	if err != nil {
		log.Fatal(err)
	}
	dockerClient.up()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		dockerClient.down()
		os.Exit(1)
	}()

	dbDriver := openDB()
	if isDatabaseUp() == false {
		log.Print("Timeout reached when waiting for database to come up. Aborting.")
		os.Exit(1)
	}

	checkForUpdatesWrapped := func() { checkForUpdates(&localStatus, &repoStatus, &dockerClient) }

	c := cron.New()
	c.AddFunc("@hourly", checkForUpdatesWrapped)
	c.Start()
	checkForUpdatesWrapped()

	router := createRouter(dbDriver)
	log.Fatal(http.ListenAndServe(":8080", router))
}
