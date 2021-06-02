package cmd

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

type AwsCliConfig struct {
	awsSharedCredentialsFilePath string
	currentProfile string
	configFile *ini.File
}

func NewAwsCliConfig() (*AwsCliConfig, error) {
	a := &AwsCliConfig{}
	profile := a.GetCurrentProfile()
	path, err := a.GetAwsSharedCredentialFilePath()
	if err != nil {
		return nil, err
	}

	a.currentProfile = profile
	a.awsSharedCredentialsFilePath = path

	return a, nil
}

func (a *AwsCliConfig) GetAwsSharedCredentialFilePath() (string, error) {
	awsConfigFile := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if awsConfigFile == "" {
		awsConfigFile = fmt.Sprintf("%v/%v", home, defaultAwsConfigFile)
	}

	return awsConfigFile, err
}

func (a *AwsCliConfig) GetCurrentProfile() string {
	awsProfile := os.Getenv("AWS_PROFILE")
	if awsProfile == "" {
		awsProfile = "default"
	}

	return awsProfile
}

func (a *AwsCliConfig) SetNewIAMCredentials(newCredentials IAMCredentials) (*ini.File, error) {
	a.configFile.Section(a.currentProfile).Key(awsAccessKeyId).SetValue(newCredentials.accessKeyId)
	a.configFile.Section(a.currentProfile).Key(awsSecretAccessKey).SetValue(newCredentials.secretAccessKey)
	if err := a.configFile.SaveTo(a.awsSharedCredentialsFilePath); err != nil {
		return nil, err
	}

	return a.configFile, nil
}

func (a AwsCliConfig) SaveConfig() error {
	return a.configFile.SaveTo(a.awsSharedCredentialsFilePath)
}

func (a *AwsCliConfig) ParseConfig() (*ini.File, error) {
	cfg, err := ini.Load(a.awsSharedCredentialsFilePath)
	if err != nil {
		return nil, err
	}

	// Remove DEFAULT section, somehow it is parsed by default even if
	// it does not exist
	cfg.DeleteSection("DEFAULT")
	a.configFile = cfg

	return cfg, nil
}