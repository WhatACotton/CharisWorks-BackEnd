package handler

import (
	"log"
	"net/http"
	"os"
	"unify/internal/database"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

// 仮登録を行う。ここでの登録内容はUserIDと作成日時だけ。
func TemporarySignUp(c *gin.Context) {
	CustomerReqPayload := new(validation.CustomerReqPayload)
	if CustomerReqPayload.VerifyCustomer(c) {
		Cart := new(database.Cart)
		Cart.SessionKey = validation.GetCartSessionKey(c)
		if Cart.SessionKey == "new" {
			log.Print("don't have sessionKey")
			Cart.CartID = validation.GetUUID()
		}
		log.Print("CartID: ", Cart.CartID)
		_, NewSessionKey := validation.CustomerSessionStart(c)
		database.CustomerSignUp(*CustomerReqPayload, NewSessionKey, Cart.CartID)
		file, err := os.OpenFile("accountlog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		logger2 := log.New(file, "", log.Ldate|log.Ltime)
		logger2.SetOutput(file)
		logger2.Println("UserID :", CustomerReqPayload.UserID, " created.")
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func SignUp(c *gin.Context) {
	//登録情報変更処理
	//bodyにアカウントの詳細情報が入っている。
	_, UserID := GetDatafromSessionKey(c)
	h := new(validation.CustomerRegisterPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	if h.InspectCusromerRegisterPayload() {
		if h.InspectFirstRegisterPayload() {
			database.CustomerRegister(UserID, *h)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "未入力欄があります。"})
		}
	} else {

		c.JSON(http.StatusBadRequest, gin.H{"message": "変更できませんでした。"})
	}
}

func LogIn(c *gin.Context) {
	UserReqPayload := new(validation.CustomerReqPayload)
	if UserReqPayload.VerifyCustomer(c) {
		log.Print("UserID : ", UserReqPayload.UserID)
		if UserReqPayload.EmailVerified {
			database.CustomerEmailVerified(1, UserReqPayload.UserID)
		} else {
			database.CustomerEmailVerified(0, UserReqPayload.UserID)
		}
		Email := database.GetEmail(UserReqPayload.UserID)
		log.Print(Email)
		if Email != UserReqPayload.Email {
			database.CustomerChangeEmail(UserReqPayload.UserID, UserReqPayload.Email)
		}
		GetDatafromSessionKey(c)
		file, err := os.OpenFile("accountlog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		logger2 := log.New(file, "", log.Ldate|log.Ltime)
		logger2.SetOutput(file)
		logger2.Println("UserID :", UserReqPayload.UserID, " logined.")
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}

}

func GetCustomer(c *gin.Context) {
	_, UserID := GetDatafromSessionKey(c)
	if UserID != "" {
		Customer := new(database.Customer)
		Customer.GetCustomer(UserID)
		c.JSON(http.StatusOK, Customer)
	}
}

func ModifyCustomer(c *gin.Context) {
	//登録情報変更処理
	//bodyにアカウントの詳細情報が入っている。
	_, UserID := GetDatafromSessionKey(c)
	h := new(validation.CustomerRegisterPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	if h.InspectCusromerRegisterPayload() {
		database.CustomerRegister(UserID, *h)
		log.Print("CustomerData was modified")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "変更できませんでした。"})
	}
}

func LogOut(c *gin.Context) {
	_, UserID := GetDatafromSessionKey(c)
	//c.JSON(http.StatusOK, "SuccessFully Logouted!!")
	file, err := os.OpenFile("accountlog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger2 := log.New(file, "", log.Ldate|log.Ltime)
	logger2.SetOutput(file)
	logger2.Println("UserID :", UserID, " logouted.")
	//ログアウト処理
	OldSessionKey := validation.CustomerSessionEnd(c)

	log.Print("SessionKey was :", OldSessionKey)

}

func DeleteCustomer(c *gin.Context) {
	//アカウントの削除
	_, UserID := GetDatafromSessionKey(c)
	database.CustomerDelete(UserID)
	database.CustomerDeleteSession(UserID)
	file, err := os.OpenFile("accountlog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger2 := log.New(file, "", log.Ldate|log.Ltime)
	logger2.SetOutput(file)
	logger2.Println("UserID :", UserID, " deleted.")
	c.JSON(http.StatusOK, gin.H{"message": "アカウントを削除しました。"})

}
