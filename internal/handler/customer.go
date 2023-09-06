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
		//新しいアカウントの構造体を作成
		newCustomer := new(models.CustomerRequestPayload)

		newCustomer.UID = user.UID
		newCustomer.Email = user.Email
		log.Printf(newCustomer.UID, newCustomer.Email)
		_, NewSessionKey := validation.SessionStart(c)
		log.Print(NewSessionKey)

		//アカウント登録
		Cart_List := new(database.Cart_List)
		Cart_List.Session_Key = validation.Get_Cart_Session(c)
		if Cart_List.Session_Key == "new" {
			log.Print("don't have sessionKey")
			Cart_List.Cart_ID = validation.GetUUID()
		} else {
			err := Cart_List.Get_Cart_ID_from_SessionKey()
			if err != nil {
				log.Fatal(err)
			}
			database.Delete_Cart_List(Cart_List.Cart_ID)
		}
		Cart_List.Session_Key = validation.GetUUID()
		Cart_List.Create_Cart_List()
		log.Print("Cart_ID: ", Cart_List.Cart_ID)
		res := database.SignUp_Customer(*newCustomer, NewSessionKey, Cart_List.Cart_ID)
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
			Cart_ID, err := database.Get_Cart_ID(user.UID)
			if err != nil {
				log.Fatal(err)
			}
			log.Print("Cart_ID:", Cart_ID)
			if Cart_ID == "" {
				Cart_List := new(database.Cart_List)
				Cart_List.Session_Key = validation.Get_Cart_Session(c)
				if Cart_List.Session_Key == "new" {
					log.Print("don't have sessionKey")
					Cart_List.Cart_ID = validation.GetUUID()
				} else {
					err := Cart_List.Get_Cart_ID_from_SessionKey()
					if err != nil {
						log.Fatal(err)
					}
					database.Delete_Cart_List(Cart_List.Cart_ID)
				}
				Cart_List.Session_Key = validation.GetUUID()
				Cart_List.Create_Cart_List()
				database.Set_Cart_ID(user.UID, Cart_List.Cart_ID)
				validation.CartSessionEnd(c)
				log.Print("Cart_ID: ", Cart_List.Cart_ID)
			}
		} else {
			c.JSON(http.StatusOK, user)
		}
		if user.Email_Verified {
			err := database.Email_Verified(user.UID)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Print(OldSessionKey)
		log.Print(NewSessionKey)
		email, err := database.Get_Email(user.UID)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(email)
		if email != user.Email {
			database.Change_Email(user.UID, user.Email)
		}

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
		log.Print("UID : ", UID)
		Customer := new(database.Customer)
		Customer.LogIn_Customer(UID, NewSessionKey)

		//c.JSON(http.StatusOK, "SuccessFully Logined!!")
	}
	log.Print("OldSessionKey : ", OldSessionKey)
	log.Print("NewSessionKey : ", NewSessionKey)
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

func LogOut(c *gin.Context) {
	//ログアウト処理
	OldSessionKey := validation.SessionEnd(c)
	database.Invalid(OldSessionKey)
	c.JSON(http.StatusOK, "SuccessFully Logouted!!")
}
func Delete_Customer(c *gin.Context) {
	//アカウントの削除
	user := new(validation.UserReqPayload)
	err := c.BindJSON(&user)
	if err != nil {
		log.Print(err)
	}
	if user.Verify(c) { //認証
		database.Delete_Customer(user.UID)
		database.Delete_Session(user.UID)
		c.JSON(http.StatusOK, gin.H{"message": "アカウントを削除しました。"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}
}
func Change_Email(c *gin.Context) {
	//アカウントの変更
	user := new(validation.UserReqPayload)
	err := c.BindJSON(&user)
	if err != nil {
		log.Print(err)
	}
	if user.Verify(c) { //認証
		database.Change_Email(user.UID, user.Email)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "ログインできませんでした。"})
	}
}
