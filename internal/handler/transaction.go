package handler

import (
	"log"
	"net/http"
	"unify/internal/database"
	"unify/internal/models"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func BuyItem(c *gin.Context) {
	user := new(validation.User)
	Cart_List := new(database.Cart_List)
	UID := c.Query("uid")
	OldSessionKey, NewSessionKey := validation.SessionStart(c)
	Customer := new(database.Customer)
	if OldSessionKey != "new" {
		if database.Verify_Customer(UID, OldSessionKey) {
			database.Invalid(OldSessionKey)
			c.JSON(http.StatusOK, "SuccessFully Logined!!")
			Customer.LogIn_Customer(UID, NewSessionKey)
			if Customer.Register {
				if user.Verify(c, UID) {
					//ここで購入処理
					Cart_List.Refresh_Cart_List()
					//ここからデータベースの処理
					Bill := new(models.Bill)
					InspectedCarts := new([]database.Cart)
					Carts, err := database.Get_Cart_Info(Cart_List.Cart_ID)
					if err != nil {
						log.Fatal(err)
					}
					//カートから購入可能な商品のみを抽出
					for _, Cart := range Carts {
						if Cart.Status == "Available" {
							*InspectedCarts = append(*InspectedCarts, Cart)
						}
					}
					if InspectedCarts == &Carts {
						//購入可能な商品のみを購入履歴に追加
						database.PostTransaction(*InspectedCarts, Cart_List.Cart_ID)
						//初期化
						Transactions := new([]models.Transaction)
						Transaction := new(models.Transaction)
						TotalPrice := 0
						TotalCount := 0
						//購入履歴を作成
						for _, Cart := range Carts {
							Transaction = new(models.Transaction)
							//Transaction.InfoId = Cart.InfoId
							Transaction.CartId = Cart_List.Cart_ID
							Transaction.Quantity = Cart.Quantity
							*Transactions = append(*Transactions, *Transaction)
							//Price := database.GetPrice(Cart.InfoId)
							//TotalPrice += Price * Cart.Quantity
							TotalCount += Cart.Quantity
						}
						Bill.Transactions = *Transactions
						Bill.TotalPrice = TotalPrice
						Bill.TotalCount = TotalCount
						Bill.CartId = Cart_List.Cart_ID
						Bill.UID = UID
						Bill.TransactionDate = database.GetDate()
						Bill.Address = Customer.Address
						Bill.Name = Customer.Name
						Bill.PhoneNumber = Customer.PhoneNumber
						Bill.Address = Customer.Address
						database.PostTransactionList(Cart_List.Cart_ID, UID)
						c.JSON(http.StatusOK, Bill)
					} else {
						c.JSON(http.StatusBadRequest, gin.H{"message": "カートに購入不可の商品が含まれています。"})
					}
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なアクセスです。"})
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "本登録が完了していません。"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "不正なアクセスです"})
		}
	} else {
		c.JSON(http.StatusOK, "未ログインです")
	}

}
