package handler

import (
	"log"
	"net/http"

	"github.com/WhatACotton/go-backend-test/cashing"
	"github.com/WhatACotton/go-backend-test/internal/database"
	"github.com/WhatACotton/go-backend-test/validation"
	"github.com/gin-gonic/gin"
)

// 商品の購入リクエストを作成。
func BuyItem(c *gin.Context) {
	log.Print("Creating PaymentIntent...")
	log.Print(c)
	UserID := GetDatafromSessionKey(c)
	CartContents := new(CartRequestPayloads)
	err := c.BindJSON(&CartContents)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("CartContents:", CartContents)
	if UserID != "" {
		Customer := new(database.Customer)
		Customer.CustomerGet(UserID)
		log.Print("Customer:", Customer.CustomerName)
		if Customer.IsRegistered && Customer.IsEmailVerified {
			TotalPrice := CartContents.inspectCart()
			if TotalPrice != 0 {
				stripeInfo, err := cashing.Purchase(TotalPrice)
				if err != nil {
					log.Fatal(err)
				}
				TransactionID := validation.GetUUID()

				database.TransactionPost(*Customer, stripeInfo, TransactionID, ConstructCart(*CartContents))

				c.JSON(http.StatusOK, gin.H{"message": "購入リンクが発行されました。", "url": stripeInfo.URL})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "カートの中身に購入不可能な商品が含まれています。"})
			}

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "本登録・またはEmail認証が完了していません。", "errcode": "401"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインが完了していません。", "errcode": "401"})
	}

}
func ConstructCart(Cart CartRequestPayloads) (CartContents database.CartContents) {
	for _, CartContent := range Cart {
		Item := new(database.Item)
		Item.ItemGet(CartContent.ItemID)
		CartContent := new(database.CartContent)
		CartContent.Price = Item.Price
		CartContent.Status = Item.Status
		CartContents = append(CartContents, *CartContent)
	}
	return CartContents
}

// カートの中身を確認し、購入可能かどうかを判定する。
func (carts *CartRequestPayloads) inspectCart() int {
	price := 0
	if len(*carts) == 0 {
		return 0
	}
	for _, Cart := range *carts {
		Item := new(database.Item)
		Item.ItemGet(Cart.ItemID)
		if Item.Status != "Available" {
			return 0
		}
		flag, stock := database.IsItemExist(Cart.ItemID)
		if !flag {
			return 0
		}
		if Cart.Quantity <= 0 {
			return 0
		}
		if stock < Cart.Quantity {
			return 0
		}
	}
	if hasDuplicates(*carts) {
		for _, Cart := range *carts {
			Item := new(database.Item)
			Item.ItemGet(Cart.ItemID)
			price += Item.Price * Cart.Quantity
		}
		return price
	}
	return 0
}

func hasDuplicates(slice CartRequestPayloads) bool {
	encountered := make(map[string]bool)
	for _, item := range slice {
		if encountered[item.ItemID] {
			return false
		}
		encountered[item.ItemID] = true
	}
	return true
}

// 購入履歴を取得する。
func GetTransaction(c *gin.Context) {
	TransactionContentsList := new([]database.TransactionDetails)
	UserID := GetDatafromSessionKey(c)

	log.Print("UserID:", UserID)
	Transactions := database.TransactionGet(UserID)
	for _, Transaction := range Transactions {
		TransactionContents := Transaction.TransactionDetailsGet()
		*TransactionContentsList = append(*TransactionContentsList, TransactionContents)
	}
	c.JSON(http.StatusOK, gin.H{"TransactionLists": Transactions, "Transactions": *TransactionContentsList})
}

// stripeからのwebhookを受け取る
func Webhook(c *gin.Context) {
	ID, status, err := cashing.PaymentComplete(c.Writer, c.Request)
	if err != nil {
		c.JSON(400, gin.H{"message": "error"})
	}
	if ID != "" {
		completePayment(ID, status)
	}
}

// 購入を完了させる
func completePayment(ID string, status string) (err error) {
	database.TransactionSetStatus(status, ID)
	Transaction := new(database.Transaction)
	Transaction.TransactionID = database.TransactionGetID(ID)
	TransactionDetails := Transaction.TransactionDetailsGet()
	log.Print(len(TransactionDetails))
	for _, TransactionDetail := range TransactionDetails {
		log.Print("TransactionContent: ", TransactionDetail)
		database.Purchased(TransactionDetail)
		Item := new(database.Item)
		Item.ItemGet(TransactionDetail.ItemID)
		amount := float64(Item.Price) * float64(TransactionDetail.Quantity) * 0.97 * 0.964
		cashing.Transfer(amount, Item.StripeAccountID, Item.Name)
	}
	UserID := database.TransactionGetUserIDfromStripeID(ID)
	database.ClearCart(UserID)
	return nil
}

// 返金処理
func Refund(c *gin.Context) {
	ID := c.Query("ID")
	status := database.TransactionGetStatus(ID)
	if status == "返金待ち" {
		cashing.Refund(ID)
	}
}
