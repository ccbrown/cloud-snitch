package report

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	jsoniter "github.com/json-iterator/go"

	"github.com/ccbrown/go-geoip"
)

type AWSCloudTrailLog struct {
	Records []AWSCloudTrailRecord
}

type AWSCloudTrailRecord struct {
	UserIdentity    *AWSCloudTrailUserIdentity
	EventTime       time.Time
	EventSource     string
	EventName       string
	EventCategory   string
	EventType       string
	SourceIPAddress string
	UserAgent       string
	ErrorCode       string
}

func (r *AWSCloudTrailRecord) EventSummaryKey() string {
	return r.EventSource + ":" + r.EventName
}

type AWSCloudTrailUserIdentityType string

const (
	AWSCloudTrailUserIdentityTypeAssumedRole     AWSCloudTrailUserIdentityType = "AssumedRole"
	AWSCloudTrailUserIdentityTypeRole            AWSCloudTrailUserIdentityType = "Role"
	AWSCloudTrailUserIdentityTypeIAMUser         AWSCloudTrailUserIdentityType = "IAMUser"
	AWSCloudTrailUserIdentityTypeAWSService      AWSCloudTrailUserIdentityType = "AWSService"
	AWSCloudTrailUserIdentityTypeAWSAccount      AWSCloudTrailUserIdentityType = "AWSAccount"
	AWSCloudTrailUserIdentityTypeWebIdentityUser AWSCloudTrailUserIdentityType = "WebIdentityUser"
)

func (t AWSCloudTrailUserIdentityType) PrincipalType() PrincipalType {
	switch t {
	case AWSCloudTrailUserIdentityTypeAssumedRole:
		return PrincipalTypeAWSAssumedRole
	case AWSCloudTrailUserIdentityTypeRole:
		return PrincipalTypeAWSRole
	case AWSCloudTrailUserIdentityTypeIAMUser:
		return PrincipalTypeAWSIAMUser
	case AWSCloudTrailUserIdentityTypeAWSService:
		return PrincipalTypeAWSService
	case AWSCloudTrailUserIdentityTypeAWSAccount:
		return PrincipalTypeAWSAccount
	case AWSCloudTrailUserIdentityTypeWebIdentityUser:
		return PrincipalTypeWebIdentityUser
	default:
		return PrincipalTypeUnknown
	}
}

type AWSCloudTrailUserIdentity struct {
	Type             AWSCloudTrailUserIdentityType
	PrincipalId      string
	SessionContext   *AWSCloudTrailUserIdentitySessionContext
	ARN              string
	InvokedBy        string
	IdentityProvider string
}

func (i *AWSCloudTrailUserIdentity) PrincipalName() string {
	if i.SessionContext != nil && i.SessionContext.SessionIssuer != nil {
		if name := i.SessionContext.SessionIssuer.PrincipalName(); name != "" {
			return name
		}
	}
	if i.IdentityProvider != "" {
		return i.IdentityProvider
	}
	if i.ARN != "" {
		return i.ARN
	}
	return i.InvokedBy
}

func (i *AWSCloudTrailUserIdentity) PrincipalARN() string {
	if i.SessionContext != nil && i.SessionContext.SessionIssuer != nil {
		if arn := i.SessionContext.SessionIssuer.PrincipalARN(); arn != "" {
			return arn
		}
	}
	if i.IdentityProvider != "" {
		return i.IdentityProvider
	}
	if i.ARN != "" {
		return i.ARN
	}
	return ""
}

func (i *AWSCloudTrailUserIdentity) PrincipalKey() string {
	if i.SessionContext != nil && i.SessionContext.SessionIssuer != nil {
		if key := i.SessionContext.SessionIssuer.PrincipalKey(); key != "" {
			return key
		}
	}
	if i.IdentityProvider != "" {
		return i.IdentityProvider
	}
	if i.PrincipalId != "" {
		return i.PrincipalId
	}
	return i.InvokedBy
}

type AWSCloudTrailUserIdentitySessionContext struct {
	SessionIssuer *AWSCloudTrailUserIdentity
}

func (r *Report) ImportAWSCloudTrailRecords(records []AWSCloudTrailRecord) {
	for _, record := range records {
		r.ImportAWSCloudTrailRecord(&record)
	}
}

func (r *Report) AddIPAddressLocation(ip net.IP) {
	if _, ok := r.IPAddressNetworks[ip.String()]; ok {
		return
	}
	if r.IPAddressNetworks == nil {
		r.IPAddressNetworks = make(map[string]*string)
	}
	record, network := geoip.Lookup(ip)
	if record == nil {
		r.IPAddressNetworks[ip.String()] = nil
	} else {
		network := network.String()
		r.IPAddressNetworks[ip.String()] = &network
		if r.NetworkLocations == nil {
			r.NetworkLocations = make(map[string]*Location)
		}
		subdivisionNames := make([]string, 0, len(record.Subdivisions))
		for _, subdivision := range record.Subdivisions {
			subdivisionNames = append(subdivisionNames, subdivision.Names.En)
		}
		r.NetworkLocations[network] = &Location{
			Latitude:         record.Location.Latitude,
			Longitude:        record.Location.Longitude,
			CountryCode:      record.Country.ISOCode,
			CountryName:      record.Country.Names.En,
			CityName:         record.City.Names.En,
			SubdivisionNames: subdivisionNames,
		}
	}
}

func (r *Report) ImportAWSCloudTrailRecord(record *AWSCloudTrailRecord) {
	if record.EventCategory != "Management" {
		return
	}

	if !r.StartTime.IsZero() {
		if record.EventTime.Before(r.StartTime) || !record.EventTime.Before(r.StartTime.Add(r.Duration())) {
			return
		}
	}

	principalKey := record.UserIdentity.PrincipalKey()
	principal, ok := r.Principals[principalKey]
	if !ok {
		principal = &Principal{
			Name:   record.UserIdentity.PrincipalName(),
			Type:   record.UserIdentity.Type.PrincipalType(),
			ARN:    record.UserIdentity.PrincipalARN(),
			Events: make(map[string]*EventSummary),
		}
		if r.Principals == nil {
			r.Principals = make(map[string]*Principal)
		}
		r.Principals[principalKey] = principal
	}

	if ip := net.ParseIP(record.SourceIPAddress); ip != nil {
		if principal.IPAddresses == nil {
			principal.IPAddresses = make(map[string]int)
		}
		principal.IPAddresses[ip.String()]++
		r.AddIPAddressLocation(ip)
	}

	if agent := strings.TrimSpace(record.UserAgent); agent != "" {
		if principal.UserAgents == nil {
			principal.UserAgents = make(map[string]int)
		}
		principal.UserAgents[agent]++
	}

	eventSummaryKey := record.EventSummaryKey()
	eventSummary, ok := principal.Events[eventSummaryKey]
	if !ok {
		eventSummary = &EventSummary{
			Name:   record.EventName,
			Source: record.EventSource,
		}
		principal.Events[eventSummaryKey] = eventSummary
	}

	eventSummary.Count++

	if record.ErrorCode != "" {
		if eventSummary.ErrorCodes == nil {
			eventSummary.ErrorCodes = make(map[string]int)
		}
		eventSummary.ErrorCodes[record.ErrorCode]++
	}
}

func (r *Report) ImportCompressedAWSCloudTrailLog(f io.Reader) error {
	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	return r.ImportAWSCloudTrailLogJSON(gz)
}

func (r *Report) ImportAWSCloudTrailLogJSON(f io.Reader) error {
	var log AWSCloudTrailLog
	if err := jsoniter.NewDecoder(f).Decode(&log); err != nil {
		return fmt.Errorf("failed to decode log: %w", err)
	}
	r.ImportAWSCloudTrailRecords(log.Records)

	return nil
}

type AmazonS3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

type ImportAWSCloudTrailLogBucketConfig struct {
	BucketName string
	S3         AmazonS3API
}

func getS3Subdirectories(ctx context.Context, client s3.ListObjectsV2APIClient, bucket, prefix string) ([]string, error) {
	var ret []string

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket:    &bucket,
		Prefix:    &prefix,
		Delimiter: aws.String("/"),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}
		for _, commonPrefix := range page.CommonPrefixes {
			parts := strings.Split(*commonPrefix.Prefix, "/")
			ret = append(ret, parts[len(parts)-2])
		}
	}

	return ret, nil
}

type AWSCloudTrailLogBucketAccountRegion struct {
	AccountsPrefix string
	AccountId      string
	Region         string
}

type ScanAWSCloudTrailLogBucketConfig struct {
	S3         AmazonS3API
	BucketName string
	KeyPrefix  string
}

func ScanAWSCloudTrailLogBucket(ctx context.Context, config ScanAWSCloudTrailLogBucketConfig) ([]AWSCloudTrailLogBucketAccountRegion, error) {
	type Account struct {
		AccountId      string
		AccountsPrefix string
	}
	var accounts []Account

	awsLogsPrefix := config.KeyPrefix + "AWSLogs/"

	subdirectories, err := getS3Subdirectories(ctx, config.S3, config.BucketName, awsLogsPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWSLogs/ subdirectories: %w", err)
	}

	for _, name := range subdirectories {
		if strings.HasPrefix(name, "o-") {
			prefix := awsLogsPrefix + name + "/"
			accountIds, err := getS3Subdirectories(ctx, config.S3, config.BucketName, prefix)
			if err != nil {
				return nil, fmt.Errorf("failed to get subdirectories for organization logs: %w", err)
			}
			for _, accountId := range accountIds {
				accounts = append(accounts, Account{
					AccountId:      accountId,
					AccountsPrefix: prefix,
				})
			}
		} else {
			accounts = append(accounts, Account{
				AccountId:      name,
				AccountsPrefix: awsLogsPrefix,
			})
		}
	}

	var ret []AWSCloudTrailLogBucketAccountRegion

	for _, account := range accounts {
		prefix := account.AccountsPrefix + account.AccountId + "/CloudTrail/"

		regions, err := getS3Subdirectories(ctx, config.S3, config.BucketName, prefix)
		if err != nil {
			return nil, fmt.Errorf("failed to get region directories: %w", err)
		}

		for _, region := range regions {
			ret = append(ret, AWSCloudTrailLogBucketAccountRegion{
				AccountsPrefix: account.AccountsPrefix,
				AccountId:      account.AccountId,
				Region:         region,
			})
		}
	}

	return ret, nil
}

func (r *Report) ImportAWSCloudTrailLogBucket(ctx context.Context, config ImportAWSCloudTrailLogBucketConfig) error {
	accountRegions, err := ScanAWSCloudTrailLogBucket(ctx, ScanAWSCloudTrailLogBucketConfig{
		BucketName: config.BucketName,
		S3:         config.S3,
	})
	if err != nil {
		return fmt.Errorf("failed to scan bucket: %w", err)
	}

	for _, accountRegion := range accountRegions {
		if err := r.ImportAWSCloudTrailLogsForAccountRegion(ctx, ImportAWSCloudTrailLogsForAccountRegionConfig{
			S3:             config.S3,
			BucketName:     config.BucketName,
			AccountId:      accountRegion.AccountId,
			AccountsPrefix: accountRegion.AccountsPrefix,
			Region:         accountRegion.Region,
		}); err != nil {
			return fmt.Errorf("failed to import logs for account region: %w", err)
		}
	}

	return nil
}

type ImportAWSCloudTrailLogsForAccountRegionConfig struct {
	S3             AmazonS3API
	BucketName     string
	AccountsPrefix string
	AccountId      string
	Region         string

	// If non-zero, we won't look at more than this many bytes of log files.
	MaxSourceBytes int64
}

func (r *Report) ImportAWSCloudTrailLogsForAccountRegion(ctx context.Context, config ImportAWSCloudTrailLogsForAccountRegionConfig) error {
	prefix := config.AccountsPrefix + config.AccountId + "/CloudTrail/"

	timePadding := 5 * time.Minute
	lastDay := r.StartTime.Add(r.Duration() + timePadding).Truncate(24 * time.Hour)

	for day := r.StartTime.Add(-timePadding).Truncate(24 * time.Hour); !day.After(lastDay); day = day.AddDate(0, 0, 1) {
		regionPrefix := prefix + config.Region + "/"
		dayPrefix := regionPrefix + day.Format("2006/01/02/")
		paginator := s3.NewListObjectsV2Paginator(config.S3, &s3.ListObjectsV2Input{
			Bucket: &config.BucketName,
			Prefix: aws.String(dayPrefix + config.AccountId + "_CloudTrail_" + config.Region + "_" + day.Format("20060102")),
		})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return fmt.Errorf("failed to list objects: %w", err)
			}
			for _, object := range page.Contents {
				if !strings.HasSuffix(*object.Key, ".json.gz") {
					continue
				}
				filename := strings.TrimPrefix(*object.Key, dayPrefix)
				parts := strings.Split(filename, "_")
				if len(parts) < 5 {
					continue
				}
				timestamp, err := time.Parse("20060102T1504Z", parts[3])
				if err != nil {
					continue
				}
				if timestamp.Before(r.StartTime.Add(-timePadding)) || !timestamp.Before(r.StartTime.Add(r.Duration()+timePadding)) {
					continue
				}

				// This is an in-scope object.

				if config.MaxSourceBytes > 0 && object.Size != nil && r.SourceBytes+*object.Size > config.MaxSourceBytes {
					r.IsIncomplete = true
					return nil
				}

				if err := r.ImportAWSCloudTrailLogBucketObject(ctx, ImportAWSCloudTrailLogBucketObjectConfig{
					S3:         config.S3,
					BucketName: config.BucketName,
					ObjectKey:  *object.Key,
				}); err != nil {
					return fmt.Errorf("failed to import log object: %w", err)
				}
			}
		}
	}

	return nil
}

type ImportAWSCloudTrailLogBucketObjectConfig struct {
	S3         AmazonS3API
	BucketName string
	ObjectKey  string
}

func (r *Report) ImportAWSCloudTrailLogBucketObject(ctx context.Context, config ImportAWSCloudTrailLogBucketObjectConfig) error {
	resp, err := config.S3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &config.BucketName,
		Key:    &config.ObjectKey,
	})
	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}
	defer resp.Body.Close()
	if resp.ContentLength != nil {
		r.SourceBytes += *resp.ContentLength
	}
	return r.ImportCompressedAWSCloudTrailLog(resp.Body)
}
