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
		log.Print("Database update triggered...")
		update(localStatus, repoStatus)
	}

}

func update(localStatus *Status, repoStatus *RepositoryStatus) error {
	checksum, err := downloadZip(repoStatus.ZipURL, localStatus.Directory)
	if err != nil {
		log.Printf("Update operation failed with error: %s", err)
		return err
	}

	localStatus.LastCheckedDate = time.Now()
	localStatus.ModifiedDate = repoStatus.ModifiedDate
	localStatus.ZipChecksum = checksum
	localStatus.save()

	return nil
}
