package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	_ "github.com/snowflakedb/gosnowflake"
)

const tmpl = `AWSTemplateFormatVersion: '2010-09-09'
Resources:
  SnowflakeClientRole:
    Type: 'AWS::IAM::Role'
    Properties: 
      RoleName: '{{ .ClientName }}'
      AssumeRolePolicyDocument: 
        Version: '2012-10-17'
        Statement: 
          - Effect: 'Allow'
            Principal: 
              AWS: '{{ .IAMUserARN }}'
            Action: 'sts:AssumeRole'
            Condition: 
              StringEquals: 
                sts:ExternalId: 
                  - '{{ .ExternalID }}'
                  - '{{ .SNS_EXternalID }}'
      Policies: 
        - PolicyName: 'SnowflakeClientPolicy'
          PolicyDocument: 
            Version: '2012-10-17'
            Statement: 
              - Effect: 'Allow'
                Action: 
                  - 's3:PutObject'
                  - 's3:GetObject'
                  - 's3:GetObjectVersion'
                  - 's3:DeleteObject'
                  - 's3:DeleteObjectVersion'
                Resource: 'arn:aws:s3:::{{ .Bucket }}/{{ .ClientName }}/*'
              - Effect: 'Allow'
                Action: 
                  - 's3:ListBucket'
                  - 's3:GetBucketLocation'
                Resource: 'arn:aws:s3:::{{ .Bucket }}'
              - Effect: 'Allow'
                Action:
                  - 'sns:Publish'
                Resource: 'arn:aws:sns:*::'
`

type TemplateData struct {
	ClientName     string
	Bucket         string
	IAMUserARN     string
	ExternalID     string
	SNS_EXternalID string
}

func saveTemplateToFile(filename string, data TemplateData) error {
	data.ClientName = strings.ToLower(data.ClientName)
	t := template.Must(template.New("cloudformation").Parse(tmpl))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return t.Execute(file, data)
}

func ClientDetails(user, password, account, clientName string, Bucketname string) {
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, account, ""+strings.ToUpper(clientName)+"")
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`USE SCHEMA`)
	if err != nil {
		log.Fatal(err)
	}

	//Retrieve the properties of the Storage Integration
	descQuery := fmt.Sprintf("DESC INTEGRATION %s", strings.ToUpper(clientName))
	rows, err := db.Query(descQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var property, propertyType, value, defaultValue string
	var externalID, iamUserARN string
	for rows.Next() {
		if err := rows.Scan(&property, &propertyType, &value, &defaultValue); err != nil {
			log.Fatal(err)
		}
		if property == "STORAGE_AWS_EXTERNAL_ID" {
			externalID = value
		} else if property == "STORAGE_AWS_IAM_USER_ARN" {
			iamUserARN = value
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("STORAGE_AWS_EXTERNAL_ID: %s", externalID)
	log.Printf("STORAGE_AWS_IAM_USER_ARN: %s", iamUserARN)

	descQuery1 := fmt.Sprintf("DESC INTEGRATION %s_SNS_ERROR_NOTIFICATION", strings.ToUpper(clientName))
	rows1, err := db.Query(descQuery1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows1.Close()
	var sns_externalid string

	for rows1.Next() {
		if err := rows1.Scan(&property, &propertyType, &value, &defaultValue); err != nil {
			log.Fatal(err)
		}
		if property == "SF_AWS_EXTERNAL_ID" {
			sns_externalid = value
		}
	}

	if err := rows1.Err(); err != nil {
		log.Fatal(err)
	}

	// Generate and save the IAM Role CloudFormation template
	data := TemplateData{
		ClientName:     clientName,
		Bucket:         Bucketname,
		IAMUserARN:     iamUserARN,
		ExternalID:     externalID,
		SNS_EXternalID: sns_externalid,
	}
	err = saveTemplateToFile("/tmp/cf.yaml", data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("CloudFormation template saved to cf.yaml")
}
