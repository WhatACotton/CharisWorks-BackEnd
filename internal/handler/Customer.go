package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/WhatACotton/go-backend-test/internal/database"
	"github.com/WhatACotton/go-backend-test/validation"
	"github.com/gin-gonic/gin"
)

// 初回登録
func Register(c *gin.Context) {
	//登録情報変更処理
	//bodyにアカウントの詳細情報が入っている。
	UserID := GetDatafromSessionKey(c)
	h := new(validation.CustomerRegisterPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	if h.InspectCusromerRegisterPayload() {
		//初回登録なので、記入漏れがないかの確認
		if h.InspectFirstRegisterPayload() {
			database.CustomerRegister(UserID, *h)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "未入力欄があります。"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "変更できませんでした。"})
	}
}

// ログイン
func LogIn(c *gin.Context) {
	UserReqPayload := new(validation.CustomerReqPayload)
	//ユーザ認証
	if UserReqPayload.VerifyCustomer(c) {
		//Eamil認証
		if UserReqPayload.EmailVerified {
			database.CustomerEmailVerified(1, UserReqPayload.UserID)
		} else {
			database.CustomerEmailVerified(0, UserReqPayload.UserID)
		}
		Customer := new(database.Customer)
		Customer.CustomerGet(UserReqPayload.UserID)
		if Customer.UserID == "not found" {
			log.Print("New User")
			database.CustomerSignUp(*UserReqPayload)
		} else {
			//Emailが変更されているかどうかの処理
			Email := database.GetEmail(UserReqPayload.UserID)
			if Email != UserReqPayload.Email {
				log.Print("Email was changed")
				database.CustomerChangeEmail(UserReqPayload.UserID, UserReqPayload.Email)
			}

		}
		_, NewSessionKey := validation.CustomerSessionStart(c)
		database.CustomerLogIn(UserReqPayload.UserID, NewSessionKey)
		validation.LoginLogging(UserReqPayload.UserID + " logined")
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}
}

// 顧客情報の取得
func GetCustomer(c *gin.Context) {
	UserID := GetDatafromSessionKey(c)
	if UserID != "" {
		Customer := new(database.Customer)
		Customer.CustomerGet(UserID)
		c.JSON(http.StatusOK, gin.H{"Customer": Customer})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未ログインです。"})
	}
}

// 顧客情報の更新
func ModifyCustomer(c *gin.Context) {
	//登録情報変更処理
	//bodyにアカウントの詳細情報が入っている。
	UserID := GetDatafromSessionKey(c)
	h := new(validation.CustomerRegisterPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	log.Print(h)
	if h.InspectCusromerRegisterPayload() {
		database.CustomerRegister(UserID, *h)
		log.Print("CustomerData was modified")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "変更できませんでした。"})
	}
}

// ログアウト
func LogOut(c *gin.Context) {
	UserID := GetDatafromSessionKey(c)
	//c.JSON(http.StatusOK, "SuccessFully Logouted!!")
	validation.LoginLogging(UserID + " logouted")
	//ログアウト処理
	OldSessionKey := validation.CustomerSessionEnd(c)
	log.Print("SessionKey was :", OldSessionKey)
}

// アカウントの削除
func DeleteCustomer(c *gin.Context) {
	//アカウントの削除
	UserID := GetDatafromSessionKey(c)
	database.CustomerDelete(UserID)
	validation.LoginLogging(UserID + "deleted")

	c.JSON(http.StatusOK, gin.H{"message": "アカウントを削除しました。"})
}

// セッションキーから、UserIDとCartを取得
// 同時にログイン中ならログイン状態の更新も行う
func GetDatafromSessionKey(c *gin.Context) (UserID string) {
	CustomerSessionKey := validation.GetCustomerSessionKey(c)
	//次にCustomerSessionKeyを持っているならログイン情報の取得
	if CustomerSessionKey != "new" {
		//ログインしている
		log.Print("logined")
		//CustomerSessionKeyからUserIDを取得
		UserID = database.GetUserID(CustomerSessionKey)
		//この場合はログイン中なので、もしCartSessionKeyを持っていたら、それを削除
		_, CustomerSessionKey = validation.CustomerSessionStart(c)
		database.CustomerLogIn(UserID, CustomerSessionKey)
	} else {
		log.Print("not logined")
		//cookieに保存

	}
	return UserID
}

type CartRequestPayload struct {
	ItemID   string `json:"ItemID"`
	Quantity int    `json:"Quantity"`
}
type CartRequestPayloads []CartRequestPayload

func Cart(c *gin.Context) {
	UserID := GetDatafromSessionKey(c)
	if UserID != "" {
		CartContents := new(CartRequestPayloads)
		err := c.BindJSON(&CartContents)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(CartContents)
		if CartContents.inspectCart() != 0 {
			jsonBytes, err := json.Marshal(CartContents)
			if err != nil {
				log.Fatal(err)
			}
			database.CartSave(UserID, string(jsonBytes))
		}
		c.JSON(http.StatusOK, gin.H{"Cart": CartContents})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未ログインです。"})
	}
}
