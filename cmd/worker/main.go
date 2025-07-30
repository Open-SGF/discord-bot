package main

import (
	"context"
	"discord-bot/pkg/worker"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"time"
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
