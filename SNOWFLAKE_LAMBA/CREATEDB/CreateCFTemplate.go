package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func CreateCFTemplate(region string, clientName string) {
	// Initialize a session that the SDK will use to load credentials
	// from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Fatal(err)
	}
	// Create a CloudFormation client with the session
	cfClient := cloudformation.New(sess)
	// Pass in the template file name and region
	templateFileName := "/tmp/cf.yaml"
	templateBody, err := os.ReadFile(templateFileName)
	if err != nil {
		log.Fatalf("failed to read template file: %v", err)
	}

	// Create the CloudFormation stack
	input := &cloudformation.CreateStackInput{
		StackName:    aws.String(fmt.Sprintf("-%s", clientName)),
		TemplateBody: aws.String(string(templateBody)),
		Parameters:   []*cloudformation.Parameter{},
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
			aws.String("CAPABILITY_NAMED_IAM"),
		},
		Tags: []*cloudformation.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(clientName),
			},
			{
				Key:   aws.String("Environment"),
				Value: aws.String("Production"),
			},
		},
	}
	_, err = cfClient.CreateStack(input)
	if err != nil {
		log.Fatalf("failed to create stack: %v", err)
	}
	log.Println("Stack creation initiated successfully")

	// Wait for the stack to be created
	err = cfClient.WaitUntilStackCreateComplete(&cloudformation.DescribeStacksInput{
		StackName: aws.String(fmt.Sprintf("<ROLE>", clientName)),
	})
	if err != nil {
		log.Fatalf("failed to wait for stack creation: %v", err)
	}
	log.Println("Stack created successfully")

}
