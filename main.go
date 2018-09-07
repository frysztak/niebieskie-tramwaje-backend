package main

import (
	"log"
	"net/http"
)

func main() {
	driver := openDB()
	router := createRouter(driver)
	log.Fatal(http.ListenAndServe(":8080", router))
}
