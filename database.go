package main

import (
	"database/sql"
	"fmt"
	_ "github.com/bitnine-oss/agensgraph-golang"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

const (
	DB_USER     = "sebastian"
	DB_PASSWORD = ""
	DB_NAME     = "wroclaw_gtfs"
	GRAPH_NAME  = "wroclaw"
)

func openDB() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable",
		DB_USER, DB_NAME)
	log.Print("Opening database...")
	return sql.Open("postgres", dbinfo)
}

func setGraphPath(db *sql.DB) error {
	q := fmt.Sprintf(`SET graph_path = %s`, GRAPH_NAME)
	_, err := db.Exec(q)
	return err
}

func getAllStopNames(db *sql.DB) ([]string, error) {
	rows, err := db.Query(getAllStopNamesQuery)
	if err != nil {
		return nil, err
	}

	stop_names := make([]string, 0)
	for rows.Next() {
		var stop_name string
		err = rows.Scan(&stop_name)
		if err != nil {
			return nil, err
		}

		stop_name = strings.Replace(stop_name, `"`, ``, -1)
		stop_names = append(stop_names, stop_name)
	}
	log.Printf(`Received %d stop names`, len(stop_names))
	return stop_names, nil
}

type Route struct {
	ID    string
	IsBus bool
}

func getAllRouteIDs(db *sql.DB) ([]Route, error) {
	rows, err := db.Query(getAllRouteIDsQuery)
	if err != nil {
		return nil, err
	}

	routes := make([]Route, 0)
	for rows.Next() {
		var route_id string
		var is_bus bool
		err = rows.Scan(&route_id, &is_bus)
		if err != nil {
			return nil, err
		}
		route_id = strings.Replace(route_id, `"`, ``, -1)
		routes = append(routes, Route{route_id, is_bus})
	}
	log.Printf(`Received %d route IDs`, len(routes))
	log.Println(routes[0])
	return routes, nil
}

type RouteVariant struct {
	FirstStop string
	LastStop  string
	TripIDs   []string
}

func getVariantsForRouteID(db *sql.DB, routeID string) ([]RouteVariant, error) {
	q := fmt.Sprintf(getVariantsForRouteIDQuery, routeID)
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	stopNamesMap := make(map[string][]string)
	for rows.Next() {
		var route_id string
		var stop_names string
		err = rows.Scan(&route_id, &stop_names)
		if err != nil {
			return nil, err
		}

		route_id = strings.Replace(route_id, `"`, ``, -1)
		if val, ok := stopNamesMap[stop_names]; ok {
			stopNamesMap[stop_names] = append(val, route_id)
		} else {
			stopNamesMap[stop_names] = []string{route_id}
		}
	}

	var variants []RouteVariant
	for stop_names, tripIDs := range stopNamesMap {
		stop_names = stop_names[1 : len(stop_names)-2] // remove trailing '[' ']'
		s := strings.Split(stop_names, ",")

		processedStopNames := []string{}
		for _, stop_name := range s {
			stop_name = strings.Replace(strings.TrimSpace(stop_name), `"`, ``, -1)
			processedStopNames = append(processedStopNames, stop_name)
		}
		first := processedStopNames[0]
		last := processedStopNames[len(processedStopNames)-1]
		variants = append(variants, RouteVariant{first, last, tripIDs})
	}

	log.Printf(`Received %d variants for routeID "%s"`, len(variants), routeID)
	return variants, nil
}
