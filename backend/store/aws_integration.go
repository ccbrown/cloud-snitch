package store

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

	"github.com/ccbrown/cloud-snitch/backend/model"
)

type IndexedAWSIntegration struct {
	*model.AWSIntegration

	PrimaryIndex
	ByteByteIndex1
}

func (s *Store) PutAWSIntegration(ctx context.Context, integration *model.AWSIntegration) error {
	return s.put(ctx, &IndexedAWSIntegration{
		AWSIntegration: integration,
		PrimaryIndex: PrimaryIndex{
			HashKey:  []byte("aws_integration:" + integration.Id),
			RangeKey: []byte("_"),
		},
		ByteByteIndex1: ByteByteIndex1{
			HashKey:  []byte("aws_integrations:" + integration.TeamId),
			RangeKey: []byte(integration.Id),
		},
	})
}

func (s *Store) GetAWSIntegrationById(ctx context.Context, id model.Id) (*model.AWSIntegration, error) {
	return getByPrimaryKey[model.AWSIntegration](ctx, s, []byte("aws_integration:"+id), ConsistencyEventual)
}

type AWSIntegrationPatch struct {
	Name *string
}

func (p *AWSIntegrationPatch) Apply(update expression.UpdateBuilder) expression.UpdateBuilder {
	if p.Name != nil {
		update = update.Set(expression.Name("Name"), expression.Value(p.Name))
	}
	return update
}

func (s *Store) PatchAWSIntegrationById(ctx context.Context, id model.Id, patch *AWSIntegrationPatch) (*model.AWSIntegration, error) {
	update := patch.Apply(expression.UpdateBuilder{})
	return updateByPrimaryKey[model.AWSIntegration](ctx, s, []byte("aws_integration:"+id), update)
}

func (s *Store) GetAWSIntegrationsByTeamId(ctx context.Context, teamId model.Id) ([]*model.AWSIntegration, error) {
	return getAllByHashKey[model.AWSIntegration](ctx, s, "_bb1", "_bb1h", []byte("aws_integrations:"+teamId))
}

func (s *Store) DeleteAWSIntegrationById(ctx context.Context, id model.Id) error {
	return deleteByPrimaryKey(ctx, s, []byte("aws_integration:"+id))
}

type IndexedAWSIntegrationRecon struct {
	*model.AWSIntegrationRecon

	PrimaryIndex
	ByteByteIndex1

	TTL
}

func (s *Store) PutAWSIntegrationRecon(ctx context.Context, recon *model.AWSIntegrationRecon) error {
	return s.put(ctx, &IndexedAWSIntegrationRecon{
		AWSIntegrationRecon: recon,
		PrimaryIndex: PrimaryIndex{
			HashKey:  []byte("aws_integration_recon:" + recon.AWSIntegrationId),
			RangeKey: []byte("_"),
		},
		ByteByteIndex1: ByteByteIndex1{
			HashKey:  []byte("aws_integration_recons:" + recon.TeamId),
			RangeKey: []byte(recon.AWSIntegrationId),
		},
		TTL: NewTTL(recon.ExpirationTime),
	})
}

func (s *Store) GetAWSIntegrationReconsByTeamId(ctx context.Context, teamId model.Id) ([]*model.AWSIntegrationRecon, error) {
	return getAllByHashKey[model.AWSIntegrationRecon](ctx, s, "_bb1", "_bb1h", []byte("aws_integration_recons:"+teamId))
}

func (s *Store) DeleteAWSIntegrationReconByAWSIntegrationId(ctx context.Context, id model.Id) error {
	return deleteByPrimaryKey(ctx, s, []byte("aws_integration_recon:"+id))
}
