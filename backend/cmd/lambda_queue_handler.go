package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ccbrown/cloud-snitch/backend/app"
)

type queueHandler struct {
	App *app.App
}

func (h *queueHandler) HandleEvents(ctx context.Context, event events.SQSEvent) error {
	for _, message := range event.Records {
		var attrs app.QueueMessageAttributes

		sendTimestamp, err := strconv.ParseInt(message.Attributes["SentTimestamp"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing sent timestamp: %w", err)
		}
		attrs.SendTime = time.UnixMilli(sendTimestamp)

		var msg app.QueueMessage
		if err := jsoniter.Unmarshal([]byte(message.Body), &msg); err != nil {
			zap.L().Error("failed to unmarshal queue message", zap.String("body", message.Body), zap.Error(err))
			return fmt.Errorf("error unmarshaling edge message: %w", err)
		} else if err := h.App.HandleQueueMessage(ctx, msg, attrs); err != nil {
			zap.L().Error("failed to handle queue message", zap.String("body", message.Body), zap.Error(err))
			return fmt.Errorf("error handling edge message: %w", err)
		}
	}
	return nil
}

var lambdaQueueHandlerCmd = &cobra.Command{
	Use:   "lambda-queue-handler",
	Short: "processes messages from the message queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := app.New(rootConfig.App)
		if err != nil {
			return err
		}

		h := &queueHandler{App: a}
		lambda.Start(h.HandleEvents)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lambdaQueueHandlerCmd)
}
