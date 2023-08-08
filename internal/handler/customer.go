package handler

import (
	"log"
	"net/http"
	"unify/internal/database"
	"unify/internal/models"
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
		OldCartSessionKey, NewSessionKey := validation.SessionStart(c)
		if OldCartSessionKey != "new" {
			database.SignUpCart(OldCartSessionKey, uid)
		}
		log.Printf(NewSessionKey)
		//新しいアカウントの構造体を作成
		newCustomer := new(models.CustomerRequestPayload)

		newCustomer.UID = user.Userdata.UID
		newCustomer.Email = user.Userdata.Email
		log.Printf(newCustomer.UID, newCustomer.Email)
		//アカウント登録
		res := database.SignUpCustomer(*newCustomer, NewSessionKey)
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
		database.RegisterCustomer(*user, *h)
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
			database.LogInCustomer(user.Userdata.UID, NewSessionKey)

			c.JSON(http.StatusOK, "SuccessFully Logined!!")
		} else {
			database.LogInCustomer(user.Userdata.UID, NewSessionKey)
			database.CartInvalid(OldSessionKey)

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
		//アカウント修正処理
		h := new(models.CustomerRegisterPayload)
		if err := c.BindJSON(&h); err != nil {
			return
		}
		database.ModifyCustomer(*user, *h)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func DeleteCustomer(c *gin.Context) {
	//アカウントの削除
	user := new(validation.User)
	uid := c.Query("uid")
	if user.Verify(c, uid) { //認証
		database.DeleteCustomer(user.Userdata.UID)
		user.DeleteCustomer(c, user.Userdata.UID)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}
	user.DeleteCustomer(c, uid)
}
func LogOut(c *gin.Context) {
	//ログアウト
	uid := c.Query("uid")
	OldSessionKey, _ := validation.SessionStart(c)
	if database.VerifyCustomer(uid, OldSessionKey) {
		database.Invalid(OldSessionKey)
		c.JSON(http.StatusOK, "SuccessFully Logouted!!")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})
	}
}
