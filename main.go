package main

import (
	"net/http"
	"unify/internal/database"
	"unify/internal/handler"
	"unify/validation"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	r := gin.Default()
	validation.CORS(r)
	database.TestSQL()
	user := new(validation.User)
	r.POST("/validation", user.Verify)
	r.Handle(http.MethodGet, "/customer", handler.Customer)
	r.Handle(http.MethodGet, "/transaction", handler.Transaction)
	r.Handle(http.MethodGet, "/item", handler.Item)

	r.Run(":8080") // 0.0.0.0:8080 でサーバーを立てます。

}
