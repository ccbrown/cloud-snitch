package apptest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/form"
)

type MockStripeBackend struct {
	stripe.Backend

	customers     map[string]*stripe.Customer
	subscriptions map[string]*stripe.Subscription
	m             sync.Mutex
}

var DummyStripeCard = &stripe.PaymentMethod{
	ID:   "pm_card_visa",
	Type: stripe.PaymentMethodTypeCard,
	Card: &stripe.PaymentMethodCard{
		Last4:    "1234",
		ExpMonth: 12,
		ExpYear:  3000,
	},
}

var DummyStripePriceIndividualSubscription = &stripe.Price{
	ID:         "price_individual_subscription",
	UnitAmount: 99,
	Currency:   stripe.CurrencyUSD,
	Recurring: &stripe.PriceRecurring{
		Interval: stripe.PriceRecurringIntervalMonth,
	},
	Type: stripe.PriceTypeRecurring,
	Product: &stripe.Product{
		ID:   "prod_individual_subscription",
		Name: "Individual Subscription",
	},
	Metadata: map[string]string{
		"use_account_quantity": "true",
	},
}

var DummyStripePriceTeamSubscription = &stripe.Price{
	ID:         "price_team_subscription",
	UnitAmount: 999,
	Currency:   stripe.CurrencyUSD,
	Recurring: &stripe.PriceRecurring{
		Interval: stripe.PriceRecurringIntervalMonth,
	},
	Type: stripe.PriceTypeRecurring,
	Product: &stripe.Product{
		ID:   "prod_team_subscription",
		Name: "Team Subscription",
	},
	Metadata: map[string]string{
		"use_account_quantity": "true",
	},
}

var DummyStripePriceCustomSubscription = &stripe.Price{
	ID:         "price_custom_subscription",
	UnitAmount: 100000,
	Currency:   stripe.CurrencyUSD,
	Product: &stripe.Product{
		ID:   "prod_custom_subscription",
		Name: "My Custom Subscription",
	},
}

func (b *MockStripeBackend) Call(method, path, key string, paramsContainer stripe.ParamsContainer, v stripe.LastResponseSetter) error {
	formValues := &form.Values{}
	form.AppendTo(formValues, paramsContainer)
	var params *stripe.Params
	if !reflect.ValueOf(paramsContainer).IsNil() {
		params = paramsContainer.GetParams()
	}
	return b.CallRaw(method, path, key, formValues, params, v)
}

func metadataFromValues(values *form.Values, prefix string) map[string]string {
	ret := make(map[string]string)
	for key, value := range values.ToValues() {
		if strings.HasPrefix(key, prefix) {
			key := strings.TrimSuffix(strings.TrimPrefix(key, prefix+"["), "]")
			ret[key] = value[0]
		}
	}
	return ret
}

func (b *MockStripeBackend) CallRaw(method, path, key string, body *form.Values, params *stripe.Params, v stripe.LastResponseSetter) error {
	b.m.Lock()
	defer b.m.Unlock()

	stringParam := func(name string) string {
		if values := body.Get(name); len(values) > 0 {
			return values[0]
		}
		return ""
	}

	respondOkay := func(body any) {
		buf, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(buf, v); err != nil {
			panic(err)
		}
		v.SetLastResponse(&stripe.APIResponse{
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			RawJSON:    buf,
			StatusCode: http.StatusOK,
		})
	}

	if path == "/v1/customers" && method == "POST" {
		c := &stripe.Customer{
			ID: fmt.Sprintf("cus_%d", len(b.customers)+1),
			Address: &stripe.Address{
				Country:    stringParam("address[country]"),
				PostalCode: stringParam("address[postal_code]"),
				Line1:      stringParam("address[line1]"),
				Line2:      stringParam("address[line2]"),
				City:       stringParam("address[city]"),
				State:      stringParam("address[state]"),
			},
		}
		if id := stringParam("invoice_settings[default_payment_method]"); id != "" {
			if id == DummyStripeCard.ID {
				c.InvoiceSettings = &stripe.CustomerInvoiceSettings{
					DefaultPaymentMethod: DummyStripeCard,
				}
			}
		}
		if b.customers == nil {
			b.customers = make(map[string]*stripe.Customer)
		}
		b.customers[c.ID] = c
		respondOkay(c)
	} else if strings.HasPrefix(path, "/v1/customers/") && method == "GET" {
		id := strings.TrimPrefix(path, "/v1/customers/")
		respondOkay(b.customers[id])
	} else if strings.HasPrefix(path, "/v1/customers/") && method == "POST" {
		id := strings.TrimPrefix(path, "/v1/customers/")
		c := b.customers[id]

		if newAddr := (stripe.Address{
			Country:    stringParam("address[country]"),
			PostalCode: stringParam("address[postal_code]"),
			Line1:      stringParam("address[line1]"),
			Line2:      stringParam("address[line2]"),
			City:       stringParam("address[city]"),
			State:      stringParam("address[state]"),
		}); newAddr != (stripe.Address{}) {
			c.Address = &newAddr
		}

		if id := stringParam("invoice_settings[default_payment_method]"); id != "" {
			if id == DummyStripeCard.ID {
				c.InvoiceSettings = &stripe.CustomerInvoiceSettings{
					DefaultPaymentMethod: DummyStripeCard,
				}
			}
		}

		respondOkay(c)
	} else if strings.HasPrefix(path, "/v1/payment_methods/") && method == "GET" {
		id := strings.TrimPrefix(path, "/v1/payment_methods/")
		if id == DummyStripeCard.ID {
			respondOkay(DummyStripeCard)
		}
	} else if strings.HasPrefix(path, "/v1/products/") && method == "GET" {
		id := strings.TrimPrefix(path, "/v1/products/")
		switch id {
		case DummyStripePriceIndividualSubscription.Product.ID:
			respondOkay(DummyStripePriceIndividualSubscription.Product)
		case DummyStripePriceTeamSubscription.Product.ID:
			respondOkay(DummyStripePriceTeamSubscription.Product)
		case DummyStripePriceCustomSubscription.Product.ID:
			respondOkay(DummyStripePriceCustomSubscription.Product)
		default:
			panic("unexpected product ID: " + id)
		}
	} else if path == "/v1/setup_intents" && method == "POST" {
		respondOkay(struct{}{})
	} else if path == "/v1/entitlements/active_entitlements" && method == "GET" {
		var entitlements []*stripe.EntitlementsActiveEntitlement
		customer := stringParam("customer")
		for _, s := range b.subscriptions {
			if s.Customer.ID != customer {
				continue
			}
			switch s.Items.Data[0].Price.ID {
			case DummyStripePriceIndividualSubscription.ID:
				entitlements = append(entitlements, &stripe.EntitlementsActiveEntitlement{
					ID:        "ent_individual_features",
					LookupKey: "individual-features",
				})
			case DummyStripePriceTeamSubscription.ID:
				entitlements = append(entitlements, &stripe.EntitlementsActiveEntitlement{
					ID:        "ent_team_features",
					LookupKey: "team-features",
				}, &stripe.EntitlementsActiveEntitlement{
					ID:        "ent_individual_features",
					LookupKey: "individual-features",
				})
			}
		}
		respondOkay(&stripe.EntitlementsActiveEntitlementList{
			Data: entitlements,
		})
	} else if path == "/v1/subscriptions" && method == "GET" {
		subscriptions := make([]*stripe.Subscription, 0, len(b.subscriptions))
		customer := stringParam("customer")
		for _, s := range b.subscriptions {
			if customer == "" || s.Customer.ID == customer {
				subscriptions = append(subscriptions, s)
			}
		}
		respondOkay(&stripe.SubscriptionList{
			Data: subscriptions,
		})
	} else if path == "/v1/subscriptions" && method == "POST" {
		s := &stripe.Subscription{
			ID: fmt.Sprintf("sub_%d", len(b.subscriptions)+1),
			Customer: &stripe.Customer{
				ID: stringParam("customer"),
			},
			Items:    &stripe.SubscriptionItemList{},
			Metadata: metadataFromValues(body, "metadata"),
		}
		for i := 0; stringParam(fmt.Sprintf("items[%d][price]", i)) != ""; i++ {
			priceId := stringParam(fmt.Sprintf("items[%d][price]", i))
			var price *stripe.Price
			switch priceId {
			case DummyStripePriceIndividualSubscription.ID:
				price = DummyStripePriceIndividualSubscription
			case DummyStripePriceTeamSubscription.ID:
				price = DummyStripePriceTeamSubscription
			case DummyStripePriceCustomSubscription.ID:
				price = DummyStripePriceCustomSubscription
			default:
				panic("unexpected price ID: " + priceId)
			}
			s.Items.Data = append(s.Items.Data, &stripe.SubscriptionItem{
				ID:    fmt.Sprintf("si_%d", len(s.Items.Data)+1),
				Price: price,
				Plan: &stripe.Plan{
					Product: price.Product,
				},
			})
		}
		if b.subscriptions == nil {
			b.subscriptions = make(map[string]*stripe.Subscription)
		}
		b.subscriptions[s.ID] = s
		respondOkay(s)
	} else if strings.HasPrefix(path, "/v1/subscriptions/") && method == "POST" {
		id := strings.TrimPrefix(path, "/v1/subscriptions/")
		s := b.subscriptions[id]

		for i := 0; stringParam(fmt.Sprintf("items[%d][id]", i)) != ""; i++ {
			itemId := stringParam(fmt.Sprintf("items[%d][id]", i))
			var item *stripe.SubscriptionItem
			for _, it := range s.Items.Data {
				if it.ID == itemId {
					item = it
					break
				}
			}

			if priceId := stringParam(fmt.Sprintf("items[%d][price]", i)); priceId != "" {
				var newPrice *stripe.Price
				switch priceId {
				case DummyStripePriceIndividualSubscription.ID:
					newPrice = DummyStripePriceIndividualSubscription
				case DummyStripePriceTeamSubscription.ID:
					newPrice = DummyStripePriceTeamSubscription
				case DummyStripePriceCustomSubscription.ID:
					newPrice = DummyStripePriceCustomSubscription
				default:
					panic("unexpected price ID: " + priceId)
				}
				item.Price = newPrice
				item.Plan = &stripe.Plan{
					Product: newPrice.Product,
				}
			}
		}

		respondOkay(s)
	} else {
		panic("unexpected call to Stripe API: " + method + " " + path)
	}

	return nil
}
