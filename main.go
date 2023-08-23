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
	r.GET("/Login", handler.LogIn)
	// 仮登録
	r.POST("/SignUp", handler.Temporary_SignUp)
	// ログアウト
	r.GET("/Logout", handler.Log_Out)
	// 本登録
	r.POST("/Registration", handler.SignUp)
	// 登録内容の修正
	r.POST("/Modify", handler.Modify_Customer)
	// アカウントの削除
	r.DELETE("/DeleteCustomer", handler.Delete_Customer)
	// ログイン状態の継続
	r.GET("/SessionStart", handler.Continue_LogIn)

	// カート機能
	// 商品の登録・修正・削除
	r.POST("/PostCart", handler.Post_Cart)
	// カートの取得
	r.GET("/GetCart", handler.Get_Cart)

	// 購入処理
	r.POST("/Transaction", handler.BuyItem)

	r.Run(":8080") // 0.0.0.0:8080 でサーバーを立てます。
}
