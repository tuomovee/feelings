package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tuomovee/feelings/pkg/db"
	"log"
	"os"
)

type pollResponse struct {
	Feeling db.Feeling `json:"feeling"`
}

var errorLog log.Logger

func init() {
	errorLog = *log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
}

// Handler for PUT
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var response pollResponse
	if err := json.Unmarshal([]byte(request.Body), &response); err != nil {
		errorLog.Printf("Failed to parse request: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, nil
	}

	if err := db.InsertPollResult(response.Feeling); err != nil {
		errorLog.Printf("Failed to insert poll response: %d", response.Feeling)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
