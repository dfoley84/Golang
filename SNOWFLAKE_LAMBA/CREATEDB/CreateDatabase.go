package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

func CreateDatabase(user, password, account, clientName string) {
	// Construct the DSN (Data Source Name) for Snowflake
	dsn := user + ":" + password + "@" + account
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	clientName = strings.ToUpper(clientName)
	// Create a Database
	DBQuery := fmt.Sprintf(`CREATE OR REPLACE DATABASE "<>"`, clientName)
	_, err = db.Exec(DBQuery)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(`Database created successfully: %`, clientName)
}
