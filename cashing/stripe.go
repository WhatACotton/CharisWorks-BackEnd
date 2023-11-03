package cashing

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/WhatACotton/go-backend-test/validation"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/account"
	"github.com/stripe/stripe-go/v75/accountlink"
	"github.com/stripe/stripe-go/v75/checkout/session"
	"github.com/stripe/stripe-go/v75/refund"
	"github.com/stripe/stripe-go/v75/transfer"
	"github.com/stripe/stripe-go/v75/webhook"
)

type StripeInfo struct {
	URL         string
	AmountTotal int64
	ID          string
}

// (第一引数).urlが飛んでほしいリンク先
func Purchase(amount int) (StripeInfo, error) {
	amount_str := strconv.Itoa(amount)
	amount_int, _ := strconv.ParseInt(amount_str, 0, 0)
	log.Printf("amount: %v\n", amount_int)
	return createCheckoutSession(amount_int)
}
func createCheckoutSession(amount int64) (StripeInfo, error) {
	stripe.Key = "sk_test_51Nj1urA3bJzqElthx8UK5v9CdaucJOZj3FwkOHZ8KjDt25IAvplosSab4uybQOyE2Ne6xxxI4Rnh8pWEbYUwPoPG00wvseAHzl"

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("jpy"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("購入金額"),
					},
					UnitAmount: stripe.Int64(amount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("http://" + os.Getenv("IP_ADDRESS") + ":80/mypage"),
		CancelURL:  stripe.String("http://" + os.Getenv("IP_ADDRESS") + ":80/signin"),
	}
	s, _ := session.New(params)
	log.Print(s.ID)

	stripe_info := StripeInfo{
		URL:         s.URL,
		AmountTotal: s.AmountTotal,
		ID:          s.ID,
	}
	return stripe_info, nil
}
func PaymentComplete(w http.ResponseWriter, req *http.Request) (string, error) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return "", err
	}

	// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key
	// You can find your endpoint's secret in your webhook settings
	endpointSecret := os.Getenv("STRIPE_KEY")
	event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return "", err
	}

	// Handle the checkout.session.completed event
	if event.Type == "checkout.session.completed" {
		var sessions stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &sessions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return "", err
		}

		params := &stripe.CheckoutSessionParams{}
		params.AddExpand("line_items")
		log.Print(sessions.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return "", err
		}
		return sessions.ID, nil
	}

	w.WriteHeader(http.StatusOK)
	return "", nil
}
func CreateStripeAccount(email string) (stripeID string, URL string) {
	stripe.Key = "sk_test_51Nj1urA3bJzqElthx8UK5v9CdaucJOZj3FwkOHZ8KjDt25IAvplosSab4uybQOyE2Ne6xxxI4Rnh8pWEbYUwPoPG00wvseAHzl"

	params := &stripe.AccountParams{
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(true),
			},
		},
		Country: stripe.String("JP"),
		Email:   stripe.String(email),
		Type:    stripe.String("express"),
	}

	a, _ := account.New(params)
	log.Print("CreatedStripeAccount. ID : ", a.ID)
	URL = createAccountLink(a.ID)
	return a.ID, URL
}
func createAccountLink(ID string) string {
	params := &stripe.AccountLinkParams{
		Account:    stripe.String(ID),
		RefreshURL: stripe.String("http://" + os.Getenv("IP_ADDRESS") + "/mypage"),
		ReturnURL:  stripe.String("http://" + os.Getenv("IP_ADDRESS") + "/mypage"),
		Type:       stripe.String("account_onboarding"),
		Collect:    stripe.String("eventually_due"),
	}
	result, err := accountlink.New(params)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(result.URL)
	return result.URL
}
func Transfer(amount float64, stripeID string, ItemName string) {
	stripe.Key = "sk_test_51Nj1urA3bJzqElthx8UK5v9CdaucJOZj3FwkOHZ8KjDt25IAvplosSab4uybQOyE2Ne6xxxI4Rnh8pWEbYUwPoPG00wvseAHzl"
	log.Print("Transfering... \n amount: ", amount, "\n stripeID: ", stripeID, "\n ItemName: ", ItemName)
	params := &stripe.TransferParams{
		Amount:      stripe.Int64(int64(amount)),
		Currency:    stripe.String(string(stripe.CurrencyJPY)),
		Destination: stripe.String(stripeID),
		Description: stripe.String(ItemName),
	}
	tr, _ := transfer.New(params)
	log.Print(tr.ID)
	validation.TransferLogging(tr.ID + stripeID + strconv.FormatFloat(amount, 'f', 2, 64) + ItemName)

}
func Refund(ID string) {
	stripe.Key = "sk_test_51Nj1urA3bJzqElthx8UK5v9CdaucJOZj3FwkOHZ8KjDt25IAvplosSab4uybQOyE2Ne6xxxI4Rnh8pWEbYUwPoPG00wvseAHzl"

	params := &stripe.RefundParams{PaymentIntent: stripe.String("ID")}
	result, err := refund.New(params)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(result)
}
