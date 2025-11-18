package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws/aws-sdk-go/service/dynamodb"
	"github.com/sorabhlahoti/golang-serverless-project/pkg/handlers"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {

	//Aws dynamodn setup
	region := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})

	if err != nil {
		return
	}

	dynaClient = dynamodb.New(sess)

	//register lambda handler to process requests(events)
	lambda.Start(handler)

}

const tableName = "go-serverless-project"

func handler(req events.APIGatewayProxyRequest) (*event.APIGatewayProxyResponse, error) {

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
