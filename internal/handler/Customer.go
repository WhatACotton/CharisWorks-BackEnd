package handler

import (
	"log"
	"net/http"

	"github.com/WhatACotton/go-backend-test/internal/database"
	"github.com/WhatACotton/go-backend-test/validation"
	"github.com/gin-gonic/gin"
)

// アカウント作成
func SignUp(c *gin.Context) {
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
		validation.LoginLogging(CustomerReqPayload.UserID + "created")
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

// 初回登録
func Register(c *gin.Context) {
	//登録情報変更処理
	//bodyにアカウントの詳細情報が入っている。
	_, UserID := GetDatafromSessionKey(c)
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
		log.Print("UserID : ", UserReqPayload.UserID)
		//Eamil認証
		if UserReqPayload.EmailVerified {
			database.CustomerEmailVerified(1, UserReqPayload.UserID)
		} else {
			database.CustomerEmailVerified(0, UserReqPayload.UserID)
		}
		//Emailが変更されているかどうかの処理
		Email := database.GetEmail(UserReqPayload.UserID)
		log.Print(Email)
		if Email != UserReqPayload.Email {
			database.CustomerChangeEmail(UserReqPayload.UserID, UserReqPayload.Email)
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
	_, UserID := GetDatafromSessionKey(c)
	if UserID != "" {
		Customer := new(database.Customer)
		Customer.CustomerGet(UserID)
		c.JSON(http.StatusOK, Customer)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未ログインです。"})
	}
}

// 顧客情報の更新
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

// ログアウト
func LogOut(c *gin.Context) {
	_, UserID := GetDatafromSessionKey(c)
	//c.JSON(http.StatusOK, "SuccessFully Logouted!!")
	validation.LoginLogging(UserID + " logouted")
	//ログアウト処理
	OldSessionKey := validation.CustomerSessionEnd(c)
	log.Print("SessionKey was :", OldSessionKey)
}

// アカウントの削除
func DeleteCustomer(c *gin.Context) {
	//アカウントの削除
	_, UserID := GetDatafromSessionKey(c)
	database.CustomerDelete(UserID)
	validation.LoginLogging(UserID + "deleted")

	c.JSON(http.StatusOK, gin.H{"message": "アカウントを削除しました。"})
}
