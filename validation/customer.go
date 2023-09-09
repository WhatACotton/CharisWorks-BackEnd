package validation

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CustomerSessionStart(c *gin.Context) (OldSessionKey string, NewSessionKey string) {
	session := sessions.DefaultMany(c, "SessionKey")
	if session.Get("SessionKey") == nil {
		SessionKey := GetUUID()
		session.Set("SessionKey", SessionKey)
		err := session.Save()
		if err != nil {
			log.Fatal(err)
		}
		return "new", SessionKey
	} else {
		SessionKey := session.Get("SessionKey")
		NewSessionKey := GetUUID()
		session.Set("SessionKey", NewSessionKey)
		session.Save()
		return SessionKey.(string), NewSessionKey
	}
}

func CustomerSessionEnd(c *gin.Context) (OldSessionKey string) {
	session := sessions.DefaultMany(c, "SessionKey")
	if session.Get("SessionKey") != nil {
		session.Options(sessions.Options{MaxAge: -1})
		OldSessionKey = session.Get("SessionKey").(string)
		session.Clear()
		session.Save()
	}
	return OldSessionKey
}

func GetCustomerSessionKey(c *gin.Context) string {
	session := sessions.DefaultMany(c, "SessionKey")
	if session.Get("SessionKey") != nil {
		SessionKey := session.Get("SessionKey").(string)

		return SessionKey
	} else {
		return "new"
	}
}
