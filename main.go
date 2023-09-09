package main

import (
	"unify/internal/database"
	"unify/internal/handler"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	validation.CORS(r)
	validation.SessionConfig(r)
	database.TestSQL()
	// アカウント管理
	//ログイン
	r.POST("/Login", handler.LogIn)
	// 仮登録
	r.POST("/SignUp", handler.TemporarySignUp)
	// 本登録
	r.POST("/Registration", handler.SignUp)
	// 登録内容の修正
	r.POST("/Modify", handler.ModifyCustomer)
	// アカウントの削除
	r.DELETE("/DeleteCustomer", handler.DeleteCustomer)
	// ログアウト cookie clear
	r.POST("/SessionEnd", handler.LogOut)
	// アカウント情報の取得
	r.GET("/GetCustomer", handler.GetCustomer)
	// 購入履歴の取得
	r.GET("/GetTransactions", handler.GetTransaction)
	// カート機能
	// 商品の登録・修正・削除
	r.POST("/PostCart", handler.PostCart)
	// カートの取得
	r.GET("/GetCart", handler.GetCart)

	// 購入処理
	r.POST("/Transaction", handler.BuyItem)
	r.POST("/stripe", handler.Webhook)
	//商品API
	r.GET("/item/top", handler.Top)
	r.GET("/item/all", handler.ALL)
	r.GET("/item/details", handler.ItemDetails)
	r.GET("/item/category/:category", handler.Category)
	r.GET("/item/color/:color", handler.Color)

	r.Run(":8080") // 0.0.0.0:8080 でサーバーを立てます。
}

//TODO
//購入履歴APIの実装。Transaction_Listも複数のレコードを持つことを考慮しないといけない。
