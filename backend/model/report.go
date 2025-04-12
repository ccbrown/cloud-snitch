package model

import "time"

type ReportRetention string

const (
	ReportRetentionOneWeek  ReportRetention = "1w"
	ReportRetentionTwoWeeks ReportRetention = "2w"
)

func (r ReportRetention) Duration() time.Duration {
	switch r {
	case ReportRetentionOneWeek:
		return 7 * 24 * time.Hour
	default:
		panic("invalid report retention")
	}
}

func NewReportId() Id {
	return NewId("r")
}

type Report struct {
	Id             Id
	CreationTime   time.Time
	ExpirationTime time.Time

	TeamId           Id
	AWSIntegrationId Id

	Scope       ReportScope
	Location    ReportLocation
	DownloadURL string

	Size               int
	SourceBytes        int
	IsIncomplete       bool
	GenerationDuration time.Duration
}

type ReportScope struct {
	StartTime time.Time
	Duration  time.Duration
	AWS       ReportScopeAWS
}

type ReportScopeAWS struct {
	AccountId string
	Region    string
}

type ReportLocation struct {
	AWSRegion string
	S3Bucket  string
	Key       string
}
