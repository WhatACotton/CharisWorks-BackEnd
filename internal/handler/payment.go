package handler

import (
	"log"
	"unify/cashing"
	"unify/internal/database"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func Webhook(c *gin.Context) {
	ID, err := cashing.PaymentComplete(c.Writer, c.Request)
	if err != nil {
		c.JSON(400, gin.H{"message": "error"})
	}
	CompletePayment(ID)

}
func CompletePayment(ID string) {
	database.ChangeTransactionStatus("決済完了", ID)
	UID, err := database.GetUIDfromStripeID(ID)
	if err != nil {
		panic(err)
	}
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
	}
	CartID, err := database.GetCartID(UID)
	if err != nil {
		panic(err)
	}
	CartContents, err := database.GetCartContents(CartID)
	if err != nil {
		panic(err)
	}
	for _, CartContent := range CartContents {
		database.DeleteCartContent(CartID, CartContent.ItemID)
	}
	err = database.DeleteCart(CartID)
	if err != nil {
		panic(err)
	}
	err = database.DeleteCartforTransaction(CartID)
	if err != nil {
		panic(err)
	}
	database.SetCartID(UID, validation.GetUUID())
}
