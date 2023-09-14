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
	Cart, UID := GetDatafromSessionKey(c)
	if UID != "" {
		Customer := new(database.Customer)
		Customer.GetCustomer(UID)
		log.Print("Customer:", Customer.Name)
		if Customer.Register && Customer.EmailVerified {
			CartContents, err := database.GetCartContents(Cart.CartID)
			if err != nil {
				log.Fatal(err)
			}
			if inspectCart(CartContents) {
				TotalPrice := totalPrice(CartContents)
				stripeInfo, err := cashing.Purchase(TotalPrice)
				if err != nil {
					log.Fatal(err)
				}
				TransactionID := validation.GetUUID()
				database.PostTransaction(Cart, *Customer, stripeInfo, TransactionID, CartContents)

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
	CustomerSessionKey := new(string)
	*CustomerSessionKey = validation.GetCustomerSessionKey(c)
	TransactionContentsList := new([]database.TransactionContents)
	UID, err := database.GetUID(*CustomerSessionKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("UID:", UID)
	Transactions, err := database.GetTransactions(UID)
	if err != nil {
		log.Fatal(err)
	}
	for _, Transaction := range Transactions {
		TransactionContents, err := Transaction.GetTransactionContents()
		if err != nil {
			log.Fatal(err)
		}
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
	database.SetTransactionStatus("決済完了", ID)
	Transaction := new(database.Transaction)
	Transaction.TransactionID, err = database.GetTransactionID(ID)
	if err != nil {
		log.Fatal(err)
	}
	TransactionContents, err := Transaction.GetTransactionContents()
	if err != nil {
		log.Fatal(err)
	}
	for _, TransactionContent := range TransactionContents {
		log.Print("TransactionContent: ", TransactionContent)
		database.Purchased(TransactionContent)
		Itemdetails, err := database.GetItemDetails(TransactionContent.InfoID)
		if err != nil {
			panic(err)
		}
		amount := float64(Itemdetails.Price) * float64(TransactionContent.Quantity) * 0.97 * 0.964
		cashing.Transfer(amount, Itemdetails.Madeby, Itemdetails.ItemName)
	}
	UID, err := database.GetUIDfromStripeID(ID)
	if err != nil {
		panic(err)
	}
	CartID, err := database.GetCartID(UID)
	if err != nil {
		panic(err)
	}
	database.DeleteCartContentforTransaction(CartID)
	err = database.DeleteCart(CartID)
	if err != nil {
		panic(err)
	}
	err = database.DeleteCartContentforTransaction(CartID)
	if err != nil {
		panic(err)
	}
	database.SetCartID(UID, validation.GetUUID())
	return nil
}

// 返金処理
func Refund(c *gin.Context) {
	ID := c.Query("ID")
	refund(ID)
}
func refund(ID string) {
	status := database.GetTransactionStatus(ID)
	if status == "返金待ち" {
		cashing.Refund(ID)
	}
}
