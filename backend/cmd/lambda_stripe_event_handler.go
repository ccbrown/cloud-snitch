package cmd

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/stripe/stripe-go/v81"

	"github.com/ccbrown/cloud-snitch/backend/app"
)

type stripeEventHandler struct {
	App *app.App
}

func (h *stripeEventHandler) HandleEvents(ctx context.Context, event events.EventBridgeEvent) error {
	var stripeEvent stripe.Event
	if err := jsoniter.Unmarshal(event.Detail, &stripeEvent); err != nil {
		return err
	}
	return h.App.HandleStripeEvent(ctx, &stripeEvent)
}

var lambdaStripeEventHandlerCmd = &cobra.Command{
	Use:   "lambda-stripe-event-handler",
	Short: "processes events from the stripe event bus",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := app.New(rootConfig.App)
		if err != nil {
			return err
		}

		h := &stripeEventHandler{App: a}
		lambda.Start(h.HandleEvents)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lambdaStripeEventHandlerCmd)
}
