package cmd

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"strings"
)

type IAMCredentials struct {
	secretAccessKey string
	accessKeyId     string
	username        string
}

type AWSUtils struct {
	iamClient *iam.Client
	stsClient *sts.Client
}

func NewAWSUtils() (*AWSUtils, error) {
	stsClient, err := getStsClient()
	if err != nil {
		return nil, err
	}
	iamClient, err := getIamClient()
	if err != nil {
		return nil, err
	}
	return &AWSUtils{iamClient: iamClient, stsClient: stsClient}, err
}

func getStsClient() (*sts.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return sts.NewFromConfig(defaultConfig), nil
}

func getIamClient() (*iam.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return iam.NewFromConfig(defaultConfig), nil
}

func (a *AWSUtils) GetCurrentUsername() (string, error) {
	params := &sts.GetCallerIdentityInput{}
	identity, err := a.stsClient.GetCallerIdentity(context.Background(), params)
	if err != nil {
		return "", err
	}

	return strings.Split(*identity.Arn, "/")[1], nil
}

func (a *AWSUtils) GetNewKeys(username string) (IAMCredentials, error) {
	params := &iam.CreateAccessKeyInput{UserName: aws.String(username)}
	credentials, err := a.iamClient.CreateAccessKey(context.Background(), params)
	if err != nil {
		return IAMCredentials{}, nil
	}

	return IAMCredentials{
		secretAccessKey: *credentials.AccessKey.SecretAccessKey,
		accessKeyId:     *credentials.AccessKey.AccessKeyId,
		username:        *credentials.AccessKey.UserName,
	}, err
}

func (a *AWSUtils) DeactivateOldKeys(accessKeyId string, username string) error {
	params := &iam.UpdateAccessKeyInput{
		AccessKeyId: aws.String(accessKeyId),
		Status:      "Inactive",
		UserName:    aws.String(username),
	}

	if _, err := a.iamClient.UpdateAccessKey(context.Background(), params); err != nil {
		return err
	}

	return nil
}
