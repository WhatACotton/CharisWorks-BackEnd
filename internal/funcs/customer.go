package funcs

import (
	"net/http"

	"unify/internal/database"
	"unify/internal/models"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

var customers []models.Customer

func SignUpCustomer(usr validation.User, c *gin.Context) {
	//アカウント登録処理

	//新しいアカウントの構造体を作成
	var newCustomer models.CustomerRequestPayload

	newCustomer.UID = usr.Userdata.UID
	newCustomer.Email = usr.Userdata.Email

	//アカウント登録
	res := database.SignUpCustomer(newCustomer)

	c.IndentedJSON(http.StatusOK, res)
}

func RegisterCustomer(usr validation.User, c *gin.Context) {
	//アカウント本登録処理
	//2回構造体を作るのは冗長かも知れないが、bindしている以上、
	//インジェクションされて予期しない場所が変更される可能性がある。
	h := new(models.CustomerRegisterPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}

	database.RegisterCustomer(usr, *h)
}

func ModifyCustomer(usr validation.User, c *gin.Context) {
	//アカウント修正処理
	h := new(models.CustomerRegisterPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	database.ModifyCustomer(usr, *h)
}

func DeleteCustomer(usr validation.User, c *gin.Context) {
	database.Delete("user", "uid", usr.Userdata.UID)
	if GetCustomer(c) == 404 {
		c.JSON(http.StatusOK, "Customer was successfully deleted.")
	} else {
		c.JSON(http.StatusBadRequest, "Could not delete customer")
	}
}

func LogIn(usr validation.User, c *gin.Context) {
	c.JSON(http.StatusOK, database.LogInCustomer(usr.Userdata.UID))
}

func GetCustomer(c *gin.Context) (err int) {
	uid := c.Query("uid")
	var response = database.GetCustomer(uid)
	if response.UID == "" {
		return http.StatusNotFound
	}
	c.JSON(http.StatusOK, response)
	return http.StatusOK
}
