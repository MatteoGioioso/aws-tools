package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const awsAccessKeyId = "aws_access_key_id"
const awsSecretAccessKey = "aws_secret_access_key"

var rootCmd = &cobra.Command{
	Use:   "aws-key-rotation",
	Short: "Automatically rotate your local AWS credentials",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() error {
	// TODO remove this
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "test_credentials")

	awsCliConfig, err := NewAwsCliConfig()
	if err != nil {
		return err
	}

	awsUtils, err := NewAWSUtils()
	if err != nil {
		return err
	}

	username, err := awsUtils.GetCurrentUsername()
	if err != nil {
		return err
	}
	fmt.Println(username)

	configFile, err := awsCliConfig.ParseConfig()
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", configFile.Section("default").Key(awsAccessKeyId))

	fmt.Println(awsCliConfig.GetCurrentProfile())

	return rootCmd.Execute()
}


