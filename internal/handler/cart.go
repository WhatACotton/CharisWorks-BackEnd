package handler

import (
	"log"
	"net/http"
	"unify/internal/database"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func Cart_List_Session_Start(c *gin.Context) (database.Cart_List, string) {
	Cart_List := new(database.Cart_List)
	Cart_Session_Key := validation.Get_Cart_Session(c)
	if Cart_Session_Key != "new" {
		Cart_List.Session_Key = Cart_Session_Key
		Cart_List.Refresh_Cart_List()
	} else {
		Cart_List.Session_Key = validation.GetUUID()
		Cart_List.Cart_ID = validation.GetUUID()
		Cart_List.Create_Cart_List()
	}
	validation.Set_Cart_Session(c, Cart_List.Session_Key)

	return *Cart_List, Cart_Session_Key
}

func Post_Cart(c *gin.Context) {
	Cart_List, _ := Cart_List_Session_Start(c)
	NewCartReq := new(database.Cart_Request_Payload)
	err := c.BindJSON(&NewCartReq)
	if err != nil {
		log.Fatal(err)
	}
	err = NewCartReq.Cart(Cart_List.Cart_ID)
	if err != nil {
		log.Fatal(err)
	}
	Carts, err := database.Get_Cart_Info(Cart_List.Cart_ID)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, Carts)
}

func Get_Cart(c *gin.Context) {
	Cart_List, _ := Cart_List_Session_Start(c)
	Carts, err := database.Get_Cart_Info(Cart_List.Cart_ID)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, Carts)
}

func Update_Cart(c *gin.Context) {
	Cart_List, _ := Cart_List_Session_Start(c)
	NewCartReq := new(database.Cart_Request_Payload)
	err := c.BindJSON(&NewCartReq)
	if err != nil {
		log.Fatal(err)
	}
	err = NewCartReq.Update_Cart(Cart_List.Cart_ID)
	if err != nil {
		log.Fatal(err)
	}
	Carts, err := database.Get_Cart_Info(Cart_List.Cart_ID)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, Carts)
}
