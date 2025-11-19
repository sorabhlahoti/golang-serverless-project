package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sorabhlahoti/golang-serverless-project/pkg/user"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

var ErrorMethodNotAllowed = "method not allowed"

func GetUser(req events.APIGatewayProxyRequest, table string, ddb *dynamodb.Client) (events.APIGatewayProxyResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	email := req.QueryStringParameters["email"]
	if email != "" {
		res, err := user.FetchUser(ctx, email, table, ddb)
		if err != nil {
			s := err.Error()
			return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: &s})
		}
		return apiResponse(http.StatusOK, res)
	}
	res, err := user.FetchUsers(ctx, table, ddb)
	if err != nil {
		s := err.Error()
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: &s})
	}
	return apiResponse(http.StatusOK, res)
}

func CreateUser(req events.APIGatewayProxyRequest, table string, ddb *dynamodb.Client) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := user.CreateUser(ctx, req.Body, table, ddb)
	if err != nil {
		s := err.Error()
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: &s})
	}
	return apiResponse(http.StatusCreated, res)
}

func UpdateUser(req events.APIGatewayProxyRequest, table string, ddb *dynamodb.Client) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := user.UpdateUser(ctx, req.Body, table, ddb)
	if err != nil {
		s := err.Error()
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: &s})
	}
	return apiResponse(http.StatusOK, res)
}

func DeleteUser(req events.APIGatewayProxyRequest, table string, ddb *dynamodb.Client) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := user.DeleteUser(ctx, req.QueryStringParameters["email"], table, ddb); err != nil {
		s := err.Error()
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: &s})
	}
	return apiResponse(http.StatusOK, map[string]string{"status": "deleted"})
}

func UnhandledMethod() (events.APIGatewayProxyResponse, error) {
	return apiResponse(405, ErrorMethodNotAllowed)
}
