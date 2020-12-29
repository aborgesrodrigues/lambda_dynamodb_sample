package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

/*type MyEvent struct {
	Name string `json:"name"`
}*/

type DragonStatsRequest struct{
	Name string `json:"name"`
}

type DragonStats struct{
	Name string `json:"dragon_name"`
	Damage int `json:"damage"`
	Description string `json:"description"`
	Family string `json:"family"`
	LocationCity string `json:"location_city"`
	LocationCountry string `json:"location_country"`
	LocationNeighborhood string `json:"location_neighborhood"`
	LocationState string `json:"location_state"`
	Protection int `json:"protection"`
}

func getResponse(body string, statusCode int)events.APIGatewayProxyResponse{

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: body,
	}
}

func getItem(body string, svc *dynamodb.DynamoDB)(string, error){
	dragonStatsRequest := DragonStatsRequest{}

	err := json.Unmarshal([]byte(body), &dragonStatsRequest)
	if err != nil {
		return "", err
	}

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("DragonStatsTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"dragon_name": {
				S: aws.String(dragonStatsRequest.Name),
			},
		},
	})

	if err != nil {
		return "", err
	}

	if result.Item == nil {
		msg := fmt.Sprintf("Could not find %s dragon", dragonStatsRequest.Name)

		return "", errors.New(msg)
	}

	dragonStats := DragonStats{}

	fmt.Println(result.Item)
	err = dynamodbattribute.UnmarshalMap(result.Item, &dragonStats)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(dragonStats)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func scan(svc *dynamodb.DynamoDB)(string, error){
	result, err := svc.Scan(
		&dynamodb.ScanInput{
			TableName:                 aws.String("DragonStatsTable"),
		},
	)

	var dragons []DragonStats

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dragons)

	if err != nil {
		return "", err
	}

	b, err := json.Marshal(dragons)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func query(body string, svc *dynamodb.DynamoDB)(string, error){
	dragonStatsRequest := DragonStatsRequest{}

	err := json.Unmarshal([]byte(body), &dragonStatsRequest)
	if err != nil {
		return "", err
	}

	result, err := svc.Query(
		&dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":v1": {
					S: aws.String(dragonStatsRequest.Name),
				},
			},
			ExpressionAttributeNames: map[string]*string{
				"#family": aws.String("family"),
			},
			ProjectionExpression: 	aws.String("dragon_name, #family"),
			KeyConditionExpression: aws.String("dragon_name = :v1"),
			TableName:            	aws.String("DragonStatsTable"),
		},
	)

	var dragons []DragonStats

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dragons)

	if err != nil {
		return "", err
	}

	b, err := json.Marshal(dragons)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	//body, err := getItem(request.Body, svc)
	//body, err := scan(svc)
	body, err := query(request.Body, svc)

	if err != nil{
		return getResponse(err.Error(), 400), err
	}

	return getResponse(body, 200), err
}

func main() {
	lambda.Start(HandleRequest)
}