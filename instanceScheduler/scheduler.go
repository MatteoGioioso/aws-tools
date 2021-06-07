package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hirvitek/aws-tools/instanceScheduler/libs"
)

func HandleRequest() error {
	ssm, err := libs.NewSSM()
	if err != nil {
		return err
	}

	config, err := ssm.GetConfig()
	if err != nil {
		return err
	}

	factory, err := libs.NewResourceClientsFactory()
	if err != nil {
		return err
	}

	schedulerConfigClient := libs.NewSchedulerConfigClient(config.Period, config.TimeZone)
	status, err := libs.SchedulerService(config, schedulerConfigClient, factory)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", status)

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
