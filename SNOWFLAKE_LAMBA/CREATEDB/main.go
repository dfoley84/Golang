package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Event struct {
	ClientName string `json:"ClientName"`
	BucketName string `json:"BucketName"`
	SecertName string `json:"SecertName"`
	Action     string `json:"Action"`
}

func GetSecretValue(ctx context.Context, cfg aws.Config, secretName string) (string, error) {
	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}
	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("unable to get secret value for: %v", err)
	}

	if result.SecretString == nil {
		return "", fmt.Errorf("secret value for %s is nil", secretName)
	}

	return *result.SecretString, nil
}

func handler(ctx context.Context, event Event) (map[string]interface{}, error) {
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, configuring AWS SDK: %v", err)
	}
	secretValue, err := GetSecretValue(ctx, cfg, event.SecertName)
	if err != nil {
		log.Fatalf("unable to get secret value: %v", err)
		return nil, err
	}
	var secretData map[string]string
	if err := json.Unmarshal([]byte(secretValue), &secretData); err != nil {
		log.Fatalf("unable to unmarshal secret value: %v", err)
		return nil, err
	}
	user := secretData["user"]
	password := secretData["password"]
	account := secretData["account"]

	switch event.Action {
	case "create":
		CreateDatabase(user, password, account, event.ClientName)
	case "schema":
		CreateSchema(user, password, account, event.ClientName)
	case "storage":
		CreateStorage(user, password, account, event.ClientName, event.BucketName)
	case "snsTopic":
		CreateSNSTopic(user, password, account, event.ClientName, region)
	case "IAMRole":
		ClientDetails(user, password, account, event.ClientName, event.BucketName)
		CreateCFTemplate(region, event.ClientName)
	}

	return map[string]interface{}{
		"status": "success",
	}, nil
}
func main() {
	lambda.Start(handler)
}
