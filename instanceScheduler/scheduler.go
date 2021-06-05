package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest() error {
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}