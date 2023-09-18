package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"unify/cashing"
	"unify/internal/database"

	"github.com/gin-gonic/gin"
)

func CreateStripeAccount(c *gin.Context) {
	_, UserID := GetDatafromSessionKey(c)
	StripeAccountID := database.GetStripeAccountID(UserID)
	if StripeAccountID != "allow" {
		email := database.GetEmail(UserID)
		AccountID, URL := cashing.CreateStripeAccount(email)
		database.CustomerCreateStripeAccount(UserID, AccountID)

		c.JSON(http.StatusOK, gin.H{"message": "アカウントが作成されました。", "URL": URL})
	} else {
		if StripeAccountID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "出品者登録が完了していません。", "errcode": "401"})
		}
		c.JSON(http.StatusOK, gin.H{"message": "アカウントが作成されています。"})
	}
}

func isMaker(c *gin.Context) bool {
	_, UserID := GetDatafromSessionKey(c)
	StripeAccountID := database.GetStripeAccountID(UserID)
	if StripeAccountID != "allow" && StripeAccountID != "" {
		return true
	}
	return false
}
func PostItem(c *gin.Context) {
	if isMaker(c) {
		postItem(c)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "出品者登録が完了していません。", "errcode": "401"})
	}
}
func postItem(c *gin.Context) {

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
