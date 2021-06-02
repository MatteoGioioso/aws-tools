package cmd

import (
	"github.com/onsi/gomega"
	"testing"
)

func Test_cacheOldCredentials(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	t.Run("should cache the credentials correctly", func(t *testing.T) {
		file, err := parseAWSCredentialsFile("../test_credentials")
		if err != nil {
			t.Error(err)
		}

		if err := cacheOldCredentials(file); err != nil {
			t.Error(err)
		}

		g.Expect(len(temporaryCacheCredentials)).To(gomega.Equal(2))
		g.Expect(temporaryCacheCredentials["default"].accessKeyId).To(gomega.Equal("SOMEFAK3KE4CCESSKEY"))
		g.Expect(temporaryCacheCredentials["default"].secretAccessKey).To(gomega.Equal("8sOm3rAndom+secre7K3yDonotTry1t"))
	})

}
