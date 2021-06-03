package cmd

import (
	"github.com/spf13/cobra"
	"log"
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
	oldIamCredentials := awsCliConfig.GetIAMCredentials()

	log.Printf(
		"[LKR] old credentials parsed with profile: %v and access key id: %v \n",
		awsCliConfig.currentProfile,
		oldIamCredentials.accessKeyId,
	)

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
	log.Printf(
		"[LKR] New IAM keys created with user: %v and access key: %v \n",
		username,
		newIamCredentials.accessKeyId,
	)

	if os.Getenv("LKR_BACKUP_OLD_KEYS") == "no" {
		log.Println("[LKR] old keys NOT backed up")
	} else {
		// Save a copy of the credentials, just in case something happens
		_, dir, err := awsCliConfig.StashOldCredentials()
		if err != nil {
			return err
		}

		log.Printf("[LKR] old keys backed up in: %v \n", dir)
	}

	if _, err := awsCliConfig.SetIAMCredentials(newIamCredentials); err != nil {
		return err
	}

	if err := awsCliConfig.SaveConfig(); err != nil {
		return err
	}
	log.Println("[LKR] New credentials saved")

	if os.Getenv("LKR_DELETE_OLD_KEYS") == "yes" {
		if err := awsUtils.DeleteCredentials(oldIamCredentials.accessKeyId, username); err != nil {
			return err
		}
		log.Println("[LKR] Old keys deleted")
	} else {
		if err := awsUtils.DeactivateOldKeys(oldIamCredentials.accessKeyId, username); err != nil {
			return err
		}
		log.Println("[LKR] Old keys deactivated")
	}

	return rootCmd.Execute()
}
