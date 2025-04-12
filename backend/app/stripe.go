package app

import (
	"context"

	"github.com/stripe/stripe-go/v81"
)

func (a *App) HandleStripeEvent(ctx context.Context, event *stripe.Event) error {
	switch event.Type {
	case stripe.EventTypeEntitlementsActiveEntitlementSummaryUpdated:
		/*
			{
				"object": {
					"object": "entitlements.active_entitlement_summary",
					"customer": "cus_S32XBi6oZUNKYl",
					"entitlements": {
						"object": "list",
						"data": [
						{
							"id": "ent_test_61SLGNaa9kBQSLR7l412ejpbHZUu9FB2",
							"object": "entitlements.active_entitlement",
							"feature": "feat_test_61SJKP40IdgRzvC5b412ejpbHZUu9SFM",
							"livemode": false,
							"lookup_key": "team-features"
						},
						{
							"id": "ent_test_61SLGNapRzSIJnhMs412ejpbHZUu9QXg",
							"object": "entitlements.active_entitlement",
							"feature": "feat_test_61SJKOxUKvmitloN0412ejpbHZUu9Pdg",
							"livemode": false,
							"lookup_key": "individual-features"
						}
						],
						"has_more": false,
						"url": "/v1/customer/cus_S32XBi6oZUNKYl/entitlements"
					},
					"livemode": false
				},
				"previous_attributes": {
					"entitlements": {
						"data": [
						{
							"id": "ent_test_61SLBtN3ZTkG92lqj412ejpbHZUu9OfQ",
							"object": "entitlements.active_entitlement",
							"feature": "feat_test_61SJKOxUKvmitloN0412ejpbHZUu9Pdg",
							"livemode": false,
							"lookup_key": "individual-features"
						}
						]
					}
				}
			}
		*/
		return a.RefreshTeamEntitlements(ctx, RefreshTeamEntitlementsInput{
			StripeCustomerId: event.Data.Object["customer"].(string),
		})
	}
	return nil
}
