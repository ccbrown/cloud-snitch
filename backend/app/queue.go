package app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type QueueMessage struct {
	QueueReportGeneration              *QueueReportGenerationInput        `json:",omitempty"`
	QueueTeamReportGeneration          *QueueTeamReportGenerationInput    `json:",omitempty"`
	GenerateAWSCloudTrailReport        *GenerateAWSCloudTrailReportInput  `json:",omitempty"`
	QueueTeamStripeSubscriptionUpdates *struct{}                          `json:",omitempty"`
	UpdateTeamStripeSubscription       *UpdateTeamStripeSubscriptionInput `json:",omitempty"`
	QueueTeamEntitlementRefreshes      *struct{}                          `json:",omitempty"`
	RefreshTeamEntitlements            *RefreshTeamEntitlementsInput      `json:",omitempty"`
}

type OutgoingQueueMessage struct {
	Message QueueMessage
	Delay   time.Duration
}

type QueueMessageAttributes struct {
	SendTime time.Time
}

func (a *App) HandleQueueMessage(ctx context.Context, message QueueMessage, attrs QueueMessageAttributes) error {
	if message.QueueReportGeneration != nil {
		input := *message.QueueReportGeneration
		if input.StartTime.IsZero() {
			input.StartTime = attrs.SendTime.Truncate(input.Duration).Add(-input.Duration)
		}
		if err := a.QueueReportGeneration(ctx, input); err != nil {
			return fmt.Errorf("failed to queue report generation: %w", err)
		}
	}
	if message.QueueTeamReportGeneration != nil {
		if err := a.QueueTeamReportGeneration(ctx, *message.QueueTeamReportGeneration); err != nil {
			return fmt.Errorf("failed to queue team report generation: %w", err)
		}
	}
	if message.GenerateAWSCloudTrailReport != nil {
		if _, err := a.GenerateAWSCloudTrailReport(ctx, *message.GenerateAWSCloudTrailReport); err != nil {
			return fmt.Errorf("failed to generate aws cloudtrail report: %w", err)
		}
	}
	if message.QueueTeamStripeSubscriptionUpdates != nil {
		if err := a.QueueTeamStripeSubscriptionUpdates(ctx); err != nil {
			return fmt.Errorf("failed to queue team stripe subscription updates: %w", err)
		}
	}
	if message.UpdateTeamStripeSubscription != nil {
		if err := a.UpdateTeamStripeSubscription(ctx, *message.UpdateTeamStripeSubscription); err != nil {
			return fmt.Errorf("failed to update team stripe subscription: %w", err)
		}
	}
	if message.QueueTeamEntitlementRefreshes != nil {
		if err := a.QueueTeamEntitlementRefreshes(ctx); err != nil {
			return fmt.Errorf("failed to queue team entitlement refreshes: %w", err)
		}
	}
	if message.RefreshTeamEntitlements != nil {
		if err := a.RefreshTeamEntitlements(ctx, *message.RefreshTeamEntitlements); err != nil {
			return fmt.Errorf("failed to refresh team entitlements: %w", err)
		}
	}
	return nil
}

const MaxQueueDelay = 15 * time.Minute

func (a *App) QueueMessages(ctx context.Context, messagesByQueueRegion map[string][]OutgoingQueueMessage) error {
	const maxBatchSize = 10
	for queueRegion, messages := range messagesByQueueRegion {
		sqsAPI := a.sqs[queueRegion]
		queueURL := "https://sqs." + queueRegion + ".amazonaws.com/" + a.config.AWSAccountId + "/" + a.config.SQSQueueName
		for batchStart := 0; batchStart < len(messages); batchStart += maxBatchSize {
			batch := messages[batchStart:min(batchStart+maxBatchSize, len(messages))]
			batchEntries := make([]sqstypes.SendMessageBatchRequestEntry, len(batch))
			for i, message := range batch {
				buf, err := jsoniter.Marshal(message.Message)
				if err != nil {
					return fmt.Errorf("failed to marshal message: %w", err)
				}
				batchEntries[i] = sqstypes.SendMessageBatchRequestEntry{
					Id:           aws.String(fmt.Sprintf("%d", i)),
					MessageBody:  aws.String(string(buf)),
					DelaySeconds: int32(message.Delay.Seconds()),
				}
			}

			entriesToSend := batchEntries
			for attempt := 1; ; attempt++ {
				output, err := sqsAPI.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
					QueueUrl: &queueURL,
					Entries:  entriesToSend,
				})
				if err != nil {
					return fmt.Errorf("failed to send message batch: %w", err)
				}
				if len(output.Failed) == 0 {
					break
				} else if attempt >= 5 {
					return fmt.Errorf("failed to send message batch after 5 attempts")
				}

				zap.L().Warn("failed to send message batch entries", zap.Int("successful", len(output.Successful)), zap.Int("failed", len(output.Failed)))

				entriesToSend = make([]sqstypes.SendMessageBatchRequestEntry, 0, len(output.Failed))
				for _, failed := range output.Failed {
					idx, err := strconv.Atoi(*failed.Id)
					if err != nil {
						return fmt.Errorf("failed to parse failed message id: %w", err)
					}
					entriesToSend = append(entriesToSend, batchEntries[idx])
				}

				time.Sleep(time.Duration(attempt*attempt) * time.Second)
			}
		}
	}

	return nil
}
