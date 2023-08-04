package main

import (
	"net/http"
	"unify/internal/database"
	"unify/internal/handler"
	"unify/validation"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	validation.CORS(r)
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

	r.POST("/LoggedInCart", handler.LoggedInPostCart)
	r.POST("/PostSessionCart", handler.PostCart)
	r.GET("/GetSessionCart", handler.GetCart)
	//r.GET("/DeleteSession", handler.DeleteLoginSession)
	r.Handle(http.MethodGet, "/transaction", handler.Transaction)

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/hello", func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("hello") != "world" {
			session.Set("hello", "world")
			session.Save()
		}

		c.JSON(200, gin.H{"hello": session.Get("hello")})
	})
	authorized.Handle(http.MethodGet, "/customer", handler.CustomerAuthorized)
	authorized.Handle(http.MethodGet, "/transaction", handler.TransactionAuthorized)
	authorized.Handle(http.MethodGet, "/item", handler.ItemAuthorized)
	r.Run(":8080") // 0.0.0.0:8080 でサーバーを立てます。
}
