package funcs

import (
	"net/http"
	"unify/internal/database"

	"github.com/gin-gonic/gin"
)

func GetItemAuth(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		GetItemList(c)
	} else {
		GetItem(c, id)
	}
}

func GetItemList(c *gin.Context) { c.JSON(http.StatusOK, database.GetItemList()) }

func GetItem(c *gin.Context, id string) { c.JSON(http.StatusOK, database.GetItem(id)) }

func DeleteItem(c *gin.Context) {
	id := c.Query("id")
	database.Delete("itemlist", "ID", id)
	c.JSON(http.StatusOK, gin.H{"message": id + " is deleted"})
}
