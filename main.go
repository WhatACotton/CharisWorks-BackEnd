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
	r.Handle(http.MethodGet, "/customer", handler.Customer)
	r.Handle(http.MethodGet, "/transaction", handler.Transaction)
	r.Handle(http.MethodGet, "/item", handler.Item)

	authorized.Handle(http.MethodGet, "/customer", handler.Customer)
	authorized.Handle(http.MethodGet, "/transaction", handler.Transaction)
	authorized.Handle(http.MethodGet, "/item", handler.Item)
	r.Run(":8080") // 0.0.0.0:8080 でサーバーを立てます。

}
