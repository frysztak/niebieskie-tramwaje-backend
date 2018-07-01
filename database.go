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

func openDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable",
		DB_USER, DB_NAME)
	log.Print("Opening database...")
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	return db
}

func setGraphPath(db *sql.DB) {
	q := fmt.Sprintf(`SET graph_path = %s`, GRAPH_NAME)
	_, err := db.Exec(q)
	checkErr(err)
}

func getAllStopNames(db *sql.DB) []string {
	rows, err := db.Query(getAllStopNamesQuery)
	checkErr(err)

	stop_names := make([]string, 0)
	for rows.Next() {
		var stop_name string
		err = rows.Scan(&stop_name)
		checkErr(err)
		stop_names = append(stop_names, stop_name)
	}
	log.Printf(`Received %d stop names`, len(stop_names))
	return stop_names
}

type Route struct {
	id    string
	isBus bool
}

func getAllRouteIDs(db *sql.DB) []Route {
	rows, err := db.Query(getAllRouteIDsQuery)
	checkErr(err)

	routes := make([]Route, 0)
	for rows.Next() {
		var route_id string
		var is_bus bool
		err = rows.Scan(&route_id, &is_bus)
		checkErr(err)
		routes = append(routes, Route{route_id, is_bus})
	}
	log.Printf(`Received %d route IDs`, len(routes))
	return routes
}

type RouteVariant struct {
	firstStop string
	lastStop  string
	tripIDs   []string
}

func getVariantsForRouteID(db *sql.DB, routeID string) []RouteVariant {
	q := fmt.Sprintf(getVariantsForRouteIDQuery, routeID)
	rows, err := db.Query(q)

	stopNamesMap := make(map[string][]string)
	for rows.Next() {
		var route_id string
		var stop_names string
		err = rows.Scan(&route_id, &stop_names)
		checkErr(err)

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
	return variants
}
