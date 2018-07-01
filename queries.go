package main

const getAllStopNamesQuery = `
	MATCH (s:stop)
	RETURN DISTINCT s.stop_name
	ORDER BY stop_name;
`

const getAllRouteIDsQuery = `
	MATCH (t:trip)
	RETURN DISTINCT t.route_id,
	CASE t.vehicle_id
		WHEN 1  THEN true
		WHEN 2  THEN true
		WHEN 8  THEN true
		WHEN 13 THEN true
		ELSE false
	END AS is_bus
	ORDER BY route_id;
`

const getVariantsForRouteIDQuery = `
	MATCH (t:trip { route_id: '%s' })-[:stops_at]-(st: stop_time)-[:is_stop]-(s:stop)
	WITH t, st.stop_sequence AS stop_sequence, s
	ORDER BY stop_sequence ASC
	RETURN t.trip_id, collect(s.stop_name) AS stop_names;
`
