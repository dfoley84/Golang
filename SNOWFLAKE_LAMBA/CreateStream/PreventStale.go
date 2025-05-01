package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

func CreatePreventStale(user, password, account, clientName string, tbName string) {
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, account, ""+strings.ToUpper(clientName)+"")
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
		log.Fatal(err)
	}
	println("Creating Task")
	query := fmt.Sprintf(`
        CREATE OR REPLACE TASK ..PREVENT_STALE_STREAM_%[1]s
        WAREHOUSE=
        SCHEDULE='1 HOUR'
        WHEN SYSTEM$STREAM_HAS_DATA('STREAM_%[1]s') = FALSE
        AS BEGIN
          SELECT * FROM STREAM_%[1]s;
          SELECT SYSTEM$STREAM_HAS_DATA('STREAM_%[1]s');
        END;`,
		strings.ToUpper(tbName))

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	println("Stage created successfully")
}
