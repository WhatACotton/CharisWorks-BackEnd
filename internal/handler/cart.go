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
	err = NewCartReq.Cart(Cart.CartID)
	if err != nil {
		log.Print(err)
	}
	Carts, err := database.GetCartContents(Cart.CartID)
	if err != nil {
		log.Print(err)
	}
	c.JSON(http.StatusOK, Carts)
}

func GetCart(c *gin.Context) {
	Cart, _ := GetDatafromSessionKey(c)
	log.Print(Cart.CartID)
	if Cart.SessionKey != "new" {
		CartContents, err := database.GetCartContents(Cart.CartID)
		if err != nil {
			log.Fatal(err)
		}
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

func GetDatafromSessionKey(c *gin.Context) (Cart database.Cart, UID string) {
	log.Print("Getting CartID...")
	CustomerSessionKey := validation.GetCustomerSessionKey(c)
	CartSessionKey := validation.GetCartSessionKey(c)
	if CartSessionKey != "new" {
		log.Print("have CartSessionKey")
		Cart.SessionKey = CartSessionKey
		Cart.GetCartIDfromCartSessionKey()
	}
	if CustomerSessionKey != "new" {
		log.Print("logined")
		UID, _ = database.GetUID(CustomerSessionKey)
		CartIDfromCustomer, err := database.GetCartID(UID)
		if err != nil {
			log.Fatal(err)
		}
		CartContents, err := database.GetCartContents(CartIDfromCustomer)
		if err != nil {
			log.Fatal(err)
		}
		if CartContents != nil {
			log.Print("have CartContents from CustomerData")
			Cart.CartID = CartIDfromCustomer
		}
		if Cart.CartID == "" {
			log.Print("don't have CartID in any place")
			Cart.CartID = validation.GetUUID()
		}
		validation.CartSessionEnd(c)
		LogInToDB(c)
		database.SetCartID(UID, Cart.CartID)
	} else {
		log.Print("not logined")
		if Cart.CartID == "" {
			log.Print("don't have CartID in any place")
			Cart.CartID = validation.GetUUID()
		}
		Cart.SessionKey = validation.GetUUID()
		Cart.CreateCart()
		log.Print("Cart with sesssion. SessionKey : ", Cart.SessionKey, " CartID : ", Cart.CartID)
		validation.SetCartSessionKey(c, Cart.SessionKey)
	}
	log.Print("CartID:", Cart.CartID)
	return Cart, UID
}
