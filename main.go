package main

import (
	"unify/internal/database"
	"unify/internal/handler"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	validation.CORS(r)
	validation.SessionConfig(r)
	database.TestSQL()
	r.GET("/item", handler.GetItem)
	r.GET("/itemlist", handler.GetItemList)

	r.GET("/Login", handler.LogIn)
	r.POST("/SignUp", handler.TemporarySignUp)
	r.GET("/Logout", handler.LogOut)
	r.POST("/Registration", handler.SignUp)
	r.POST("/Modify", handler.ModifyCustomer)
	r.DELETE("/DeleteCustomer", handler.DeleteCustomer)

	r.GET("/SessionStart", handler.ContinueLogIn)
	r.PATCH("/UpdateCart", handler.UpdateCart)
	r.POST("/PostCart", handler.PostCart)
	r.GET("/GetCart", handler.GetCart)

	r.POST("/Transaction", handler.BuyItem)
	r.Run(":8080") // 0.0.0.0:8080 でサーバーを立てます。
}
