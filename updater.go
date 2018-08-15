package main

import (
	"log"
	"time"
)

func checkForUpdates(localStatus *Status, repoStatus *RepositoryStatus, client *DockerClient) {
	log.Print("Checking for updates...")
	err := getRepositoryStatus(repoStatus)
	if err != nil {
		log.Printf("Getting repository status failed. %s", err)
		return
	}

	if repoStatus.ModifiedDate.After(localStatus.ModifiedDate) {
		log.Print("Database update triggered...")
		update(localStatus, repoStatus, client)
	} else {
		log.Print("Nothing to update")
	}
}

func update(localStatus *Status, repoStatus *RepositoryStatus, client *DockerClient) error {
	checksum, err := downloadZip(repoStatus.ZipURL, localStatus.Directory)
	if err != nil {
		log.Printf("Update operation failed with error: %s", err)
		return err
	}

	if localStatus.ZipChecksum == checksum {
		log.Print("Checksums match; file didn't actually change. Aborting.")
		localStatus.LastCheckedDate = time.Now()
		localStatus.save()
		return nil
	}

	localStatus.LastCheckedDate = time.Now()
	localStatus.ModifiedDate = repoStatus.ModifiedDate
	localStatus.ZipChecksum = checksum
	localStatus.save()

	return nil
}
