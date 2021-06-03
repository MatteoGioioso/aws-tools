package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

const awsAccessKeyId = "aws_access_key_id"
const awsSecretAccessKey = "aws_secret_access_key"

var rootCmd = &cobra.Command{
	Use:   "aws-key-rotation",
	Short: "Automatically rotate your local AWS credentials",
	Long:  "Automatically rotate your local AWS credentials locally",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() error {
	// Initialize AWS Config
	awsCliConfig, err := NewAwsCliConfig()
	if err != nil {
		return err
	}
	// Parse the file
	_, err = awsCliConfig.ParseConfig()
	if err != nil {
		return err
	}

	if os.Getenv("LKR_BACKUP_OLD_KEYS") == "no" {

	} else {
		// Save a copy of the credentials, just in case something happens
		if _, _, err := awsCliConfig.StashOldCredentials(); err != nil {
			return err
		}
	}

	oldIamCredentials := awsCliConfig.GetIAMCredentials()

	awsUtils, err := NewAWSUtils()
	if err != nil {
		return err
	}

	username, err := awsUtils.GetCurrentUsername()
	if err != nil {
		return err
	}

	newIamCredentials, err := awsUtils.GetNewKeys(username)
	if err != nil {
		return err
	}

	if _, err := awsCliConfig.SetIAMCredentials(newIamCredentials); err != nil {
		return err
	}

	if err := awsCliConfig.SaveConfig(); err != nil {
		return err
	}

	if os.Getenv("LKR_DELETE_OLD_KEYS") == "yes" {
		if err := awsUtils.DeleteCredentials(oldIamCredentials.accessKeyId, username); err != nil {
			return err
		}
	} else {
		if err := awsUtils.DeactivateOldKeys(oldIamCredentials.accessKeyId, username); err != nil {
			return err
		}
	}

	return rootCmd.Execute()
}
