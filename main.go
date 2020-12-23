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
	Damage int `json:"damage"`
	Description string `json:"description"`
	Family string `json:"family"`
	LocationCity string `json:"location_city"`
	LocationCountry string `json:"location_country"`
	LocationNeighborhood string `json:"location_neighborhood"`
	LocationState string `json:"location_state"`
	Protection int `json:"protection"`
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

	dragonStatsRequest := DragonStatsRequest{}

	err := json.Unmarshal([]byte(request.Body), &dragonStatsRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 410,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: "1-" + err.Error(),
		}, err
	}

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("DragonStatsTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"dragon_name": {
				S: aws.String(dragonStatsRequest.Name),
			},
		},
	})
	fmt.Println("5")
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 411,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: "2-" + err.Error(),
		}, err
	}

	if result.Item == nil {
		msg := fmt.Sprintf("Could not find %s dragon", dragonStatsRequest.Name)

		return events.APIGatewayProxyResponse{
			StatusCode: 412,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: "3-" + msg,
		}, errors.New(msg)
	}

	dragonStats := DragonStats{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &dragonStats)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 413,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: "4-" + err.Error(),
		}, err
	}

	b, err := json.Marshal(dragonStats)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 414,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: "5-" + err.Error(),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(b),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}