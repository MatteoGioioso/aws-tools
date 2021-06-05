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
	for id, resource := range config.Resources {
		wakeup, err := schedulerConfigClient.ShouldWakeup()
		if err != nil {
			return err
		}

		client := factory[resource.Type]
		if wakeup {
			state, err := client.WakeUp(libs.ResourceClientArgs{
				Identifiers:  []string{id},
			})
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", state)
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}