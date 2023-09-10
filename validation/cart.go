package validation

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetCartSessionKey(c *gin.Context, CartSessionKey string) {
	session := sessions.DefaultMany(c, "CartSessionKey")
	session.Set("CartSessionKey", CartSessionKey)
	session.Save()
}

func GetCartSessionKey(c *gin.Context) (CartSessionKey string) {
	session := sessions.DefaultMany(c, "CartSessionKey")
	if session.Get("CartSessionKey") == nil {
		return "new"
	} else {
		CartSessionKey = session.Get("CartSessionKey").(string)
		return CartSessionKey
	}
}
func CartSessionEnd(c *gin.Context) (OldSessionKey string) {
	session := sessions.DefaultMany(c, "CartSessionKey")
	if session.Get("CartSessionKey") != nil {
		OldSessionKey = session.Get("CartSessionKey").(string)
		session.Delete("CartSessionKey")
		session.Save()
	}
	return OldSessionKey
}
