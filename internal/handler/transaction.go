package handler

import (
	"log"
	"net/http"
	"unify/cashing"
	"unify/internal/database"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

// 商品の購入リクエストを作成。
func BuyItem(c *gin.Context) {
	log.Print("Creating PaymentIntent...")
	Cart, UserID := GetDatafromSessionKey(c)
	if UserID != "" {
		Customer := new(database.Customer)
		Customer.GetCustomer(UserID)
		log.Print("Customer:", Customer.Name)
		if Customer.IsRegistered && Customer.IsEmailVerified {
			CartContents := database.GetCartContents(Cart.CartID)
			if inspectCart(CartContents) {
				TotalPrice := totalPrice(CartContents)
				stripeInfo, err := cashing.Purchase(TotalPrice)
				if err != nil {
					log.Fatal(err)
				}
				TransactionID := validation.GetUUID()
				database.TransactionPost(Cart, *Customer, stripeInfo, TransactionID, CartContents)

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

// カートの中身を確認し、購入可能かどうかを判定する。
func inspectCart(carts database.CartContents) bool {
	if len(carts) == 0 {
		return false
	}
	for _, Cart := range carts {
		if Cart.Status != "Available" {
			return false
		}
	}
	return true
}

// カートの合計金額を計算する。
func totalPrice(Carts database.CartContents) (TotalPrice int) {
	for _, Cart := range Carts {
		TotalPrice += Cart.Price * Cart.Quantity
	}
	return TotalPrice
}

// 購入履歴を取得する。
func GetTransaction(c *gin.Context) {
	TransactionContentsList := new([]database.TransactionDetails)
	_, UserID := GetDatafromSessionKey(c)

	log.Print("UserID:", UserID)
	Transactions := database.TransactionGet(UserID)
	for _, Transaction := range Transactions {
		TransactionContents := Transaction.TransactionGetContents()
		*TransactionContentsList = append(*TransactionContentsList, TransactionContents)
	}
	c.JSON(http.StatusOK, gin.H{"TransactionLists": Transactions, "Transactions": *TransactionContentsList})
}

func Webhook(c *gin.Context) {
	ID, err := cashing.PaymentComplete(c.Writer, c.Request)
	if err != nil {
		c.JSON(400, gin.H{"message": "error"})
	}
	if ID != "" {
		completePayment(ID)
	}
}

func completePayment(ID string) (err error) {
	database.TransactionSetStatus("決済完了", ID)
	Transaction := new(database.Transaction)
	Transaction.TransactionID = database.TransactionGetID(ID)
	if err != nil {
		log.Fatal(err)
	}
	TransactionContents := Transaction.TransactionGetContents()
	if err != nil {
		log.Fatal(err)
	}
	for _, TransactionContent := range TransactionContents {
		log.Print("TransactionContent: ", TransactionContent)
		database.Purchased(TransactionContent)
		Item := new(database.Item)
		Item.ItemGet(TransactionContent.ItemID)
		amount := float64(Item.Price) * float64(TransactionContent.Quantity) * 0.97 * 0.964
		cashing.Transfer(amount, Item.MadeBy, Item.ItemName)
	}
	UserID := database.TransactionGetUserIDfromStripeID(ID)
	if err != nil {
		panic(err)
	}
	CartID := database.GetCartID(UserID)
	if err != nil {
		panic(err)
	}
	database.DeleteCartContentforTransaction(CartID)
	database.DeleteCart(CartID)
	database.CustomerSetCartID(UserID, validation.GetUUID())
	return nil
}

// 返金処理
func Refund(c *gin.Context) {
	ID := c.Query("ID")
	refund(ID)
}
func refund(ID string) {
	status := database.TransactionGetStatus(ID)
	if status == "返金待ち" {
		cashing.Refund(ID)
	}
}
