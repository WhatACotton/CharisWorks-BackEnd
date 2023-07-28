package handler

import (
	"net/http"
	"unify/internal/funcs"

	"github.com/gin-gonic/gin"
)

func Customer(c *gin.Context) {
	requestMethod := http.MethodGet
	switch request := requestMethod; request {
	case "GET":
		funcs.GetCustomer(c)
	case "POST":
		funcs.PostCustomer(c)
	case "PATCH":
		funcs.PatchCustomer(c)
	case "DELETE":
		funcs.DeleteCustomer(c)
	}
}

func Item(c *gin.Context) {
	requestMethod := http.MethodGet
	switch request := requestMethod; request {
	case "GET":
		funcs.GetItem(c)
	case "POST":
		funcs.PostItem(c)
	case "PATCH":
		funcs.PatchItem(c)
	case "DELETE":
		funcs.DeleteItem(c)
	}
}

func Transaction(c *gin.Context) {
	requestMethod := http.MethodGet
	switch request := requestMethod; request {
	case "GET":
		funcs.GetTransaction(c)
	case "POST":
		funcs.PostTransaction(c)
	case "PATCH":
		funcs.PatchTransaction(c)
	case "DELETE":
		funcs.DeleteTransaction(c)
	}
}
