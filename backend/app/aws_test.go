package app

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMostSimilarKnownAWSRegion(t *testing.T) {
	assert.True(t, strings.HasPrefix(MostSimilarKnownAWSRegion("us-east-99"), "us-east-"))
	assert.True(t, strings.HasPrefix(MostSimilarKnownAWSRegion("us-foo-99"), "us-"))
}

func TestClosestAvailableAWSRegion(t *testing.T) {
	available := []string{
		"us-east-1",
		"us-west-2",
		"eu-central-1",
	}
	for region, expected := range map[string]string{
		"af-south-1":     "eu-central-1",
		"ap-east-1":      "eu-central-1",
		"ap-northeast-1": "us-west-2",
		"ap-northeast-2": "us-west-2",
		"ap-northeast-3": "us-west-2",
		"ap-south-1":     "eu-central-1",
		"ap-south-2":     "eu-central-1",
		"ap-southeast-1": "eu-central-1",
		"ap-southeast-2": "us-west-2",
		"ap-southeast-3": "eu-central-1",
		"ap-southeast-4": "us-west-2",
		"ap-southeast-5": "eu-central-1",
		"ap-southeast-7": "eu-central-1",
		"ca-central-1":   "us-east-1",
		"ca-west-1":      "us-west-2",
		"cn-north-1":     "eu-central-1",
		"cn-northwest-1": "eu-central-1",
		"eu-central-1":   "eu-central-1",
		"eu-central-2":   "eu-central-1",
		"eu-north-1":     "eu-central-1",
		"eu-south-1":     "eu-central-1",
		"eu-south-2":     "eu-central-1",
		"eu-west-1":      "eu-central-1",
		"eu-west-2":      "eu-central-1",
		"eu-west-3":      "eu-central-1",
		"il-central-1":   "eu-central-1",
		"me-central-1":   "eu-central-1",
		"me-south-1":     "eu-central-1",
		"mx-central-1":   "us-east-1",
		"sa-east-1":      "us-east-1",
		"us-east-1":      "us-east-1",
		"us-east-2":      "us-east-1",
		"us-gov-east-1":  "us-east-1",
		"us-gov-west-1":  "us-west-2",
		"us-west-1":      "us-west-2",
		"us-west-2":      "us-west-2",
	} {
		t.Run(region, func(t *testing.T) {
			assert.Equal(t, expected, ClosestAvailableAWSRegion(region, available))
		})
	}
}
