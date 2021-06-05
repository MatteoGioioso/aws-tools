package aws

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hirvitek/aws-tools/instanceScheduler/libs"
)

type EC2 struct {
	client *ec2.Client
}

func NewEC2() (*EC2, error) {
	e := &EC2{}
	client, err := e.getEC2Client()
	if err != nil {
		return nil, err
	}
	e.client = client

	return e, nil
}

func (e *EC2) getEC2Client() (*ec2.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return ec2.NewFromConfig(defaultConfig), nil
}

func (e *EC2) Sleep(identifiers []string) (libs.ResourcesState, error) {
	params := &ec2.StopInstancesInput{
		InstanceIds: identifiers,
	}
	instances, err := e.client.StopInstances(context.Background(), params)
	if err != nil {
		return nil, err
	}

	resourcesState := make(libs.ResourcesState)
	for _, instance := range instances.StoppingInstances {
		resourcesState[*instance.InstanceId] = libs.ResourceState{
			State: string(instance.CurrentState.Name),
		}
	}

	return resourcesState, nil
}

func (e *EC2) WakeUp(identifiers []string) (libs.ResourcesState, error) {
	params := &ec2.StartInstancesInput{
		InstanceIds: identifiers,
	}
	instances, err := e.client.StartInstances(context.Background(), params)
	if err != nil {
		return nil, err
	}

	resourcesState := make(libs.ResourcesState)
	for _, instance := range instances.StartingInstances {
		resourcesState[*instance.InstanceId] = libs.ResourceState{
			State: string(instance.CurrentState.Name),
		}
	}

	return resourcesState, nil
}

type RDS struct {
	client *rds.Client
}

func NewRDS() (*RDS, error) {
	r := &RDS{}
	client, err := r.getRDSClient()
	if err != nil {
		return nil, err
	}
	r.client = client
	return r, err
}

func (r *RDS) getRDSClient() (*rds.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return rds.NewFromConfig(defaultConfig), nil
}

func (r *RDS) Sleep(identifier string, resourceType string) (libs.ResourcesState, error) {
	if resourceType == "RDSCluster" {
		params := &rds.StopDBClusterInput{
			DBClusterIdentifier: aws.String(identifier),
		}
		instance, err := r.client.StopDBCluster(context.Background(), params)
		if err != nil {
			return nil, err
		}

		resourcesState := make(libs.ResourcesState)
		resourcesState[*instance.DBCluster.DBClusterIdentifier] = libs.ResourceState{State: *instance.DBCluster.Status}

		return resourcesState, nil
	} else {
		params := &rds.StopDBInstanceInput{
			DBInstanceIdentifier: aws.String(identifier),
		}
		instance, err := r.client.StopDBInstance(context.Background(), params)
		if err != nil {
			return nil, err
		}

		resourcesState := make(libs.ResourcesState)
		resourcesState[*instance.DBInstance.DBInstanceIdentifier] = libs.ResourceState{State: *instance.DBInstance.DBInstanceStatus}

		return resourcesState, nil
	}
}

func (r *RDS) WakeUp(identifier string, resourceType string) (libs.ResourcesState, error) {
	if resourceType == "RDSCluster" {
		params := &rds.StartDBClusterInput{
			DBClusterIdentifier: aws.String(identifier),
		}
		instance, err := r.client.StartDBCluster(context.Background(), params)
		if err != nil {
			return nil, err
		}

		resourcesState := make(libs.ResourcesState)
		resourcesState[*instance.DBCluster.DBClusterIdentifier] = libs.ResourceState{State: *instance.DBCluster.Status}

		return resourcesState, nil
	} else {
		params := &rds.StartDBInstanceInput{
			DBInstanceIdentifier: aws.String(identifier),
		}
		instance, err := r.client.StartDBInstance(context.Background(), params)
		if err != nil {
			return nil, err
		}

		resourcesState := make(libs.ResourcesState)
		resourcesState[*instance.DBInstance.DBInstanceIdentifier] = libs.ResourceState{State: *instance.DBInstance.DBInstanceStatus}

		return resourcesState, nil
	}
}

type SSM struct {
	client *ssm.Client
}

func NewSSM() (*SSM, error) {
	s := &SSM{}
	client, err := s.getSSMParameterStoreClient()
	if err != nil {
		return nil, err
	}
	s.client = client
	return s, err
}

func (s *SSM) getSSMParameterStoreClient() (*ssm.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return ssm.NewFromConfig(defaultConfig), nil
}

func (s SSM) GetConfig() (*libs.SchedulerConfig, error) {
	params := &ssm.GetParameterInput{
		Name:           aws.String("INSTANCE_SCHEDULER_CONFIG_NAME"),
	}
	parameter, err := s.client.GetParameter(context.Background(), params)
	if err != nil {
		return &libs.SchedulerConfig{}, err
	}

	sc := &libs.SchedulerConfig{}

	if err := json.Unmarshal([]byte(*parameter.Parameter.Value), sc); err != nil {
		return &libs.SchedulerConfig{}, err
	}

	return sc, err
}
