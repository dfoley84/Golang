package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func CreateSNSTopic(user, password, account, clientName string, region string) {
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, account, ""+strings.ToUpper(clientName)+"")
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`USE SCHEMA FUSION`)
	if err != nil {
		log.Fatal(err)
	}
	// Create the SNS Topic Integration
	query := fmt.Sprintf(`CREATE OR REPLACE NOTIFICATION INTEGRATION %s_SNS_ERROR_NOTIFICATION
	ENABLED = TRUE
	DIRECTION = OUTBOUND
	TYPE = QUEUE
	NOTIFICATION_PROVIDER = AWS_SNS
	AWS_SNS_TOPIC_ARN = 'arn:aws:sns:%s::'
	AWS_SNS_ROLE_ARN = 'arn:aws:iam:::role/SNOWFLAKE-%s-Role'`,
		strings.ToUpper(clientName), region, strings.ToUpper(clientName))
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Storage Integration created successfully: %s_FUSION_S3_INT", strings.ToUpper(clientName))
}
