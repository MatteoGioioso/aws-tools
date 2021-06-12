package libs

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/smithy-go"
	"log"
	"os"
	"strconv"
	"strings"
)

type ResourceState struct {
	State        string `json:"state"`
	ResourceType string `json:"type"`
}
type ResourcesState map[string]ResourceState

type ResourceClientArgs struct {
	Identifiers    []string
	ResourcesState *ResourcesState
}

type ResourcesClient interface {
	Sleep(args ResourceClientArgs) (ResourcesState, error)
	WakeUp(args ResourceClientArgs) (ResourcesState, error)
}

type ResourceClientsFactory map[string]ResourcesClient

func NewResourceClientsFactory() (ResourceClientsFactory, error) {
	newEC2, err := NewEC2()
	if err != nil {
		return nil, err
	}

	newRDS, err := NewRDS()
	if err != nil {
		return nil, err
	}

	newAurora, err := NewAurora()
	if err != nil {
		return nil, err
	}

	newFargate, err := NewFargate()
	if err != nil {
		return nil, err
	}

	factory := make(ResourceClientsFactory)
	factory[elasticComputeCloud] = newEC2
	factory[relationalDatabase] = newRDS
	factory[aurora] = newAurora
	factory[fargate] = newFargate

	return factory, err
}

// ========================================== EC2 ==================================== //

type EC2TypeClient interface {
	StopInstances(ctx context.Context, params *ec2.StopInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StopInstancesOutput, error)
	StartInstances(ctx context.Context, params *ec2.StartInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
}

type EC2 struct {
	client EC2TypeClient
}

func NewEC2() (*EC2, error) {
	e := &EC2{}
	client, err := e.getClient()
	if err != nil {
		return nil, err
	}
	e.client = client

	return e, nil
}

func (e *EC2) getClient() (*ec2.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return ec2.NewFromConfig(defaultConfig), nil
}

func (e *EC2) Sleep(args ResourceClientArgs) (ResourcesState, error) {
	params := &ec2.StopInstancesInput{
		InstanceIds: args.Identifiers,
	}
	instances, err := e.client.StopInstances(context.Background(), params)
	if err != nil {
		return checkError(err, args)
	}

	for _, instance := range instances.StoppingInstances {
		(*args.ResourcesState)[*instance.InstanceId] = ResourceState{
			State:        string(instance.CurrentState.Name),
			ResourceType: elasticComputeCloud,
		}
	}

	return *args.ResourcesState, nil
}

func (e *EC2) WakeUp(args ResourceClientArgs) (ResourcesState, error) {
	params := &ec2.StartInstancesInput{
		InstanceIds: args.Identifiers,
	}
	instances, err := e.client.StartInstances(context.Background(), params)
	if err != nil {
		return checkError(err, args)
	}

	for _, instance := range instances.StartingInstances {
		(*args.ResourcesState)[*instance.InstanceId] = ResourceState{
			State:        string(instance.CurrentState.Name),
			ResourceType: elasticComputeCloud,
		}
	}

	return *args.ResourcesState, nil
}

// ========================================== RDS ==================================== //

type RDSTypeClient interface {
	StopDBInstance(ctx context.Context, params *rds.StopDBInstanceInput, optFns ...func(*rds.Options)) (*rds.StopDBInstanceOutput, error)
	StartDBInstance(ctx context.Context, params *rds.StartDBInstanceInput, optFns ...func(*rds.Options)) (*rds.StartDBInstanceOutput, error)
}

type RDS struct {
	client RDSTypeClient
}

func NewRDS() (*RDS, error) {
	r := &RDS{}
	client, err := r.getClient()
	if err != nil {
		return nil, err
	}
	r.client = client
	return r, err
}

func (r *RDS) getClient() (*rds.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return rds.NewFromConfig(defaultConfig), nil
}

func (r *RDS) Sleep(args ResourceClientArgs) (ResourcesState, error) {
	params := &rds.StopDBInstanceInput{
		DBInstanceIdentifier: aws.String(args.Identifiers[0]),
	}
	instance, err := r.client.StopDBInstance(context.Background(), params)
	if err != nil {
		return checkError(err, args)
	}

	(*args.ResourcesState)[*instance.DBInstance.DBInstanceIdentifier] = ResourceState{
		State:        *instance.DBInstance.DBInstanceStatus,
		ResourceType: relationalDatabase,
	}

	return *args.ResourcesState, nil
}

func (r *RDS) WakeUp(args ResourceClientArgs) (ResourcesState, error) {
	params := &rds.StartDBInstanceInput{
		DBInstanceIdentifier: aws.String(args.Identifiers[0]),
	}
	instance, err := r.client.StartDBInstance(context.Background(), params)
	if err != nil {
		return checkError(err, args)
	}

	(*args.ResourcesState)[*instance.DBInstance.DBInstanceIdentifier] = ResourceState{
		State:        *instance.DBInstance.DBInstanceStatus,
		ResourceType: relationalDatabase,
	}

	return *args.ResourcesState, nil
}


// ========================================== AURORA ==================================== //

type AuroraTypeClient interface {
	StopDBCluster(ctx context.Context, params *rds.StopDBClusterInput, optFns ...func(*rds.Options)) (*rds.StopDBClusterOutput, error)
	StartDBCluster(ctx context.Context, params *rds.StartDBClusterInput, optFns ...func(*rds.Options)) (*rds.StartDBClusterOutput, error)
}
type Aurora struct {
	client AuroraTypeClient
}

func NewAurora() (*Aurora, error) {
	r := &Aurora{}
	client, err := r.getClient()
	if err != nil {
		return nil, err
	}
	r.client = client
	return r, err
}

func (r *Aurora) getClient() (*rds.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return rds.NewFromConfig(defaultConfig), nil
}

func (r *Aurora) Sleep(args ResourceClientArgs) (ResourcesState, error) {
	params := &rds.StopDBClusterInput{
		DBClusterIdentifier: aws.String(args.Identifiers[0]),
	}
	instance, err := r.client.StopDBCluster(context.Background(), params)
	if err != nil {
		return checkError(err, args)
	}

	(*args.ResourcesState)[*instance.DBCluster.DBClusterIdentifier] = ResourceState{
		State:        *instance.DBCluster.Status,
		ResourceType: aurora,
	}

	return *args.ResourcesState, nil
}

func (r *Aurora) WakeUp(args ResourceClientArgs) (ResourcesState, error) {
	params := &rds.StartDBClusterInput{
		DBClusterIdentifier: aws.String(args.Identifiers[0]),
	}
	instance, err := r.client.StartDBCluster(context.Background(), params)
	if err != nil {
		return checkError(err, args)
	}

	resourcesState := make(ResourcesState)
	resourcesState[*instance.DBCluster.DBClusterIdentifier] = ResourceState{
		State:        *instance.DBCluster.Status,
		ResourceType: aurora,
	}

	return resourcesState, nil
}

// ========================================== SSM ==================================== //

type SSMTypeClient interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

type SSM struct {
	client SSMTypeClient
}

func NewSSM() (*SSM, error) {
	s := &SSM{}
	client, err := s.getClient()
	if err != nil {
		return nil, err
	}
	s.client = client
	return s, err
}

func (s *SSM) getClient() (*ssm.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return ssm.NewFromConfig(defaultConfig), nil
}

func (s SSM) GetConfig() (*SchedulerConfig, error) {
	params := &ssm.GetParameterInput{
		Name: aws.String(os.Getenv("INSTANCE_SCHEDULER_CONFIG_NAME")),
	}
	parameter, err := s.client.GetParameter(context.Background(), params)
	if err != nil {
		return &SchedulerConfig{}, err
	}

	sc := &SchedulerConfig{}

	if err := json.Unmarshal([]byte(*parameter.Parameter.Value), sc); err != nil {
		return &SchedulerConfig{}, err
	}

	return sc, err
}

// ========================================== Fargate ==================================== //

type FargateTypeClient interface {
	UpdateService(ctx context.Context, params *ecs.UpdateServiceInput, optFns ...func(*ecs.Options)) (*ecs.UpdateServiceOutput, error)
}

type Fargate struct {
	client FargateTypeClient
}

func NewFargate() (*Fargate, error) {
	c := &Fargate{}
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}
	c.client = client
	return c, err
}

func (e Fargate) getClient() (*ecs.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return ecs.NewFromConfig(defaultConfig), nil
}

func (e Fargate) WakeUp(args ResourceClientArgs) (ResourcesState, error) {
	return e.update(args, 1)
}

func (e Fargate) Sleep(args ResourceClientArgs) (ResourcesState, error) {
	return e.update(args, 0)
}

func (e Fargate) update(args ResourceClientArgs, desiredCount int32) (ResourcesState, error) {
	clusterService := strings.Split(args.Identifiers[0], ":")
	params := &ecs.UpdateServiceInput{
		Service:            aws.String(clusterService[1]),
		Cluster:            aws.String(clusterService[0]),
		DesiredCount:       aws.Int32(desiredCount),
		ForceNewDeployment: true,
	}
	service, err := e.client.UpdateService(context.Background(), params)
	if err != nil {
		return nil, err
	}

	(*args.ResourcesState)[*service.Service.ServiceName] = ResourceState{
		State:        strconv.Itoa(int(service.Service.RunningCount)),
		ResourceType: fargate,
	}

	return *args.ResourcesState, nil
}

// ============================================== SNS =========================================== //

type MessageBus interface {
	Send(message string) error
}

type SNS struct {
	client *sns.Client
}

func NewSNS() (*SNS, error) {
	c := &SNS{}
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}
	c.client = client
	return c, err
}

func (e SNS) getClient() (*sns.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return sns.NewFromConfig(defaultConfig), nil
}

func (e SNS) Send(message string) error {
	params := &sns.PublishInput{
		Message:  aws.String(message),
		Subject:  aws.String("AWS Scheduler report status"),
		TopicArn: aws.String(os.Getenv("TOPIC_ARN")),
	}
	if _, err := e.client.Publish(context.Background(), params); err != nil {
		return err
	}

	return nil
}

func checkError(err error, args ResourceClientArgs) (ResourcesState, error) {
	var oe *smithy.OperationError
	if errors.As(err, &oe) {
		log.Printf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		if strings.Contains(oe.Operation(), "Stop") {
			(*args.ResourcesState)[args.Identifiers[0]] = ResourceState{
				State:        "Stopped",
				ResourceType: relationalDatabase,
			}
		}

		if strings.Contains(oe.Operation(), "Start") {
			(*args.ResourcesState)[args.Identifiers[0]] = ResourceState{
				State:        "Started",
				ResourceType: relationalDatabase,
			}
		}

		return *args.ResourcesState, nil
	}

	return nil, err
}