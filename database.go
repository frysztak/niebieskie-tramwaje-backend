package main

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
)

const (
	URL = "bolt://neo4j:krowa@localhost:7687"
)

func openDB() bolt.Driver {
	log.Print("Creating driver...")
	return bolt.NewDriver()
}

type Stop struct {
	Name      string
	Latitude  float32
	Longitude float32
}

func getAllStops(driver bolt.Driver) ([]Stop, error) {
	conn, err := driver.OpenNeo(URL)
	stops := make([]Stop, 0)

	if err != nil {
		return stops, err
	}
	defer conn.Close()

	rows, err := conn.QueryNeo(getAllStopNamesQuery, nil)
	if err != nil {
		return stops, err
	}

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return stops, err
		} else if err != io.EOF {
			name := row[0].(string)
			lat := row[1].(float64)
			long := row[2].(float64)
			stops = append(stops, Stop{name, float32(lat), float32(long)})
		}
	}

	log.Printf(`Received %d stops`, len(stops))
	return stops, nil
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
			routeID := row[0].(string)
			routeType := row[1].(string)
			isBus := strings.Contains(routeType, "bus")
			route := Route{routeID, isBus}
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
	TripIDs   []int
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
			routeType := row[1].(string)
			firstStopName := row[2].(string)
			lastStopName := row[3].(string)
			tripIDs := row[4].([]interface{})

			s := make([]int, len(tripIDs))
			for i, v := range tripIDs {
				s[i] = int(v.(int64))
			}
			isBus := strings.Contains(routeType, "bus")
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
			routeType := row[1].(string)
			firstStopName := row[2].(string)
			lastStopName := row[3].(string)
			tripIDs := row[4].([]interface{})

			s := make([]int, len(tripIDs))
			for i, v := range tripIDs {
				s[i] = int(v.(int64))
			}
			isBus := strings.Contains(routeType, "bus")
			variants = append(variants, RouteVariant{routeID, isBus, firstStopName, lastStopName, s})
		}
	}

	log.Printf(`Received %d variants for stop name "%s"`, len(variants), stopName)
	return variants, nil
}

type TimeTableEntry struct {
	TripID        int
	ArrivalTime   string
	DepartureTime string
}

type TimeTable struct {
	RouteID   string
	StopName  string
	Direction string
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

func normaliseTime(time string) string {
	mapping := map[string]string{
		"24:": "00:",
		"25:": "01:",
		"26:": "02:",
		"27:": "03:",
		"28:": "04:",
		"29:": "05:",
	}

	for original, replacement := range mapping {
		if strings.Contains(time, original) {
			return strings.Replace(time, original, replacement, -1)
		}
	}

	return time
}

func (entry *TimeTableEntry) normalise() {
	entry.ArrivalTime = normaliseTime(entry.ArrivalTime)
	entry.DepartureTime = normaliseTime(entry.DepartureTime)
}

func (tt TimeTable) normalise() {
	normaliseDay := func(slice []TimeTableEntry) {
		for idx, _ := range slice {
			(&slice[idx]).normalise()
		}
	}
	normaliseDay(tt.Weekdays)
	normaliseDay(tt.Saturdays)
	normaliseDay(tt.Sundays)
}

func getTimetable(driver bolt.Driver, routeID string, stopName string, direction string) (TimeTable, error) {
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
		"routeID":   routeID,
		"stopName":  stopName,
		"direction": direction})
	if err != nil {
		return TimeTable{}, err
	}

	var timeTable TimeTable
	timeTable.RouteID = routeID
	timeTable.StopName = stopName
	timeTable.Direction = direction
	timeTable.Weekdays = []TimeTableEntry{}
	timeTable.Saturdays = []TimeTableEntry{}
	timeTable.Sundays = []TimeTableEntry{}

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return TimeTable{}, err
		} else if err != io.EOF {
			tripID := int(row[0].(int64))
			arrivalTime := row[1].(string)
			departureTime := row[2].(string)

			entry := TimeTableEntry{tripID, arrivalTime, departureTime}

			tripIDString := strconv.Itoa(tripID)
			switch prefix := tripIDString[0]; prefix {
			case '2': // Monday-Thursday
				fallthrough
			case '6':
				timeTable.Weekdays = append(timeTable.Weekdays, entry)

			case '8': // Friday
				fallthrough
			case '1': // it's actual '10', but we can cheat a little
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
	timeTable.normalise()

	log.Printf(`Received %d time table entries for route ID "%s", stop name "%s" and direction "%s"`, len(timeTable.Weekdays)+len(timeTable.Saturdays)+len(timeTable.Sundays), routeID, stopName, direction)
	return timeTable, nil
}

type RouteInfo struct {
	RouteID     string
	RouteType   string
	IsBus       bool
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
			routeType := row[1].(string)
			validFrom := row[2].(string)
			validUntil := row[3].(string)
			agencyName := row[4].(string)
			agencyUrl := row[5].(string)
			agencyPhone := row[6].(string)

			isBus := strings.Contains(routeType, "bus")
			routeInfo = RouteInfo{routeID, routeType, isBus, validFrom, validUntil, agencyName, agencyUrl, agencyPhone}
		}
	}

	log.Printf(`Received route info for route id "%s"`, routeID)
	return routeInfo, nil
}

type RouteDirections struct {
	RouteID    string
	Directions []string
}

func getRouteDirections(driver bolt.Driver, routeID string) (RouteDirections, error) {
	var routeDirections RouteDirections
	routeDirections.RouteID = routeID

	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return routeDirections, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(getRouteDirectionsQuery)
	if err != nil {
		return routeDirections, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"routeID": routeID})
	if err != nil {
		return routeDirections, err
	}

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return routeDirections, err
		} else if err != io.EOF {
			direction := row[0].(string)
			routeDirections.Directions = append(routeDirections.Directions, direction)
		}
	}

	log.Printf(`Received %d route directions for route id "%s"`, len(routeDirections.Directions), routeID)
	return routeDirections, nil
}

func getRouteDirectionsThroughStop(driver bolt.Driver, routeID string, stopName string) (RouteDirections, error) {
	var routeDirections RouteDirections
	routeDirections.RouteID = routeID

	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return routeDirections, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(getRouteDirectionsThroughStopQuery)
	if err != nil {
		return routeDirections, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"routeID":  routeID,
		"stopName": stopName})
	if err != nil {
		return routeDirections, err
	}

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return routeDirections, err
		} else if err != io.EOF {
			direction := row[0].(string)
			routeDirections.Directions = append(routeDirections.Directions, direction)
		}
	}

	log.Printf(`Received %d route directions for route id "%s" going through stop "%s"`, len(routeDirections.Directions), routeID, stopName)
	return routeDirections, nil
}

type TripTimelineEntry struct {
	StopName      string
	DepartureTime string
	OnDemand      bool
}

type TripTimeline struct {
	TripID   int
	Timeline []TripTimelineEntry
}

func getTripTimeline(driver bolt.Driver, tripID int) (TripTimeline, error) {
	var timeline TripTimeline
	timeline.TripID = tripID

	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return timeline, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(getTripTimelineQuery)
	if err != nil {
		return timeline, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"tripID": tripID})
	if err != nil {
		return timeline, err
	}

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return timeline, err
		} else if err != io.EOF {
			stopName := row[0].(string)
			departureTime := normaliseTime(row[1].(string))
			onDemand := row[2].(bool)
			timeline.Timeline = append(timeline.Timeline, TripTimelineEntry{stopName, departureTime, onDemand})
		}
	}

	log.Printf(`Received %d timeline entries for trip id "%d"`, len(timeline.Timeline), tripID)
	return timeline, nil
}

type StopsForRoute struct {
	RouteID   string
	StopNames []string
}

func getStopsForRouteID(driver bolt.Driver, routeID string) (StopsForRoute, error) {
	var data StopsForRoute
	data.RouteID = routeID

	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return data, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(getStopsForRouteIDQuery)
	if err != nil {
		return data, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"routeID": routeID})
	if err != nil {
		return data, err
	}

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return data, err
		} else if err != io.EOF {
			stopName := row[0].(string)
			data.StopNames = append(data.StopNames, stopName)
		}
	}

	log.Printf(`Received %d stop names for route id "%s"`, len(data.StopNames), routeID)
	return data, nil
}

// key -> shapeID
// value -> list of trip IDs
type ShapeMap map[int][]int

func getShapeIDs(driver bolt.Driver, routeID, direction, stopName string) (ShapeMap, error) {
	data := ShapeMap{}

	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return data, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(getShapeIDsQuery)
	if err != nil {
		return data, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"routeID":   routeID,
		"direction": direction,
		"stopName":  stopName,
	})
	if err != nil {
		return data, err
	}

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return data, err
		} else if err != io.EOF {
			shapeID := int(row[0].(int64))
			tripIDs_ := row[1].([]interface{})
			tripIDs := make([]int, len(tripIDs_))
			for i, v := range tripIDs_ {
				tripIDs[i] = int(v.(int64))
			}
			data[shapeID] = tripIDs
		}
	}

	log.Printf(`Received %d shape-map entries for route id "%s", direction "%s" and stopName "%s"`, len(data), routeID, direction, stopName)
	return data, nil
}

type ShapePoint struct {
	ShapeID       int `json:"-"`
	ShapeSequence int `json:"-"`
	Latitude      float32
	Longitude     float32
}

type ShapePoints []ShapePoint

func getShapePoints(driver bolt.Driver, shapeID int) (ShapePoints, error) {
	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(getShapeQuery)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"shapeID": shapeID,
	})
	if err != nil {
		return nil, err
	}

	var data ShapePoints
	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return data, err
		} else if err != io.EOF {
			shapeID := int(row[0].(int64))
			shapeSeq := int(row[1].(int64))
			lat := float32(row[2].(float64))
			lon := float32(row[3].(float64))
			data = append(data, ShapePoint{shapeID, shapeSeq, lat, lon})
		}
	}

	log.Printf(`Received %d shape-points for shape id %d`, len(data), shapeID)
	return data, nil
}

type StopOnDemand struct {
	Stop
	OnDemand bool
}

func getStopsForTripID(driver bolt.Driver, tripID int) ([]StopOnDemand, error) {
	conn, err := driver.OpenNeo(URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(getTripStopsQuery)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"tripID": tripID,
	})
	if err != nil {
		return nil, err
	}

	var stops []StopOnDemand
	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return stops, err
		} else if err != io.EOF {
			name := row[0].(string)
			lat := row[1].(float64)
			long := row[2].(float64)
			onDemand := row[3].(bool)
			stops = append(stops, StopOnDemand{Stop{name, float32(lat), float32(long)}, onDemand})
		}
	}

	log.Printf(`Received %d stops for trip id %d`, len(stops), tripID)
	return stops, nil
}

type Shape struct {
	ShapeID int
	Points  ShapePoints
}

type StopOnMap struct {
	StopOnDemand
	FirstOrLast bool
}

type MapData struct {
	Shapes []Shape
	Stops  []StopOnMap
}

func getMapData(driver bolt.Driver, routeID, direction, stopName string) (MapData, error) {
	var data MapData

	shapeMap, err := getShapeIDs(driver, routeID, direction, stopName)
	if err != nil {
		return data, err
	}

	// get all shape IDs
	// at the same time build a set of canonical tripIDs (one for each shape ID). they'll be used later.
	shapeIDs := make([]int, 0, len(shapeMap))
	tripIDs := make([]int, 0, len(shapeMap))
	for k := range shapeMap {
		shapeIDs = append(shapeIDs, k)
		tripIDs = append(tripIDs, shapeMap[k][0])
	}

	for _, shapeID := range shapeIDs {
		// TODO: use pipeline for concurrent querying
		points, err := getShapePoints(driver, shapeID)
		if err != nil {
			return data, err
		}
		data.Shapes = append(data.Shapes, Shape{shapeID, points})
	}

	// value (bool) marks first or last stop in the trip
	stopsMap := map[StopOnDemand]bool{}
	for _, tripID := range tripIDs {
		newStops, err := getStopsForTripID(driver, tripID)
		if err != nil {
			return data, err
		}

		for idx, newStop := range newStops {
			if _, ok := stopsMap[newStop]; ok {
				// we already have such a stop. do nothing.
			} else {
				stopsMap[newStop] = false
			}

			// mark either first or last stop in the trip
			if idx == 0 || idx == len(newStops)-1 {
				stopsMap[newStop] = true
			}
		}
	}

	stops := make([]StopOnMap, 0, len(stopsMap))
	for stop, firstOrLast := range stopsMap {
		stops = append(stops, StopOnMap{stop, firstOrLast})
	}
	data.Stops = stops

	log.Printf(`Received %d shapes and %d stops`, len(data.Shapes), len(data.Stops))
	return data, nil
}
