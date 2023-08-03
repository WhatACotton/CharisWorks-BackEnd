package handler

import (
	"net/http"
	"unify/internal/funcs"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func TemporarySignUp(c *gin.Context) {
	//signup処理
	//仮登録を行う。ここでの登録内容はUIDと作成日時だけ。
	user := new(validation.User)
	if user.Verify(c) { //認証
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func SignUp(c *gin.Context) {
	//本登録処理
	//本登録を行う。bodyにアカウントの詳細情報が入っている。
	user := new(validation.User)
	if user.Verify(c) { //認証
		funcs.RegisterCustomer(*user, c)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func LogIn(c *gin.Context) {
	//LogIn処理
	user := new(validation.User)
	if user.Verify(c) { //認証
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}
}

func ModifyCustomer(c *gin.Context) {
	//登録情報変更処理
	//bodyにアカウントの詳細情報が入っている。
	user := new(validation.User)
	if user.Verify(c) { //認証
		funcs.ModifyCustomer(*user, c)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}

func DeleteCustomer(c *gin.Context) {
	//アカウントの削除
	user := new(validation.User)
	if user.Verify(c) { //認証
		funcs.DeleteCustomer(*user, c)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}
}
