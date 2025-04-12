package store

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/ccbrown/cloud-snitch/backend/model"
)

type Store struct {
	client *dynamodb.Client
	config *Config
}

func New(cfg Config) (*Store, error) {
	if cfg.DynamoDB.TableName == "" {
		return nil, fmt.Errorf("a dynamodb table name must be configured")
	}

	awsConfig, err := cfg.DynamoDB.AWSConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading aws config: %w", err)
	}

	return &Store{
		client: dynamodb.NewFromConfig(awsConfig),
		config: &cfg,
	}, nil
}

func (s *Store) Client() *dynamodb.Client {
	return s.client
}

// Deletes the underlying table if it exists and creates a new one.
func (s *Store) Recreate(ctx context.Context) error {
	return RecreateTable(ctx, s.client, s.config.DynamoDB.TableName)
}

// Deletes the table if it exists and creates a new one.
func RecreateTable(ctx context.Context, client *dynamodb.Client, name string) error {
	if _, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(name),
	}); err == nil {
		dynamodb.NewTableNotExistsWaiter(client).Wait(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(name),
		}, 2*time.Minute)
	}

	if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("_hk"),
				AttributeType: types.ScalarAttributeTypeB,
			}, {
				AttributeName: aws.String("_rk"),
				AttributeType: types.ScalarAttributeTypeB,
			}, {
				AttributeName: aws.String("_bb1h"),
				AttributeType: types.ScalarAttributeTypeB,
			}, {
				AttributeName: aws.String("_bb1r"),
				AttributeType: types.ScalarAttributeTypeB,
			}, {
				AttributeName: aws.String("_bb2h"),
				AttributeType: types.ScalarAttributeTypeB,
			}, {
				AttributeName: aws.String("_bb2r"),
				AttributeType: types.ScalarAttributeTypeB,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("_hk"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("_rk"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("_bb1"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("_bb1h"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("_bb1r"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
			{
				IndexName: aws.String("_bb2"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("_bb2h"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("_bb2r"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},
		TableName:   &name,
		BillingMode: types.BillingModePayPerRequest,
	}); err != nil {
		return err
	}

	_, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(name),
		TimeToLiveSpecification: &types.TimeToLiveSpecification{
			AttributeName: aws.String("_ttl"),
			Enabled:       aws.Bool(true),
		},
	})
	return err
}

func (s *Store) put(ctx context.Context, items ...any) error {
	if len(items) == 1 {
		item := items[0]
		attrs, err := attributevalue.MarshalMap(item)
		if err != nil {
			return err
		}
		_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
			Item:      attrs,
			TableName: &s.config.DynamoDB.TableName,
		})
		return err
	} else {
		var ops []types.TransactWriteItem
		for _, item := range items {
			attrs, err := attributevalue.MarshalMap(item)
			if err != nil {
				return err
			}
			ops = append(ops, types.TransactWriteItem{
				Put: &types.Put{
					Item:      attrs,
					TableName: &s.config.DynamoDB.TableName,
				},
			})
		}
		_, err := s.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: ops,
		})
		return err
	}
}

func createOrUpdateByPrimaryKey[T any](ctx context.Context, s *Store, hk []byte, update expression.UpdateBuilder) (*T, error) {
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, err
	}

	if output, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"_hk": &types.AttributeValueMemberB{
				Value: hk,
			},
			"_rk": &types.AttributeValueMemberB{
				Value: []byte("_"),
			},
		},
		TableName:                 &s.config.DynamoDB.TableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueAllNew,
	}); err != nil {
		return nil, err
	} else if len(output.Attributes) == 0 {
		return nil, nil
	} else {
		var ret T
		if err := attributevalue.UnmarshalMap(output.Attributes, &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

func updateByPrimaryKey[T any](ctx context.Context, s *Store, hk []byte, update expression.UpdateBuilder) (*T, error) {
	if reflect.ValueOf(update).IsZero() {
		// there's nothing to update, so this just becomes a get operation
		return getByPrimaryKey[T](ctx, s, hk, ConsistencyStrongInRegion)
	}

	expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(expression.AttributeExists(expression.Name("_hk"))).Build()
	if err != nil {
		return nil, err
	}

	if output, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"_hk": &types.AttributeValueMemberB{
				Value: hk,
			},
			"_rk": &types.AttributeValueMemberB{
				Value: []byte("_"),
			},
		},
		TableName:                 &s.config.DynamoDB.TableName,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueAllNew,
	}); err != nil {
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			return nil, nil
		}
		return nil, err
	} else if len(output.Attributes) == 0 {
		return nil, nil
	} else {
		var ret T
		if err := attributevalue.UnmarshalMap(output.Attributes, &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

func countByPrimaryHashKey(ctx context.Context, s *Store, hk []byte) (int, error) {
	var ret int
	var exclusiveStartKey map[string]types.AttributeValue
	for {
		if output, err := s.client.Query(ctx, &dynamodb.QueryInput{
			ExpressionAttributeNames: map[string]string{
				"#attr": "_hk",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v": &types.AttributeValueMemberB{
					Value: hk,
				},
			},
			ExclusiveStartKey:      exclusiveStartKey,
			KeyConditionExpression: aws.String("#attr = :v"),
			TableName:              &s.config.DynamoDB.TableName,
			Select:                 types.SelectCount,
		}); err != nil {
			return 0, err
		} else {
			ret += int(output.Count)
			if output.LastEvaluatedKey == nil {
				break
			}
			exclusiveStartKey = output.LastEvaluatedKey
		}
	}
	return ret, nil
}

func getByPrimaryKey[T any](ctx context.Context, s *Store, hk []byte, consistency Consistency) (*T, error) {
	if output, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"_hk": &types.AttributeValueMemberB{
				Value: hk,
			},
			"_rk": &types.AttributeValueMemberB{
				Value: []byte("_"),
			},
		},
		TableName:      &s.config.DynamoDB.TableName,
		ConsistentRead: aws.Bool(consistency == ConsistencyStrongInRegion),
	}); err != nil {
		return nil, err
	} else if len(output.Item) == 0 {
		return nil, nil
	} else {
		var ret T
		if err := attributevalue.UnmarshalMap(output.Item, &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

func getByPrimaryKeys[T any](ctx context.Context, s *Store, hks ...[]byte) ([]*T, error) {
	const maxBatchSize = 100

	keys := make([]map[string]types.AttributeValue, len(hks))
	for i, hk := range hks {
		keys[i] = map[string]types.AttributeValue{
			"_hk": &types.AttributeValueMemberB{
				Value: hk,
			},
			"_rk": &types.AttributeValueMemberB{
				Value: []byte("_"),
			},
		}
	}

	var ret []*T

	for len(keys) > 0 {
		batch := keys
		if len(batch) > maxBatchSize {
			batch = batch[len(batch)-maxBatchSize:]
		}
		keys = keys[:len(keys)-len(batch)]

		if output, err := s.client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				s.config.DynamoDB.TableName: {
					Keys: batch,
				},
			},
		}); err != nil {
			return nil, err
		} else {
			for _, item := range output.Responses[s.config.DynamoDB.TableName] {
				var v T
				if err := attributevalue.UnmarshalMap(item, &v); err != nil {
					return nil, err
				}
				ret = append(ret, &v)
			}
			keys = append(keys, output.UnprocessedKeys[s.config.DynamoDB.TableName].Keys...)
		}
	}

	return ret, nil
}

func attributeValue(value any) (types.AttributeValue, error) {
	switch v := value.(type) {
	case []byte:
		return &types.AttributeValueMemberB{
			Value: v,
		}, nil
	case string:
		return &types.AttributeValueMemberS{
			Value: v,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected attribute value type: %T", value)
	}
}

func getAllByHashKey[T any](ctx context.Context, s *Store, index, name string, value any) ([]*T, error) {
	var indexNamePtr *string
	if index != "" {
		indexNamePtr = &index
	}

	var exclusiveStartKey map[string]types.AttributeValue

	attrValue, err := attributeValue(value)
	if err != nil {
		return nil, err
	}

	var ret []*T
	for {
		if output, err := s.client.Query(ctx, &dynamodb.QueryInput{
			ExpressionAttributeNames: map[string]string{
				"#attr": name,
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v": attrValue,
			},
			KeyConditionExpression: aws.String("#attr = :v"),
			IndexName:              indexNamePtr,
			TableName:              &s.config.DynamoDB.TableName,
			ExclusiveStartKey:      exclusiveStartKey,
		}); err != nil {
			return nil, err
		} else {
			for _, item := range output.Items {
				var v T
				if err := attributevalue.UnmarshalMap(item, &v); err != nil {
					return nil, err
				}
				ret = append(ret, &v)
			}
			if len(output.LastEvaluatedKey) == 0 {
				break
			}
			exclusiveStartKey = output.LastEvaluatedKey
		}
	}
	return ret, nil
}

func getAllByHashAndRangeKey[T any](ctx context.Context, s *Store, index string, hashKey, rangeKey any) ([]*T, error) {
	indexNamePtr := &index

	var exclusiveStartKey map[string]types.AttributeValue

	hashKeyAttrValue, err := attributeValue(hashKey)
	if err != nil {
		return nil, err
	}

	rangeKeyAttrValue, err := attributeValue(rangeKey)
	if err != nil {
		return nil, err
	}

	var ret []*T
	for {
		if output, err := s.client.Query(ctx, &dynamodb.QueryInput{
			ExpressionAttributeNames: map[string]string{
				"#hattr": index + "h",
				"#rattr": index + "r",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":hk": hashKeyAttrValue,
				":rk": rangeKeyAttrValue,
			},
			KeyConditionExpression: aws.String("#hattr = :hk AND #rattr = :rk"),
			IndexName:              indexNamePtr,
			TableName:              &s.config.DynamoDB.TableName,
			ExclusiveStartKey:      exclusiveStartKey,
		}); err != nil {
			return nil, err
		} else {
			for _, item := range output.Items {
				var v T
				if err := attributevalue.UnmarshalMap(item, &v); err != nil {
					return nil, err
				}
				ret = append(ret, &v)
			}
			if len(output.LastEvaluatedKey) == 0 {
				break
			}
			exclusiveStartKey = output.LastEvaluatedKey
		}
	}
	return ret, nil
}

func getAllByHashKeyWithMinRangeKey[T any](ctx context.Context, s *Store, index, hkname string, value any, rkname string, min any) ([]*T, error) {
	var indexNamePtr *string
	if index != "" {
		indexNamePtr = &index
	}

	var exclusiveStartKey map[string]types.AttributeValue

	hkAttrValue, err := attributeValue(value)
	if err != nil {
		return nil, err
	}

	rkAttrValue, err := attributeValue(min)
	if err != nil {
		return nil, err
	}

	var ret []*T
	for {
		if output, err := s.client.Query(ctx, &dynamodb.QueryInput{
			ExpressionAttributeNames: map[string]string{
				"#hk": hkname,
				"#rk": rkname,
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":hkv": hkAttrValue,
				":rkv": rkAttrValue,
			},
			KeyConditionExpression: aws.String("#hk = :hkv AND #rk >= :rkv"),
			IndexName:              indexNamePtr,
			TableName:              &s.config.DynamoDB.TableName,
			ExclusiveStartKey:      exclusiveStartKey,
		}); err != nil {
			return nil, err
		} else {
			for _, item := range output.Items {
				var v T
				if err := attributevalue.UnmarshalMap(item, &v); err != nil {
					return nil, err
				}
				ret = append(ret, &v)
			}
			if len(output.LastEvaluatedKey) == 0 {
				break
			}
			exclusiveStartKey = output.LastEvaluatedKey
		}
	}
	return ret, nil
}

func deleteByPrimaryKey(ctx context.Context, s *Store, hk []byte) error {
	if _, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"_hk": &types.AttributeValueMemberB{
				Value: hk,
			},
			"_rk": &types.AttributeValueMemberB{
				Value: []byte("_"),
			},
		},
		TableName: &s.config.DynamoDB.TableName,
	}); err != nil {
		return err
	}
	return nil
}

func deleteByPrimaryKeys(ctx context.Context, s *Store, hks ...[]byte) error {
	const batchSize = 25

	for len(hks) > 0 {
		batch := hks
		if len(batch) > batchSize {
			batch = batch[:batchSize]
		}
		hks = hks[len(batch):]

		var ops []types.WriteRequest
		for _, hk := range batch {
			ops = append(ops, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						"_hk": &types.AttributeValueMemberB{
							Value: hk,
						},
						"_rk": &types.AttributeValueMemberB{
							Value: []byte("_"),
						},
					},
				},
			})
		}
		if _, err := s.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				s.config.DynamoDB.TableName: ops,
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

type Consistency int

const (
	ConsistencyEventual Consistency = iota
	ConsistencyStrongInRegion
)

type TTL struct {
	Time *attributevalue.UnixTime `dynamodbav:"_ttl,omitempty"`
}

func NewTTL(t time.Time) TTL {
	if t.IsZero() {
		return TTL{}
	}
	unix := attributevalue.UnixTime(t)
	return TTL{
		Time: &unix,
	}
}

func prefixIds(prefix string, ids []model.Id) [][]byte {
	m := make(map[model.Id]struct{}, len(ids))
	ret := make([][]byte, 0, len(ids))
	for _, id := range ids {
		if _, ok := m[id]; !ok {
			m[id] = struct{}{}
			ret = append(ret, []byte(prefix+id.String()))
		}
	}
	return ret
}
