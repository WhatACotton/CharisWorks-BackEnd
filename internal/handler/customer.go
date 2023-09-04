package handler

import (
	"log"
	"net/http"
	"unify/internal/database"
	"unify/internal/models"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func Temporary_SignUp(c *gin.Context) {
	//signup処理
	//仮登録を行う。ここでの登録内容はUIDと作成日時だけ。
	user := new(validation.UserReqPayload)
	err := c.BindJSON(&user)
	if err != nil {
		log.Print(err)
	}

	if user.Verify(c) { //認証
		log.Printf(user.Email)
		_, NewSessionKey := validation.SessionStart(c)

		log.Print(NewSessionKey)
		//新しいアカウントの構造体を作成
		newCustomer := new(models.CustomerRequestPayload)

		newCustomer.UID = user.UID
		newCustomer.Email = user.Email
		log.Printf(newCustomer.UID, newCustomer.Email)
		//アカウント登録
		res := database.SignUp_Customer(*newCustomer, NewSessionKey)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func SignUp(c *gin.Context) {
	//本登録処理
	//本登録を行う。bodyにアカウントの詳細情報が入っている。
	user := new(validation.UserReqPayload)
	err := c.BindJSON(&user)
	if err != nil {
		log.Print(err)
	}
	log.Print(c)
	if user.Verify(c) { //認証
		//アカウント本登録処理
		//2回構造体を作るのは冗長かも知れないが、bindしている以上、
		//インジェクションされて予期しない場所が変更される可能性がある。
		h := new(models.CustomerRegisterPayload)
		if err := c.BindJSON(&h); err != nil {
			return
		}
		database.Register_Customer(*user, *h)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func LogIn(c *gin.Context) {
	//LogIn処理
	user := new(validation.UserReqPayload)
	err := c.BindJSON(&user)
	if err != nil {
		log.Print(err)
	}
	Customer := new(database.Customer)
	if user.Verify(c) { //認証
		OldSessionKey, NewSessionKey := validation.SessionStart(c)
		Customer.LogIn_Customer(user.UID, NewSessionKey)
		if OldSessionKey == "new" {
			c.JSON(http.StatusOK, "SuccessFully Logined!!")
		} else {
			c.JSON(http.StatusOK, user)
		}
		log.Print(Customer)
		log.Print(OldSessionKey)
		log.Print(NewSessionKey)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}

}

func Continue_LogIn(c *gin.Context) {
	OldSessionKey, NewSessionKey := validation.SessionStart(c)
	if OldSessionKey == "new" {
		c.JSON(http.StatusOK, "未ログインです")
	} else {
		UID, err := database.Get_UID(OldSessionKey)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(UID)
		Customer := new(database.Customer)
		Customer.LogIn_Customer(UID, NewSessionKey)
		log.Print(Customer)

		c.JSON(http.StatusOK, "SuccessFully Logined!!")
	}
	log.Print(OldSessionKey)
	log.Print(NewSessionKey)
}

func Modify_Customer(c *gin.Context) {
	//登録情報変更処理
	//bodyにアカウントの詳細情報が入っている。
	user := new(validation.UserReqPayload)
	err := c.BindJSON(&user)
	if err != nil {
		log.Print(err)
	}
	if user.Verify(c) { //認証
		//アカウント修正処理
		h := new(models.CustomerRegisterPayload)
		if err := c.BindJSON(&h); err != nil {
			return
		}
		database.Modify_Customer(*user, *h)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func Delete_Customer(c *gin.Context) {
	//アカウントの削除
	user := new(validation.UserReqPayload)
	err := c.BindJSON(&user)
	if err != nil {
		log.Print(err)
	}
	if user.Verify(c) { //認証
		user.DeleteCustomer(c)
		database.Delete_Customer(user.UID)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}
}
