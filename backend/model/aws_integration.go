package model

import "time"

func NewAWSIntegrationId() Id {
	return NewId("aws")
}

type AWSIntegration struct {
	Id           Id
	CreationTime time.Time

	TeamId Id
	Name   string

	RoleARN string

	GetAccountNamesFromOrganizations bool
	ManageSCPs                       bool
	CloudTrailTrail                  *AWSIntegrationCloudTrailTrail
}

type AWSIntegrationCloudTrailTrail struct {
	S3BucketName string
	S3KeyPrefix  string
}

// Contains information gathered by an integration as report generations are queued.
type AWSIntegrationRecon struct {
	AWSIntegrationId Id
	TeamId           Id

	Time           time.Time
	ExpirationTime time.Time
	CanManageSCPs  bool

	Accounts []AWSIntegrationAccountRecon
}

type AWSIntegrationAccountRecon struct {
	Id   string
	Name string
}

type AWSSCP struct {
	Content string
}

type AWSAccessReport struct {
	Services []AWSAccessReportService
}

type AWSAccessReportService struct {
	Name                   string
	Namespace              string
	LastAuthenticationTime time.Time
}
