package validation

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SessionStart(c *gin.Context) (OldSessionKey string, NewSessionKey string) {
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

func SessionEnd(c *gin.Context) (OldSessionKey string) {
	session := sessions.DefaultMany(c, "SessionKey")
	if session.Get("SessionKey") != nil {
		session.Options(sessions.Options{MaxAge: -1})
		OldSessionKey = session.Get("SessionKey").(string)
		session.Clear()
		session.Save()
	}
	return OldSessionKey
}

func CartSessionEnd(c *gin.Context) (OldSessionKey string) {
	session := sessions.DefaultMany(c, "CartSessionKey")
	if session.Get("CartSessionKey") != nil {
		session.Options(sessions.Options{MaxAge: -1})
		OldSessionKey = session.Get("CartSessionKey").(string)
		session.Clear()
		session.Save()
	}
	return OldSessionKey
}

func SessionConfig(r *gin.Engine) {
	store := cookie.NewStore([]byte(GenerateRandomKey()))
	cookies := []string{"CartSessionKey", "SessionKey"}
	r.Use(sessions.SessionsMany(cookies, store))
}

func LogInStatus(c *gin.Context) bool {
	session := sessions.DefaultMany(c, "SessionKey")
	if session.Get("SessionKey") == nil {
		return false
	} else {
		return true
	}
}

func Customer_Get_SessionKey(c *gin.Context) string {
	session := sessions.DefaultMany(c, "SessionKey")
	if session.Get("SessionKey") != nil {
		SessionKey := session.Get("SessionKey").(string)

		return SessionKey
	} else {
		return "new"
	}
}
