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
	r.POST("/SignUp", handler.Temporary_SignUp)
	// 本登録
	r.POST("/Registration", handler.SignUp)
	// 登録内容の修正
	r.POST("/Modify", handler.Modify_Customer)
	// アカウントの削除
	r.DELETE("/DeleteCustomer", handler.Delete_Customer)
	// ログイン状態の継続
	r.POST("/SessionContinue", handler.Continue_LogIn)
	// ログアウト cookie clear
	r.POST("/SessionEnd", handler.LogOut)

	// カート機能
	// 商品の登録・修正・削除
	r.POST("/PostCart", handler.Post_Cart)
	// カートの取得
	r.GET("/GetCart", handler.Get_Cart)

	// 購入処理
	r.POST("/Transaction", handler.BuyItem)

	//商品API
	r.GET("/item/top", handler.Top)
	r.GET("/item/all", handler.ALL)
	r.GET("/item/details", handler.Item_Details)
	r.GET("/item/category/:category", handler.Category)
	r.GET("/item/color/:color", handler.Color)

	r.Run(":8080") // 0.0.0.0:8080 でサーバーを立てます。
}
