USING PERIODIC COMMIT
LOAD CSV WITH HEADERS FROM "file:///trips.txt" AS row
CREATE (:Trip {
    tripID:    toInteger(replace(row.trip_id, "_", "")),
    routeID:   row.route_id,
    vehicleID: toInteger(row.vehicle_id),
    variantID: toInteger(row.variant_id),
    headsign:  row.trip_headsign
});

USING PERIODIC COMMIT
LOAD CSV WITH HEADERS FROM "file:///stops.txt" AS row
CREATE (:Stop {
    stopID:    toInteger(row.stop_id),
    code:      toInteger(row.stop_code),
    name:      row.stop_name,
    latitude:  toFloat(row.stop_lat),
    longitude: toFloat(row.stop_lon)
});

USING PERIODIC COMMIT
LOAD CSV WITH HEADERS FROM "file:///stop_times.txt" AS row
CREATE (:StopTime {
    tripID:        toInteger(replace(row.trip_id, "_", "")),
    arrivalTime:   substring(row.arrival_time, 0, 5),
    departureTime: substring(row.departure_time, 0, 5),
    stopID:        toInteger(row.stop_id),
    stopSequence:  toInteger(row.stop_sequence),
    onDemand:      toInteger(row.drop_off_type) = 3
});

USING PERIODIC COMMIT
LOAD CSV WITH HEADERS FROM "file:///vehicle_types.txt" AS row
CREATE (:VehicleType {
    vehicleID:   toInteger(row.vehicle_type_id),
    name:        row.vehicle_type_name,
    description: row.vehicle_type_description,
    symbol:      row.vehicle_type_symbol
});

USING PERIODIC COMMIT
LOAD CSV WITH HEADERS FROM "file:///routes.txt" AS row
CREATE (:Route {
    routeID:    row.route_id,
    agencyID:   toInteger(row.agency_id),
    typeID:     toInteger(row.route_type2_id),
    validFrom:  row.valid_from,
    validUntil: row.valid_until
});

USING PERIODIC COMMIT
LOAD CSV WITH HEADERS FROM "file:///route_types.txt" AS row
CREATE (:RouteType {
    typeID: toInteger(row.route_type2_id),
    name:   row.route_type2_name
});

USING PERIODIC COMMIT
LOAD CSV WITH HEADERS FROM "file:///agency.txt" AS row
CREATE (:Agency {
    agencyID: toInteger(row.agency_id),
    name:     row.agency_name,
    url:      row.agency_url,
    phone:    row.agency_phone
});

USING PERIODIC COMMIT
LOAD CSV WITH HEADERS FROM "file:///feed_info.txt" AS row
CREATE (:FeedInfo {
    publisherName: row.feed_publisher_name,
    publisherURL:  row.feed_publisher_url,
    lang:          row.feed_lang,
    startDate:     row.feed_start_date,
    endDate:       row.feed_end_date
});

USING PERIODIC COMMIT
LOAD CSV WITH HEADERS FROM "file:///calendar_dates.txt" AS row
CREATE (:CalendarDate {
    serviceID:     toInteger(row.service_id),
    date:          row.date,
    lang:          row.feed_lang,
    expectionType: row.exception_type
});

CREATE INDEX ON :Trip(tripID);
CREATE INDEX ON :Stop(stopID);
CREATE INDEX ON :StopTime(tripID);
CREATE INDEX ON :VehicleType(vehicleID);
CREATE INDEX ON :Route(routeID);
CREATE INDEX ON :RouteType(typeID);
CREATE INDEX ON :Agency(agencyID);

CALL db.awaitIndexes();

// create `start_at` relationship
CALL apoc.periodic.iterate(
"MATCH (trip:Trip), (stopTime: StopTime{stopSequence: 0})
WHERE trip.tripID = stopTime.tripID
RETURN trip, stopTime",
"CREATE (trip)-[:starts_at]->(stopTime)",
{batchSize:100, parallel:true});

// create `ends_at` relationship
CALL apoc.periodic.iterate(
"MATCH (stopTime: StopTime)
WITH stopTime.tripID as tID, max(stopTime.stopSequence) as lastStopSeq
MATCH (trip: Trip{tripID: tID}), (stopTime: StopTime{tripID: tID, stopSequence: lastStopSeq})
RETURN trip, stopTime",
"CREATE (trip)-[:ends_at]->(stopTime)",
{batchSize:100, parallel:true});

// create `happens_at` relationship
CALL apoc.periodic.iterate(
"MATCH (st:StopTime), (stop:Stop {stopID: st.stopID})
RETURN st, stop",
"CREATE (st)-[:happens_at]->(stop)",
{batchSize:200, parallel:false});

// create `next` relationship
CALL apoc.periodic.iterate(
"MATCH (st:StopTime)
WITH st
ORDER BY st.stopSequence
RETURN st.tripID as tripID, COLLECT(st) AS sts",
"FOREACH(i in RANGE(0, size(sts)-2) |
  FOREACH(st1 in [sts[i]] |
    FOREACH(st2 in [sts[i+1]] |
      CREATE (st1)-[:next]->(st2))))",
{batchSize:200, parallel:false});


// create `is_type` relationship
MATCH (route:Route), (routeType: RouteType{typeID: route.typeID})
CREATE (route)-[:is_type]->(routeType);

