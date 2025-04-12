package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v81"
)

type Currency = stripe.Currency

type CurrencyAmount struct {
	Currency Currency
	Amount   int64
}

func (ca CurrencyAmount) String() string {
	if ca.Amount < 0 {
		return "-" + CurrencyAmount{
			Currency: ca.Currency,
			Amount:   -ca.Amount,
		}.String()
	}

	switch ca.Currency {
	case stripe.CurrencyUSD:
		return fmt.Sprintf("$%01d.%02d", ca.Amount/100, ca.Amount%100)
	default:
		return strconv.FormatInt(ca.Amount, 10)
	}
}

type TeamBillingProfile struct {
	Name    string
	Address TeamBillingAddress
	Balance *CurrencyAmount
}

type TeamBillingAddress struct {
	City       *string
	Country    string
	Line1      *string
	Line2      *string
	PostalCode string
	State      *string
}

type TeamPaymentMethodCard struct {
	Last4Digits     string
	ExpirationMonth int
	ExpirationYear  int
}

type TeamPaymentMethodUSBankAccount struct {
	Last4Digits string
}

type TeamPaymentMethod struct {
	Card          *TeamPaymentMethodCard
	USBankAccount *TeamPaymentMethodUSBankAccount
}

type TeamSubscription struct {
	Name     string
	Accounts int
	Price    *TeamSubscriptionPrice
}

type TeamSubscriptionPrice struct {
	AccountMonth *CurrencyAmount
}

type TeamBillableAccount struct {
	Id             string
	TeamId         Id
	ExpirationTime time.Time
}
