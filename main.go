package main

import _ "log"

func main() {
	db := openDB()
	defer db.Close()

	setGraphPath(db)
	getAllStopNames(db)
	getAllRouteIDs(db)
	getVariantsForRouteID(db, "122")

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
