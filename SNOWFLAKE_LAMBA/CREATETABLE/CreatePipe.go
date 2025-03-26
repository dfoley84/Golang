package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

func CreatePipe(user, password, account, clientName, tableName string) {
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, account, ""+strings.ToUpper(clientName)+"")
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("USE SCHEMA")
	if err != nil {
		log.Fatal(err)
	}

	query := fmt.Sprintf("SHOW COLUMNS IN TABLE %s", strings.ToUpper(tableName))
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var columns []string
	for rows.Next() {
		var tblName, schemaName, colName, dataType, nullFlag, defVal, kind, expr, comment, dbName, autoInc, policyName interface{}
		if err := rows.Scan(&tblName, &schemaName, &colName, &dataType, &nullFlag, &defVal, &kind, &expr, &comment, &dbName, &autoInc, &policyName); err != nil {
			log.Fatal(err)
		}
		columns = append(columns, fmt.Sprintf("%v", colName))
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if len(columns) == 0 {
		log.Fatal("No columns found for table: ", tableName)
	}

	// Print out the column names
	fmt.Println("Columns: ", columns)

	columnList := strings.Join(columns, `", "`)
	columnList = fmt.Sprintf(`"%s"`, columnList)

	var mappedColumns []string
	for i, col := range columns {
		if col == "ExportDate" {
			mappedColumns = append(mappedColumns, `CURRENT_TIMESTAMP AS "ExportDate"`)
		} else {
			mappedColumns = append(mappedColumns, fmt.Sprintf(`t.$%d AS "%s"`, i+1, col))
		}
	}
	mappedColumnsQuery := strings.Join(mappedColumns, ",\n")

	fmt.Println("Creating Pipe")

	pipeQuery := fmt.Sprintf(`
    CREATE OR REPLACE PIPE %s_SNOW_PIPE_%s
    AUTO_INGEST = TRUE
    ERROR_INTEGRATION = '%s_SNS_ERROR_NOTIFICATION'
    AS
    COPY INTO %s (%s)
    FROM (
        SELECT
            %s
        FROM @PROD_%s_FUSION_%s_S3_STAGE t
    )
    FILE_FORMAT = (TYPE = 'CSV', FIELD_DELIMITER = ',', RECORD_DELIMITER = '\n', SKIP_HEADER = 1, ERROR_ON_COLUMN_COUNT_MISMATCH = TRUE)`,

		strings.ToUpper(clientName), strings.ToUpper(tableName), strings.ToUpper(clientName),
		strings.ToUpper(tableName), columnList,
		mappedColumnsQuery,
		strings.ToUpper(clientName), strings.ToUpper(tableName))

	// Log the generated query for debugging
	fmt.Println("Generated Pipe Query:")
	fmt.Println(pipeQuery)

	_, err = db.Exec(pipeQuery)
	if err != nil {
		log.Fatal("Error executing pipe creation query: ", err)
	}
	fmt.Println("Pipe created successfully")
}
