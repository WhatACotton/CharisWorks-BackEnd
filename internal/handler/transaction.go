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
	Cart_List, UID := Get_Cart_ID(c)
	if UID != "" {
		Customer := new(database.Customer)
		Customer.Get_Customer(UID)
		log.Print("Customer:", Customer)
		if Customer.Register {
			Carts, err := database.Get_Cart_Info(Cart_List.Cart_ID)
			if err != nil {
				log.Fatal(err)
			}
			if inspect_Cart(Carts) {
				Total_Price := total_Price(Carts)
				stripe_info, err := cashing.Purchase(Total_Price)
				if err != nil {
					log.Fatal(err)
				}
				Transaction_ID := validation.GetUUID()
				Transaction_List := new(database.Transaction_List)
				Transaction_List.Construct_Transaction_List(Cart_List, *Customer, stripe_info, Transaction_ID)
				Transaction_List.Post_Transaction_List()
				Transaction := new(database.Transaction)
				for _, Cart := range Carts {
					log.Print("Cart:", Cart)
					Transaction.Construct_Transaction(Cart, Cart_List, Transaction_ID)
					Transaction.Post_Transaction()
				}
				c.JSON(http.StatusOK, gin.H{"message": "購入リンクが発行されました。", "url": stripe_info.URL})
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

func inspect_Cart(carts []database.Cart) bool {
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

func total_Price(Carts []database.Cart) (Total_Price int) {
	for _, Cart := range Carts {
		Total_Price += Cart.Price
	}
	return Total_Price
}

func Get_Transaction(c *gin.Context) {
	Customer_SessionKey := new(string)
	*Customer_SessionKey = validation.Customer_Get_SessionKey(c)
	Return_Transactions := new([]database.Transaction)
	UID, err := database.Get_UID(*Customer_SessionKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("UID:", UID)
	Transaction_Lists, err := database.Get_Transaction_Lists(UID)
	if err != nil {
		log.Fatal(err)
	}
	for _, Transaction_List := range Transaction_Lists {
		Transactions, err := Transaction_List.Get_Transactions()
		if err != nil {
			log.Fatal(err)
		}
		*Return_Transactions = append(*Return_Transactions, Transactions...)
	}
	log.Println("Transaction_Lists:", Transaction_Lists)
	c.JSON(http.StatusOK, gin.H{"Transaction_Lists": Transaction_Lists, "Transactions": *Return_Transactions})
}

//ログイン状態の確認
//email認証・本登録の確認

//商品・価格の取得・購入までの処理
//UIDからCart_ID,Name,Address,Email,Phone_Numberを取得
//カート処理
//Cart_IDからCartを取得
//CartからItem_IDを取得
//Item_IDからInfo_IDを取得
//Info_IDからPriceを取得
//Priceを合計し、stripeに渡す

//購入履歴処理
//Cart_IDとInfo_IDを紐付け、Transactionに追加
//Cart_ID,UID,TransactionDate,TotalPrice,Address,Name,Phone_NumberをTransaction_Listに追加

//購入後処理
//UIDに紐付けられたCart_IDを更新
//Cart.dbからCart_IDに紐付けられたCartを削除
//Cart_List.dbからCart_IDを削除
