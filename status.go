package main

import (
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"github.com/tevino/abool"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Status struct {
	ModifiedDate    time.Time         `json:"modifiedDate"`
	LastCheckedDate time.Time         `json:"lastCheckedDate"`
	ZipChecksum     string            `json:"checksum"`
	Directory       string            `json:"-"`
	DatabaseBusy    *abool.AtomicBool `json:"-"`
}

const (
	statusFileName = "status.json"
)

func (s Status) save() error {
	json, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	statusPath := filepath.Join(s.Directory, statusFileName)
	return ioutil.WriteFile(statusPath, json, 0644)
}

func createDefaultStatusFile(dir string) error {
	status := Status{time.Unix(0, 0), time.Unix(0, 0), "", dir, nil}
	return status.save()
}

func loadStatus(status *Status) error {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Print("Homedir not found. %s", err)
		return err
	}

	mpkDir := filepath.Join(homeDir, "MPK")

	if _, err := os.Stat(mpkDir); os.IsNotExist(err) {
		log.Print("MPK dir doesn't exist")
		return err
	}

	statusPath := filepath.Join(mpkDir, statusFileName)
	log.Printf("Looking for file %s...", statusPath)

	if _, err := os.Stat(statusPath); os.IsNotExist(err) {
		log.Print("Status file doesn't exist, creating one")
		err = createDefaultStatusFile(mpkDir)
		if err != nil {
			return err
		}
	}

	jsonFile, err := os.Open(statusPath)
	if err != nil {
		log.Printf("Can't open status file. %s", err)
		return err
	}
	defer jsonFile.Close()
	log.Print("Found it!")

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Print("Failed to read status file. %s", err)
		return err
	}

	err = json.Unmarshal(byteValue, status)
	if err != nil {
		log.Print("Failed to parse status file. %s", err)
		return err
	}

	status.Directory = mpkDir
	status.DatabaseBusy = abool.New()

	return nil
}
