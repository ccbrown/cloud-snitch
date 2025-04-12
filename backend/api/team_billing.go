package api

import (
	"context"

	"github.com/ccbrown/cloud-snitch/backend/api/apispec"
	"github.com/ccbrown/cloud-snitch/backend/app"
	"github.com/ccbrown/cloud-snitch/backend/model"
)

func CurrencyAmountFromModel(ca model.CurrencyAmount) apispec.CurrencyAmount {
	return apispec.CurrencyAmount{
		Text: ca.String(),
	}
}

func TeamBillingAddressFromModel(addr *model.TeamBillingAddress) apispec.TeamBillingAddress {
	return apispec.TeamBillingAddress{
		Line1:      addr.Line1,
		Line2:      addr.Line2,
		City:       addr.City,
		State:      addr.State,
		Country:    addr.Country,
		PostalCode: addr.PostalCode,
	}
}

func TeamBillingProfileFromModel(profile *model.TeamBillingProfile) apispec.TeamBillingProfile {
	ret := apispec.TeamBillingProfile{
		Name:    profile.Name,
		Address: TeamBillingAddressFromModel(&profile.Address),
	}
	if profile.Balance != nil {
		ret.Balance = pointer(CurrencyAmountFromModel(*profile.Balance))
	}
	return ret
}

func TeamBillingAddressFromSpec(addr apispec.TeamBillingAddress) model.TeamBillingAddress {
	return model.TeamBillingAddress{
		Line1:      addr.Line1,
		Line2:      addr.Line2,
		City:       addr.City,
		State:      addr.State,
		Country:    addr.Country,
		PostalCode: addr.PostalCode,
	}
}

func (api *API) GetTeamBillingProfile(ctx context.Context, request apispec.GetTeamBillingProfileRequestObject) (apispec.GetTeamBillingProfileResponseObject, error) {
	sess := ctxSession(ctx)
	teamId := model.Id(request.TeamId)

	if profile, err := sess.GetTeamBillingProfileById(ctx, teamId); err != nil {
		return nil, err
	} else if profile == nil {
		return nil, app.NotFoundError("No such billing profile.")
	} else {
		return apispec.GetTeamBillingProfile200JSONResponse(TeamBillingProfileFromModel(profile)), nil
	}
}

func (api *API) CreateTeamBillingProfile(ctx context.Context, request apispec.CreateTeamBillingProfileRequestObject) (apispec.CreateTeamBillingProfileResponseObject, error) {
	sess := ctxSession(ctx)

	input := app.CreateTeamBillingProfileInput{
		Name:    request.Body.Name,
		Address: TeamBillingAddressFromSpec(request.Body.Address),
	}

	if profile, err := sess.CreateTeamBillingProfileById(ctx, model.Id(request.TeamId), input); err != nil {
		return nil, err
	} else {
		return apispec.CreateTeamBillingProfile200JSONResponse(TeamBillingProfileFromModel(profile)), nil
	}
}

func (api *API) UpdateTeamBillingProfile(ctx context.Context, request apispec.UpdateTeamBillingProfileRequestObject) (apispec.UpdateTeamBillingProfileResponseObject, error) {
	sess := ctxSession(ctx)

	patch := app.TeamBillingProfilePatch{
		Name: request.Body.Name,
	}
	if request.Body.Address != nil {
		patch.Address = pointer(TeamBillingAddressFromSpec(*request.Body.Address))
	}

	if profile, err := sess.PatchTeamBillingProfileById(ctx, model.Id(request.TeamId), patch); err != nil {
		return nil, err
	} else if profile == nil {
		return nil, app.NotFoundError("No such billing profile.")
	} else {
		return apispec.UpdateTeamBillingProfile200JSONResponse(TeamBillingProfileFromModel(profile)), nil
	}
}

func TeamPaymentMethodFromModel(method *model.TeamPaymentMethod) apispec.TeamPaymentMethod {
	var ret apispec.TeamPaymentMethod

	if card := method.Card; card != nil {
		ret.FromTeamPaymentMethodCard(apispec.TeamPaymentMethodCard{
			Last4Digits:     card.Last4Digits,
			ExpirationMonth: card.ExpirationMonth,
			ExpirationYear:  card.ExpirationYear,
		})
	} else if account := method.USBankAccount; account != nil {
		ret.FromTeamPaymentMethodUSBankAccount(apispec.TeamPaymentMethodUSBankAccount{
			Last4Digits: account.Last4Digits,
		})
	} else {
		ret.FromTeamPaymentMethodOther(apispec.TeamPaymentMethodOther{})
	}

	return ret
}

func (api *API) GetTeamPaymentMethod(ctx context.Context, request apispec.GetTeamPaymentMethodRequestObject) (apispec.GetTeamPaymentMethodResponseObject, error) {
	sess := ctxSession(ctx)

	if method, err := sess.GetTeamPaymentMethodById(ctx, model.Id(request.TeamId)); err != nil {
		return nil, err
	} else if method == nil {
		return nil, app.NotFoundError("No such payment method.")
	} else {
		return apispec.GetTeamPaymentMethod200JSONResponse(TeamPaymentMethodFromModel(method)), nil
	}
}

func (api *API) PutTeamPaymentMethod(ctx context.Context, request apispec.PutTeamPaymentMethodRequestObject) (apispec.PutTeamPaymentMethodResponseObject, error) {
	sess := ctxSession(ctx)

	input := app.PutTeamPaymentMethodInput{
		StripePaymentMethodId: request.Body.StripePaymentMethodId,
	}
	if r := ctxRequest(ctx); r != nil {
		input.IPAddress = api.httpRequestIPAddress(r)
		input.UserAgent = r.UserAgent()
	}

	if method, err := sess.PutTeamPaymentMethodById(ctx, model.Id(request.TeamId), input); err != nil {
		return nil, err
	} else if method == nil {
		return nil, app.NotFoundError("No such payment method.")
	} else {
		return apispec.PutTeamPaymentMethod200JSONResponse(TeamPaymentMethodFromModel(method)), nil
	}
}

func TeamSubscriptionFromModel(sub *model.TeamSubscription) apispec.TeamSubscription {
	ret := apispec.TeamSubscription{
		Name:     sub.Name,
		Accounts: float32(sub.Accounts),
	}
	if price := sub.Price; price != nil {
		ret.Price = &apispec.TeamSubscriptionPrice{}
		if price.AccountMonth != nil {
			ret.Price.AccountMonth = &apispec.CurrencyAmount{
				Text: price.AccountMonth.String(),
			}
		}
	}
	return ret
}

func (api *API) GetTeamSubscription(ctx context.Context, request apispec.GetTeamSubscriptionRequestObject) (apispec.GetTeamSubscriptionResponseObject, error) {
	sess := ctxSession(ctx)

	if subscription, err := sess.GetTeamSubscriptionById(ctx, model.Id(request.TeamId)); err != nil {
		return nil, err
	} else if subscription == nil {
		return nil, app.NotFoundError("No such subscription.")
	} else {
		return apispec.GetTeamSubscription200JSONResponse(TeamSubscriptionFromModel(subscription)), nil
	}
}

func TeamSubscriptionTierFromApp(tier apispec.TeamSubscriptionTier) app.TeamSubscriptionTier {
	switch tier {
	case apispec.INDIVIDUAL:
		return app.TeamSubscriptionTierIndividual
	case apispec.TEAM:
		return app.TeamSubscriptionTierTeam
	default:
		panic("unknown subscription tier")
	}
}

func (api *API) CreateTeamSubscription(ctx context.Context, request apispec.CreateTeamSubscriptionRequestObject) (apispec.CreateTeamSubscriptionResponseObject, error) {
	sess := ctxSession(ctx)

	input := app.CreateTeamSubscriptionInput{
		Tier: TeamSubscriptionTierFromApp(request.Body.Tier),
	}

	if subscription, err := sess.CreateTeamSubscriptionById(ctx, model.Id(request.TeamId), input); err != nil {
		return nil, err
	} else {
		return apispec.CreateTeamSubscription200JSONResponse(TeamSubscriptionFromModel(subscription)), nil
	}
}

func (api *API) UpdateTeamSubscription(ctx context.Context, request apispec.UpdateTeamSubscriptionRequestObject) (apispec.UpdateTeamSubscriptionResponseObject, error) {
	sess := ctxSession(ctx)

	input := app.UpdateTeamSubscriptionInput{}
	if request.Body.Tier != nil {
		input.Tier = pointer(TeamSubscriptionTierFromApp(*request.Body.Tier))
	}

	if subscription, err := sess.UpdateTeamSubscriptionById(ctx, model.Id(request.TeamId), input); err != nil {
		return nil, err
	} else if subscription == nil {
		return nil, app.NotFoundError("No such subscription.")
	} else {
		return apispec.UpdateTeamSubscription200JSONResponse(TeamSubscriptionFromModel(subscription)), nil
	}
}
