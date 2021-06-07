package libs

import (
	"fmt"
	"log"
)

func SchedulerService(
	config *SchedulerConfig,
	schedulerConfigClient *SchedulerConfigClient,
	factory ResourceClientsFactory,
) (ResourcesState, error) {
	resourcesState := make(ResourcesState)
	for id, resource := range config.Resources {
		wakeup, err := schedulerConfigClient.ShouldWakeup()
		if err != nil {
			return nil, err
		}
		log.Printf("Should resource %v %v wake up? %v", resource.Type, id, wakeup)

		client := factory[resource.Type]
		if wakeup {
			_, err = client.WakeUp(ResourceClientArgs{
				Identifiers:    []string{id},
				ResourcesState: &resourcesState,
			})
			if err != nil {
				return nil, err
			}
		} else {
			_, err = client.Sleep(ResourceClientArgs{
				Identifiers:    []string{id},
				ResourcesState: &resourcesState,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return resourcesState, nil
}

func PrintStatusReport(state ResourcesState) string {
	report := fmt.Sprintf("AWS Scheduler daily status report: \n\n")
	for id, resourceState := range state {
		if resourceState.ResourceType == fargate {
			report += fmt.Sprintf(
				"%v with id %v has tasks running count of: %v \n",
				resourceState.ResourceType,
				id,
				resourceState.State,
			)
		} else {
			report += fmt.Sprintf(
				"%v with id %v has status: %v \n",
				resourceState.ResourceType,
				id,
				resourceState.State,
			)
		}
	}

	fmt.Println(report)

	return report
}
