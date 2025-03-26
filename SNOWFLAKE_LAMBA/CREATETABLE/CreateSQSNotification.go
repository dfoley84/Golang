package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	_ "github.com/snowflakedb/gosnowflake"
)

func CreateSQS(ctx context.Context, cfg aws.Config, user, password, account, clientName string, tableName string, bucketName string, s3Path string) {
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

	descQuery := fmt.Sprintf("DESC STAGE %%s_S3_STAGE",
		strings.ToUpper(clientName), strings.ToUpper(tableName))
	fmt.Println("Executing query:", descQuery) // Debugging statement
	rows, err := db.Query(descQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Columns:", columns)

	var parentProperty, property, propertyType, propertyValue, propertyDefault string
	var DIRECTORY_NOTIFICATION_CHANNEL string
	for rows.Next() {
		if err := rows.Scan(&parentProperty, &property, &propertyType, &propertyValue, &propertyDefault); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ParentProperty: %s, Property: %s, PropertyType: %s, PropertyValue: %s, PropertyDefault: %s\n",
			parentProperty, property, propertyType, propertyValue, propertyDefault) // Debugging statement
		if strings.HasSuffix(property, "DIRECTORY_NOTIFICATION_CHANNEL") {
			DIRECTORY_NOTIFICATION_CHANNEL = propertyValue
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	if DIRECTORY_NOTIFICATION_CHANNEL == "" {
		log.Fatal("DIRECTORY_NOTIFICATION_CHANNEL not found")
	}
	fmt.Println("Notification Channel: ", DIRECTORY_NOTIFICATION_CHANNEL)
	s3client := s3.NewFromConfig(cfg)
	clientPath := fmt.Sprintf("%s%s", clientName, clientName, tableName)

	_, err = s3client.PutBucketNotificationConfiguration(ctx, &s3.PutBucketNotificationConfigurationInput{
		Bucket: aws.String(bucketName),
		NotificationConfiguration: &s3types.NotificationConfiguration{
			QueueConfigurations: []s3types.QueueConfiguration{
				{
					Id:       aws.String("SQSNotification"),
					QueueArn: aws.String(DIRECTORY_NOTIFICATION_CHANNEL),
					Events: []s3types.Event{
						s3types.EventS3ObjectCreatedPut,
						s3types.EventS3ObjectCreatedPost,
						s3types.EventS3ObjectCreatedCompleteMultipartUpload,
						s3types.EventS3ObjectCreatedCopy,
					},
					Filter: &s3types.NotificationConfigurationFilter{
						Key: &s3types.S3KeyFilter{
							FilterRules: []s3types.FilterRule{
								{
									Name:  s3types.FilterRuleNamePrefix,
									Value: aws.String(clientPath),
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to create SQS notification: %v", err)
	}
}
