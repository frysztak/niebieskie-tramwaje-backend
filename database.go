package main

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"io"
	"log"
	"sort"
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

type TimeTableEntry struct {
	TripID        string
	ArrivalTime   string
	DepartureTime string
}

type TimeTable struct {
	RouteID   string
	AtStop    string
	FromStop  string
	ToStop    string
	Weekdays  []TimeTableEntry
	Saturdays []TimeTableEntry
	Sundays   []TimeTableEntry
}

func (tt TimeTable) sort() {
	type Predicate func(i, j int) bool
	predFactory := func(slice []TimeTableEntry) Predicate {
		return func(i, j int) bool {
			return slice[i].ArrivalTime < slice[j].ArrivalTime
		}
	}
	sort.Slice(tt.Weekdays, predFactory(tt.Weekdays))
	sort.Slice(tt.Saturdays, predFactory(tt.Saturdays))
	sort.Slice(tt.Sundays, predFactory(tt.Sundays))
}

func getTimetable(driver bolt.Driver, routeID string, atStopName string, fromStopName string, toStopName string) (TimeTable, error) {
	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return TimeTable{}, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(getTimetableQuery)
	if err != nil {
		return TimeTable{}, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"routeID":      routeID,
		"atStopName":   atStopName,
		"fromStopName": fromStopName,
		"toStopName":   toStopName})
	if err != nil {
		return TimeTable{}, err
	}

	var timeTable TimeTable
	timeTable.RouteID = routeID
	timeTable.AtStop = atStopName
	timeTable.FromStop = fromStopName
	timeTable.ToStop = toStopName

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return TimeTable{}, err
		} else if err != io.EOF {
			tripID := row[0].(string)
			arrivalTime := row[1].(string)
			departureTime := row[2].(string)

			entry := TimeTableEntry{tripID, arrivalTime, departureTime}

			switch prefix := tripID[0]; prefix {
			case '6': // Monday-Thursday
				timeTable.Weekdays = append(timeTable.Weekdays, entry)
			case '8': // Friday
				// ignore.
				// for some reason, in Wrocław GTFS they make distinction between
				// Mondays-Thurdays and Fridays. To the best of my knowledge,
				// there is no difference whatsoever.
			case '3':
				timeTable.Saturdays = append(timeTable.Saturdays, entry)
			case '4':
				timeTable.Sundays = append(timeTable.Sundays, entry)
			default:
				panic(fmt.Sprintf("Unknown prefix: %d", prefix))
			}
		}
	}

	timeTable.sort()

	log.Printf(`Received %d time table entries for route ID "%s" and stop name "%s"`, len(timeTable.Weekdays)+len(timeTable.Saturdays)+len(timeTable.Sundays), routeID, atStopName)
	return timeTable, nil
}

type RouteInfo struct {
	RouteID     string
	TypeID      int
	ValidFrom   string
	ValidUntil  string
	AgencyName  string
	AgencyUrl   string
	AgencyPhone string
}

func getRouteInfo(driver bolt.Driver, routeID string) (RouteInfo, error) {
	var routeInfo RouteInfo

	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return routeInfo, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(getRouteInfoQuery)
	if err != nil {
		return routeInfo, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"routeID": routeID})
	if err != nil {
		return routeInfo, err
	}

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return routeInfo, err
		} else if err != io.EOF {
			routeID := row[0].(string)
			typeID := row[1].(int64)
			validFrom := row[2].(string)
			validUntil := row[3].(string)
			agencyName := row[4].(string)
			agencyUrl := row[5].(string)
			agencyPhone := row[6].(string)

			routeInfo = RouteInfo{routeID, int(typeID), validFrom, validUntil, agencyName, agencyUrl, agencyPhone}
		}
	}

	log.Printf(`Received route info for route id "%s"`, routeID)
	return routeInfo, nil
}
