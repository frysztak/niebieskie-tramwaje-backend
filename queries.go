package main

const getAllStopNamesQuery = `
	MATCH (s:Stop)
    RETURN DISTINCT s.name
    ORDER BY s.name;
`

const getAllRouteIDsQuery = `
	MATCH (t:Trip)
    RETURN DISTINCT t.routeID,
    CASE t.vehicleID
        WHEN 1  THEN true
        WHEN 2  THEN true
        WHEN 8  THEN true
        WHEN 13 THEN true
        ELSE false
    END AS is_bus
    ORDER BY t.routeID;
`

const getRouteVariantsByRouteIDQuery = `
	MATCH (trip:Trip{routeID: {routeID}})-[:starts_at]->(st:StopTime)-[:happens_at]->(stop:Stop)
	WITH trip.tripID as tripID, stop.name as firstStopName
	MATCH (trip:Trip{tripID: tripID})-[:ends_at]->(st:StopTime)-[:happens_at]->(stop:Stop)
	RETURN
		trip.routeID as routeID,
		CASE trip.vehicleID
    	    WHEN 1  THEN true
    	    WHEN 2  THEN true
    	    WHEN 8  THEN true
    	    WHEN 13 THEN true
    	    ELSE false
    	END AS isBus,
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
	RETURN
		trip.routeID as routeID,
		CASE trip.vehicleID
    	    WHEN 1  THEN true
    	    WHEN 2  THEN true
    	    WHEN 8  THEN true
    	    WHEN 13 THEN true
    	    ELSE false
    	END AS isBus,
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
	MATCH (route:Route {routeID: {routeID}})
	WITH route
	MATCH (agency:Agency {agencyID: route.agencyID})
	RETURN
		route.routeID as routeID,
	    route.typeID as typeID,
	    route.validFrom as validFrom,
	    route.validUntil as validUntil,
	    agency.name as agencyName,
	    agency.url as agencyUrl,
	    agency.phone as agencyPhone;
`

const getTripTimelineQuery = `
	MATCH p=(t:Trip {tripID: {tripID}})-[:starts_at]-(:StopTime)-[:next*]-(:StopTime)-[:ends_at]-(t)
	WITH filter(n in nodes(p) WHERE EXISTS(n.stopID)) as nodes
	WITH extract(n in nodes | [n.stopID, n.departureTime]) AS stopIDsAndTime
	UNWIND stopIDsAndTime as tuples
	WITH tuples[0] as stopID, tuples[1] as departureTime

	MATCH (s:Stop {stopID: stopID})
	RETURN s.name as stopName, departureTime
`

const getStopsForRouteIDQuery = `
	MATCH (t:Trip {routeID: {routeID}})
	WITH collect(t.tripID) as tripIDs
	MATCH (st: StopTime)-[:happens_at]-(s:Stop)
	WHERE st.tripID in tripIDs
	WITH DISTINCT s.name as name, count(s.name) as cnt
	RETURN name ORDER BY cnt DESC
`
