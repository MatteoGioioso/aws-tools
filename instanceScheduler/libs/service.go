package libs

import (
	"fmt"
	"log"
)

type SchedulerService struct {
	Config                *SchedulerConfig
	SchedulerConfigClient *SchedulerConfigClient
	Factory               ResourceClientsFactory
	MessageBus            MessageBus
}

func (s SchedulerService) Execute() (ResourcesState, error) {
	resourcesState := make(ResourcesState)
	for id, resource := range s.Config.Resources {
		wakeup, err := s.SchedulerConfigClient.ShouldWakeup()
		if err != nil {
			return nil, err
		}
		log.Printf("Should resource %v %v wake up? %v", resource.Type, id, wakeup)

		client := s.Factory[resource.Type]
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

	if _, err := s.sendReport(resourcesState); err != nil {
		return nil, err
	}

	return resourcesState, nil
}

func (s SchedulerService) printStatusReport(state ResourcesState) string {
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

	return report
}

func (s SchedulerService) sendReport(state ResourcesState) (string, error) {
	ok, err := s.SchedulerConfigClient.ShouldSendReport()
	if err != nil {
		return "", err
	}

	var report string
	if ok {
		report = s.printStatusReport(state)
		if err := s.MessageBus.Send(report); err != nil {
			return "", err
		}
	}

	return report, nil
}
