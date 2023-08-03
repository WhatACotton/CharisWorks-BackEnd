package handler

import (
	"net/http"
	"unify/internal/auth"
	"unify/internal/funcs"

	"github.com/gin-gonic/gin"
)

func GetItem(c *gin.Context) {
	id := c.Query("id")
	funcs.GetItem(c, id)
}

func GetItemList(c *gin.Context) {
	funcs.GetItemList(c)
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
