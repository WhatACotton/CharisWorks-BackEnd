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
	user := new(validation.User)
	uid := c.Query("uid")

	if user.Verify(c, uid) { //認証
		log.Printf(user.Userdata.Email)
		_, NewSessionKey := validation.SessionStart(c)

		log.Print(NewSessionKey)
		//新しいアカウントの構造体を作成
		newCustomer := new(models.CustomerRequestPayload)

		newCustomer.UID = user.Userdata.UID
		newCustomer.Email = user.Userdata.Email
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
	user := new(validation.User)
	uid := c.Query("uid")
	if user.Verify(c, uid) { //認証
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
	user := new(validation.User)
	uid := c.Query("uid")
	Customer := new(database.Customer)
	if user.Verify(c, uid) { //認証
		log.Printf(user.Userdata.Email)
		OldSessionKey, NewSessionKey := validation.SessionStart(c)
		Customer.LogIn_Customer(user.Userdata.UID, NewSessionKey)
		if OldSessionKey == "new" {
			c.JSON(http.StatusOK, "SuccessFully Logined!!")
		} else {
			database.Invalid(OldSessionKey)
			c.JSON(http.StatusOK, user)
		}
		log.Print(OldSessionKey)
		log.Print(NewSessionKey)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}

}

func Continue_LogIn(c *gin.Context) {
	OldSessionKey, NewSessionKey := validation.SessionStart(c)
	uid := c.Query("uid")
	if OldSessionKey == "new" {
		c.JSON(http.StatusOK, "未ログインです")
	} else {
		if database.Verify_Customer(uid, OldSessionKey) {
			database.LogIn_Log(uid, NewSessionKey)
			database.Invalid(OldSessionKey)
			c.JSON(http.StatusOK, "SuccessFully Logined!!")
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})

		}
	}
	log.Println(uid)
	log.Print(OldSessionKey)
	log.Print(NewSessionKey)
}

func Modify_Customer(c *gin.Context) {
	//登録情報変更処理
	//bodyにアカウントの詳細情報が入っている。
	uid := c.Query("uid")
	user := new(validation.User)
	if user.Verify(c, uid) { //認証
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
	Log_Out(c)
	//アカウントの削除
	user := new(validation.User)
	uid := c.Query("uid")
	if user.Verify(c, uid) { //認証
		user.DeleteCustomer(c, uid)
		database.Delete_Customer(user.Userdata.UID)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}
}

func Log_Out(c *gin.Context) {
	//ログアウト
	uid := c.Query("uid")
	OldSessionKey := validation.SessionEnd(c)
	database.Invalid(OldSessionKey)
	if database.Verify_Customer(uid, OldSessionKey) {
		c.JSON(http.StatusOK, "SuccessFully Loggedout!!")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})
	}
}
