package storetest

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/require"

	"github.com/ccbrown/cloud-snitch/backend/store"
)

// Initializes a database and returns a store configuration that can be used to connect to it.
func NewStoreConfig(t *testing.T) store.Config {
	keyBytes := make([]byte, 20)
	if _, err := rand.Read(keyBytes); err != nil {
		t.Fatal(err)
	}
	key := base64.RawURLEncoding.EncodeToString(keyBytes)

	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	tableNameBytes := make([]byte, 20)
	if _, err := rand.Read(tableNameBytes); err != nil {
		t.Fatal(err)
	}
	tableName := base64.RawURLEncoding.EncodeToString(tableNameBytes)

	config := store.Config{
		DynamoDB: store.DynamoDBConfig{
			TableName: tableName,
			Endpoint:  endpoint,
			Region:    "us-east-1",
			StaticCredentials: &store.DynamoDBStaticCredentials{
				AccessKeyId:     key,
				SecretAccessKey: key,
			},
		},
	}

	s, err := store.New(config)
	require.NoError(t, err)

	if _, err := s.Client().ListTables(context.Background(), &dynamodb.ListTablesInput{}); err != nil {
		t.Skip("no dynamodb server available. to start one: docker run -p 8000:8000 --rm -it amazon/dynamodb-local")
	}

	if err := s.Recreate(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		s.Client().DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
	})

	require.NoError(t, dynamodb.NewTableExistsWaiter(s.Client()).Wait(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}, 2*time.Minute))

	return config
}
