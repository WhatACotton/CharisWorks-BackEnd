package handler

import (
	"log"
	"net/http"
	"unify/cashing"
	"unify/internal/database"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func BuyItem(c *gin.Context) {
	Cart, UID := GetCartID(c)
	if UID != "" {
		Customer := new(database.Customer)
		Customer.GetCustomer(UID)
		log.Print("Customer:", Customer)
		if Customer.Register {
			CartContents, err := database.GetCartInfo(Cart.CartID)
			if err != nil {
				log.Fatal(err)
			}
			if inspectCart(CartContents) {
				TotalPrice := totalPrice(CartContents)
				stripeInfo, err := cashing.Purchase(TotalPrice)
				if err != nil {
					log.Fatal(err)
				}
				TransactionID := validation.GetUUID()
				database.PostTransaction(Cart, *Customer, stripeInfo, TransactionID, CartContents)

				c.JSON(http.StatusOK, gin.H{"message": "購入リンクが発行されました。", "url": stripeInfo.URL})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "カートの中身に購入不可能な商品が含まれています。"})
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "本登録が完了していません。"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未ログインです。"})
	}

}

func inspectCart(carts database.CartContents) bool {
	if len(carts) == 0 {
		return false
	}
	for _, Cart := range carts {
		if Cart.Status != "Available" {
			return false
		}
	}
	return true
}

func totalPrice(Carts database.CartContents) (TotalPrice int) {
	for _, Cart := range Carts {
		TotalPrice += Cart.Price
	}
	return TotalPrice
}

func GetTransaction(c *gin.Context) {
	CustomerSessionKey := new(string)
	*CustomerSessionKey = validation.GetCustomerSessionKey(c)
	TransactionContentsList := new([]database.TransactionContents)
	UID, err := database.GetUID(*CustomerSessionKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("UID:", UID)
	Transactions, err := database.GetTransactions(UID)
	if err != nil {
		log.Fatal(err)
	}
	for _, Transaction := range Transactions {
		TransactionContents, err := Transaction.GetTransactionContents()
		if err != nil {
			log.Fatal(err)
		}
		*TransactionContentsList = append(*TransactionContentsList, TransactionContents)
	}
	c.JSON(http.StatusOK, gin.H{"TransactionLists": Transactions, "Transactions": *TransactionContentsList})
}

//ログイン状態の確認
//email認証・本登録の確認

//商品・価格の取得・購入までの処理
//UIDからCartID,Name,Address,Email,PhoneNumberを取得
//カート処理
//CartIDからCartを取得
//CartからItemIDを取得
//ItemIDからInfoIDを取得
//InfoIDからPriceを取得
//Priceを合計し、stripeに渡す

//購入履歴処理
//CartIDとInfoIDを紐付け、Transactionに追加
//CartID,UID,TransactionDate,TotalPrice,Address,Name,PhoneNumberをTransactionListに追加

//購入後処理
//UIDに紐付けられたCartIDを更新
//Cart.dbからCartIDに紐付けられたCartを削除
//CartList.dbからCartIDを削除
