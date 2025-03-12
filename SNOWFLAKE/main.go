package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	user := ""
	password := ""
	account := "-"

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please provide an operation type:")
	fmt.Println("1) create_db")
	fmt.Println("2) create_table")
	fmt.Print("Enter choice: ")
	operation, _ := reader.ReadString('\n')
	operation = strings.TrimSpace(operation)

	switch operation {
	//Switch Statement to Create the DB and Schema
	case "1", "create_db":
		fmt.Print("Please provide the Client Name")
		clientName, _ := reader.ReadString('\n')
		clientName = strings.TrimSpace(clientName)
		fmt.Print("Please provide the S3 Bucket Name: ")
		bucketName, _ := reader.ReadString('\n')
		bucketName = strings.TrimSpace(bucketName)
		CreateDatabase(user, password, account, clientName)
		CreateSchema(user, password, account, clientName)
		//CreateStorage(user, password, account, clientName, bucketName)

	//Switch Statement to Create the DB Table
	case "2", "create_table":
		fmt.Print("Please provide the Client Name: ")
		clientName, _ := reader.ReadString('\n')
		clientName = strings.TrimSpace(clientName)

		fmt.Print("Please provide the Table Name: ")
		tableName, _ := reader.ReadString('\n')
		tableName = strings.TrimSpace(tableName)

		fmt.Print("Please provide the CSV file name: ")
		csvFileName, _ := reader.ReadString('\n')
		csvFileName = strings.TrimSpace(csvFileName)

		CreateTable(user, password, account, clientName, tableName, csvFileName)
		CreateStage(user, password, account, clientName, tableName)
		CreatePipe(user, password, account, clientName, tableName)
	default:
		log.Fatal("Operation not recognized. Please provide a valid operation type: \n1) create_db \n2) create_table \n")
	}
}
