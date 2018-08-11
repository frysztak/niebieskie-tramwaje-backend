package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func wrapJSON(name string, item interface{}) ([]byte, error) {
	wrapped := map[string]interface{}{
		name: item,
	}
	return json.Marshal(wrapped)
}

type Handler func(w http.ResponseWriter, r *http.Request)

func StopsHandler(driver bolt.Driver) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func StopsAndRoutesHandler(driver bolt.Driver) Handler {
	return func(w http.ResponseWriter, r *http.Request) {

		type StopsAndRoutes struct {
			Stops  []string
			Routes []Route
		}

		stops, err := getAllStopNames(driver)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		routes, err := getAllRouteIDs(driver)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := StopsAndRoutes{stops, routes}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

func RoutesHandler(driver bolt.Driver) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func RoutesVariantsByIdHandler(driver bolt.Driver) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		routeID, err := url.QueryUnescape(vars["routeID"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data, err := getRouteVariantsForRouteID(driver, routeID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		wrappedData := map[string][]RouteVariant{"RouteVariants": data}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
		json.NewEncoder(w).Encode(wrappedData)
	}
}

func RoutesVariantsByStopNameHandler(driver bolt.Driver) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		stopName, err := url.QueryUnescape(vars["stopName"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

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
	}
}

func RoutesTimeTableHandler(driver bolt.Driver) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		routeID, err := url.QueryUnescape(vars["routeID"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		stopName, err := url.QueryUnescape(vars["stopName"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		direction, err := url.QueryUnescape(vars["direction"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data, err := getTimetable(driver, routeID, stopName, direction)
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
	}
}

func RouteInfoHandler(driver bolt.Driver) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		routeID, err := url.QueryUnescape(vars["routeID"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data, err := getRouteInfo(driver, routeID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

func RouteDirectionsHandler(driver bolt.Driver) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		routeID, err := url.QueryUnescape(vars["routeID"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data, err := getRouteDirections(driver, routeID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

func RouteDirectionsThroughStopHandler(driver bolt.Driver) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		routeID, err := url.QueryUnescape(vars["routeID"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		stopName, err := url.QueryUnescape(vars["stopName"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data, err := getRouteDirectionsThroughStop(driver, routeID, stopName)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

func RouteStopsHandler(driver bolt.Driver) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		routeID, err := url.QueryUnescape(vars["routeID"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data, err := getStopsForRouteID(driver, routeID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

func TripTimelineHandler(driver bolt.Driver) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tripIDString, err := url.QueryUnescape(vars["tripID"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		tripID, err := strconv.ParseInt(tripIDString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data, err := getTripTimeline(driver, int(tripID))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Expires", "Wed, 21 Oct 2020 07:28:00 GMT") //TODO: dynamically read actual date from DB
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

func main() {
	driver := openDB()
	router := mux.NewRouter().UseEncodedPath()
	router.HandleFunc("/stops", StopsHandler(driver))
	router.HandleFunc("/stops/and/routes", StopsAndRoutesHandler(driver))
	router.HandleFunc("/routes", RoutesHandler(driver))
	router.HandleFunc("/routes/variants/id/{routeID}", RoutesVariantsByIdHandler(driver))
	router.HandleFunc("/routes/variants/stop/{stopName}", RoutesVariantsByStopNameHandler(driver))
	router.HandleFunc("/route/{routeID}/timetable/at/{stopName}/direction/{direction}", RoutesTimeTableHandler(driver))
	router.HandleFunc("/route/{routeID}/info", RouteInfoHandler(driver))
	router.HandleFunc("/route/{routeID}/directions", RouteDirectionsHandler(driver))
	router.HandleFunc("/route/{routeID}/directions/through/{stopName}", RouteDirectionsThroughStopHandler(driver))
	router.HandleFunc("/route/{routeID}/stops", RouteStopsHandler(driver))
	router.HandleFunc("/trip/{tripID}/timeline", TripTimelineHandler(driver))
	http.ListenAndServe(":8080", router)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
