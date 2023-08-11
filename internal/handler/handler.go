package handler

import (
	"net/http"
	"unify/internal/database"

	"github.com/gin-gonic/gin"
)

func GetItem(c *gin.Context) {
	id := c.Query("id")
	c.JSON(http.StatusOK, database.GetItem(id))
}

func GetItemList(c *gin.Context) {
	c.JSON(http.StatusOK, database.GetItemList())

}
