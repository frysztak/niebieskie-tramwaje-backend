package main

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func StopsHandler(db *sql.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := getAllStopNames(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	})
}

func RoutesHandler(db *sql.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := getAllRouteIDs(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	})
}

func RoutesVariantsHandler(db *sql.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		routeID := p.ByName("routeID")
		data, err := getVariantsForRouteID(db, routeID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	})
}

func main() {
	db, err := openDB()
	if err != nil {
		panic(err)
	}
	err = setGraphPath(db)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	router := httprouter.New()
	router.GET("/stops", StopsHandler(db))
	router.GET("/routes", RoutesHandler(db))
	router.GET("/routes/:routeID/variants", RoutesVariantsHandler(db))
	http.ListenAndServe(":8080", router)

	//getAllStopNames(db)
	//getAllRouteIDs(db)
	//getVariantsForRouteID(db, "122")

	/*
	   fmt.Println("# Inserting values")

	   var lastInsertId int
	   err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
	   checkErr(err)
	   fmt.Println("last inserted id =", lastInsertId)

	   fmt.Println("# Updating")
	   stmt, err := db.Prepare("update userinfo set username=$1 where uid=$2")
	   checkErr(err)

	   res, err := stmt.Exec("astaxieupdate", lastInsertId)
	   checkErr(err)

	   affect, err := res.RowsAffected()
	   checkErr(err)

	   fmt.Println(affect, "rows changed")

	   fmt.Println("# Querying")
	   rows, err := db.Query("SELECT * FROM userinfo")
	   checkErr(err)

	   for rows.Next() {
	       var uid int
	       var username string
	       var department string
	       var created time.Time
	       err = rows.Scan(&uid, &username, &department, &created)
	       checkErr(err)
	       fmt.Println("uid | username | department | created ")
	       fmt.Printf("%3v | %8v | %6v | %6v\n", uid, username, department, created)
	   }

	   fmt.Println("# Deleting")
	   stmt, err = db.Prepare("delete from userinfo where uid=$1")
	   checkErr(err)

	   res, err = stmt.Exec(lastInsertId)
	   checkErr(err)

	   affect, err = res.RowsAffected()
	   checkErr(err)

	   fmt.Println(affect, "rows changed")
	*/
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
