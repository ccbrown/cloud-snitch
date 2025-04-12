package api

import (
	"context"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func AWSRegionFromApp(id string, region *app.AWSRegion) apispec.AWSRegion {
	return apispec.AWSRegion{
		Id:                 id,
		Name:               region.Name,
		GeolocationCountry: region.GeolocationCountry,
		GeolocationRegion:  region.GeolocationRegion,
		Partition:          region.Partition,
		Latitude:           float32(region.Latitude),
		Longitude:          float32(region.Longitude),
	}
}

func (api *API) GetAWSRegions(ctx context.Context, request apispec.GetAWSRegionsRequestObject) (apispec.GetAWSRegionsResponseObject, error) {
	ret := make([]apispec.AWSRegion, 0, len(app.KnownAWSRegions))
	for id, region := range app.KnownAWSRegions {
		ret = append(ret, AWSRegionFromApp(id, &region))
	}
	return apispec.GetAWSRegions200JSONResponse(ret), nil
}

func (api *API) GetAWSAccountsByTeamId(ctx context.Context, request apispec.GetAWSAccountsByTeamIdRequestObject) (apispec.GetAWSAccountsByTeamIdResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if recons, err := sess.GetAWSIntegrationReconsByTeamId(ctx, teamId); err != nil {
		return nil, err
	} else {
		accounts := map[string]apispec.AWSAccount{}
		for _, recon := range recons {
			for _, account := range recon.Accounts {
				if existing, ok := accounts[account.Id]; ok {
					if account.Name != "" {
						existing.Name = &account.Name
					}
				} else {
					accounts[account.Id] = apispec.AWSAccount{
						Id:   account.Id,
						Name: nilIfEmpty(account.Name),
					}
				}
			}
		}
		ret := make([]apispec.AWSAccount, 0, len(accounts))
		for _, account := range accounts {
			ret = append(ret, account)
		}
		return apispec.GetAWSAccountsByTeamId200JSONResponse(ret), nil
	}
}
