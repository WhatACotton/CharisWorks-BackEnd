package handler

import (
	"log"
	"net/http"
	"unify/internal/database"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func PostCart(c *gin.Context) {
	Cart, _ := GetCartID(c)
	NewCartReq := new(database.CartRequestPayload)
	err := c.BindJSON(&NewCartReq)
	if err != nil {
		log.Print(err)
	}
	err = NewCartReq.Cart(Cart.CartID)
	if err != nil {
		log.Print(err)
	}
	Carts, err := database.GetCartInfo(Cart.CartID)
	if err != nil {
		log.Print(err)
	}
	c.JSON(http.StatusOK, Carts)
}

func GetCart(c *gin.Context) {
	Cart, _ := GetCartID(c)
	log.Print(Cart.CartID)
	if Cart.SessionKey != "new" {
		Carts, err := database.GetCartInfo(Cart.CartID)
		if err != nil {
			log.Fatal(err)
		}
		if Carts == nil {
			c.JSON(http.StatusOK, "There is no Cart")
		} else {
			c.JSON(http.StatusOK, Carts)
		}
		log.Print(Carts)
	} else {
		c.JSON(http.StatusOK, "未ログインです")
	}
}

func GetCartID(c *gin.Context) (Cart database.Cart, UID string) {
	CartID := new(string)
	CustomerSessionKey := new(string)
	*CustomerSessionKey = validation.GetCustomerSessionKey(c)
	UID, err := database.GetUID(*CustomerSessionKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(UID)
	if UID != "" {
		*CartID, err = database.GetCartID(UID)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("CartID:", *CartID)
		if *CartID == "" {
			if *CustomerSessionKey != "new" {
				Cart := new(database.Cart)
				Cart.SessionKey = *CustomerSessionKey
				err := Cart.GetCartIDfromSessionKey()
				if err != nil {
					log.Fatal(err)
				}
				*CartID = Cart.CartID
			} else {
				*CartID = "new"
			}
		}
	} else {
		*CartID = "new"
	}

	if *CartID == "new" {
		log.Print("not login")
		Cart.SessionKey = validation.GetCartSessionKey(c)
		if Cart.SessionKey == "new" {
			log.Print("don't have sessionKey")
			Cart.CartID = validation.GetUUID()
		} else {
			err := Cart.GetCartIDfromSessionKey()
			if err != nil {
				log.Fatal(err)
			}
			database.DeleteCartList(Cart.CartID)
		}
		Cart.SessionKey = validation.GetUUID()
		Cart.CreateCartList()
		validation.SetCartSessionKey(c, Cart.SessionKey)
	} else {
		log.Print("logined")
		validation.CartSessionEnd(c)
		Cart.CartID = *CartID
		LogInToDB(c)
	}
	return Cart, UID
}
