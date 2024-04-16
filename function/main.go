package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	var bucketname = os.Getenv("BUCKET_NAME")
	fmt.Printf("Bucket name = %s.\n", bucketname)

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf(" ++ %s: %s\n", key, value)
	}

	return events.LambdaFunctionURLResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(handler)
}
