package main

import (
	"encoding/json"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/julienschmidt/httprouter"
	"net/http"
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
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
		w.WriteHeader(http.StatusOK)
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
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
		w.WriteHeader(http.StatusOK)
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
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
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
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
		w.WriteHeader(http.StatusOK)
		w.Write(wrappedData)
	})
}

func RoutesTimeTableHandler(driver bolt.Driver) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		routeID := p.ByName("routeID")
		atStopName := p.ByName("atStopName")
		fromStopName := p.ByName("fromStopName")
		toStopName := p.ByName("toStopName")

		data, err := getTimetable(driver, routeID, atStopName, fromStopName, toStopName)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})
}

func main() {
	driver := openDB()
	router := httprouter.New()
	router.GET("/stops", StopsHandler(driver))
	router.GET("/routes", RoutesHandler(driver))
	router.GET("/routes/variants/id/:routeID", RoutesVariantsByIdHandler(driver))
	router.GET("/routes/variants/stop/:stopName", RoutesVariantsByStopNameHandler(driver))
	router.GET("/route/:routeID/timetable/at/:atStopName/from/:fromStopName/to/:toStopName", RoutesTimeTableHandler(driver))
	http.ListenAndServe(":8080", router)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
