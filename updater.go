package main

import (
	"log"
	"time"
)

func checkForUpdates(localStatus *Status, repoStatus *RepositoryStatus) {
	err := getRepositoryStatus(repoStatus)
	if err != nil {
		log.Printf("Getting repository status failed. %s", err)
		return
	}

	if repoStatus.ModifiedDate.After(localStatus.ModifiedDate) {
		log.Print("DB update triggered...")
		err := downloadZip(repoStatus.ZipURL, localStatus.Directory)
		if err != nil {
			log.Fatal(err)
		}

		localStatus.LastCheckedDate = time.Now()
		localStatus.ModifiedDate = repoStatus.ModifiedDate
		// save localStatus
	}

}
