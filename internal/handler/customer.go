package handler

import (
	"net/http"

	"unify/internal/database"
	"unify/internal/models"

	"github.com/gin-gonic/gin"
)

var customers []models.Customer

func PostCustomer(c *gin.Context) {
	var CreatedDate = database.GetDate()
	var newCustomer models.CustomerRequestPayload
	if err := c.BindJSON(&newCustomer); err != nil {
		return
	}
	newCustomer.CreatedDate = CreatedDate
	res := database.PostCustomer(newCustomer)
	c.IndentedJSON(http.StatusOK, res)
}

func UpdateCustomerCustomer(c *gin.Context) {
	h := new(database.PatchRequestPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	h.Patch("user", false, "uid")
	GetCustomer(c)
}

func GetCustomer(c *gin.Context) {
	uid := c.Query("uid")
	var response = database.GetCustomer(uid)
	if response.UID == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "customer not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func DeleteCustomer(c *gin.Context) {
	uid := c.Query("uid")
	database.Delete("user", "uid", uid)
	GetCustomer(c)
}
