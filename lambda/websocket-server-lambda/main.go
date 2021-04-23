package main

import (
	"context"
	"log"

	"websocket-server-lambda/connection"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	rc := event.RequestContext
	log.Println(rc.RouteKey)
	switch rk := rc.RouteKey; rk {
	case "$connect":
		log.Println("Connecting...")
		connectionId := event.RequestContext.ConnectionID
		err := connection.HandleConnect(connectionId)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}, nil
		}
	case "$disconnect":
		log.Println("Disconnecting...")
		connectionId := event.RequestContext.ConnectionID
		err := connection.HandleDisconnect(connectionId)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}, nil
		}
	case "subscribe":
		log.Println("Subscribing...")
		connectionId := event.RequestContext.ConnectionID
		err := connection.HandleSubscription(connectionId, event, true)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}, nil
		}
	case "unsubscribe":
		log.Println("Unsubscribing...")
		connectionId := event.RequestContext.ConnectionID
		err := connection.HandleSubscription(connectionId, event, false)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}, nil
		}
	case "$default":
		log.Printf("Unknown RouteKey %v", rk)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(handler)
}