package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tuomovee/feelings/pkg/db"
	"log"
	"os"
	"time"
)

var errorLog log.Logger

func init() {
	errorLog = *log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
}

// Handler for GET {date}
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var date time.Time
	if dateResource := request.PathParameters["date"]; dateResource != "" {
		var err error
		date, err = time.Parse(db.DateLayout, dateResource)
		if err != nil {
			errorLog.Printf("Failed to parse date query parameter %s: %v ", dateResource, err)
			return events.APIGatewayProxyResponse{
				Body:       fmt.Sprintf("Error: %v", err),
				StatusCode: 400,
			}, nil
		}
	} else {
		// Something is wrong in API gateway configuration since we got a request
		// without date path parameter
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: 400,
		}, nil
	}

	r, err := db.GetPollResult(date)
	if err != nil {
		errorLog.Printf("Failed to fetch results: %v", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error: %v", err),
			StatusCode: 503,
		}, err
	}

	resultJson, err := json.Marshal(r)
	if err != nil {
		errorLog.Printf("Failed to serialize poll results %v: %v", resultJson, err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error: %v", err),
			StatusCode: 500,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(resultJson),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
