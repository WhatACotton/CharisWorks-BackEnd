package handler

import (
	"net/http"
	"unify/internal/auth"
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

func Transaction(c *gin.Context) {
	requestMethod := http.MethodGet
	switch request := requestMethod; request {
	case "GET":
		auth.GetTransaction(c)
	case "POST":
		auth.PostTransaction(c)
	}
}
