package main

import (
	"encoding/json"
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var ( // TODO: make const somehow
	maxAge            = 60 * 60 * 24 * 7 // 7 days
	cacheControlValue = fmt.Sprintf("max-age:%d, public", maxAge)
)

func wrapJSON(name string, item interface{}) ([]byte, error) {
	wrapped := map[string]interface{}{
		name: item,
	}
	return json.Marshal(wrapped)
}

func StopsHandler(driver bolt.Driver) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := getAllStopNames(driver)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		wrappedData, err := wrapJSON("stopNames", data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", cacheControlValue)
		w.WriteHeader(http.StatusCreated)
		w.Write(wrappedData)
	})
}

func RoutesHandler(driver bolt.Driver) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := getAllRouteIDs(driver)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// TODO: wrap data
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", cacheControlValue)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	})
}

func RoutesVariantsByIdHandler(driver bolt.Driver) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		routeID := p.ByName("routeID")
		data, err := getRouteVariantsForRouteID(driver, routeID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		wrappedData := map[string][]RouteVariant{"RouteVariants": data}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", cacheControlValue)
		json.NewEncoder(w).Encode(wrappedData)
	})
}

func RoutesVariantsByStopNameHandler(driver bolt.Driver) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		stopName := p.ByName("stopName")
		data, err := getRouteVariantsByStopName(driver, stopName)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		wrappedData, err := wrapJSON("routeVariants", data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", cacheControlValue)
		w.WriteHeader(http.StatusCreated)
		w.Write(wrappedData)
	})
}

func main() {
	driver := openDB()
	router := httprouter.New()
	router.GET("/stops", StopsHandler(driver))
	router.GET("/routes", RoutesHandler(driver))
	router.GET("/routes/variants/id/:routeID", RoutesVariantsByIdHandler(driver))
	router.GET("/routes/variants/stop/:stopName", RoutesVariantsByStopNameHandler(driver))
	http.ListenAndServe(":8080", router)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
