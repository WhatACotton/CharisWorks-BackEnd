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

	//アカウント作成日時を取得
	var CreatedDate = database.GetDate()
	//新しいアカウントの構造体を作成
	var newCustomer models.CustomerRequestPayload

	newCustomer.UID = usr.Userdata.UID
	newCustomer.Email = usr.Userdata.Email
	newCustomer.CreatedDate = CreatedDate

	//アカウント登録
	res := database.SignUpCustomer(newCustomer)

	c.IndentedJSON(http.StatusOK, res)
}

func RegisterCustomer(usr validation.User, c *gin.Context) {
	//アカウント本登録処理
	//2回構造体を作るのは冗長かも知れないが、bindしている以上、
	//インジェクションされて予期しない場所が変更される可能性がある。
	h := new(models.Customer)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	//本登録の日時を取得
	var RegisteredDate = database.GetDate()
	var RegisterCustomer models.Customer

	RegisterCustomer.UID = usr.Userdata.UID
	RegisterCustomer.RegisteredDate = RegisteredDate
	RegisterCustomer.Name = h.Name
	RegisterCustomer.Address = h.Address
	RegisterCustomer.Email = h.Email
	RegisterCustomer.PhoneNumber = h.PhoneNumber

	database.RegisterCustomer(RegisterCustomer)
}

func ModifyCustomer(usr validation.User, c *gin.Context) {
	//アカウント修正処理
	//2回構造体を作るのは冗長かも知れないが、bindしている以上、
	//インジェクションされて予期しない場所が変更される可能性がある。
	h := new(models.Customer)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	//本登録の日時を取得
	var ModifiedDate = database.GetDate()
	var ModifyCustomer models.Customer

	ModifyCustomer.UID = usr.Userdata.UID
	ModifyCustomer.ModifiedDate = ModifiedDate
	ModifyCustomer.Name = h.Name
	ModifyCustomer.Address = h.Address
	ModifyCustomer.Email = h.Email
	ModifyCustomer.PhoneNumber = h.PhoneNumber

	database.ModifyCustomer(ModifyCustomer)
}

func PatchCustomer(c *gin.Context) {
	h := new(database.PatchRequestPayload)
	if err := c.BindJSON(&h); err != nil {
		return
	}
	h.Patch("user", "uid")
	GetCustomer(c)
}

func GetCustomer(c *gin.Context) {
	uid := c.Query("uid")
	var response = database.GetCustomer(uid)
	if response.UID == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "customer not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func DeleteCustomer(c *gin.Context) {
	uid := c.Query("uid")
	database.Delete("user", "uid", uid)
	GetCustomer(c)
}
