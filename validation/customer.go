package validation

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CustomerSessionStart(c *gin.Context) (OldSessionKey string, NewSessionKey string) {
	session := sessions.Default(c)
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

// CustomerSessionKeyの取得
func GetCustomerSessionKey(c *gin.Context) string {
	session := sessions.Default(c)
	if session.Get("SessionKey") != nil {
		SessionKey := session.Get("SessionKey").(string)

		return SessionKey
	} else {
		return "new"
	}
}
func CustomerSessionEnd(c *gin.Context) (OldSessionKey string) {
	session := sessions.Default(c)
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
	Address1    string `json:"Address1"`
	Address2    string `json:"Address2"`
	Address3    string `json:"Address3"`
	PhoneNumber string `json:"PhoneNumber"`
}

func (c *CustomerRegisterPayload) InspectCusromerRegisterPayload() bool {
	if c.Name == "" || c.Address1 == "" || c.Address2 == "" || c.PhoneNumber == "" || c.ZipCode == "" {
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
	if c.Name == "default" || c.Address1 == "default" || c.Address2 == "default" || c.Address3 == "default" || c.PhoneNumber == "default" || c.ZipCode == "default" {
		return false
	}
	return true
}
