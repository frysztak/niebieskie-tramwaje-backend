package main

import (
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
)

func main() {
	driver := openDB()
	router := createRouter(driver)
	log.Fatal(http.ListenAndServe(":8080", gziphandler.GzipHandler(router)))
}
