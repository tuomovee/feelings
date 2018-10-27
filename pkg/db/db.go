// Package db implements storage and retreival of Feelings poll results and
// related data structures
package db

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"
import "github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
import "github.com/aws/aws-sdk-go-v2/aws/external"
import "github.com/aws/aws-sdk-go-v2/aws"
import "time"
import "log"
import "fmt"
import "os"

// Date format string for formatting or parsing yyyy-mm-dd dates
const DateLayout = "2006-01-02"

// Our scale of emotions
type Feeling int

const (
	VeryBad Feeling = iota
	Bad
	Good
	VeryGood
)

func (f Feeling) String() string {
	switch f {
	case VeryBad:
		return "very_bad"
	case Bad:
		return "bad"
	case Good:
		return "good"
	case VeryGood:
		return "very_good"
	}
	errorLog.Printf("No string value for feeling %d", f)
	return "failure"
}

// Poll statistics for a day
type PollResult struct {
	Date     string `json:"date"`
	VeryBad  int    `json:"very_bad"`
	Bad      int    `json:"bad"`
	Good     int    `json:"good"`
	VeryGood int    `json:"very_good"`
}

var (
	dbClient  *dynamodb.DynamoDB
	errorLog  log.Logger
	tableName *string
)

// Initialize package scope DynamoDB client and error loggers.
// Panics on failure since without proper config from env this package is useless.
func init() {
	errorLog = *log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)

	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("Unable to load AWS config:" + err.Error())
	}
	dbClient = dynamodb.New(config)

	if tableName = aws.String(os.Getenv("TABLE_NAME")); *tableName == "" {
		panic("Environment variable TABLE_NAME not set")
	}
}

// GetPollResult fetches poll results for a specific date
func GetPollResult(t time.Time) (PollResult, error) {
	var result PollResult

	key := t.Format(DateLayout)
	item, err := dbClient.GetItemRequest(&dynamodb.GetItemInput{
		TableName: tableName,
		Key: map[string]dynamodb.AttributeValue{
			"date": {
				S: &key,
			},
		},
	}).Send()
	if err != nil {
		errorLog.Printf("GetPollResult: GetItem failed: %v", err)
		return result, err
	}

	if err := dynamodbattribute.UnmarshalMap(item.Item, &result); err != nil {
		errorLog.Printf("GetPollResult: Failed to unmarshal DynamoDB item: %v", err)
	}
	return result, err
}

// InsertPollResult saves a poll response
func InsertPollResult(f Feeling) error {
	if f < VeryBad || f > VeryGood {
		return fmt.Errorf("Feeling out of range: %d", f)
	}

	key := time.Now().Format(DateLayout)
	updateExpression := fmt.Sprintf("ADD %s :v", f.String())
	_, err := dbClient.UpdateItemRequest(&dynamodb.UpdateItemInput{
		TableName:        tableName,
		UpdateExpression: &updateExpression,
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":v": {
				N: aws.String("1"),
			},
		},
		Key: map[string]dynamodb.AttributeValue{
			"date": {
				S: &key,
			},
		},
	}).Send()
	if err != nil {
		errorLog.Printf("InsertPollResult: Failed to update %s with feeling %v: %v", *tableName, f, err)
		return err
	}
	return nil
}
