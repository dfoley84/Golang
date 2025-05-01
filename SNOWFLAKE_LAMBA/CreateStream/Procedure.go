package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

func CreateProcedure(user, password, account, clientName string, tbName string) {
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

	fmt.Println("Creating Procedure")
	query := fmt.Sprintf(`
		CREATE OR REPLACE PROCEDURE ..%[1]s()
		RETURNS VARCHAR
		LANGUAGE JAVASCRIPT
		EXECUTE AS OWNER
		AS
		$$
			var result = snowflake.execute({
				sqlText: "SELECT TO_CHAR(CURRENT_TIMESTAMP, 'YYYYMMDD_HH24MISS')"
			});
			result.next();
			var ts = result.getColumnValue(1);
			
			var filename = 'STREAM_%[1]s' + ts + '.csv';
			
			var sql_command = "COPY INTO @..STAGE_%[1]s/" + filename + 
				" FROM ..VIEW_STREAM_%[1]s " +
				"FILE_FORMAT = (" +
					"TYPE = 'CSV', " +
					"FIELD_DELIMITER = ',', " +
					"RECORD_DELIMITER = '\\n', " +
					"FIELD_OPTIONALLY_ENCLOSED_BY = '\"', " +
					"COMPRESSION = NONE" +
				") " +
				"HEADER = TRUE " +
				"SINGLE = TRUE " +
				"MAX_FILE_SIZE = 5368709120";
			
			snowflake.execute({ sqlText: sql_command });
			return 'Copy command executed successfully with file ' + filename;
		$$;
	`, strings.ToUpper(tbName))

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating procedure: %v", err)
	}
	fmt.Println("Procedure created successfully")
}
