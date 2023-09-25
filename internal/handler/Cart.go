package handler

import (
	"log"
	"net/http"

	"github.com/WhatACotton/go-backend-test/internal/database"
	"github.com/WhatACotton/go-backend-test/validation"
	"github.com/gin-gonic/gin"
)

// カートの追加・変更・削除
func PostCart(c *gin.Context) {
	Cart, _ := GetDatafromSessionKey(c)
	NewCartReq := new(database.CartContentRequestPayload)
	err := c.BindJSON(&NewCartReq)
	if err != nil {
		log.Print(err)
	}
	NewCartReq.Cart(Cart.CartID)
	Carts, err := database.GetCartContents(Cart.CartID)
	if err != nil {
		log.Print(err)
	}
	c.JSON(http.StatusOK, Carts)
}

// カートの取得
func GetCart(c *gin.Context) {
	Cart, _ := GetDatafromSessionKey(c)
	log.Print(Cart.CartID)
	CartContents, err := database.GetCartContents(Cart.CartID)
	if err != nil {
		log.Print(err)
	}
	if CartContents == nil {
		c.JSON(http.StatusOK, "There is no Cart")
	} else {
		c.JSON(http.StatusOK, CartContents)
	}
	log.Print(CartContents)
}

// セッションキーから、UserIDとCartを取得
// 同時にログイン中ならログイン状態の更新も行う
func GetDatafromSessionKey(c *gin.Context) (Cart database.Cart, UserID string) {
	log.Print("Getting CartID...")
	CustomerSessionKey := validation.GetCustomerSessionKey(c)
	CartSessionKey := validation.GetCartSessionKey(c)
	//まずはCartSessionKeyの取得情報からCartIDを取得しておく
	if CartSessionKey != "new" {
		log.Print("have CartSessionKey")
		Cart.SessionKey = CartSessionKey
		Cart.CartSessionListGetCartID()
	}
	//次にCustomerSessionKeyを持っているならログイン情報の取得
	if CustomerSessionKey != "new" {
		//ログインしている
		log.Print("logined")
		//CustomerSessionKeyからUserIDを取得
		UserID = database.GetUserID(CustomerSessionKey)
		//UserIDからCartIDを取得
		CartIDfromCustomer := database.GetCartID(UserID)
		//カートの中に商品が入っているかの処理
		CartContents, err := database.GetCartContents(CartIDfromCustomer)
		if err != nil {
			log.Print(err)
		}
		if CartContents != nil {
			log.Print("have CartContents from CustomerData")
			Cart.CartID = CartIDfromCustomer
		}
		//CartSessionKeyからカート情報を取得できず、更に顧客情報に登録されているカートにも商品が入っていなかったらCartIDを振り直す
		if Cart.CartID == "" {
			log.Print("don't have CartID in any place")
			Cart.CartID = validation.GetUUID()
		}
		//この場合はログイン中なので、もしCartSessionKeyを持っていたら、それを削除
		validation.CartSessionEnd(c)
		_, CustomerSessionKey = validation.CustomerSessionStart(c)
		database.CustomerLogIn(UserID, CustomerSessionKey)
		database.CustomerSetCartID(UserID, Cart.CartID)
	} else {
		//ログインしていない
		log.Print("not logined")
		if Cart.CartID == "" {
			//CartIDが取得失敗
			log.Print("don't have CartID in any place")
			Cart.CartID = validation.GetUUID()
		} else {
			//CartIDが取得できたので、そのCartIDに紐づくSessionKeyをリセット
			database.CartSessionListDelete(Cart.CartID)
		}
		//新しいSessionKeyの発行

		Cart.SessionKey = validation.GetUUID()
		//CaerSessionListへの登録
		Cart.CartSessionListCreate()
		log.Print("Cart with sesssion. SessionKey : ", Cart.SessionKey, " CartID : ", Cart.CartID)
		//cookieに保存
		validation.SetCartSessionKey(c, Cart.SessionKey)

	}
	log.Print("CartID:", Cart.CartID)
	return Cart, UserID
}
