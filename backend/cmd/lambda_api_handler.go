package cmd

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/spf13/cobra"

	"github.com/ccbrown/cloud-snitch/backend/api"
	"github.com/ccbrown/cloud-snitch/backend/app"
)

var lambdaAPIHandlerCmd = &cobra.Command{
	Use:   "lambda-api-handler",
	Short: "serves the api as an aws lambda handler",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := app.New(rootConfig.App)
		if err != nil {
			return err
		}

		api := api.New(a, rootConfig.API)
		adapter := httpadapter.NewV2(api)

		lambda.Start(adapter.ProxyWithContext)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lambdaAPIHandlerCmd)
}
