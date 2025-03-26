package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

func CreateStorage(user, password, account, clientName string, Bucketname string) {
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, account, ""+strings.ToUpper(clientName)+"")
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`USE SCHEMA `)
	if err != nil {
		log.Fatal(err)
	}

	// Create the Storage Integration
	query := fmt.Sprintf(`CREATE OR REPLACE STORAGE INTEGRATION %s_<>_S3_INT
	TYPE = EXTERNAL_STAGE
	STORAGE_PROVIDER = 'S3'
	ENABLED = TRUE
	STORAGE_AWS_ROLE_ARN = ''
	STORAGE_ALLOWED_LOCATIONS = ('s3://%s/%s/*');
	`, strings.ToUpper(clientName), strings.ToUpper(clientName), Bucketname, strings.ToUpper(clientName))
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Storage Integration created successfully: %<>", strings.ToUpper(clientName))
}
