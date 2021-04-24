package websocket

import (
	"errors"
	"log"
	"technical-analysis-lambda/dynamosymbol"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func ReadConnections(symbol string, symbolRate *[]byte) (*[]dynamosymbol.Connection, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	filt := expression.Name(symbol).Equal(expression.Value("true"))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, errors.New("could not create filter expression")
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String("WebsocketConnectionsTable"),
	}

	result, err := svc.Scan(params)
	if err != nil {
		log.Println("Scan Failed LMAO", err)
		return nil, errors.New("DynamoDB scan failed")
	}

	connectionList := make([]dynamosymbol.Connection, *result.Count)
	if err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &connectionList); err != nil {
		return nil, errors.New("connection list unmarshaling error")
	}

	return &connectionList, nil
}