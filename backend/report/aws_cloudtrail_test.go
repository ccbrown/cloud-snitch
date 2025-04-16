package report

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReport_ImportCompressedAWSCloudTrailLog(t *testing.T) {
	r := &Report{
		StartTime:       time.Date(2025, 3, 6, 2, 25, 0, 0, time.UTC),
		DurationSeconds: 60 * 60,
	}

	f, err := os.Open("testdata/aws-cloudtrail-logs/AWSLogs/o-1234abcde/222222222222/CloudTrail/us-east-1/2025/03/06/222222222222_CloudTrail_us-east-1_20250306T0230Z_QF6qcgyiaVnSueXa.json.gz")
	require.NoError(t, err)
	defer f.Close()
	require.NoError(t, r.ImportCompressedAWSCloudTrailLog(f))

	actualJSON, err := json.MarshalIndent(r, "", "    ")
	require.NoError(t, err)

	expectedJSON := `{
		"startTime": "2025-03-06T02:25:00Z",
		"durationSeconds": 3600,
		"networkLocations": {
			"123.12.0.0/17": {
				"latitude": 34.7472,
				"longitude": 113.625,
				"countryCode": "CN",
				"countryName": "China",
				"cityName": "Zhengzhou",
				"subdivisionNames": [
					"Henan"
				]
			},
			"44.223.86.0/23": {
				"latitude": 38.9547,
				"longitude": -77.4043,
				"countryCode": "US",
				"countryName": "United States",
				"cityName": "Herndon",
				"subdivisionNames": [
					"Virginia"
				]
			},
			"98.80.0.0/17": {
				"latitude": 39.0438,
				"longitude": -77.4874,
				"countryCode": "US",
				"countryName": "United States",
				"cityName": "Ashburn",
				"subdivisionNames": [
					"Virginia"
				]
			}
		},
		"ipAddressNetworks": {
			"123.12.3.4": "123.12.0.0/17",
			"44.223.86.2": "44.223.86.0/23",
			"98.80.15.110": "98.80.0.0/17"
		},
		"principals": {
			"AIDAJCEX7SE6A3IUMPJEO": {
				"name": "arn:aws:iam::222222222222:user/chris",
				"type": "AWSIAMUser",
				"arn": "arn:aws:iam::222222222222:user/chris",
				"userAgents": {
					"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3 Safari/605.1.15": 1
				},
				"ipAddresses": {
					"123.12.3.4": 1
				},
				"events": {
					"health.amazonaws.com:DescribeEventAggregates": {
						"name": "DescribeEventAggregates",
						"source": "health.amazonaws.com",
						"count": 1
					}
				}
			},
			"AROAIXNNZA45TGR7PBZQ2": {
				"name": "arn:aws:iam::222222222222:role/ecs-cluster-instance",
				"type": "AWSAssumedRole",
				"arn": "arn:aws:iam::222222222222:role/ecs-cluster-instance",
				"userAgents": {
					"aws-sdk-go/1.55.5 (go1.22.7; linux; amd64) amazon-ssm-agent/": 1
				},
				"ipAddresses": {
					"98.80.15.110": 1
				},
				"events": {
					"ssm.amazonaws.com:UpdateInstanceInformation": {
						"name": "UpdateInstanceInformation",
						"source": "ssm.amazonaws.com",
						"count": 1
					}
				}
			},
			"AROAJPOWK32OXQMNTD5A2": {
				"name": "arn:aws:iam::222222222222:role/my-LambdaFunctionRole-54321GFDSX",
				"type": "AWSAssumedRole",
				"arn": "arn:aws:iam::222222222222:role/my-LambdaFunctionRole-54321GFDSX",
				"userAgents": {
					"aws-sdk-java/2.30.21 md/io#async md/http#NettyNio ua/2.1 os/Linux#5.10.234-225.895.amzn2.x86_64 lang/java#17.0.14 md/OpenJDK_64-Bit_Server_VM#17.0.14+7-LTS md/vendor#Amazon.com_Inc. md/en_US m/E": 1
				},
				"ipAddresses": {
					"44.223.86.2": 1
				},
				"events": {
					"kms.amazonaws.com:Decrypt": {
						"name": "Decrypt",
						"source": "kms.amazonaws.com",
						"count": 1
					}
				}
			},
			"AROAJSSCJRRGHVOV2IMRO": {
				"name": "arn:aws:iam::222222222222:role/aws-ec2-spot-fleet-tagging-role",
				"type": "AWSAssumedRole",
				"arn": "arn:aws:iam::222222222222:role/aws-ec2-spot-fleet-tagging-role",
				"userAgents": {
					"spotfleet.amazonaws.com": 1
				},
				"events": {
					"ec2.amazonaws.com:DescribeInstanceStatus": {
						"name": "DescribeInstanceStatus",
						"source": "ec2.amazonaws.com",
						"count": 1
					}
				}
			},
			"cloudtrail.amazonaws.com": {
				"name": "cloudtrail.amazonaws.com",
				"type": "AWSService",
				"userAgents": {
					"cloudtrail.amazonaws.com": 14
				},
				"events": {
					"s3.amazonaws.com:GetBucketAcl": {
						"name": "GetBucketAcl",
						"source": "s3.amazonaws.com",
						"count": 14
					}
				}
			},
			"logs.amazonaws.com": {
				"name": "logs.amazonaws.com",
				"type": "AWSService",
				"userAgents": {
					"logs.amazonaws.com": 1
				},
				"events": {
					"sts.amazonaws.com:AssumeRole": {
						"name": "AssumeRole",
						"source": "sts.amazonaws.com",
						"count": 1
					}
				}
			},
			"spotfleet.amazonaws.com": {
				"name": "spotfleet.amazonaws.com",
				"type": "AWSService",
				"userAgents": {
					"spotfleet.amazonaws.com": 1
				},
				"events": {
					"sts.amazonaws.com:AssumeRole": {
						"name": "AssumeRole",
						"source": "sts.amazonaws.com",
						"count": 1
					}
				}
			}
		}
	}`

	assert.JSONEq(t, expectedJSON, string(actualJSON), string(actualJSON))
}

type MockAmazonS3API struct {
	T *testing.T
}

func (api MockAmazonS3API) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	file, err := os.Open("testdata/" + *params.Bucket + "/" + *params.Key)
	require.NoError(api.T, err)
	return &s3.GetObjectOutput{
		Body: file,
	}, nil
}

func (api MockAmazonS3API) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	testdataDir := "testdata/" + *params.Bucket + "/"

	if params.Delimiter != nil {
		entries, err := os.ReadDir(testdataDir + *params.Prefix)
		if !os.IsNotExist(err) {
			require.NoError(api.T, err)
		}
		commonPrefixes := make([]s3types.CommonPrefix, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				commonPrefixes = append(commonPrefixes, s3types.CommonPrefix{
					Prefix: aws.String(*params.Prefix + entry.Name() + "/"),
				})
			}
		}
		return &s3.ListObjectsV2Output{
			CommonPrefixes: commonPrefixes,
		}, nil
	}

	prefixDir := *params.Prefix
	if !strings.HasSuffix(prefixDir, "/") {
		prefixDir = filepath.Dir(prefixDir)
	}
	var objects []s3types.Object
	err := filepath.Walk(testdataDir+prefixDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			key := strings.TrimPrefix(path, testdataDir)
			if !info.IsDir() && strings.HasPrefix(key, *params.Prefix) {
				objects = append(objects, s3types.Object{
					Key: aws.String(key),
				})
			}
			return nil
		})
	if !os.IsNotExist(err) {
		require.NoError(api.T, err)
	}
	return &s3.ListObjectsV2Output{
		Contents: objects,
	}, nil
}

func TestReport_ImportCompressedAWSCloudTrailLogBucket(t *testing.T) {
	r := &Report{
		StartTime:       time.Date(2025, 3, 6, 2, 25, 0, 0, time.UTC),
		DurationSeconds: 60 * 60,
	}

	require.NoError(t, r.ImportAWSCloudTrailLogBucket(context.Background(), ImportAWSCloudTrailLogBucketConfig{
		S3: MockAmazonS3API{
			T: t,
		},
		BucketName: "aws-cloudtrail-logs",
	}))

	assert.Len(t, r.Principals, 8)
}

func TestReport_AWSAccountEvent(t *testing.T) {
	rawEvent := `{
		"eventVersion": "1.08",
		"userIdentity": {
			"type": "AWSAccount",
			"principalId": "AROA3ISBVQKHHSGTEMQ45:cloud-snitch-us-east-1-dev-QueueHandler49A28051-kJGBH9cs5Vvs",
			"accountId": "774305579662"
		},
		"eventTime": "2025-03-19T00:30:20Z",
		"eventSource": "sts.amazonaws.com",
		"eventName": "AssumeRole",
		"awsRegion": "us-east-1",
		"sourceIPAddress": "18.232.97.18",
		"userAgent": "aws-sdk-go-v2/1.36.3 ua/2.1 os/linux lang/go#1.24.1 md/GOOS#linux md/GOARCH#arm64 exec-env/AWS_Lambda_Image api/sts#1.33.17 m/g",
		"requestParameters": {
			"roleArn": "arn:aws:iam::123412341234:role/cloud-snitch-dev-CloudSnitchIntegrationRole03668212-LXdv1ZrGg24c",
			"roleSessionName": "cloud_snitch",
			"externalId": "t-fZFf5o79ZxjFhK8XnfLRCD"
		},
		"responseElements": {
			"credentials": {
				"accessKeyId": "ASIATGOPDG7TNXVGU6TJ",
				"sessionToken": "",
				"expiration": "Mar 19, 2025, 1:30:20 AM"
			},
			"assumedRoleUser": {
				"assumedRoleId": "AROATGOPDG7THYFSRMFNE:cloud_snitch",
				"arn": "arn:aws:sts::123412341234:assumed-role/cloud-snitch-dev-CloudSnitchIntegrationRole03668212-LXdv1ZrGg24c/cloud_snitch"
			}
		},
		"additionalEventData": {
			"RequestDetails": {
				"awsServingRegion": "us-east-1",
				"endpointType": "regional"
			}
		},
		"requestID": "9d068207-15fb-4fb7-bd59-754e7e798828",
		"eventID": "1c94cbc7-9c84-3f25-9481-d5172c6ee15e",
		"readOnly": true,
		"resources": [
			{
				"accountId": "123412341234",
				"type": "AWS::IAM::Role",
				"ARN": "arn:aws:iam::123412341234:role/cloud-snitch-dev-CloudSnitchIntegrationRole03668212-LXdv1ZrGg24c"
			}
		],
		"eventType": "AwsApiCall",
		"managementEvent": true,
		"recipientAccountId": "123412341234",
		"sharedEventID": "6e410184-19dd-4d7d-a8f0-c17c584957a0",
		"eventCategory": "Management",
		"tlsDetails": {
			"tlsVersion": "TLSv1.3",
			"cipherSuite": "TLS_AES_128_GCM_SHA256",
			"clientProvidedHostHeader": "sts.us-east-1.amazonaws.com"
		}
	}`

	var record AWSCloudTrailRecord
	require.NoError(t, json.Unmarshal([]byte(rawEvent), &record))

	r := &Report{}
	r.ImportAWSCloudTrailRecord(&record)

	require.Len(t, r.Principals, 1)
}

func TestReport_AWSInternalServiceEvent(t *testing.T) {
	rawEvent := `{
		"eventVersion": "1.09",
		"userIdentity": {
			"accountId": "774305579662",
			"invokedBy": "AWS Internal"
		},
		"eventTime": "2025-03-13T03:22:13Z",
		"eventSource": "ecr.amazonaws.com",
		"eventName": "PolicyExecutionEvent",
		"awsRegion": "us-east-1",
		"sourceIPAddress": "AWS Internal",
		"userAgent": "AWS Internal",
		"requestParameters": null,
		"responseElements": null,
		"eventID": "a185533f-657a-4e17-831d-b7743921c04e",
		"readOnly": true,
		"resources": [
			{
				"accountId": "774305579662",
				"type": "AWS::ECR::Repository",
				"ARN": "arn:aws:ecr:us-east-1:774305579662:repository/cdk-hnb659fds-container-assets-774305579662-us-east-1"
			}
		],
		"eventType": "AwsServiceEvent",
		"managementEvent": true,
		"recipientAccountId": "774305579662",
		"serviceEventDetails": {
			"repositoryName": "cdk-hnb659fds-container-assets-774305579662-us-east-1",
			"lifecycleEventPolicy": {
				"lifecycleEventRules": [
					{
						"rulePriority": 1,
						"description": "Untagged images should not exist, but expire any older than one year",
						"lifecycleEventSelection": {
							"tagStatus": "Untagged",
							"tagPrefixList": [],
							"tagPatternList": [],
							"countType": "Time since image pushed",
							"countUnit": "Days",
							"countNumber": 365
						},
						"action": "expire"
					}
				],
				"lastEvaluatedAt": 1741782216108,
				"policyVersion": 1,
				"policyId": "d5175d14-57c1-4b11-9f06-610483b24692"
			},
			"lifecycleEventImageActions": [],
			"lifecycleEventFailureDetails": []
		},
		"eventCategory": "Management"
	}`

	var record AWSCloudTrailRecord
	require.NoError(t, json.Unmarshal([]byte(rawEvent), &record))

	r := &Report{}
	r.ImportAWSCloudTrailRecord(&record)

	require.Len(t, r.Principals, 1)
}

func TestReport_WebIdentityUserEvent(t *testing.T) {
	rawEvent := `{
		"eventVersion": "1.08",
		"userIdentity": {
			"type": "WebIdentityUser",
			"principalId": "arn:aws:iam::774305579662:oidc-provider/token.actions.githubusercontent.com:sts.amazonaws.com:repo:ccbrown/cloud-snitch:ref:refs/heads/main",
			"userName": "repo:ccbrown/cloud-snitch:ref:refs/heads/main",
			"identityProvider": "arn:aws:iam::774305579662:oidc-provider/token.actions.githubusercontent.com"
		},
		"eventTime": "2025-04-16T02:46:56Z",
		"eventSource": "sts.amazonaws.com",
		"eventName": "AssumeRoleWithWebIdentity",
		"awsRegion": "us-east-1",
		"sourceIPAddress": "52.154.133.36",
		"userAgent": "aws-sdk-nodejs/2.1112.0 linux/v20.19.0 configure-aws-credentials-for-github-actions promise",
		"requestParameters": {
			"roleArn": "arn:aws:iam::774305579662:role/cloud-snitch-github-actio-GithubActionsRoleF5CC769F-MoeCHSu77MYt",
			"roleSessionName": "GitHubActions",
			"durationSeconds": 3600
		},
		"responseElements": {
			"credentials": {
				"accessKeyId": "ASIA3ISBVQKHJOCMRQJH",
				"sessionToken": "",
				"expiration": "Apr 16, 2025, 3:46:56 AM"
			},
			"subjectFromWebIdentityToken": "repo:ccbrown/cloud-snitch:ref:refs/heads/main",
			"assumedRoleUser": {
				"assumedRoleId": "AROA3ISBVQKHLCPMWHJG5:GitHubActions",
				"arn": "arn:aws:sts::774305579662:assumed-role/cloud-snitch-github-actio-GithubActionsRoleF5CC769F-MoeCHSu77MYt/GitHubActions"
			},
			"provider": "arn:aws:iam::774305579662:oidc-provider/token.actions.githubusercontent.com",
			"audience": "sts.amazonaws.com"
		},
		"additionalEventData": {
			"identityProviderConnectionVerificationMethod": "IAMTrustStore",
			"RequestDetails": {
				"awsServingRegion": "us-east-1",
				"endpointType": "regional"
			}
		},
		"requestID": "6936cb5f-948b-4139-8939-06118942c0cc",
		"eventID": "ef0464cd-918e-4528-9ea3-1a081d39b393",
		"readOnly": true,
		"resources": [
			{
				"accountId": "774305579662",
				"type": "AWS::IAM::Role",
				"ARN": "arn:aws:iam::774305579662:role/cloud-snitch-github-actio-GithubActionsRoleF5CC769F-MoeCHSu77MYt"
			}
		],
		"eventType": "AwsApiCall",
		"managementEvent": true,
		"recipientAccountId": "774305579662",
		"eventCategory": "Management",
		"tlsDetails": {
			"tlsVersion": "TLSv1.3",
			"cipherSuite": "TLS_AES_128_GCM_SHA256",
			"clientProvidedHostHeader": "sts.us-east-1.amazonaws.com"
		}
	}`

	var record AWSCloudTrailRecord
	require.NoError(t, json.Unmarshal([]byte(rawEvent), &record))

	r := &Report{}
	r.ImportAWSCloudTrailRecord(&record)

	require.Len(t, r.Principals, 1)
	for id, p := range r.Principals {
		assert.Equal(t, PrincipalTypeWebIdentityUser, p.Type)
		assert.Equal(t, "arn:aws:iam::774305579662:oidc-provider/token.actions.githubusercontent.com", p.Name)
		assert.Equal(t, "arn:aws:iam::774305579662:oidc-provider/token.actions.githubusercontent.com", p.ARN)
		assert.Equal(t, "arn:aws:iam::774305579662:oidc-provider/token.actions.githubusercontent.com", id)
		break
	}
}
