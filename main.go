package main

import (
	"encoding/json"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func StopsHandler(driver bolt.Driver) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := getAllStopNames(driver)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	})
}

func RoutesHandler(driver bolt.Driver) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := getAllRouteIDs(driver)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	})
}

func RoutesVariantsHandler(driver bolt.Driver) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		routeID := p.ByName("routeID")
		data, err := getVariantsForRouteID(driver, routeID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	})
}

func main() {
	driver := openDB()
	router := httprouter.New()
	router.GET("/stops", StopsHandler(driver))
	router.GET("/routes", RoutesHandler(driver))
	router.GET("/routes/:routeID/variants", RoutesVariantsHandler(driver))
	http.ListenAndServe(":8080", router)

	//getAllStopNames(db)
	//getAllRouteIDs(db)
	//getVariantsForRouteID(db, "122")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
