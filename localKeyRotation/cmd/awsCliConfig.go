package cmd

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"time"
)

const defaultAwsConfigFile = "/.aws/credentials"

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
		awsConfigFile = fmt.Sprintf("%v%v", home, defaultAwsConfigFile)
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

// StashOldCredentials Stash old credentials just in case we might needed it after de-activation
func (a *AwsCliConfig) StashOldCredentials() (string, error) {
	a.configFile.Section(a.currentProfile)

	newFile := ini.Empty()
	if err := newFile.NewSections(a.currentProfile); err != nil {
		return "", err
	}

	keyId := a.configFile.Section(a.currentProfile).Key(awsAccessKeyId).String()
	secretKey := a.configFile.Section(a.currentProfile).Key(awsSecretAccessKey).String()

	newFile.Section(a.currentProfile).Key(awsAccessKeyId).SetValue(keyId)
	newFile.Section(a.currentProfile).Key(awsSecretAccessKey).SetValue(secretKey)

	dir, _ := filepath.Split(a.awsSharedCredentialsFilePath)
	join := filepath.Join(dir, fmt.Sprintf("inactive-config-%v", time.Now().Unix()))

	if err := newFile.SaveTo(join); err != nil {
		return "", err
	}

	return join, nil
}