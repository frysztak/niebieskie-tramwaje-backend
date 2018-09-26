package main

import (
	"./GTFS"
	"log"
	"net/http"
)

func main() {
	driver := GTFS.OpenDB()
	router := GTFS.CreateRouter(driver)
	log.Fatal(http.ListenAndServe(":8080", router))
}
