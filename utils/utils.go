package utils

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
	"net/http"
	"io/ioutil"
)

var (
	logger   = logrus.WithField("module", "utils")
)

type Item struct {
	Node_type string`json:"node_type"`
	Ip string`json:"ip"`
}


func SetManagerDbInfo(db_sess *session.Session, lockTable, localIp string) (bool, error){
	logger.Debug("Attempting to update DB record")
	svc := dynamodb.New(db_sess)

	item := Item{
		Node_type: "primary_manager",
		Ip: localIp,
	}

	av, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		logger.Fatalf("Failed to marshal record, %v", err)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(lockTable),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		return false, err
	}

	return false, nil

}

func FetchToken(ip string, tokenType string) (string, error) {
	host := fmt.Sprintf("http://%s:9024/token/%s", ip, tokenType)

	logger.Debugf("Token from %s", host)
	response, err := http.Get(host)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(contents), nil

}

func FetchManagerIp(db_sess *session.Session, lockTable string) (string, error) {
	logger.Debug("Fetching manager info from DB")
	svc := dynamodb.New(db_sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(lockTable),
		Key: map[string]*dynamodb.AttributeValue{
			"node_type":
			{
				S: aws.String("primary_manager"),
			},
		},
	})

	if err != nil {
		logger.Fatalf("DynamoDB error, %v ", err)
	}

	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		logger.Fatalf("Failed to unmarshal Record, %v", err)
	}

	logger.Debugf("Found %v in DB", item)

	return item.Ip, nil

}