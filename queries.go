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
