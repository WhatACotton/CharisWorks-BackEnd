package handler

import (
	"log"
	"unify/cashing"
	"unify/internal/database"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func Webhook(c *gin.Context) {
	ID, err := cashing.Payment_complete(c.Writer, c.Request)
	if err != nil {
		c.JSON(400, gin.H{"message": "error"})
	}

	Complete_Payment(ID)

}
func Complete_Payment(ID string) {
	log.Print("ID: ", ID)
	database.Change_Transaction_Status("決済完了", ID)
	UID, err := database.Get_UID_from_Stripe_ID(ID)
	if err != nil {
		panic(err)
	}
	Cart_ID, err := database.Get_Cart_ID(UID)
	if err != nil {
		panic(err)
	}
	err = database.Delete_Cart_List(Cart_ID)
	if err != nil {
		panic(err)
	}
	err = database.Delete_Cart_for_Transaction(Cart_ID)
	if err != nil {
		panic(err)
	}
	database.Set_Cart_ID(UID, validation.GetUUID())
}
