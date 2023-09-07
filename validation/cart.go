package validation

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Cart_Session_Start(c *gin.Context) (NewSessionKey string) {
	session := sessions.DefaultMany(c, "CartSessionKey")
	Session_Key := session.Get("CartSessionKey")
	if Session_Key == nil {
		SessionKey := GetUUID()
		session.Set("CartSessionKey", SessionKey)
		err := session.Save()
		if err != nil {
			log.Fatal(err)
		}
		return SessionKey
	} else {
		SessionKey := GetUUID()
		session.Set("CartSessionKey", NewSessionKey)
		session.Save()
		return SessionKey
	}
}

func Set_Cart_Session(c *gin.Context, Cart_Session_Key string) {
	session := sessions.DefaultMany(c, "CartSessionKey")
	session.Set("CartSessionKey", Cart_Session_Key)
	session.Save()
}

func Get_Cart_Session(c *gin.Context) (Cart_Session_Key string) {
	session := sessions.DefaultMany(c, "CartSessionKey")
	if session.Get("CartSessionKey") == nil {
		return "new"
	} else {
		Cart_Session_Key = session.Get("CartSessionKey").(string)
		return Cart_Session_Key
	}
}
