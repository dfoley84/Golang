package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/transfer"
)

type Event struct {
	UserName      string `json:"user_name"`
	Role          string `json:"role"`
	HomeDirectory string `json:"home_directory"`
	ServerId      string `json:"server_id"`
	BucketName    string `json:"bucket_name"`
}

func handler(ctx context.Context, event Event) error {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return fmt.Errorf("AWS_REGION environment variable is not set")
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("AWS config error: %v", err)
	}

	SFTPClient := transfer.NewFromConfig(cfg)
	resp, err := SFTPClient.CreateUser(ctx, &transfer.CreateUserInput{
		UserName:          &event.UserName,
		Role:              &event.Role,
		HomeDirectory:     &event.HomeDirectory,
		HomeDirectoryType: "LOGICAL",
		HomeDirectoryMappings: []transfer.HomeDirectoryMapEntry{
			{
				Entry:  "/",
				Target: "/" + event.BucketName + "/" + event.HomeDirectory + "/",
				Type:   "DIRECTORY",
			},
		},
		ServerId: &event.ServerId,
	})
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
		return err
	}
	fmt.Printf("User created successfully: %v\n", resp)
	return nil
}

func main() {
	lambda.Start(handler)
}
