package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

func CreateStage(user, password, account, clientName string, tbName string) {
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, account, +strings.ToUpper(clientName)+)
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("USE SCHEMA ")
	if err != nil {
		log.Fatal(err)
	}
	println("Creating Stage")
	// Create the Stage
	query := fmt.Sprintf(`CREATE OR REPLACE STAGE %s_%s_S3_STAGE
		STORAGE_INTEGRATION = S3_INT
		URL = 's3:///%s////'`, clientName, tbName, clientName, clientName, tbName)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	println("Stage created successfully")
}
