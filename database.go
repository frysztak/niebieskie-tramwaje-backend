package main

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"io"
	"log"
)

const (
	URL = "bolt://neo4j:krowa@localhost:7687"
)

func openDB() bolt.Driver {
	log.Print("Creating driver...")
	return bolt.NewDriver()
}

func getAllStopNames(driver bolt.Driver) ([]string, error) {
	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryNeo(getAllStopNamesQuery, nil)
	if err != nil {
		return nil, err
	}

	stop_names := make([]string, 0)
	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return nil, err
		} else if err != io.EOF {
			stop_names = append(stop_names, row[0].(string))
		}
	}

	log.Printf(`Received %d stop names`, len(stop_names))
	return stop_names, nil
}

type Route struct {
	ID    string
	IsBus bool
}

func getAllRouteIDs(driver bolt.Driver) ([]Route, error) {
	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryNeo(getAllRouteIDsQuery, nil)
	if err != nil {
		return nil, err
	}

	routes := make([]Route, 0)
	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return nil, err
		} else if err != io.EOF {
			route := Route{row[0].(string), row[1].(bool)}
			routes = append(routes, route)
		}
	}

	log.Printf(`Received %d route IDs`, len(routes))
	return routes, nil
}

type RouteVariant struct {
	RouteID   string
	IsBus     bool
	FirstStop string
	LastStop  string
	TripIDs   []string
}

func getRouteVariantsForRouteID(driver bolt.Driver, routeID string) ([]RouteVariant, error) {
	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(getRouteVariantsByRouteIDQuery)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"routeID": routeID})
	if err != nil {
		return nil, err
	}

	var variants []RouteVariant
	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return nil, err
		} else if err != io.EOF {
			routeID := row[0].(string)
			isBus := row[1].(bool)
			firstStopName := row[2].(string)
			lastStopName := row[3].(string)
			tripIDs := row[4].([]interface{})

			s := make([]string, len(tripIDs))
			for i, v := range tripIDs {
				s[i] = fmt.Sprint(v)
			}
			variants = append(variants, RouteVariant{routeID, isBus, firstStopName, lastStopName, s})
		}
	}

	log.Printf(`Received %d variants for routeID "%s"`, len(variants), routeID)
	return variants, nil
}

func getRouteVariantsByStopName(driver bolt.Driver, stopName string) ([]RouteVariant, error) {
	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(getRouteVariantsByStopNameQuery)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"stopName": stopName})
	if err != nil {
		return nil, err
	}

	var variants []RouteVariant
	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return nil, err
		} else if err != io.EOF {
			routeID := row[0].(string)
			isBus := row[1].(bool)
			firstStopName := row[2].(string)
			lastStopName := row[3].(string)
			tripIDs := row[4].([]interface{})

			s := make([]string, len(tripIDs))
			for i, v := range tripIDs {
				s[i] = fmt.Sprint(v)
			}
			variants = append(variants, RouteVariant{routeID, isBus, firstStopName, lastStopName, s})
		}
	}

	log.Printf(`Received %d variants for stop name "%s"`, len(variants), stopName)
	return variants, nil
}
