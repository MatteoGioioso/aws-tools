package libs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
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
