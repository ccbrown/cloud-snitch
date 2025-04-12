package api

import (
	"context"
	"time"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func ReportRetentionFromSpec(retention apispec.ReportRetention) model.ReportRetention {
	switch retention {
	case apispec.ONEWEEK:
		return model.ReportRetentionOneWeek
	case apispec.TWOWEEKS:
		return model.ReportRetentionTwoWeeks
	default:
		panic("unknown report retention")
	}
}

func ReportScopeFromModel(scope *model.ReportScope) apispec.ReportScope {
	return apispec.ReportScope{
		StartTime:       scope.StartTime,
		DurationSeconds: int(scope.Duration / time.Second),
		Aws:             ReportScopeAWSFromModel(&scope.AWS),
	}
}

func ReportScopeAWSFromModel(aws *model.ReportScopeAWS) apispec.ReportScopeAWS {
	return apispec.ReportScopeAWS{
		AccountId: aws.AccountId,
		Region:    aws.Region,
	}
}

func ReportFromModel(report *model.Report) apispec.Report {
	return apispec.Report{
		Id:                        report.Id.String(),
		TeamId:                    report.TeamId.String(),
		AwsIntegrationId:          report.AWSIntegrationId.String(),
		Scope:                     ReportScopeFromModel(&report.Scope),
		DownloadUrl:               report.DownloadURL,
		Size:                      report.Size,
		SourceBytes:               report.SourceBytes,
		IsIncomplete:              &report.IsIncomplete,
		GenerationDurationSeconds: int(report.GenerationDuration / time.Second),
	}
}

func (api *API) GetReportsByTeamId(ctx context.Context, request apispec.GetReportsByTeamIdRequestObject) (apispec.GetReportsByTeamIdResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if reports, err := sess.GetReportsByTeamId(ctx, teamId); err != nil {
		return nil, err
	} else {
		return apispec.GetReportsByTeamId200JSONResponse(mapSlice(reports, ReportFromModel)), nil
	}
}

func (api *API) DeleteReportById(ctx context.Context, request apispec.DeleteReportByIdRequestObject) (apispec.DeleteReportByIdResponseObject, error) {
	sess := ctxSession(ctx)
	reportId := model.Id(request.ReportId)

	if err := sess.DeleteReportById(ctx, reportId); err != nil {
		return nil, err
	} else {
		return apispec.DeleteReportById200Response{}, nil
	}
}
