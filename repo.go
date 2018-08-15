package main

import (
	"errors"
	"fmt"
	"github.com/deiu/rdf2go"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type RepositoryStatus struct {
	ModifiedDate time.Time
	ZipURL       string
}

func getRepositoryStatus(status *RepositoryStatus) error {
	uri := "https://www.wroclaw.pl/open-data/dataset/rozkladjazdytransportupublicznegoplik_data.jsonld"
	g := rdf2go.NewGraph(uri, true)

	log.Print("Querying OpenData for JSON-LD...")
	err := g.LoadURI(uri)
	if err != nil {
		return err
	}

	modifiedTerm := "http://purl.org/dc/terms/modified"

	if g.Len() == 0 {
		err = errors.New("Received graph is empty")
		return err
	}

	triple := g.One(nil, rdf2go.NewResource(modifiedTerm), nil)
	if triple == nil {
		err = errors.New("`modified` term not found")
		return err
	}

	modifiedDateS := triple.Object.String()                                                           // "2018-08-08T08:19:12.979963"^^<http://www.w3.org/2001/XMLSchema#dateTime>
	modifiedDateS = strings.Replace(modifiedDateS, `"`, "", 2)                                        // 2018-08-08T08:19:12.979963^^<http://www.w3.org/2001/XMLSchema#dateTime>
	modifiedDateS = strings.TrimRight(modifiedDateS, "^^<http://www.w3.org/2001/XMLSchema#dateTime>") // 2018-08-08T08:19:12.979963

	layout := "2006-01-02T15:04:05.999999"
	modifiedDate, err := time.Parse(layout, modifiedDateS)
	if err != nil {
		err = fmt.Errorf("Time parsing failed! Layout `%s` didn't match time `%s`", layout, modifiedDateS)
		return err
	}

	accessUrlTerm := "http://www.w3.org/ns/dcat#accessURL"
	triple = g.One(nil, rdf2go.NewResource(accessUrlTerm), nil)
	if triple == nil {
		err = errors.New("`accessURL` term not found")
		return err
	}

	accessURL := triple.Object.String()
	accessURL = strings.Replace(accessURL, "<", "", 1)
	accessURL = strings.Replace(accessURL, ">", "", 1)

	status.ModifiedDate = modifiedDate
	status.ZipURL = accessURL

	return nil
}

func downloadZip(zipUrl string, destination string) error {
	filename := filepath.Base(zipUrl)
	//path := filepath.Join(destination, filename)

	log.Printf("Downloading %s into %s", filename, destination)
	//err := downloadFile(path, zipUrl)
	//if err == nil {
	//	log.Printf("Download finished successfully")
	//}
	//return err
	return nil
}
