package auth

import (
	"net/http"
	"unify/internal/database"
	"unify/internal/models"

	"github.com/gin-gonic/gin"
)

func GetItem(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		GetItemList(c)
	} else {
		GetAItem(c, id)
	}
}

func GetItemList(c *gin.Context) { c.JSON(http.StatusOK, database.GetItemList()) }

func GetAItem(c *gin.Context, id string) { c.JSON(http.StatusOK, database.GetItem(id)) }

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
	h.Patch("itemlist", "ID")
	c.JSON(http.StatusOK, database.GetItemList())
}

func DeleteItem(c *gin.Context) {
	id := c.Query("id")
	database.Delete("itemlist", "ID", id)
	c.JSON(http.StatusOK, gin.H{"message": id + " is deleted"})
}
