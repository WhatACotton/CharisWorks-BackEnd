package handler

import (
	"log"
	"net/http"
	"unify/internal/database"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func Post_Cart(c *gin.Context) {
	Cart_List, _ := Get_Cart_ID(c)
	NewCartReq := new(database.Cart_Request_Payload)
	err := c.BindJSON(&NewCartReq)
	if err != nil {
		log.Print(err)
	}
	err = NewCartReq.Cart(Cart_List.Cart_ID)
	if err != nil {
		log.Print(err)
	}
	Carts, err := database.Get_Cart_Info(Cart_List.Cart_ID)
	if err != nil {
		log.Print(err)
	}
	c.JSON(http.StatusOK, Carts)
}

func Get_Cart(c *gin.Context) {
	Cart_List, _ := Get_Cart_ID(c)
	Carts, err := database.Get_Cart_Info(Cart_List.Cart_ID)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, Carts)
	log.Print(Carts)
}

func Get_Cart_ID(c *gin.Context) (Cart_List database.Cart_List, UID string) {
	Cart_ID := new(string)
	Customer_SessionKey := new(string)
	*Customer_SessionKey = validation.Customer_Get_SessionKey(c)
	UID, err := database.Get_UID(*Customer_SessionKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(UID)
	if UID != "" {
		*Cart_ID, err = database.Get_Cart_ID(UID)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Cart_ID:", *Cart_ID)
		if *Cart_ID == "" {
			if *Customer_SessionKey != "new" {
				Cart_List := new(database.Cart_List)
				Cart_List.Session_Key = *Customer_SessionKey
				err := Cart_List.Get_Cart_ID_from_SessionKey()
				if err != nil {
					log.Fatal(err)
				}
				*Cart_ID = Cart_List.Cart_ID
			} else {
				*Cart_ID = "new"
			}
		}
	} else {
		*Cart_ID = "new"
	}

	if *Cart_ID == "new" {
		log.Print("not login")
		Cart_List.Session_Key = validation.Get_Cart_Session(c)
		if Cart_List.Session_Key == "new" {
			log.Print("don't have sessionKey")
			Cart_List.Cart_ID = validation.GetUUID()
		} else {
			err := Cart_List.Get_Cart_ID_from_SessionKey()
			if err != nil {
				log.Fatal(err)
			}
			database.Delete_Cart_List(Cart_List.Cart_ID)
		}
		Cart_List.Session_Key = validation.GetUUID()
		Cart_List.Create_Cart_List()
		validation.Set_Cart_Session(c, Cart_List.Session_Key)
	} else {
		log.Print("logined")
		validation.CartSessionEnd(c)
		Cart_List.Cart_ID = *Cart_ID
		Continue_LogIn(c)

	}
	return Cart_List, UID
}
