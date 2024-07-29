package database

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB 
// db is a struct

func ConnectToDB(postgresURL string) {
  
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		fmt.Println("not able to connect to database")
		return
	}
	DB=db
	err = DB.Ping() 
	//Ping verifies a connection to the database is still alive, 
	//establishing a connection if necessary
	if err != nil {
		log.Fatalf("failed to ping %s ", err)
		return
	}
	fmt.Println("connected to database")
}
