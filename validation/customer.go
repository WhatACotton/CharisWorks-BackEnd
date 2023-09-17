package validation

import (
	"fmt"
	"log"
	"regexp"

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

func GetCustomerSessionKey(c *gin.Context) string {
	session := sessions.DefaultMany(c, "SessionKey")
	if session.Get("SessionKey") != nil {
		SessionKey := session.Get("SessionKey").(string)

		return SessionKey
	} else {
		return "new"
	}
}
func CustomerSessionEnd(c *gin.Context) (OldSessionKey string) {
	session := sessions.DefaultMany(c, "SessionKey")
	if session.Get("SessionKey") != nil {
		OldSessionKey = session.Get("SessionKey").(string)
		session.Delete("SessionKey")
		session.Save()
	}
	return OldSessionKey
}

type CustomerRegisterPayload struct {
	Name        string `json:"Name"`
	ZipCode     string `json:"ZipCode"`
	Address     string `json:"Address"`
	PhoneNumber string `json:"PhoneNumber"`
}

func (c *CustomerRegisterPayload) InspectCusromerRegisterPayload() bool {
	if c.Name == "" || c.Address == "" || c.ZipCode == "" {
		return false
	}

	zipCodePattern := regexp.MustCompile(`\d{3}-\d{4}`)
	matched := zipCodePattern.MatchString(c.ZipCode)
	if matched {
	} else {
		fmt.Println("郵便番号の形式が正しくありません")
		return false
	}
	return true
}
func (c *CustomerRegisterPayload) InspectFirstRegisterPayload() bool {
	if c.Name == "山田 カリス" || c.Address == "住所" || c.ZipCode == "000-0000" {
		return false
	}
	return true
}
