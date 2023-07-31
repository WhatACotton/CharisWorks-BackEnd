package funcs

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

func GetTransaction(c *gin.Context) {
	transactionId := c.Query("transactionId")
	var response = database.GetTransaction(transactionId)
	if response.UID == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "customer not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}
