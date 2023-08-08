package handler

import (
	"log"
	"net/http"
	"unify/internal/database"
	"unify/internal/models"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func PostCart(c *gin.Context) {
	CartId := new(string)
	OldSessionKey, NewSessionKey := validation.SessionStart(c)
	OldCartSessionKey, NewCartSessionKey := validation.CartSessionStart(c)
	if OldCartSessionKey == "new" {
		*CartId = validation.GetUUID()
	} else {
		if database.VerifyCart(OldCartSessionKey) {
			*CartId = database.GetCartId(OldCartSessionKey)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})
		}
	}
	if validation.LogInStatus(c) {
		UID := database.GetUID(OldSessionKey)
		database.LogInLog(UID, NewSessionKey)
		database.Invalid(OldSessionKey)
		database.LoginCart(NewCartSessionKey, *CartId, UID)
	} else {
		database.SessionCart(NewCartSessionKey, *CartId)
	}
	database.CartInvalid(OldCartSessionKey)
	NewCartReq := new(models.CartRequestPayload)
	err := c.BindJSON(&NewCartReq)
	if err != nil {
		log.Fatal(err)
	}
	Carts := database.PostCart(*NewCartReq, *CartId)
	c.JSON(http.StatusOK, Carts)
}

func GetCart(c *gin.Context) {
	CartId := new(string)
	OldCartSessionKey, NewCartSessionKey := validation.CartSessionStart(c)
	OldSessionKey, NewSessionKey := validation.SessionStart(c)
	if OldCartSessionKey != "new" {
		if database.VerifyCart(OldCartSessionKey) {
			*CartId = database.GetCartId(OldCartSessionKey)
			database.CartInvalid(OldCartSessionKey)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})
		}
		if validation.LogInStatus(c) {
			UID := database.GetUID(OldSessionKey)
			database.LogInLog(UID, NewSessionKey)
			database.Invalid(OldSessionKey)
			database.LoginCart(NewCartSessionKey, *CartId, UID)
		} else {
			database.SessionCart(NewCartSessionKey, *CartId)
		}
		Carts := database.GetCart(*CartId)
		c.JSON(http.StatusOK, Carts)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "カートが見つかりませんでした"})
	}
}

func UpdateCart(c *gin.Context) {
	CartId := new(string)
	OldCartSessionKey, NewCartSessionKey := validation.CartSessionStart(c)
	OldSessionKey, NewSessionKey := validation.SessionStart(c)
	if OldCartSessionKey != "new" {
		if database.VerifyCart(OldCartSessionKey) {
			*CartId = database.GetCartId(OldCartSessionKey)
			database.CartInvalid(OldCartSessionKey)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})
		}
		if validation.LogInStatus(c) {
			UID := database.GetUID(OldSessionKey)
			database.LogInLog(UID, NewSessionKey)
			database.Invalid(OldSessionKey)
			database.LoginCart(NewCartSessionKey, *CartId, UID)
		} else {
			database.SessionCart(NewCartSessionKey, *CartId)
		}
		NewCartReq := new(models.CartRequestPayload)
		err := c.BindJSON(&NewCartReq)
		if err != nil {
			log.Fatal(err)
		}
		database.UpdateCart(*CartId, *NewCartReq)
		Carts := database.GetCart(*CartId)
		c.JSON(http.StatusOK, Carts)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "カートが見つかりませんでした"})
	}
}
