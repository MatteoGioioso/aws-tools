package cmd

import (
	"fmt"
	"github.com/onsi/gomega"
	"gopkg.in/ini.v1"
	"os"
	"testing"
)

func cleanUp(t *testing.T) {
	if err := os.Unsetenv("AWS_SHARED_CREDENTIALS_FILE"); err != nil {
		t.Error(err)
	}

	if err := os.Unsetenv("AWS_PROFILE"); err != nil {
		t.Error(err)
	}

	if err := os.Unsetenv("AWS_CONFIG_FILE"); err != nil {
		t.Error(err)
	}
}

func TestAwsCliConfig_StashOldCredentials(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	t.Run("should stash the credentials correctly", func(t *testing.T) {
		if err := os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "../test_credentials"); err != nil {
			t.Error(err)
		}

		config, err := NewAwsCliConfig()
		if err != nil {
			t.Error(err)
		}

		if _, err := config.ParseConfig(); err != nil {
			t.Error(err)
		}

		dir, err := config.StashOldCredentials()
		if err != nil {
			t.Error(err)
		}

		file, err := ini.Load(dir)
		if err != nil {
			t.Error(err)
		}

		g.Expect(file.Section("default").Key(awsAccessKeyId).String()).
			To(gomega.Equal("SOMEFAK3KE4CCESSKEY"))
		g.Expect(file.Section("default").Key(awsSecretAccessKey).String()).
			To(gomega.Equal("8sOm3rAndom+secre7K3yDonotTry1t"))

		if err := os.Remove(dir); err != nil {
			t.Error()
		}

		cleanUp(t)
	})
}

func TestAwsCliConfig_SetNewIAMCredentials(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	t.Run("should add the new credentials to a file", func(t *testing.T) {
		if err := os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "../test_credentials"); err != nil {
			t.Error(err)
		}

		config, err := NewAwsCliConfig()
		if err != nil {
			t.Error(err)
		}

		if _, err := config.ParseConfig(); err != nil {
			t.Error(err)
		}

		newCredentials := IAMCredentials{
			secretAccessKey: "123",
			accessKeyId:     "ABC",
			username:        "aws-user",
		}

		newFile, err := config.SetNewIAMCredentials(newCredentials)
		if err != nil {
			t.Error(err)
		}

		g.Expect(newFile.Section("default").Key(awsAccessKeyId).String()).To(gomega.Equal("ABC"))
		g.Expect(newFile.Section("default").Key(awsSecretAccessKey).String()).To(gomega.Equal("123"))

		cleanUp(t)
	})
}

func TestNewAwsCliConfig(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	t.Run("should get profile and path name correctly with default", func(t *testing.T) {
		config, err := NewAwsCliConfig()
		if err != nil {
			t.Error(err)
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Error(err)
		}

		g.Expect(config.currentProfile).To(gomega.Equal("default"))
		g.Expect(config.awsSharedCredentialsFilePath).
			To(gomega.Equal(fmt.Sprintf("%v%v", homeDir, "/.aws/credentials")))
	})

	t.Run("should get profile and path with non default", func(t *testing.T) {
		if err := os.Setenv("AWS_PROFILE", "second-profile"); err != nil {
			t.Error(err)
		}

		if err := os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "../test_credentials"); err != nil {
			t.Error(err)
		}

		config, err := NewAwsCliConfig()
		if err != nil {
			t.Error(err)
		}

		if _, err := config.ParseConfig(); err != nil {
			t.Error(err)
		}

		g.Expect(config.currentProfile).To(gomega.Equal("second-profile"))
		g.Expect(config.awsSharedCredentialsFilePath).To(gomega.Equal("../test_credentials"))
		g.Expect(config.configFile.Section("second-profile").Key(awsAccessKeyId).String()).
			To(gomega.Equal("SOMEFAK3KE4CCESSKEY2"))
		g.Expect(config.configFile.Section("second-profile").Key(awsSecretAccessKey).String()).
			To(gomega.Equal("8sOm3rAndom+secre7K3yDonotTry1t2"))
		cleanUp(t)
	})

}