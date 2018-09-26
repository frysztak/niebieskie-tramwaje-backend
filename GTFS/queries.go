package GTFS

const getAllStopNamesQuery = `
	MATCH (s:Stop)
    RETURN DISTINCT s.name, s.stopID, s.latitude, s.longitude
	ORDER BY s.name;
`

const getAllRouteIDsQuery = `
	MATCH (t:Trip)
	WITH t
	MATCH (:Route{routeID: t.routeID})-[:is_type]->(routeType:RouteType)
    RETURN DISTINCT
		t.routeID,
		routeType.name
	ORDER BY t.routeID;
`

const getRouteVariantsByRouteIDQuery = `
	MATCH (trip:Trip{routeID: {routeID}})-[:starts_at]->(st:StopTime)-[:happens_at]->(stop:Stop)
	WITH trip.tripID as tripID, stop.name as firstStopName
	MATCH (trip:Trip{tripID: tripID})-[:ends_at]->(st:StopTime)-[:happens_at]->(stop:Stop)
	WITH trip, stop, firstStopName
	MATCH (route:Route{routeID: trip.routeID})-[:is_type]->(routeType: RouteType)
	RETURN
	    trip.routeID as routeID,
		routeType.name as routeType,
		firstStopName,
		stop.name as lastStopName,
		collect(trip.tripID) as tripIDs
	ORDER BY routeID;
`

const getRouteVariantsByStopNameQuery = `
	MATCH (st:StopTime)-[:happens_at]->(stop:Stop{name: {stopName}})
	WITH st.tripID as tripID
	MATCH (trip:Trip{tripID: tripID})-[:starts_at]-(st:StopTime)-[:happens_at]->(stop:Stop)
	WITH trip.tripID as tripID, stop.name as firstStopName
	MATCH (trip:Trip{tripID: tripID})-[:ends_at]-(st:StopTime)-[:happens_at]->(stop:Stop)
	WITH trip, firstStopName, stop
	MATCH (:Route{routeID: trip.routeID})-[:is_type]-(routeType:RouteType)
	RETURN
		trip.routeID as routeID,
		routeType.name as routeType,
		firstStopName,
		stop.name as lastStopName,
		collect(trip.tripID) as tripIDs
	ORDER BY routeID;
`

const getRouteDirectionsQuery = `
	MATCH (t:Trip {routeID: {routeID}})
	WITH t.headsign as headsign, count(t.tripID) as cnt
	ORDER BY cnt DESC
	RETURN headsign
`

const getRouteDirectionsThroughStopQuery = `
	MATCH (st: StopTime)-[:happens_at]-(s: Stop{name: {stopName}})
	WITH collect(st.tripID) as tripIDs

	MATCH (t:Trip {routeID: {routeID}})
	WHERE t.tripID in tripIDs
	WITH t.headsign as headsign, count(t.tripID) as cnt
	ORDER BY cnt DESC
	RETURN headsign
`

const getTimetableQuery = `
	MATCH (t:Trip{routeID: {routeID}})-[:ends_at]-(:StopTime)-[:happens_at]->(stop:Stop {name: {direction}})
	WITH collect(t.tripID) as tripIDs

	MATCH (s:Stop {name: {stopName}})
	WITH collect(s.stopID) as stopIDs, tripIDs
	MATCH (st:StopTime)
	WHERE st.tripID IN tripIDs AND st.stopID in stopIDs
	RETURN
		st.tripID as tripID,
		st.arrivalTime as arrivalTime,
		st.departureTime as departureTime;
`

const getRouteInfoQuery = `
	MATCH (route:Route {routeID: {routeID}})-[:is_type]->(routeType:RouteType)
	WITH route, routeType
	MATCH (agency:Agency {agencyID: route.agencyID})
	RETURN
        route.routeID as routeID,
        routeType.name as routeType,
        route.validFrom as validFrom,
        route.validUntil as validUntil,
        agency.name as agencyName,
        agency.url as agencyUrl,
        agency.phone as agencyPhone;
`

const getTripTimelineQuery = `
	MATCH p=(t:Trip {tripID: {tripID}})-[:starts_at]-(:StopTime)-[:next*]-(:StopTime)-[:ends_at]-(t)
	WITH filter(n in nodes(p) WHERE EXISTS(n.stopID)) as nodes
	WITH extract(n in nodes | [n.stopID, n.departureTime, n.onDemand]) AS tuple
	UNWIND tuple as tuples
	WITH tuples[0] as stopID, tuples[1] as departureTime, tuples[2] as onDemand

	MATCH (s:Stop {stopID: stopID})
	RETURN s.name as stopName, departureTime, onDemand
`

const getStopsForRouteIDQuery = `
	MATCH (t:Trip {routeID: {routeID}})
	WITH collect(t.tripID) as tripIDs
	MATCH (st: StopTime)-[:happens_at]-(s:Stop)
	WHERE st.tripID in tripIDs
	WITH DISTINCT s.name as name, count(s.name) as cnt
	RETURN name ORDER BY cnt DESC
`

const getShapeIDsQuery = `
	MATCH (t:Trip{routeID: {routeID}})-[:ends_at]-(:StopTime)-[:happens_at]->(stop:Stop {name: {direction}})
    WITH collect(t.tripID) as tripIDs

    MATCH (s:Stop {name: {stopName}})
    WITH collect(s.stopID) as stopIDs, tripIDs
    MATCH (st:StopTime)
    WHERE st.tripID IN tripIDs AND st.stopID in stopIDs
    WITH collect(st.tripID) as tripIDs

    MATCH (t: Trip)
    WHERE t.tripID in tripIDs
    WITH collect(t.tripID) as tripID, t.shapeID as shapeID
    RETURN shapeID, tripID
`

const getShapeQuery = `
	MATCH (s: ShapePoint {shapeID: {shapeID}})
	RETURN s.shapeID, s.shapeSequence, s.latitude, s.longitude
	ORDER BY s.shapeSequence
`

const getTripStopsQuery = `
	MATCH p=(t:Trip {tripID: {tripID}})-[:starts_at]-(:StopTime)-[:next*]-(:StopTime)-[:ends_at]-(t)
    WITH filter(n in nodes(p) WHERE EXISTS(n.stopID)) as nodes
    WITH extract(n in nodes | [n.stopID, n.stopSequence, n.onDemand]) AS tuple
    UNWIND tuple as tuples
    WITH tuples[0] as stopID, tuples[1] as stopSequence, tuples[2] as onDemand

    MATCH (s:Stop {stopID: stopID})
    RETURN s.name, s.stopID, s.latitude, s.longitude, onDemand
    ORDER BY stopSequence
`

const getShapeForTripIDQuery = `
	MATCH (t: Trip{tripID: {tripID}})
	WITH t.shapeID as shapeID
	MATCH (s: ShapePoint {shapeID: shapeID})
	RETURN s.shapeID, s.shapeSequence, s.latitude, s.longitude
	ORDER BY s.shapeSequence
`

const getUpcomingDeparturesQuery = `
	MATCH (stop:Stop {name: {stopName}})<-[:happens_at]-(st: StopTime)
	WITH stop, st
	MATCH (t:Trip {tripID: st.tripID})
	RETURN stop.stopID, stop.name, stop.latitude, stop.longitude, st.tripID, st.departureTime, st.onDemand, t.routeID, t.headsign
	ORDER BY st.departureTime
`
