package handler

import (
	"net/http"
	"unify/internal/funcs"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func LoggedInPostCart(c *gin.Context) {
	user := new(validation.User)
	if user.Verify(c) { //認証
		funcs.PostCartLoggedIn(*user, c)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}
}
func PostCartWithSession(c *gin.Context) {

}
