package handler

import (
	"net/http"
	"unify/internal/auth"
	"unify/internal/funcs"

	"unify/validation"

	"github.com/gin-gonic/gin"
)

func TemporaryRegistration(c *gin.Context) {
	//signup処理
	//仮登録を行う。ここでの登録内容はUIDと作成日時だけ。
	user := new(validation.User)
	if user.Verify(c) { //認証
		funcs.SignUpCustomer(*user, c)
	}
}

func UserRegistration(c *gin.Context) {
	//本登録処理
	//本登録を行う。bodyにアカウントの詳細情報が入っている。
	user := new(validation.User)
	if user.Verify(c) { //認証
		funcs.RegisterCustomer(*user, c)
	}
}

func GetItem(c *gin.Context) {
	id := c.Query("id")
	funcs.GetItem(c, id)

}
func GetItemList(c *gin.Context) {
	funcs.GetItemList(c)
}

func Transaction(c *gin.Context) {
	requestMethod := http.MethodGet
	switch request := requestMethod; request {
	case "GET":
		auth.GetTransaction(c)
	case "POST":
		auth.PostTransaction(c)
	case "PATCH":
		auth.PatchTransaction(c)
	case "DELETE":
		auth.DeleteTransaction(c)
	}
}
