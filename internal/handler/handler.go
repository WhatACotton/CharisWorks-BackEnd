package handler

import "github.com/gin-gonic/gin"

func Customer(c *gin.Context) {
	requestMethod := c.Request.Method
	switch request := requestMethod; request {
	case "GET":
		GetCustomer(c)
	case "POST":
		PostCustomer(c)
	case "PATCH":
		PatchCustomer(c)
	case "DELETE":
		DeleteCustomer(c)
	}
}

func Item(c *gin.Context) {
	requestMethod := c.Request.Method
	switch request := requestMethod; request {
	case "GET":
		GetItem(c)
	case "POST":
		PostItem(c)
	case "PATCH":
		PatchItem(c)
	case "DELETE":
		DeleteItem(c)
	}
}

func Transaction(c *gin.Context) {
	requestMethod := c.Request.Method
	switch request := requestMethod; request {
	case "GET":
		GetTransaction(c)
	case "POST":
		PostTransaction(c)
	case "PATCH":
		PatchTransaction(c)
	case "DELETE":
		DeleteTransaction(c)
	}
}
