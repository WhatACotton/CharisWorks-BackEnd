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

	r.POST("/go/Cart", handler.Cart)

	r.GET("/go/GetCart", handler.GetCart)
	// 購入処理
	r.POST("/go/Transaction", func(ctx *gin.Context) {handler.BuyItem(ctx,handler.GetUserIDTestimpl{})})
	r.POST("/go/stripe", handler.Webhook)

	//商品API
	r.GET("/go/item/top", func(ctx *gin.Context){handler.Top(ctx,handler.TopItemimpl{})})
	r.GET("/go/item/all", handler.ALL)
	r.GET("/go/item/details/:ItemID", handler.ItemDetails)
	r.GET("/go/item/category/:category", handler.Category)
	r.GET("/go/item/color/:color", handler.Color)
	r.GET("/go/item/maker/:MakerName", handler.ItemMakerGet)
	r.GET("/go/item/maker/id/:StripeAccountID", handler.ItemMakerIDGet)
	r.POST("/go/item/CartDetails", handler.CartDetails)

	//MakerAPI
	r.POST("/go/Maker/AccountCreate", handler.MakerStripeAccountCreate)
	r.POST("/go/Maker/ItemMainCreate", handler.MakerItemMainCreate)
	r.POST("/go/Maker/ItemDetailCreate", handler.MakerItemDetailCreate)
	r.POST("/go/Maker/ItemDetailModyfy", handler.MakerItemDetailModyfy)
	r.GET("/go/Maker/Details", handler.MakerDetailsGet)
	r.POST("/go/Maker/DetailsRegister", handler.MakerAccountRegister)
	r.GET("/go/Maker/GetItem", handler.MakerGetItem)
	r.Run(os.Getenv("PORT")) // 0.0.0.0:8080 でサーバーを立てます。
}

//TODO
//購入履歴APIの実装。Transaction_Listも複数のレコードを持つことを考慮しないといけない。
