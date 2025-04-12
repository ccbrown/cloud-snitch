package app

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/stripe/stripe-go/v81"

	"github.com/ccbrown/cloud-snitch/backend/model"
	"github.com/ccbrown/cloud-snitch/backend/store"
)

type TeamBillingProfilePatch struct {
	Name    *string
	Address *model.TeamBillingAddress
}

func ValidateBillingAddress(address model.TeamBillingAddress) UserFacingError {
	if address.Country == "" {
		return NewUserError("A country is required.")
	} else if address.PostalCode == "" {
		return NewUserError("A postal code is required.")
	}
	return nil
}

func teamBillingProfileFromCustomer(c *stripe.Customer) *model.TeamBillingProfile {
	if c == nil {
		return nil
	}

	ret := &model.TeamBillingProfile{
		Name: c.Name,
		Address: model.TeamBillingAddress{
			Country:    c.Address.Country,
			PostalCode: c.Address.PostalCode,
			Line1:      nilIfEmpty(c.Address.Line1),
			Line2:      nilIfEmpty(c.Address.Line2),
			State:      nilIfEmpty(c.Address.State),
			City:       nilIfEmpty(c.Address.City),
		},
	}
	if c.Balance != 0 {
		ret.Balance = &model.CurrencyAmount{
			Currency: c.Currency,
			Amount:   c.Balance,
		}
	}
	return ret
}

func (s *Session) GetTeamBillingProfileById(ctx context.Context, teamId model.Id) (*model.TeamBillingProfile, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyStrongInRegion)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if team == nil {
		return nil, NotFoundError("Team not found.")
	} else if team.StripeCustomerId == "" {
		return nil, nil
	}

	customer, err := s.app.stripe.Customers.Get(team.StripeCustomerId, nil)
	if err != nil {
		return nil, s.SanitizedError(fmt.Errorf("unable to get stripe customer: %w", err))
	}
	return teamBillingProfileFromCustomer(customer), nil
}

type CreateTeamBillingProfileInput struct {
	Name    string
	Address model.TeamBillingAddress
}

func (s *Session) CreateTeamBillingProfileById(ctx context.Context, teamId model.Id, input CreateTeamBillingProfileInput) (*model.TeamBillingProfile, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyStrongInRegion)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if team == nil {
		return nil, NotFoundError("Team not found.")
	} else if team.StripeCustomerId != "" {
		return nil, NewUserError("A billing profile already exists for this team.")
	}

	if len(input.Name) == 0 {
		return nil, NewUserError("A name is required.")
	} else if err := ValidateBillingAddress(input.Address); err != nil {
		return nil, err
	}

	params := &stripe.CustomerParams{
		Name: &input.Name,
		Metadata: map[string]string{
			"team_id":   team.Id.String(),
			"team_name": team.Name,
		},
		Address: &stripe.AddressParams{
			City:       input.Address.City,
			Country:    &input.Address.Country,
			Line1:      input.Address.Line1,
			Line2:      input.Address.Line2,
			PostalCode: &input.Address.PostalCode,
			State:      input.Address.State,
		},
	}

	result, err := s.app.stripe.Customers.New(params)
	if IsStripeBadRequestError(err) {
		return nil, NewUserError("Request rejected by Stripe. Please double check your billing information.")
	} else if err != nil {
		return nil, s.SanitizedError(fmt.Errorf("unable to create stripe customer: %w", err))
	}
	if _, err := s.app.store.PatchTeamById(ctx, teamId, &store.TeamPatch{
		StripeCustomerId: &result.ID,
	}); err != nil {
		return nil, s.SanitizedError(err)
	}
	return teamBillingProfileFromCustomer(result), nil
}

func (s *Session) PatchTeamBillingProfileById(ctx context.Context, teamId model.Id, patch TeamBillingProfilePatch) (*model.TeamBillingProfile, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyStrongInRegion)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if team == nil {
		return nil, NotFoundError("Team not found.")
	} else if team.StripeCustomerId == "" {
		return nil, nil
	}

	params := &stripe.CustomerParams{
		Metadata: map[string]string{
			"team_id":   team.Id.String(),
			"team_name": team.Name,
		},
	}

	if patch.Name != nil {
		if len(*patch.Name) == 0 {
			return nil, NewUserError("A name is required.")
		}
		params.Name = patch.Name
	}

	if patch.Address != nil {
		if err := ValidateBillingAddress(*patch.Address); err != nil {
			return nil, err
		}

		params.Address = &stripe.AddressParams{
			City:       patch.Address.City,
			Country:    &patch.Address.Country,
			Line1:      patch.Address.Line1,
			Line2:      patch.Address.Line2,
			PostalCode: &patch.Address.PostalCode,
			State:      patch.Address.State,
		}
	}

	if result, err := s.app.stripe.Customers.Update(team.StripeCustomerId, params); err != nil {
		if IsStripeBadRequestError(err) {
			return nil, NewUserError("Request rejected by Stripe. Please double check your billing information.")
		} else {
			return nil, s.SanitizedError(fmt.Errorf("unable to update stripe customer: %w", err))
		}
	} else {
		return teamBillingProfileFromCustomer(result), nil
	}
}

func teamPaymentMethodFromStripePaymentMethod(method *stripe.PaymentMethod) *model.TeamPaymentMethod {
	ret := &model.TeamPaymentMethod{}

	if method.Type == stripe.PaymentMethodTypeCard {
		card := method.Card
		ret.Card = &model.TeamPaymentMethodCard{
			Last4Digits:     card.Last4,
			ExpirationMonth: int(card.ExpMonth),
			ExpirationYear:  int(card.ExpYear),
		}
	}

	if method.Type == stripe.PaymentMethodTypeUSBankAccount {
		acct := method.USBankAccount
		ret.USBankAccount = &model.TeamPaymentMethodUSBankAccount{
			Last4Digits: acct.Last4,
		}
	}

	return ret
}

func (s *Session) GetTeamPaymentMethodById(ctx context.Context, teamId model.Id) (*model.TeamPaymentMethod, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyStrongInRegion)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if team == nil {
		return nil, NotFoundError("Team not found.")
	} else if team.StripeCustomerId == "" {
		return nil, nil
	}

	if customer, err := s.app.stripe.Customers.Get(team.StripeCustomerId, &stripe.CustomerParams{
		Expand: []*string{stripe.String("invoice_settings.default_payment_method")},
	}); err != nil {
		return nil, s.SanitizedError(fmt.Errorf("unable to get stripe customer: %w", err))
	} else if customer == nil || customer.InvoiceSettings == nil || customer.InvoiceSettings.DefaultPaymentMethod == nil {
		return nil, nil
	} else {
		return teamPaymentMethodFromStripePaymentMethod(customer.InvoiceSettings.DefaultPaymentMethod), nil
	}
}

type PutTeamPaymentMethodInput struct {
	StripePaymentMethodId string
	IPAddress             string
	UserAgent             string
}

func (s *Session) PutTeamPaymentMethodById(ctx context.Context, teamId model.Id, input PutTeamPaymentMethodInput) (*model.TeamPaymentMethod, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyStrongInRegion)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if team == nil {
		return nil, NotFoundError("Team not found.")
	} else if team.StripeCustomerId == "" {
		return nil, NewUserError("Team does not have a billing profile configured.")
	}

	method, err := s.app.stripe.PaymentMethods.Get(input.StripePaymentMethodId, nil)
	if err != nil {
		return nil, s.SanitizedError(fmt.Errorf("error getting stripe payment method: %w", err))
	} else if method == nil {
		return nil, NewUserError("Payment method not found.")
	}

	params := &stripe.SetupIntentParams{
		AutomaticPaymentMethods: &stripe.SetupIntentAutomaticPaymentMethodsParams{
			Enabled:        stripe.Bool(true),
			AllowRedirects: stripe.String("never"),
		},
		Confirm:       stripe.Bool(true),
		Customer:      stripe.String(team.StripeCustomerId),
		PaymentMethod: stripe.String(input.StripePaymentMethodId),
		MandateData: &stripe.SetupIntentMandateDataParams{
			CustomerAcceptance: &stripe.SetupIntentMandateDataCustomerAcceptanceParams{
				Type: stripe.MandateCustomerAcceptanceTypeOnline,
				Online: &stripe.SetupIntentMandateDataCustomerAcceptanceOnlineParams{
					IPAddress: stripe.String(input.IPAddress),
					UserAgent: stripe.String(input.UserAgent),
				},
			},
		},
	}
	if _, err := s.app.stripe.SetupIntents.New(params); err != nil {
		return nil, s.SanitizedError(fmt.Errorf("error confirming stripe payment method: %w", err))
	}

	if _, err := s.app.stripe.Customers.Update(team.StripeCustomerId, &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(input.StripePaymentMethodId),
		},
	}); err != nil {
		return nil, s.SanitizedError(fmt.Errorf("error updating stripe customer: %w", err))
	}

	return teamPaymentMethodFromStripePaymentMethod(method), nil
}

func (a *App) getSingleStripeSubscription(customerId string) (*stripe.Subscription, error) {
	subscriptions := a.stripe.Subscriptions.List(&stripe.SubscriptionListParams{
		Customer: stripe.String(customerId),
	})
	if subscriptions.Next() {
		return subscriptions.Subscription(), nil
	}
	if err := subscriptions.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (a *App) teamSubscriptionFromStripe(subscription *stripe.Subscription) (*model.TeamSubscription, error) {
	if subscription == nil || len(subscription.Items.Data) == 0 {
		return nil, nil
	}

	ret := &model.TeamSubscription{
		Name: "Custom Subscription",
	}

	if len(subscription.Items.Data) == 1 {
		item := subscription.Items.Data[0]
		productId := item.Plan.Product.ID
		product, err := a.stripe.Products.Get(productId, nil)
		if err != nil {
			return nil, fmt.Errorf("unable to get stripe product: %w", err)
		} else if product != nil {
			ret.Name = product.Name
		}
	}

	for _, item := range subscription.Items.Data {
		if item.Price.Metadata["use_account_quantity"] == "true" {
			ret.Accounts = int(item.Quantity)

			if item.Price != nil && item.Price.Type == stripe.PriceTypeRecurring && item.Price.Recurring != nil && item.Price.Recurring.Interval == stripe.PriceRecurringIntervalMonth {
				ret.Price = &model.TeamSubscriptionPrice{
					AccountMonth: &model.CurrencyAmount{
						Currency: item.Price.Currency,
						Amount:   item.Price.UnitAmount,
					},
				}
			}

			break
		}
	}

	return ret, nil
}

type TeamSubscriptionTier int

const (
	TeamSubscriptionTierIndividual TeamSubscriptionTier = iota
	TeamSubscriptionTierTeam
)

func (s *Session) ValidateTierForTeam(ctx context.Context, tier TeamSubscriptionTier, team *model.Team) UserFacingError {
	if tier == TeamSubscriptionTierIndividual && team.Entitlements.TeamFeatures {
		// If the team is downgrading, make sure they're not using any team features.
		members, err := s.app.store.GetTeamMembershipsByTeamId(ctx, team.Id)
		if err != nil {
			return s.SanitizedError(err)
		}
		if len(members) > 1 {
			return NewUserError("This team has multiple members. Please remove them to downgrade to an individual subscription.")
		}
	}

	return nil
}

func (a *App) TeamSubscriptionTierPriceId(tier TeamSubscriptionTier) string {
	switch tier {
	case TeamSubscriptionTierIndividual:
		return a.config.Pricing.IndividualSubscriptionStripePriceId
	case TeamSubscriptionTierTeam:
		return a.config.Pricing.TeamSubscriptionStripePriceId
	default:
		panic("invalid subscription tier")
	}
}

type CreateTeamSubscriptionInput struct {
	Tier TeamSubscriptionTier
}

func (s *Session) CreateTeamSubscriptionById(ctx context.Context, teamId model.Id, input CreateTeamSubscriptionInput) (*model.TeamSubscription, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyStrongInRegion)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if team == nil {
		return nil, NotFoundError("Team not found.")
	} else if team.StripeCustomerId == "" {
		return nil, nil
	} else if subscription, err := s.app.getSingleStripeSubscription(team.StripeCustomerId); err != nil {
		return nil, s.SanitizedError(fmt.Errorf("unable to get stripe subscription: %w", err))
	} else if subscription != nil {
		return nil, NewUserError("A subscription already exists for this team.")
	} else if err := s.ValidateTierForTeam(ctx, input.Tier, team); err != nil {
		return nil, err
	}

	if customer, err := s.app.stripe.Customers.Get(team.StripeCustomerId, &stripe.CustomerParams{
		Expand: []*string{stripe.String("invoice_settings.default_payment_method")},
	}); err != nil {
		return nil, s.SanitizedError(fmt.Errorf("unable to get stripe customer: %w", err))
	} else if customer == nil || customer.InvoiceSettings == nil || customer.InvoiceSettings.DefaultPaymentMethod == nil {
		return nil, NewUserError("A payment method must be provided before starting a subscription.")
	}

	priceId := s.app.TeamSubscriptionTierPriceId(input.Tier)

	subscription, err := s.app.stripe.Subscriptions.New(&stripe.SubscriptionParams{
		Params: stripe.Params{
			IdempotencyKey: stripe.String("create_subscription:" + teamId.String()),
		},
		Customer: stripe.String(team.StripeCustomerId),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(0),
			},
		},
		Metadata: map[string]string{
			"revision": model.NewId("ssrev").String(),
		},
		OffSession:      stripe.Bool(true),
		PaymentBehavior: stripe.String("error_if_incomplete"),
	})
	if err != nil {
		if IsStripeBadRequestError(err) {
			return nil, NewUserError("Subscription rejected by Stripe. Please double check your billing information.")
		}
		return nil, s.SanitizedError(fmt.Errorf("unable to create stripe subscription: %w", err))
	}

	if err := s.app.expectTeamEntitlementChange(ctx, RefreshTeamEntitlementsInput{
		TeamId:           team.Id,
		StripeCustomerId: team.StripeCustomerId,
	}); err != nil {
		return nil, s.SanitizedError(err)
	}

	ret, err := s.app.teamSubscriptionFromStripe(subscription)
	return ret, s.SanitizedError(err)
}

type updateStripeSubscriptionInput struct {
	Tier     *TeamSubscriptionTier
	Accounts *int
}

func (a *App) updateStripeSubscriptionByCustomerId(ctx context.Context, customerId string, input updateStripeSubscriptionInput) (*stripe.Subscription, error) {
	for attempt := 0; attempt < 3; attempt++ {
		existing, err := a.getSingleStripeSubscription(customerId)
		if err != nil {
			return nil, err
		} else if existing == nil {
			return nil, nil
		}

		var item *stripe.SubscriptionItem
		for _, i := range existing.Items.Data {
			if i.Price != nil && i.Price.Metadata["use_account_quantity"] == "true" {
				item = i
				break
			}
		}
		if item == nil {
			if input.Tier != nil {
				return nil, NewUserError("You currently have a custom subscription. Please contact support if you'd like to modify it.")
			} else {
				// Don't touch the subscription the quantity shouldn't be based on accounts.
				return existing, nil
			}
		}

		itemUpdate := &stripe.SubscriptionItemsParams{
			ID: stripe.String(item.ID),
		}
		hasChange := false

		if input.Tier != nil {
			newPriceId := a.TeamSubscriptionTierPriceId(*input.Tier)
			if item.Price == nil || item.Price.ID != newPriceId {
				itemUpdate.Price = stripe.String(newPriceId)
				itemUpdate.Quantity = stripe.Int64(item.Quantity)
				hasChange = true
			}
		}

		if input.Accounts != nil && item.Quantity != int64(*input.Accounts) {
			itemUpdate.Quantity = stripe.Int64(int64(*input.Accounts))
			hasChange = true
		}

		if !hasChange {
			return existing, nil
		}

		subscription, err := a.stripe.Subscriptions.Update(existing.ID, &stripe.SubscriptionParams{
			Params: stripe.Params{
				IdempotencyKey: stripe.String("update_subscription:" + customerId + ":" + existing.Metadata["revision"]),
			},
			Items: []*stripe.SubscriptionItemsParams{itemUpdate},
			Metadata: map[string]string{
				"revision": model.NewId("ssrev").String(),
			},
			ProrationBehavior: stripe.String("always_invoice"),
		})
		if err != nil {
			if IsStripeBadRequestError(err) {
				return nil, NewUserError("Subscription rejected by Stripe. Please double check your billing information.")
			}
			return nil, fmt.Errorf("unable to create stripe subscription: %w", err)
		}
		if subscription.LastResponse != nil && subscription.LastResponse.Header.Get("Idempotent-Replayed") == "true" {
			// Our update did nothing. Re-fetch the existing subscription and try again.
			continue
		}
		return subscription, nil
	}

	return nil, fmt.Errorf("unable to update subscription after 3 attempts")
}

type UpdateTeamSubscriptionInput struct {
	Tier *TeamSubscriptionTier
}

func (s *Session) UpdateTeamSubscriptionById(ctx context.Context, teamId model.Id, input UpdateTeamSubscriptionInput) (*model.TeamSubscription, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyStrongInRegion)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if team == nil {
		return nil, NotFoundError("Team not found.")
	} else if team.StripeCustomerId == "" {
		return nil, nil
	} else if err := s.ValidateTierForTeam(ctx, *input.Tier, team); err != nil {
		return nil, err
	}

	if subscription, err := s.app.updateStripeSubscriptionByCustomerId(ctx, team.StripeCustomerId, updateStripeSubscriptionInput{
		Tier: input.Tier,
	}); err != nil {
		return nil, s.SanitizedError(err)
	} else if subscription == nil {
		return nil, NotFoundError("Subscription not found.")
	} else if err := s.app.expectTeamEntitlementChange(ctx, RefreshTeamEntitlementsInput{
		TeamId:           team.Id,
		StripeCustomerId: team.StripeCustomerId,
	}); err != nil {
		return nil, s.SanitizedError(err)
	} else {
		ret, err := s.app.teamSubscriptionFromStripe(subscription)
		return ret, s.SanitizedError(err)
	}
}

func (s *Session) GetTeamSubscriptionById(ctx context.Context, teamId model.Id) (*model.TeamSubscription, UserFacingError) {
	if err := s.RequireTeamAdministrator(ctx, teamId); err != nil {
		return nil, err
	}

	team, err := s.app.store.GetTeamById(ctx, teamId, store.ConsistencyStrongInRegion)
	if err != nil {
		return nil, s.SanitizedError(err)
	} else if team == nil {
		return nil, NotFoundError("Team not found.")
	} else if team.StripeCustomerId == "" {
		return nil, nil
	}

	if subscription, err := s.app.getSingleStripeSubscription(team.StripeCustomerId); err != nil {
		return nil, s.SanitizedError(fmt.Errorf("unable to get stripe subscription: %w", err))
	} else {
		ret, err := s.app.teamSubscriptionFromStripe(subscription)
		return ret, s.SanitizedError(err)
	}
}

func (a *App) QueueTeamStripeSubscriptionUpdates(ctx context.Context) error {
	teams, err := a.store.GetTeams(ctx)
	if err != nil {
		return fmt.Errorf("unable to get teams: %w", err)
	}

	var msgs []OutgoingQueueMessage
	for _, team := range teams {
		if team.StripeCustomerId == "" {
			continue
		}
		msgs = append(msgs, OutgoingQueueMessage{
			Delay: time.Duration(rand.Intn(int(MaxQueueDelay/time.Second))) * time.Second,
			Message: QueueMessage{
				UpdateTeamStripeSubscription: &UpdateTeamStripeSubscriptionInput{
					TeamId:           team.Id,
					StripeCustomerId: team.StripeCustomerId,
				},
			},
		})
	}

	if err := a.QueueMessages(ctx, map[string][]OutgoingQueueMessage{
		a.awsRegion: msgs,
	}); err != nil {
		return fmt.Errorf("unable to queue stripe subscription updates: %w", err)
	}

	return nil
}

type UpdateTeamStripeSubscriptionInput struct {
	TeamId           model.Id
	StripeCustomerId string
}

func (a *App) UpdateTeamStripeSubscription(ctx context.Context, input UpdateTeamStripeSubscriptionInput) error {
	count, err := a.store.GetTeamBillableAccountCountByTeamId(ctx, input.TeamId)
	if err != nil {
		return fmt.Errorf("unable to get team billing account count: %w", err)
	}
	if _, err := a.updateStripeSubscriptionByCustomerId(ctx, input.StripeCustomerId, updateStripeSubscriptionInput{
		Accounts: &count,
	}); err != nil {
		return fmt.Errorf("unable to update stripe subscription: %w", err)
	}
	return nil
}

func (a *App) QueueTeamEntitlementRefreshes(ctx context.Context) error {
	teams, err := a.store.GetTeams(ctx)
	if err != nil {
		return fmt.Errorf("unable to get teams: %w", err)
	}

	var msgs []OutgoingQueueMessage
	for _, team := range teams {
		if team.StripeCustomerId == "" {
			continue
		}
		msgs = append(msgs, OutgoingQueueMessage{
			Delay: time.Duration(rand.Intn(int(MaxQueueDelay/time.Second))) * time.Second,
			Message: QueueMessage{
				RefreshTeamEntitlements: &RefreshTeamEntitlementsInput{
					TeamId:           team.Id,
					StripeCustomerId: team.StripeCustomerId,
				},
			},
		})
	}

	if err := a.QueueMessages(ctx, map[string][]OutgoingQueueMessage{
		a.awsRegion: msgs,
	}); err != nil {
		return fmt.Errorf("unable to queue team entitlement refreshes: %w", err)
	}

	return nil
}

// We expect entitlements to change soon, so do one refresh now and queue a few more.
func (a *App) expectTeamEntitlementChange(ctx context.Context, input RefreshTeamEntitlementsInput) error {
	if err := a.RefreshTeamEntitlements(ctx, input); err != nil {
		return fmt.Errorf("unable to refresh team entitlements: %w", err)
	}

	var msgs []OutgoingQueueMessage
	for i := 1; i <= 3; i++ {
		msgs = append(msgs, OutgoingQueueMessage{
			Delay: time.Duration(i*i) * time.Second,
			Message: QueueMessage{
				RefreshTeamEntitlements: &input,
			},
		})
	}

	if err := a.QueueMessages(ctx, map[string][]OutgoingQueueMessage{
		a.awsRegion: msgs,
	}); err != nil {
		return fmt.Errorf("unable to queue team entitlement refreshes: %w", err)
	}

	return nil
}

type RefreshTeamEntitlementsInput struct {
	// If not given, we'll do some extra fetches to figure out the team id.
	TeamId model.Id

	StripeCustomerId string
}

// Pulls the latest entitlements from Stripe and updates the team.
func (a *App) RefreshTeamEntitlements(ctx context.Context, input RefreshTeamEntitlementsInput) error {
	if input.TeamId == "" {
		customer, err := a.stripe.Customers.Get(input.StripeCustomerId, nil)
		if err != nil {
			return fmt.Errorf("unable to get stripe customer: %w", err)
		}
		input.TeamId = model.Id(customer.Metadata["team_id"])
		if input.TeamId == "" {
			return nil
		}

		team, err := a.store.GetTeamById(ctx, input.TeamId, store.ConsistencyEventual)
		if err != nil {
			return fmt.Errorf("unable to get team: %w", err)
		} else if team == nil || team.StripeCustomerId != input.StripeCustomerId {
			return nil
		}
	}

	var entitlements model.TeamEntitlements
	iter := a.stripe.EntitlementsActiveEntitlements.List(&stripe.EntitlementsActiveEntitlementListParams{
		Customer: stripe.String(input.StripeCustomerId),
	})
	for iter.Next() {
		entitlement := iter.EntitlementsActiveEntitlement()
		switch entitlement.LookupKey {
		case "individual-features":
			entitlements.IndividualFeatures = true
		case "team-features":
			entitlements.TeamFeatures = true
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("unable to get stripe entitlements: %w", err)
	}

	storePatch := &store.TeamPatch{
		Entitlements: &entitlements,
	}

	if _, err := a.store.PatchTeamById(ctx, input.TeamId, storePatch); err != nil {
		return fmt.Errorf("unable to patch team entitlements: %w", err)
	}

	return nil
}
