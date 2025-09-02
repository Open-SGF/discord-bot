package main

import (
	"context"
	"discord-bot/pkg/worker"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var service *worker.Service

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newService, err := worker.InitService(ctx)

	if err != nil {
		log.Fatal(err)
	}

	service = newService
}

func main() {
	lambda.Start(service.Execute)
}
