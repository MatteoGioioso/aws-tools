package libs

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go"
	"github.com/onsi/gomega"
	"io/ioutil"
	"testing"
)

type mockSSM struct{}

func (m mockSSM) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	file, err := ioutil.ReadFile("config_sample_1.json")
	if err != nil {
		return nil, err
	}

	return &ssm.GetParameterOutput{
		Parameter: &types.Parameter{
			Value: aws.String(string(file)),
		},
	}, nil
}

func TestSSM_GetConfig(t *testing.T) {
	g := gomega.NewWithT(t)

	t.Run("should parse config from JSON", func(t *testing.T) {
		s := SSM{
			client: mockSSM{},
		}
		got, err := s.GetConfig()
		if err != nil {
			t.Error(err)
		}

		g.Expect(got.Period.Pattern).To(gomega.Equal(officeHours))
		g.Expect(got.Report.SendReport).To(gomega.Equal(true))
		g.Expect(got.Report.Hour).To(gomega.Equal(8))
		g.Expect(len(got.Resources)).To(gomega.Equal(1))
		g.Expect(got.Resources["i-123"].Type).To(gomega.Equal(elasticComputeCloud))
	})
}

func Test_checkError(t *testing.T) {
	g := gomega.NewWithT(t)

	type args struct {
		err  error
		args ResourceClientArgs
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should swallow error and print status",
			args: args{
			err: &smithy.OperationError{
				ServiceID:     "RDS",
				OperationName: "StopDBInstance",
				Err:           errors.New("test Operation error"),
			},
			args: ResourceClientArgs{
				Identifiers:    []string{"rds-instance-id"},
				ResourcesState: &ResourcesState{},
			},
		},
			want: "Stopped",
			wantErr: false,
		},
		{
			name: "should swallow error and print status",
			args: args{
			err: errors.New("generic error"),
			args: ResourceClientArgs{
				Identifiers:    []string{"rds-instance-id"},
				ResourcesState: &ResourcesState{},
			},
		},
			want: "Stopped",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkError(tt.args.err, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr {
				g.Expect(err).To(gomega.Equal(errors.New("generic error")))
				return
			}

			g.Expect(got["rds-instance-id"].State).To(gomega.Equal(tt.want))
		})
	}
}
