package handler

import (
	"log"
	"net/http"
	"unify/cashing"
	"unify/internal/database"

	"github.com/gin-gonic/gin"
)

func CreateStripeAccount(c *gin.Context) {
	_, UID := GetDatafromSessionKey(c)
	StripeAccountID, err := database.GetStripeAccountID(UID)
	if err != nil {
		log.Fatal(err)
	}
	if StripeAccountID != "allow" {
		email, err := database.GetEmail(UID)
		if err != nil {
			log.Fatal(err)
		}
		AccountID, URL := cashing.CreateStripeAccount(email)
		err = database.CreateStripeAccount(UID, AccountID)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "アカウントが作成されました。", "URL": URL})
	} else {
		if StripeAccountID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "出品者登録が完了していません。", "errcode": "401"})
		}
		c.JSON(http.StatusOK, gin.H{"message": "アカウントが作成されています。"})
	}
}

func isMaker(c *gin.Context) bool {
	_, UID := GetDatafromSessionKey(c)
	StripeAccountID, err := database.GetStripeAccountID(UID)
	if err != nil {
		log.Fatal(err)
	}
	if StripeAccountID != "allow" && StripeAccountID != "" {
		return true
	}
	return false
}
func PostItem(c *gin.Context) {
	if isMaker(c) {
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "出品者登録が完了していません。", "errcode": "401"})
	}
}
