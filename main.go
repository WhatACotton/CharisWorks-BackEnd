package main

import (
	"io"
	"os"

	"github.com/WhatACotton/go-backend-test/internal/database"
	"github.com/WhatACotton/go-backend-test/internal/handler"
	"github.com/WhatACotton/go-backend-test/validation"
	"github.com/gin-gonic/gin"
)

func main() {

	// Logging to a file.
	f, _ := os.Create("log/gin" + database.GetDate() + ".log")
	gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	// ログの出力
	r := gin.Default()
	validation.CORS(r)
	validation.SessionConfig(r)
	database.TestSQL()
	// アカウント管理
	//ログイン
	r.POST("/go/Login", handler.LogIn)
	// 仮登録
	r.POST("/go/SignUp", handler.SignUp)
	// 本登録
	r.POST("/go/Registration", handler.Register)
	// 登録内容の修正
	r.POST("/go/Modify", handler.ModifyCustomer)
	// アカウントの削除
	r.DELETE("/go/DeleteCustomer", handler.DeleteCustomer)
	// ログアウト cookie clear
	r.POST("/go/SessionEnd", handler.LogOut)
	// アカウント情報の取得
	r.GET("/go/GetCustomer", handler.GetCustomer)
	// 購入履歴の取得
	r.GET("/go/GetTransactions", handler.GetTransaction)

	// カート機能
	// 商品の登録・修正・削除
	r.POST("/go/PostCart", handler.PostCart)
	// カートの取得
	r.GET("/go/GetCart", handler.GetCart)

	// 購入処理
	r.POST("/go/Transaction", handler.BuyItem)
	r.POST("/go/stripe", handler.Webhook)
	//商品API
	r.GET("/go/item/top", handler.Top)
	r.GET("/go/item/all", handler.ALL)
	r.GET("/go/item/details", handler.ItemDetails)
	r.GET("/go/item/category/:category", handler.Category)
	r.GET("/go/item/color/:color", handler.Color)
	r.GET("/go/item/maker/:MakerName", handler.ItemMakerGet)

	r.POST("/go/Maker/AccountCreate", handler.MakerStripeAccountCreate)
	r.POST("/go/Maker/ItemMainCreate", handler.MakerItemMainCreate)
	r.POST("/go/Maker/ItemDetailCreate", handler.MakerItemDetailCreate)
	r.POST("/go/Maker/ItemDetailModyfy", handler.MakerItemDetailModyfy)
	r.GET("/go/Maker/Details", handler.MakerDetailsGet)
	r.POST("/go/Maker/DetailsRegister", handler.MakerAccountRegister)
	r.GET("/go/Maker/GetItem", handler.MakerGetItem)
	r.Run(":8080") // 0.0.0.0:8080 でサーバーを立てます。
}

//TODO
//購入履歴APIの実装。Transaction_Listも複数のレコードを持つことを考慮しないといけない。
