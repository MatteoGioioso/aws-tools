package libs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
)

type ResourceState struct {
	State string `json:"state"`
}
type ResourcesState map[string]ResourceState

func getASGClient() (*autoscaling.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return autoscaling.NewFromConfig(defaultConfig), nil
}


