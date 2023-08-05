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
	validation.SessionConfig(r)
	validation.CartConfig(r)
	authorized := validation.Basic(r)
	database.TestSQL()
	r.GET("/item", handler.GetItem)
	r.GET("/itemlist", handler.GetItemList)

	r.GET("/Login", handler.LogIn)
	//r.GET("/SessionRestart", handler.SessionRestart)
	r.POST("/SignUp", handler.TemporarySignUp)
	r.POST("/Registration", handler.SignUp)
	r.POST("/Modify", handler.ModifyCustomer)
	r.DELETE("/DeleteCustomer", handler.DeleteCustomer)

	r.GET("/SessionStart", handler.ContinueLogIn)
	r.POST("/PostCart", handler.PostCart)
	r.GET("/GetCart", handler.GetCart)
	//r.GET("/DeleteSession", handler.DeleteLoginSession)
	r.Handle(http.MethodGet, "/transaction", handler.Transaction)
	authorized.Handle(http.MethodGet, "/customer", handler.CustomerAuthorized)
	authorized.Handle(http.MethodGet, "/transaction", handler.TransactionAuthorized)
	authorized.Handle(http.MethodGet, "/item", handler.ItemAuthorized)
	r.Run(":8080") // 0.0.0.0:8080 でサーバーを立てます。
}
