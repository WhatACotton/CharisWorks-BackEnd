package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/WhatACotton/go-backend-test/cashing"
	"github.com/WhatACotton/go-backend-test/internal/database"
	"github.com/gin-gonic/gin"
)

// StripeAccountの作成
func MakerStripeAccountCreate(c *gin.Context) {
	UserID := GetDatafromSessionKey(c)
	role := database.CustomerGetStripeAccountID(UserID)
	if role == "preSeller" {
		email := database.GetEmail(UserID)
		AccountID, URL := cashing.CreateStripeAccount(email)
		database.MakerAccountCreate(UserID, AccountID)
		c.JSON(http.StatusOK, gin.H{"message": "アカウント作成のリンクが作成されました。", "URL": URL})
	} else {
		if role == "seller" {
			c.JSON(http.StatusOK, gin.H{"message": "アカウントが作成されています。"})
		}
	}
}

// 商品の登録(Status,Price,Stock,ItemName)
func MakerItemMainCreate(c *gin.Context) {
	log.Print("Creating ItemMain...")
	StripeAccountID := makerStripeAccountIDGet(c)
	if StripeAccountID != "" {
		i := new(database.ItemMain)
		if err := c.BindJSON(&i); err != nil {
			log.Print(err)
		}
		log.Print("ItemMain:", i)
		if i.Name != "" || i.Price != 0 || i.Stock != 0 || i.Status != "" || i.Description != "" || i.Color != "" || i.Series != "" || i.Size != "" {
			database.ItemMainCreate(*i, StripeAccountID)
		}
	}
}

// 商品の説明などの登録(Description,Color,Series,Size)
func MakerItemDetailCreate(c *gin.Context) {
	StripeAccountID := makerStripeAccountIDGet(c)
	if StripeAccountID != "" {
		i := new(database.ItemDetail)
		if err := c.BindJSON(&i); err != nil {
			log.Print(err)
		}
		log.Print("ItemDetail:", i)
		MadeBy := database.MakerStripeAccountIDGet(StripeAccountID)
		if i.ItemID != "" || i.Description != "" || i.Color != "" || i.Series != "" || i.Size != "" {
			database.ItemDetailCreate(*i, MadeBy)
		}
	}
}

// 商品の説明などの編集(Description,Color,Series,Size)
func MakerItemDetailModyfy(c *gin.Context) {
	StripeAccountID := makerStripeAccountIDGet(c)
	if StripeAccountID != "" {
		i := new(database.ItemDetail)
		if err := c.BindJSON(&i); err != nil {
			log.Print(err)
		}
		MadeBy := database.MakerStripeAccountIDGet(StripeAccountID)
		if i.ItemID == "" || i.Description == "" || i.Color == "" || i.Series == "" || i.Size == "" {
			database.ItemDetailCreate(*i, MadeBy)
		}
	}
}

func MakerUploadImage(c *gin.Context) {
	//画像を取得して、特定の場所に保存
	//画像のインスペクト・サイズの確認・
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保存先のファイルパスを指定します
	// ここではカレントディレクトリの"uploads"ディレクトリに保存します
	dstPath := "uploads/" + file.Filename

	// アップロードされたファイルを保存します
	if err := c.SaveUploadedFile(file, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("ファイル %s がアップロードされました", file.Filename)})

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		log.Fatal("ディレクトリの作成に失敗しました:", err)
	}
}

// StripeAccountの取得
func makerStripeAccountIDGet(c *gin.Context) (StripeAccountID string) {
	UserID := GetDatafromSessionKey(c)
	role := database.CustomerGetStripeAccountID(UserID)
	if role == "Seller" || role == "Admin" {
		StripeAccountID = database.MakerGetStripeID(UserID)
		log.Print("StripeAccountID:", StripeAccountID)
		return StripeAccountID
	} else {
		if role == "preSeller" {
			return "preSeller"
		} else {
			return ""
		}
	}
}

func MakerDetailsGet(c *gin.Context) {
	StripeAccountID := makerStripeAccountIDGet(c)
	if StripeAccountID != "" {
		if StripeAccountID == "preSeller" {
			maker := new(database.Maker)
			maker.MakerName = "preSeller"
			c.JSON(http.StatusOK, gin.H{"Maker": maker})
		} else {
			Maker := new(database.Maker)
			Maker.StripeAccountID = StripeAccountID
			Maker.MakerDetailsGet()
			c.JSON(http.StatusOK, gin.H{"Maker": Maker})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"Maker": new(database.Maker)})
	}
}
func MakerAccountRegister(c *gin.Context) {
	StripeAccountID := makerStripeAccountIDGet(c)

	if StripeAccountID != "" {
		m := new(database.Maker)
		if err := c.BindJSON(&m); err != nil {
			log.Print(err)
		}
		m.StripeAccountID = StripeAccountID
		m.MakerAccountModyfy()
		m.MakerDetailsGet()
		c.JSON(http.StatusOK, gin.H{"Maker": m})
	}
}
func MakerGetItem(c *gin.Context) {
	UserID := GetDatafromSessionKey(c)
	StribeAccountID := database.MakerNameGet(UserID)
	if StribeAccountID != "" {
		Items := database.ItemGetMaker(StribeAccountID)
		c.JSON(http.StatusOK, Items)
	}
}
