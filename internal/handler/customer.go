package handler

import (
	"log"
	"net/http"
	"unify/internal/database"
	"unify/internal/funcs"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func TemporarySignUp(c *gin.Context) {
	//signup処理
	//仮登録を行う。ここでの登録内容はUIDと作成日時だけ。
	user := new(validation.User)
	uid := c.Query("uid")

	if user.Verify(c, uid) { //認証
		log.Printf(user.Userdata.Email)
		_, NewSessionKey := validation.SessionStart(c)
		log.Printf(NewSessionKey)
		funcs.SignUpCustomer(*user, NewSessionKey, c)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func SignUp(c *gin.Context) {
	//本登録処理
	//本登録を行う。bodyにアカウントの詳細情報が入っている。
	user := new(validation.User)
	uid := c.Query("uid")
	if user.Verify(c, uid) { //認証
		funcs.RegisterCustomer(*user, c)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func LogIn(c *gin.Context) {
	//LogIn処理
	user := new(validation.User)
	uid := c.Query("uid")

	if user.Verify(c, uid) { //認証
		log.Printf(user.Userdata.Email)
		OldSessionKey, NewSessionKey := validation.SessionStart(c)
		log.Printf(OldSessionKey)
		log.Printf(NewSessionKey)
		if OldSessionKey == "new" {
			funcs.NewLogIn(*user, NewSessionKey)
			c.JSON(http.StatusOK, "SuccessFully Logined!!")
		} else {
			funcs.StoredLogIn(*user, OldSessionKey, NewSessionKey)
			c.JSON(http.StatusOK, user)
		}

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}
}

func ContinueLogIn(c *gin.Context) {
	OldSessionKey, NewSessionKey := validation.SessionStart(c)
	uid := c.Query("uid")
	log.Println(uid)

	log.Printf(OldSessionKey)
	log.Printf(NewSessionKey)
	if OldSessionKey == "new" {
		c.JSON(http.StatusOK, "未ログインです")
	} else {
		if database.VerifyCustomer(uid, OldSessionKey) {
			database.LogInLog(uid, NewSessionKey)
			database.Invalid(OldSessionKey)
			c.JSON(http.StatusOK, "SuccessFully Logined!!")
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})

		}
	}
}

func ModifyCustomer(c *gin.Context) {
	//登録情報変更処理
	//bodyにアカウントの詳細情報が入っている。
	uid := c.Query("uid")
	user := new(validation.User)
	if user.Verify(c, uid) { //認証
		funcs.ModifyCustomer(*user, c)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func DeleteCustomer(c *gin.Context) {
	//アカウントの削除
	user := new(validation.User)
	uid := c.Query("uid")
	if user.Verify(c, uid) { //認証
		funcs.DeleteCustomer(*user, c)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}
}
