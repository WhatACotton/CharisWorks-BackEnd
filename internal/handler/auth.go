package handler

import (
	"net/http"
	"unify/internal/auth"

	"github.com/gin-gonic/gin"
)

func CustomerAuthorized(c *gin.Context) {
	requestMethod := http.MethodGet
	switch request := requestMethod; request {
	case "GET":
		auth.GetCustomer(c)
	case "POST":
		auth.PostCustomer(c)
	case "PATCH":
		auth.PatchCustomer(c)
	case "DELETE":
		auth.DeleteCustomer(c)
	}
}

func ItemAuthorized(c *gin.Context) {
	requestMethod := http.MethodGet
	switch request := requestMethod; request {
	case "GET":
		auth.GetItem(c)
	case "POST":
		auth.PostItem(c)
	case "PATCH":
		auth.PatchItem(c)
	case "DELETE":
		auth.DeleteItem(c)
	}
}

func TransactionAuthorized(c *gin.Context) {
	requestMethod := http.MethodGet
	switch request := requestMethod; request {
	case "GET":
		auth.GetTransaction(c)
	case "POST":
		auth.PostTransaction(c)
	case "PATCH":
		auth.PatchTransaction(c)
	case "DELETE":
		auth.DeleteTransaction(c)
	}
}
