package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/sorabhlahoti/golang-serverless-project/pkg/handlers"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch req.HTTPMethod {

	case "GET":
		return handlers.GetUser(req, tableName, dynaClient)

	case "POST":
		return handlers.CreateUser(req, tableName, dynaClient)

	case "PUT":
		return handlers.UpdateUser(req, tableName, dynaClient)

	case "DELETE":
		return handlers.DeleteUser(req, tableName, dynaClient)

	default:
		return handlers.UnhandledMethod()

	}

}

var (
	dynaClient *dynamodb.Client
	tableName  = "go-serverless-project"
)

func main() {

	//Aws dynamodn setup
	region := os.Getenv("AWS_REGION")
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))

	if err != nil {
		return
	}

	dynaClient = dynamodb.NewFromConfig(cfg)

	//register lambda handler to process requests(events)
	lambda.Start(handler)

}
