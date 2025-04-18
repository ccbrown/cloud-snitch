package api

import (
	"context"
	"time"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func AWSIntegrationFromModel(integration *model.AWSIntegration) apispec.AWSIntegration {
	ret := apispec.AWSIntegration{
		Id:                               integration.Id.String(),
		CreationTime:                     integration.CreationTime,
		TeamId:                           integration.TeamId.String(),
		Name:                             integration.Name,
		GetAccountNamesFromOrganizations: integration.GetAccountNamesFromOrganizations,
		ManageScps:                       integration.ManageSCPs,
	}
	if trail := integration.CloudTrailTrail; trail != nil {
		ret.CloudtrailTrail = &apispec.AWSIntegrationCloudTrailTrail{
			S3BucketName: trail.S3BucketName,
			S3KeyPrefix:  nilIfEmpty(trail.S3KeyPrefix),
		}
	}
	return ret
}

func (api *API) DeleteAWSIntegration(ctx context.Context, request apispec.DeleteAWSIntegrationRequestObject) (apispec.DeleteAWSIntegrationResponseObject, error) {
	sess := ctxSession(ctx)
	integrationId := model.Id(request.IntegrationId)

	if err := sess.DeleteAWSIntegrationById(ctx, integrationId, emptyIfNil(emptyIfNil(request.Body).DeleteAssociatedData)); err != nil {
		return nil, err
	} else {
		return apispec.DeleteAWSIntegration200JSONResponse{}, nil
	}
}

func (api *API) UpdateAWSIntegration(ctx context.Context, request apispec.UpdateAWSIntegrationRequestObject) (apispec.UpdateAWSIntegrationResponseObject, error) {
	sess := ctxSession(ctx)

	patch := app.AWSIntegrationPatch{
		Name: request.Body.Name,
	}

	if integration, err := sess.PatchAWSIntegrationById(ctx, model.Id(request.IntegrationId), patch); err != nil {
		return nil, err
	} else if integration == nil {
		return nil, app.NotFoundError("No such integration.")
	} else {
		return apispec.UpdateAWSIntegration200JSONResponse(AWSIntegrationFromModel(integration)), nil
	}
}

func (api *API) GetAWSIntegrationsByTeamId(ctx context.Context, request apispec.GetAWSIntegrationsByTeamIdRequestObject) (apispec.GetAWSIntegrationsByTeamIdResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if memberships, err := sess.GetAWSIntegrationsByTeamId(ctx, teamId); err != nil {
		return nil, err
	} else {
		return apispec.GetAWSIntegrationsByTeamId200JSONResponse(mapSlice(memberships, AWSIntegrationFromModel)), nil
	}
}

func (api *API) CreateAWSIntegration(ctx context.Context, request apispec.CreateAWSIntegrationRequestObject) (apispec.CreateAWSIntegrationResponseObject, error) {
	sess := ctxSession(ctx)

	input := app.CreateAWSIntegrationInput{
		Name:                  request.Body.Name,
		TeamId:                model.Id(request.TeamId),
		RoleARN:               request.Body.RoleArn,
		QueueReportGeneration: emptyIfNil(request.Body.QueueReportGeneration),
	}
	if request.Body.GetAccountNamesFromOrganizations != nil {
		input.GetAccountNamesFromOrganizations = *request.Body.GetAccountNamesFromOrganizations
	}
	if request.Body.ManageScps != nil {
		input.ManageSCPs = *request.Body.ManageScps
	}
	if trail := request.Body.CloudtrailTrail; trail != nil {
		input.CloudTrailTrail = &app.CreateAWSIntegrationCloudTrailTrailInput{
			S3BucketName: trail.S3BucketName,
			S3KeyPrefix:  emptyIfNil(trail.S3KeyPrefix),
		}
	}

	if integration, err := sess.CreateAWSIntegration(ctx, input); err != nil {
		return nil, err
	} else {
		return apispec.CreateAWSIntegration200JSONResponse(AWSIntegrationFromModel(integration)), nil
	}
}

func (api *API) QueueAWSIntegrationReportGeneration(ctx context.Context, request apispec.QueueAWSIntegrationReportGenerationRequestObject) (apispec.QueueAWSIntegrationReportGenerationResponseObject, error) {
	sess := ctxSession(ctx)
	integrationId := model.Id(request.IntegrationId)

	if err := sess.QueueAWSIntegrationReportGeneration(ctx, app.QueueAWSIntegrationReportGenerationInput{
		IntegrationId: integrationId,
		StartTime:     request.Body.StartTime,
		Duration:      time.Second * time.Duration(request.Body.DurationSeconds),
		Retention:     ReportRetentionFromSpec(request.Body.Retention),
	}); err != nil {
		return nil, err
	} else {
		return apispec.QueueAWSIntegrationReportGeneration200Response{}, nil
	}
}

func AWSSCPFromModel(scp *model.AWSSCP) apispec.AWSSCP {
	return apispec.AWSSCP{
		Content: scp.Content,
	}
}

func (api *API) GetManagedAWSSCP(ctx context.Context, request apispec.GetManagedAWSSCPRequestObject) (apispec.GetManagedAWSSCPResponseObject, error) {
	sess := ctxSession(ctx)

	if scp, err := sess.GetManagedAWSSCPByTeamAndAccountId(ctx, model.Id(request.TeamId), request.AccountId); err != nil {
		return nil, err
	} else if scp == nil {
		return nil, app.NotFoundError("No such SCP.")
	} else {
		return apispec.GetManagedAWSSCP200JSONResponse(AWSSCPFromModel(scp)), nil
	}
}

func (api *API) PutManagedAWSSCP(ctx context.Context, request apispec.PutManagedAWSSCPRequestObject) (apispec.PutManagedAWSSCPResponseObject, error) {
	sess := ctxSession(ctx)

	if scp, err := sess.PutManagedAWSSCPByTeamAndAccountId(ctx, model.Id(request.TeamId), request.AccountId, app.PutManagedAWSSCPInput{
		Content: request.Body.Content,
	}); err != nil {
		return nil, err
	} else if scp == nil {
		return nil, app.NotFoundError("No such account.")
	} else {
		return apispec.PutManagedAWSSCP200JSONResponse(AWSSCPFromModel(scp)), nil
	}
}

func AWSAccessReportFromModel(report *model.AWSAccessReport) apispec.AWSAccessReport {
	ret := apispec.AWSAccessReport{
		Services: make([]apispec.AWSAccessReportService, 0, len(report.Services)),
	}
	for _, service := range report.Services {
		out := apispec.AWSAccessReportService{
			Name:      service.Name,
			Namespace: service.Namespace,
		}
		if !service.LastAuthenticationTime.IsZero() {
			out.LastAuthenticationTime = &service.LastAuthenticationTime
		}
		ret.Services = append(ret.Services, out)
	}
	return ret
}

func (api *API) GetAWSAccessReport(ctx context.Context, request apispec.GetAWSAccessReportRequestObject) (apispec.GetAWSAccessReportResponseObject, error) {
	sess := ctxSession(ctx)

	if report, err := sess.GetAWSAccessReportByTeamAndAccountId(ctx, model.Id(request.TeamId), request.AccountId); err != nil {
		return nil, err
	} else if report == nil {
		return nil, app.NotFoundError("No such account.")
	} else {
		return apispec.GetAWSAccessReport200JSONResponse(AWSAccessReportFromModel(report)), nil
	}
}
