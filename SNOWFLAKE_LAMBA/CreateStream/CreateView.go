package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

func CreateView(user, password, account, clientName string, tbName string) {
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, account, ""+strings.ToUpper(clientName)+"")
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`USE ROLE ACCOUNTADMIN`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("USE SCHEMA ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating View")
	query := fmt.Sprintf(`
		CREATE OR REPLACE VIEW VIEW_STREAM_%[1]s 
		AS SELECT *
		FROM STREAM_%[1]s`,
		strings.ToUpper(tbName))
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("View created successfully")

}
