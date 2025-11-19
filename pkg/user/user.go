package user

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sorabhlahoti/golang-serverless-project/pkg/validators"
)

var (
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorInvalidUserData         = "invalid user data"
	ErrorInvalidEmail            = "invalid email"
	ErrorCouldNotMarshalItem     = "could not marshal item"
	ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotDynamoPutItem   = "could not dynamo put item"
	ErrorUserAlreadyExists       = "user.User already exists"
	ErrorUserDoesNotExist        = "user.User does not exist"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func awsString(s string) *string { return &s } // helper for *string
func awsBool(b bool) *bool       { return &b } // helper for *bool

func FetchUser(ctx context.Context, email, table string, ddb *dynamodb.Client) (*User, error) {

	key, err := attributevalue.MarshalMap(map[string]string{"email": email})

	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	out, err := ddb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      &table,
		Key:            key,
		ConsistentRead: awsBool(true),
	})

	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	if out.Item == nil {
		return &User{}, nil
	}

	var u User

	if err := attributevalue.UnmarshalMap(out.Item, &u); err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return &u, nil

}

func FetchUsers(ctx context.Context, tableName string, ddb *dynamodb.Client) (*[]User, error) {

	out, err := ddb.Scan(ctx, &dynamodb.ScanInput{TableName: &tableName})

	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	var items []User

	if err := attributevalue.UnmarshalListOfMaps(out.Items, &items); err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return &items, nil

}

func CreateUser(ctx context.Context, body string, tableName string, ddb *dynamodb.Client) (*User, error) {
	var u User

	if err := json.Unmarshal([]byte(body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ := FetchUser(ctx, u.Email, tableName, ddb)

	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}

	av, err := attributevalue.MarshalMap(u)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	_, err = ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &tableName,
		Item:                av,
		ConditionExpression: awsString("attribute_not_exists(email)"),
	})

	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &u, nil
}

func UpdateUser(ctx context.Context, body string, tableName string, ddb *dynamodb.Client) (
	*User,
	error) {

	var u User

	if err := json.Unmarshal([]byte(body), &u); err != nil {

		return nil, errors.New(ErrorInvalidEmail)

	}

	currentUser, _ := FetchUser(ctx, u.Email, tableName, ddb)

	if currentUser != nil && len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	av, err := attributevalue.MarshalMap(u)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	_, err = ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      av,
	})

	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &u, nil

}

func DeleteUser(ctx context.Context, email, tableName string, ddb *dynamodb.Client) error {
	key, err := attributevalue.MarshalMap(map[string]string{"email": email})
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}
	_, err = ddb.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:           &tableName,
		Key:                 key,
		ConditionExpression: awsString("attribute_exists(email)"),
	})

	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil

}
