package stripe

import (
	"log"
	"strconv"

	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/checkout/session"
)

// (第一引数).urlが飛んでほしいリンク先
func Purchase(amount string) (stripe.CheckoutSession, error) {
	amount_int, _ := strconv.ParseInt(amount, 0, 0)
	log.Printf("amount: %v\n", amount_int)
	return createCheckoutSession(amount_int)
}

func createCheckoutSession(amount int64) (stripe.CheckoutSession, error) {
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

		SuccessURL: stripe.String("https://localhost:3000/Success"),
		CancelURL:  stripe.String("http://localhost:3000/Cancel"),
	}

	s, _ := session.New(params)
	return *s, nil
}
