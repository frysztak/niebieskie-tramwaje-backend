package main

import (
	"./GTFS"
	"./News"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"log"
	"net/http"
)

func main() {
	driver := GTFS.OpenDB()
	newsDb := News.OpenDatabase()

	router := mux.NewRouter().UseEncodedPath()
	router.HandleFunc("/stops", GTFS.StopsHandler(driver))
	router.HandleFunc("/stops/{stopNames}/departures", GTFS.StopsUpcomingDeparturesHandler(driver))
	router.HandleFunc("/stops/and/routes", GTFS.StopsAndRoutesHandler(driver))
	router.HandleFunc("/routes", GTFS.RoutesHandler(driver))
	router.HandleFunc("/routes/variants/id/{routeID}", GTFS.RoutesVariantsByIdHandler(driver))
	router.HandleFunc("/routes/variants/stop/{stopName}", GTFS.RoutesVariantsByStopNameHandler(driver))
	router.HandleFunc("/route/{routeID}/timetable/at/{stopName}/direction/{direction}", GTFS.RoutesTimeTableHandler(driver))
	router.HandleFunc("/route/{routeID}/info", GTFS.RouteInfoHandler(driver))
	router.HandleFunc("/route/{routeID}/directions", GTFS.RouteDirectionsHandler(driver))
	router.HandleFunc("/route/{routeID}/directions/through/{stopName}", GTFS.RouteDirectionsThroughStopHandler(driver))
	router.HandleFunc("/route/{routeID}/stops", GTFS.RouteStopsHandler(driver))
	router.HandleFunc("/route/{routeID}/map/at/{stopName}/direction/{direction}", GTFS.RouteMapHandler(driver))
	router.HandleFunc("/trip/{tripID}/timeline", GTFS.TripTimelineHandler(driver))
	router.HandleFunc("/trip/{tripID}/map", GTFS.TripMapHandler(driver))
	router.HandleFunc("/news/recent", News.RecentNewsHandler(newsDb))
	router.HandleFunc("/news/page/{pageNum}", News.NewsHandler(newsDb))

	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	c := cron.New()
	c.AddFunc("@every 15m", func() { News.UpdateNews(newsDb) })
	c.Start()

	News.UpdateNews(newsDb)

	select {} // "sleep" forever
}
