package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

func CreateStream(user, password, account, clientName, schema string, clientTableName string, ShareDB string) {
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, account, "+strings.ToUpper(clientName)+"")
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//USE THE ACCOUNTADMIN ROLE TO CREATE TABLES
	_, err = db.Exec(`USE ROLE ACCOUNTADMIN`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("USE SCHEMA ")
	if err != nil {
		log.Fatal("Error", err)
	}

	println("Creating Stream")
	// Create the Stage
	query := fmt.Sprintf(`CREATE OR REPLACE STREAM %s.%s.STREAM_%s ON TABLE %s.%[2]s.%s APPEND_ONLY = TRUE`,
		""+strings.ToUpper(clientName)+"", strings.ToUpper(schema), clientTableName, ShareDB, clientTableName)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	println("Stream created successfully")
}
