package auth

import (
	"net/http"
	"unify/internal/database"
	"unify/internal/models"

	"github.com/gin-gonic/gin"
)

func PostTransaction(c *gin.Context) {
	var newCustomer models.TransactionRequestPayload
	if err := c.BindJSON(&newCustomer); err != nil {
		return
	}
	res := database.PostTransaction(newCustomer)
	c.IndentedJSON(http.StatusOK, res)
}

func PatchTransaction(c *gin.Context) {
	h := new(database.PatchRequestPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	h.Patch("transaction", "transactionid")
	c.JSON(http.StatusOK, database.GetItemList())
}

func GetTransaction(c *gin.Context) {
	transactionId := c.Query("transactionId")
	var response = database.GetTransaction(transactionId)
	if response.UID == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "customer not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func DeleteTransaction(c *gin.Context) {
	transactionid := c.Query("transactionid")
	database.Delete("transaction", "ID", transactionid)
	c.JSON(http.StatusOK, gin.H{"message": "transaction was successfully deleted"})
}
