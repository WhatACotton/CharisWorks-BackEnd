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
	if validation.LogInStatus(c) {
		OldCartSessionKey, NewCartSessionKey := validation.CartSessionStart(c)
		OldSessionKey, NewSessionKey := validation.SessionStart(c)
		UID := database.GetUID(OldSessionKey)
		database.LogInLog(UID, NewSessionKey)
		database.Invalid(OldSessionKey)
		if OldCartSessionKey == "new" {
			CartId := database.NewCartLogin(NewCartSessionKey, UID)
			NewCartReq := new(models.CartRequestPayload)
			err := c.BindJSON(&NewCartReq)
			if err != nil {
				log.Fatal(err)
			}
			Carts := database.PostCart(*NewCartReq, CartId)
			c.JSON(http.StatusOK, Carts)
		} else {
			if database.VerifyCart(OldCartSessionKey) {
				CartId := database.GetCartId(OldCartSessionKey)
				database.StoredLoginCart(NewCartSessionKey, CartId, UID)
				database.CartInvalid(OldCartSessionKey)
				NewCartReq := new(models.CartRequestPayload)
				err := c.BindJSON(&NewCartReq)
				if err != nil {
					log.Fatal(err)
				}
				Carts := database.PostCart(*NewCartReq, CartId)
				c.JSON(http.StatusOK, Carts)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})

			}
		}
	} else {
		OldSessionKey, NewSessionKey := validation.CartSessionStart(c)
		if OldSessionKey == "new" {
			CartId := database.NewCartSession(NewSessionKey)
			NewCartReq := new(models.CartRequestPayload)
			err := c.BindJSON(&NewCartReq)
			if err != nil {
				log.Fatal(err)
			}
			Carts := database.PostCart(*NewCartReq, CartId)
			c.JSON(http.StatusOK, Carts)
		} else {
			if database.VerifyCart(OldSessionKey) {
				CartId := database.GetCartId(OldSessionKey)
				database.StoredCartSession(NewSessionKey, CartId)
				database.CartInvalid(OldSessionKey)
				NewCartReq := new(models.CartRequestPayload)
				err := c.BindJSON(&NewCartReq)
				if err != nil {
					log.Fatal(err)
				}
				Carts := database.PostCart(*NewCartReq, CartId)
				c.JSON(http.StatusOK, Carts)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})

			}
		}
	}

}
func GetCart(c *gin.Context) {
	if validation.LogInStatus(c) {
		OldCartSessionKey, NewCartSessionKey := validation.CartSessionStart(c)
		OldSessionKey, NewSessionKey := validation.SessionStart(c)
		UID := database.GetUID(OldSessionKey)
		database.LogInLog(UID, NewSessionKey)
		database.Invalid(OldSessionKey)
		if OldCartSessionKey == "new" {
			database.NewCartLogin(NewCartSessionKey, UID)
		} else {
			if database.VerifyCart(OldCartSessionKey) {
				CartId := database.NewCartLogin(NewCartSessionKey, UID)
				database.Invalid(OldCartSessionKey)
				Carts := database.GetCart(CartId)
				c.JSON(http.StatusOK, Carts)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})

			}
		}

	} else {
		OldCartSessionKey, NewCartSessionKey := validation.CartSessionStart(c)
		if OldCartSessionKey == "new" {
			database.NewCartSession(NewCartSessionKey)
		} else {
			if database.VerifyCart(OldCartSessionKey) {
				CartId := database.NewCartSession(NewCartSessionKey)
				database.Invalid(OldCartSessionKey)
				Carts := database.GetCart(CartId)
				c.JSON(http.StatusOK, Carts)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})

			}
		}
	}
}
