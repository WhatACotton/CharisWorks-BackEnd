package customer

import (
	"net/http"
	"unify/internal/handler"
	"unify/internal/models"

	"github.com/gin-gonic/gin"
)

func PostTransaction(c *gin.Context) {
	var newCustomer models.TransactionRequestPayload
	if err := c.BindJSON(&newCustomer); err != nil {
		return
	}
	res := handler.PostTransaction(newCustomer)
	c.IndentedJSON(http.StatusOK, res)
}
func UpdateTransaction(c *gin.Context) {
	var fixCustomer models.Transaction
	if err := c.BindJSON(&fixCustomer); err != nil {
		return
	}
	res := handler.UpdateTransaction(fixCustomer)
	c.IndentedJSON(http.StatusOK, res)

}

func GetTransaction(c *gin.Context) {
	uid := c.Query("uid")
	var response = handler.GetTransaction(uid)
	if response.UID == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "customer not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func DeleteTransaction(c *gin.Context) {
	uid := c.Query("uid")
	var response = handler.DeleteTransaction(uid)
	c.IndentedJSON(http.StatusOK, response)
}
