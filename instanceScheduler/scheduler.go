package main

import (
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

	sns, err := libs.NewSNS()
	if err != nil {
		return err
	}

	schedulerConfigClient := libs.NewSchedulerConfigClient(config)

	service := libs.SchedulerService{
		Config:                config,
		SchedulerConfigClient: schedulerConfigClient,
		Factory:               factory,
		MessageBus:            sns,
	}
	if _, err := service.Execute(); err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
