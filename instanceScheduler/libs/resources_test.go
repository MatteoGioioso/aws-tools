package libs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/onsi/gomega"
	"testing"
)

type mockSSM struct {}

func (m mockSSM) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)  {
	return &ssm.GetParameterOutput{
		Parameter:      &types.Parameter{
			Value:            nil,
		},
	}, nil
}

func TestSSM_GetConfig(t *testing.T) {
	g := gomega.NewWithT(t)

	t.Run("", func(t *testing.T) {
		s := SSM{
			client: mockSSM{},
		}
		got, err := s.GetConfig()
		if err != nil {
			t.Error(err)
		}

		g.Expect(got.Period.Pattern).To(gomega.Equal(officeHours))
	})
}
