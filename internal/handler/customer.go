package handler

import (
	"log"
	"net/http"
	"unify/internal/database"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

// 仮登録を行う。ここでの登録内容はUIDと作成日時だけ。
func TemporarySignUp(c *gin.Context) {
	CustomerReqPayload := new(validation.CustomerReqPayload)
	if CustomerReqPayload.VerifyCustomer(c) {
		Cart := new(database.Cart)
		Cart.SessionKey = validation.GetCartSessionKey(c)
		if !Cart.SessionGet() {
			log.Print("don't have sessionKey")
			Cart.CartID = validation.GetUUID()
		}
		log.Print("CartID: ", Cart.CartID)
		res := database.SignUpCustomer(*CustomerReqPayload, signUpToDB(c, CustomerReqPayload.UID), Cart.CartID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func SignUp(c *gin.Context) {
	//登録情報変更処理
	//bodyにアカウントの詳細情報が入っている。
	_, UID := GetDatafromSessionKey(c)
	h := new(validation.CustomerRegisterPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	if h.InspectCusromerRegisterPayload() {
		database.RegisterCustomer(UID, *h)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "変更できませんでした。"})
	}
}

func LogIn(c *gin.Context) {
	UserReqPayload := new(validation.CustomerReqPayload)
	if UserReqPayload.VerifyCustomer(c) {
		log.Print("UID : ", UserReqPayload.UID)
		if UserReqPayload.EmailVerified {
			err := database.EmailVerified(1, UserReqPayload.UID)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := database.EmailVerified(0, UserReqPayload.UID)
			if err != nil {
				log.Fatal(err)
			}
		}
		email, err := database.GetEmail(UserReqPayload.UID)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(email)
		if email != UserReqPayload.Email {
			database.ChangeEmail(UserReqPayload.UID, UserReqPayload.Email)
		}
		_ = signUpToDB(c, UserReqPayload.UID)
		GetDatafromSessionKey(c)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}

}

func GetCustomer(c *gin.Context) {
	UID := LogInToDB(c)
	if UID != "" {
		Customer := new(database.Customer)
		Customer.GetCustomer(UID)
		c.JSON(http.StatusOK, Customer)
	}
}

func ModifyCustomer(c *gin.Context) {
	//登録情報変更処理
	//bodyにアカウントの詳細情報が入っている。
	_, UID := GetDatafromSessionKey(c)
	h := new(validation.CustomerRegisterPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	if h.InspectCusromerRegisterPayload() {
		err := database.RegisterCustomer(UID, *h)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("CustomerData was modified")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "変更できませんでした。"})
	}
}

func LogOut(c *gin.Context) {
	//ログアウト処理
	OldSessionKey := validation.CustomerSessionEnd(c)
	database.Invalid(OldSessionKey)
	log.Print("SessionKey was :", OldSessionKey)
	//c.JSON(http.StatusOK, "SuccessFully Logouted!!")
}

func DeleteCustomer(c *gin.Context) {
	//アカウントの削除
	_, UID := GetDatafromSessionKey(c)
	database.DeleteCustomer(UID)
	database.DeleteSession(UID)
	c.JSON(http.StatusOK, gin.H{"message": "アカウントを削除しました。"})

}

func LogInToDB(c *gin.Context) (UID string) {
	OldSessionKey, NewSessionKey := validation.CustomerSessionStart(c)
	if OldSessionKey == "new" {
		validation.CustomerSessionEnd(c)
		c.JSON(http.StatusOK, "未ログインです")
		return ""
	} else {
		UID, err := database.GetUID(OldSessionKey)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("UID : ", UID)
		database.LogIn(UID, NewSessionKey)
		return UID
	}
}
func signUpToDB(c *gin.Context, UID string) (SessionKey string) {
	_, SessionKey = validation.CustomerSessionStart(c)
	log.Print("UID : ", UID)
	database.LogInLog(UID, SessionKey)
	return UID
}
