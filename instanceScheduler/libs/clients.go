package libs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func getEC2Client() (*ec2.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return ec2.NewFromConfig(defaultConfig), nil
}

func getRDSClient() (*rds.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return rds.NewFromConfig(defaultConfig), nil
}

func getASGClient() (*autoscaling.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return autoscaling.NewFromConfig(defaultConfig), nil
}

func getSSMParameterStoreClient() (*ssm.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return ssm.NewFromConfig(defaultConfig), nil
}

