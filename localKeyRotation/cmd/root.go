package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	"os"
)

const awsAccessKeyId = "aws_access_key_id"
const awsSecretAccessKey = "aws_secret_access_key"

type keys struct {
	accessKeyId string
	secretAccessKey string
}
var temporaryCacheCredentials = make(map[string]keys, 0)

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

// Get iam user name


// Temporarily cache credentials just in case
func cacheOldCredentials(credentials *ini.File) error {
	sections := credentials.Sections()
	for _, section := range sections {
		accessKeyId, err := section.GetKey(awsAccessKeyId)
		if err != nil {
			return err
		}

		secretAccessKey, err := section.GetKey(awsSecretAccessKey)
		if err != nil {
			return err
		}

		k := keys{
			accessKeyId:     accessKeyId.String(),
			secretAccessKey: secretAccessKey.String(),
		}

		temporaryCacheCredentials[section.Name()] = k
	}

	return nil
}