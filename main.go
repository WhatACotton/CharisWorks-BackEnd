package main

import (
	"net/http"
	"unify/internal/database"
	"unify/internal/handler"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	validation.CORS(r)
	authorized := validation.Basic(r)
	database.TestSQL()
	r.POST("/SignUp", handler.TemporaryRegistration)
	r.POST("/Registration", handler.UserRegistration)
	r.Handle(http.MethodGet, "/transaction", handler.Transaction)
	r.GET("/item", handler.GetItem)
	r.GET("/itemlist", handler.GetItemList)

	authorized.Handle(http.MethodGet, "/customer", handler.CustomerAuthorized)
	authorized.Handle(http.MethodGet, "/transaction", handler.TransactionAuthorized)
	authorized.Handle(http.MethodGet, "/item", handler.ItemAuthorized)
	r.Run(":8080") // 0.0.0.0:8080 でサーバーを立てます。

}
