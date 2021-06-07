package libs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/onsi/gomega"
	"strings"
	"testing"
	"time"
)

type ec2ClientMock struct{}

func (e ec2ClientMock) StopInstances(ctx context.Context, params *ec2.StopInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StopInstancesOutput, error) {
	return &ec2.StopInstancesOutput{
		StoppingInstances: []types.InstanceStateChange{
			{
				CurrentState: &types.InstanceState{
					Code: nil,
					Name: "stopped",
				},
				InstanceId:    aws.String(params.InstanceIds[0]),
				PreviousState: nil,
			},
		},
	}, nil
}

func (e ec2ClientMock) StartInstances(ctx context.Context, params *ec2.StartInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error) {
	return &ec2.StartInstancesOutput{
		StartingInstances: []types.InstanceStateChange{
			{
				CurrentState: &types.InstanceState{
					Code: nil,
					Name: "running",
				},
				InstanceId:    aws.String(params.InstanceIds[0]),
				PreviousState: nil,
			},
		},
	}, nil
}

type fargateClientMock struct{}

func (f fargateClientMock) UpdateService(ctx context.Context, params *ecs.UpdateServiceInput, optFns ...func(*ecs.Options)) (*ecs.UpdateServiceOutput, error) {
	var status string
	if *params.DesiredCount == 0 {
		status = "INACTIVE"
	} else {
		status = "ACTIVE"
	}

	return &ecs.UpdateServiceOutput{
		Service: &ecsTypes.Service{
			RunningCount: *params.DesiredCount,
			ServiceName:  params.Service,
			Status:       aws.String(status),
		},
	}, nil
}

type messageBusMock struct{}

func (b messageBusMock) Send(message string) error {
	return nil
}

func Test_scheduler(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	t.Run("should schedule resources as wake and send the report", func(t *testing.T) {
		config := &SchedulerConfig{
			Report: Report{
				SendReport: true,
				Hour:       9,
			},
			Period: Period{
				Pattern: officeHours,
			},
			TimeZone: "Europe/Helsinki",
			Resources: map[string]Resource{
				"i-123": {
					Type:       elasticComputeCloud,
					Identifier: "i-123",
				},
				"i-456": {
					Type:       elasticComputeCloud,
					Identifier: "i-456",
				},
				"cluster-name:service-name": {
					Type:       fargate,
					Identifier: "cluster-name:service-name",
				},
			},
		}
		client := &SchedulerConfigClient{
			Config: config,
			now: func() time.Time {
				return insideOfficeHourAndDay
			},
		}

		factory := ResourceClientsFactory{}
		factory[elasticComputeCloud] = &EC2{client: ec2ClientMock{}}
		factory[fargate] = &Fargate{client: fargateClientMock{}}

		service := SchedulerService{
			Config:                config,
			SchedulerConfigClient: client,
			Factory:               factory,
			MessageBus:            messageBusMock{},
		}

		got, err := service.Execute()
		if err != nil {
			t.Error(err)
			return
		}

		report, err := service.sendReport(got)
		if err != nil {
			t.Error(err)
		}

		g.Expect(len(got)).To(gomega.Equal(3))
		g.Expect(got["i-123"].State).To(gomega.Equal("running"))
		g.Expect(got["service-name"].State).To(gomega.Equal("1"))
		g.Expect(strings.Contains(report, "EC2 with id i-456 has status: running")).To(gomega.Equal(true))
		g.Expect(strings.Contains(report, "EC2 with id i-123 has status: running")).To(gomega.Equal(true))
		g.Expect(strings.Contains(report, "Fargate with id service-name has tasks running count of: 1")).To(gomega.Equal(true))
	})

	t.Run("should schedule resources as sleeping and not print the report", func(t *testing.T) {
		config := &SchedulerConfig{
			Report: Report{
				SendReport: true,
				Hour:       9,
			},
			Period: Period{
				Pattern: officeHours,
			},
			TimeZone: "Europe/Helsinki",
			Resources: map[string]Resource{
				"i-123": {
					Type:       elasticComputeCloud,
					Identifier: "i-123",
				},
				"i-456": {
					Type:       elasticComputeCloud,
					Identifier: "i-456",
				},
				"cluster-name:service-name": {
					Type:       fargate,
					Identifier: "cluster-name:service-name",
				},
			},
		}
		client := &SchedulerConfigClient{
			Config: config,
			now: func() time.Time {
				return outsideOfficeHourButNotDay
			},
		}

		factory := ResourceClientsFactory{}
		factory[elasticComputeCloud] = &EC2{client: ec2ClientMock{}}
		factory[fargate] = &Fargate{client: fargateClientMock{}}

		service := SchedulerService{
			Config:                config,
			SchedulerConfigClient: client,
			Factory:               factory,
			MessageBus:            messageBusMock{},
		}

		got, err := service.Execute()
		if err != nil {
			t.Error(err)
			return
		}

		report, err := service.sendReport(got)
		if err != nil {
			t.Error(err)
		}

		g.Expect(len(got)).To(gomega.Equal(3))
		g.Expect(got["i-123"].State).To(gomega.Equal("stopped"))
		g.Expect(got["service-name"].State).To(gomega.Equal("0"))
		g.Expect(report).To(gomega.Equal(""))
	})
}
