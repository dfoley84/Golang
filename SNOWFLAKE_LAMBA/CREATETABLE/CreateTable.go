package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/snowflakedb/gosnowflake"
)

// sanitizeHeader ensures column names are clean and formatted properly
func sanitizeHeader(header string) string {
	header = strings.TrimSpace(header)
	header = strings.ReplaceAll(header, " ", "_")
	header = strings.ReplaceAll(header, "/", "")
	header = strings.TrimLeft(header, ".") // Remove any leading dot

	// Remove any non-alphanumeric characters except underscores
	reg, _ := regexp.Compile(`[^a-zA-Z0-9_]`)
	header = reg.ReplaceAllString(header, "")

	return header
}

// CreateTable creates a table in Snowflake based on a CSV file stored in S3
func CreateTable(ctx context.Context, cfg aws.Config, user, password, account, clientName, tableName, bucketName, s3Path, csvFileName string) error {
	// Construct DSN for Snowflake
	dsn := fmt.Sprintf("%s:%s@%s/", user, password, account, strings.ToUpper(clientName))
	fmt.Printf("Connecting to database with DSN: %s\n", dsn)

	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Switch schema
	fmt.Println("Using schema ")
	_, err = db.Exec(`USE SCHEMA `)
	if err != nil {
		return fmt.Errorf("failed to use schema: %w", err)
	}

	tableName = strings.ToUpper(tableName)

	// Set up S3 client
	s3Client := s3.NewFromConfig(cfg)
	key := fmt.Sprintf("%s/%s", strings.TrimRight(s3Path, "/"), csvFileName)
	fmt.Printf("Constructed S3 Key: '%s'\n", key)

	// Check if object exists
	_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("S3 object not found: %w", err)
	}

	// Fetch the CSV file
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer result.Body.Close()

	// Read CSV headers
	reader := csv.NewReader(result.Body)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Check for Byte Order Mark (BOM) and remove it
	if strings.HasPrefix(headers[0], "\ufeff") {
		headers[0] = strings.TrimPrefix(headers[0], "\ufeff")
	}

	fmt.Println("CSV Headers:", headers)

	// Process column names
	columnDefinitions := make([]string, len(headers))
	for i, header := range headers {
		sanitizedHeader := sanitizeHeader(header)
		columnDefinitions[i] = fmt.Sprintf(`"%s" VARCHAR(16777216)`, sanitizedHeader)
	}
	columnDefinitions = append(columnDefinitions, `"ExportDate" VARCHAR(16777216)`)
	fmt.Printf("Creating table '%s' with columns: %v\n", tableName, columnDefinitions)

	// Construct and execute CREATE TABLE statement
	createTableSQL := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columnDefinitions, ", "))
	fmt.Printf("Executing SQL: %s\n", createTableSQL)

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("Table created successfully")
	return nil
}
