package store

import (
	"context"

	"github.com/ccbrown/cloud-snitch/backend/model"
)

type IndexedReport struct {
	*model.Report

	PrimaryIndex
	ByteByteIndex1

	TTL
}

func (s *Store) PutReport(ctx context.Context, report *model.Report) error {
	return s.put(ctx, &IndexedReport{
		Report: report,
		PrimaryIndex: PrimaryIndex{
			HashKey:  []byte("report:" + report.Id),
			RangeKey: []byte("_"),
		},
		ByteByteIndex1: ByteByteIndex1{
			HashKey:  []byte("reports:" + report.TeamId),
			RangeKey: []byte(report.Id),
		},
		TTL: NewTTL(report.ExpirationTime),
	})
}

func (s *Store) GetReportById(ctx context.Context, id model.Id) (*model.Report, error) {
	return getByPrimaryKey[model.Report](ctx, s, []byte("report:"+id), ConsistencyEventual)
}

func (s *Store) GetReportsByTeamId(ctx context.Context, teamId model.Id) ([]*model.Report, error) {
	return getAllByHashKey[model.Report](ctx, s, "_bb1", "_bb1h", []byte("reports:"+teamId))
}

func (s *Store) DeleteReportById(ctx context.Context, id model.Id) error {
	return s.DeleteReportsByIds(ctx, id)
}

func (s *Store) DeleteReportsByIds(ctx context.Context, ids ...model.Id) error {
	keys := make([][]byte, len(ids))
	for i, id := range ids {
		keys[i] = []byte("report:" + id)
	}
	return deleteByPrimaryKeys(ctx, s, keys...)
}
