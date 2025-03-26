package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

func CreateSchema(user, password, account, clientName string) {
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, account, ""+strings.ToUpper(clientName)+"")
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the DB Schema
	_, err = db.Exec("CREATE OR REPLACE SCHEMA ")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Schema created successfully")
}
