package handler

import (
	"net/http"
	"unify/internal/database"
	"unify/internal/models"

	"github.com/gin-gonic/gin"
)

func PostItem(c *gin.Context) {
	var newItem models.Item
	if err := c.BindJSON(&newItem); err != nil {
		return
	}
	c.JSON(http.StatusOK, database.PostItem(newItem))
}

func PatchItem(c *gin.Context) {
	h := new(database.PatchRequestPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	var isint bool
	if h.Attribute == "Price" ||
		h.Attribute == "Stonesize" ||
		h.Attribute == "Minlength" ||
		h.Attribute == "Maxlength" {
		isint = true
	} else {
		isint = false
	}
	h.Patch("itemlist", isint, "ID")
	c.JSON(http.StatusOK, database.GetItemList())
}

func GetItemList(c *gin.Context) {
	c.JSON(http.StatusOK, database.GetItemList())
}

func GetItem(c *gin.Context) {
	id := c.Query("id")
	c.JSON(http.StatusOK, database.GetItem(id))
}

func DeleteItem(c *gin.Context) {
	id := c.Query("id")
	database.Delete("itemlist", "ID", id)
	c.JSON(http.StatusOK, gin.H{"message": id + " is deleted"})
}
