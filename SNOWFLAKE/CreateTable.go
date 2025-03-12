package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

func CreateTable(user, password, account, clientName string, tableName string, csvFileName string) {
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, account, ""+.ToUpper(clientName)+"")
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`USE SCHEMA `)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(csvFileName)
	if err != nil {
		log.Fatal(err)
	}
	//Just read the Headers of the CSV file
	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Headers: ", headers)
	columnDefinitions := make([]string, len(headers))
	for i, header := range headers {
		header = strings.ReplaceAll(header, " ", "_")
		columnDefinitions[i] = fmt.Sprintf(`"%s" STRING`, header)
	}
	columnDefinitions = append(columnDefinitions, `"ExportDate" STRING`)
	createTableSQL := fmt.Sprintf(`CREATE TABLE %s (%s)`, tableName, strings.Join(columnDefinitions, ", "))
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table created successfully")
}
