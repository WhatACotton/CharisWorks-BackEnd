package handler

import (
	"net/http"
	"unify/cashing"
	"unify/internal/database"

	"github.com/gin-gonic/gin"
)

func CreateStripeAccount(c *gin.Context) {
	_, UserID := GetDatafromSessionKey(c)
	StripeAccountID := database.GetStripeAccountID(UserID)
	if StripeAccountID != "allow" {
		email := database.GetEmail(UserID)
		AccountID, URL := cashing.CreateStripeAccount(email)
		database.CustomerCreateStripeAccount(UserID, AccountID)

		c.JSON(http.StatusOK, gin.H{"message": "アカウントが作成されました。", "URL": URL})
	} else {
		if StripeAccountID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "出品者登録が完了していません。", "errcode": "401"})
		}
		c.JSON(http.StatusOK, gin.H{"message": "アカウントが作成されています。"})
	}
}

func isMaker(c *gin.Context) bool {
	_, UserID := GetDatafromSessionKey(c)
	StripeAccountID := database.GetStripeAccountID(UserID)
	if StripeAccountID != "allow" && StripeAccountID != "" {
		return true
	}
	return false
}
func PostItem(c *gin.Context) {
	if isMaker(c) {
		postItem(c)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "出品者登録が完了していません。", "errcode": "401"})
	}
}
func postItem(c *gin.Context) {

}
