package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

func CreateStage(user, password, account, clientName string, tbName string) {
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, account,strings.ToUpper(clientName))
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
	_, err = db.Exec("USE SCHEMA <>")
	if err != nil {
		log.Fatal(err)
	}
	println("Creating Stage")
	// Create the Stage
	query := fmt.Sprintf(`CREATE OR REPLACE STAGE STAGE_%s
		STORAGE_INTEGRATION = 
		URL = 's3://'`, strings.ToUpper(tbName), strings.ToLower(tbName))
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	println("Stage created successfully")
}
