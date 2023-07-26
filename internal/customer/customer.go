package customer

import (
	"net/http"

	"unify/internal/handler"
	"unify/internal/models"

	"github.com/gin-gonic/gin"
)

var customers []models.Customer

func PostCustomer(c *gin.Context) {
	var CreatedDate = handler.GetDate()
	var newCustomer models.CustomerRequestPayload
	if err := c.BindJSON(&newCustomer); err != nil {
		return
	}
	newCustomer.CreatedDate = CreatedDate
	res := handler.PostCustomer(newCustomer)
	c.IndentedJSON(http.StatusOK, res)
}
func UpdateCustomerCustomer(c *gin.Context) {
	var fixCustomer models.Customer
	if err := c.BindJSON(&fixCustomer); err != nil {
		return
	}
	res := handler.UpdateCustomer(fixCustomer)
	c.IndentedJSON(http.StatusOK, res)

}

func GetCustomer(c *gin.Context) {
	uid := c.Query("uid")
	var response = handler.GetCustomer(uid)
	if response.UID == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "customer not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func DeleteCustomer(c *gin.Context) {
	uid := c.Query("uid")
	var response = handler.DeleteCustomer(uid)
	c.IndentedJSON(http.StatusOK, response)
}
