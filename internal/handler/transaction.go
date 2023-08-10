package handler

import (
	"net/http"
	"unify/internal/database"
	"unify/internal/models"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func BuyItem(c *gin.Context) {
	user := new(validation.User)
	UID := c.Query("uid")
	if user.Verify(c, UID) {
		//ここで購入処理
		OldCartSessionKey := validation.CartSessionEnd(c)
		database.CartInvalid(OldCartSessionKey)
		CartId := database.GetCartId(OldCartSessionKey)
		//ここからデータベースの処理
		Bill := new(models.Bill)
		InspectedCarts := new([]models.Cart)
		Carts := database.GetCart(CartId)
		//カートから購入可能な商品のみを抽出
		for _, Cart := range Carts {
			if Cart.Status == "Available" {
				*InspectedCarts = append(*InspectedCarts, Cart)
			}
		}
		//購入可能な商品のみを購入履歴に追加
		database.PostTransaction(*InspectedCarts)
		//初期化
		Transactions := new([]models.Transaction)
		Transaction := new(models.Transaction)
		TotalPrice := 0
		TotalCount := 0
		//購入履歴を作成
		for _, Cart := range Carts {
			Transaction = new(models.Transaction)
			Transaction.InfoId = Cart.InfoId
			Transaction.CartId = Cart.CartId
			Transaction.Quantity = Cart.Quantity
			*Transactions = append(*Transactions, *Transaction)
			Price := database.GetPrice(Cart.InfoId)
			TotalPrice += Price * Cart.Quantity
			TotalCount += Cart.Quantity
		}
		Bill.Transactions = *Transactions
		Bill.TotalPrice = TotalPrice
		Bill.TotalCount = TotalCount
		Bill.CartId = CartId
		Bill.UID = UID
		Bill.TransactionDate = database.GetDate()
		database.PostTransactionList(CartId, UID, Bill.TransactionDate)
		c.JSON(http.StatusOK, Bill)

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
	}
}
