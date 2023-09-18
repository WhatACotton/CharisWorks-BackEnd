package handler

import (
	"log"
	"net/http"
	"unify/internal/database"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func PostCart(c *gin.Context) {
	Cart, _ := GetDatafromSessionKey(c)
	NewCartReq := new(database.CartContentRequestPayload)
	err := c.BindJSON(&NewCartReq)
	if err != nil {
		log.Print(err)
	}
	NewCartReq.Cart(Cart.CartID)
	Carts := database.GetCartContents(Cart.CartID)
	c.JSON(http.StatusOK, Carts)
}

func GetCart(c *gin.Context) {
	Cart, _ := GetDatafromSessionKey(c)
	log.Print(Cart.CartID)
	if Cart.SessionKey != "new" {
		CartContents := database.GetCartContents(Cart.CartID)
		if CartContents == nil {
			c.JSON(http.StatusOK, "There is no Cart")
		} else {
			c.JSON(http.StatusOK, CartContents)
		}
		log.Print(CartContents)
	} else {
		c.JSON(http.StatusOK, "未ログインです")
	}
}

func GetDatafromSessionKey(c *gin.Context) (Cart database.Cart, UserID string) {
	log.Print("Getting CartID...")
	CustomerSessionKey := validation.GetCustomerSessionKey(c)
	CartSessionKey := validation.GetCartSessionKey(c)
	if CartSessionKey != "new" {
		log.Print("have CartSessionKey")
		Cart.SessionKey = CartSessionKey
		Cart.CartSessionListGetCartID()
	}
	if CustomerSessionKey != "new" {
		log.Print("logined")
		UserID = database.GetUserID(CustomerSessionKey)
		CartIDfromCustomer := database.GetCartID(UserID)
		CartContents := database.GetCartContents(CartIDfromCustomer)
		if CartContents != nil {
			log.Print("have CartContents from CustomerData")
			Cart.CartID = CartIDfromCustomer
		}
		if Cart.CartID == "" {
			log.Print("don't have CartID in any place")
			Cart.CartID = validation.GetUUID()
		}
		validation.CartSessionEnd(c)
		OldSessionKey, NewSessionKey := validation.CustomerSessionStart(c)
		if OldSessionKey == "new" {
			validation.CustomerSessionEnd(c)
			c.JSON(http.StatusOK, "未ログインです")
		} else {
			UserID := database.GetUserID(OldSessionKey)
			log.Print("UserID : ", UserID)
			database.CustomerLogIn(UserID, NewSessionKey)
		}
		database.CustomerSetCartID(UserID, Cart.CartID)
	} else {
		log.Print("not logined")
		if Cart.CartID == "" {
			log.Print("don't have CartID in any place")
			Cart.CartID = validation.GetUUID()
		}
		database.CartSessionListDelete(Cart.CartID)
		Cart.SessionKey = validation.GetUUID()
		Cart.CartSessionListCreate()
		log.Print("Cart with sesssion. SessionKey : ", Cart.SessionKey, " CartID : ", Cart.CartID)
		validation.SetCartSessionKey(c, Cart.SessionKey)
	}
	log.Print("CartID:", Cart.CartID)
	return Cart, UserID
}
