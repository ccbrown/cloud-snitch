package store

import (
	"context"

	"github.com/ccbrown/cloud-snitch/backend/model"
)

type IndexedTeamBillableAccount struct {
	*model.TeamBillableAccount

	PrimaryIndex

	TTL
}

func (s *Store) PutTeamBillableAccount(ctx context.Context, account *model.TeamBillableAccount) error {
	return s.put(ctx, &IndexedTeamBillableAccount{
		TeamBillableAccount: account,
		PrimaryIndex: PrimaryIndex{
			HashKey:  []byte("billable_aws_accounts:" + account.TeamId),
			RangeKey: []byte(account.Id),
		},
		TTL: NewTTL(account.ExpirationTime),
	})
}

func (s *Store) GetTeamBillableAccountCountByTeamId(ctx context.Context, teamId model.Id) (int, error) {
	return countByPrimaryHashKey(ctx, s, []byte("billable_aws_accounts:"+teamId))
}
