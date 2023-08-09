package database

import (
	"log"
	"unify/internal/models"
)

func PostCart(req models.CartRequestPayload, CartId string) (carts []models.Cart) {
	Itemlist := GetItemList()
	if InspectItems(req.ItemId, Itemlist) {
		Carts := GetCart(CartId)
		if SearchCart(Carts, req.ItemId) {
			if req.Quantity == 0 {
				DeleteCart(CartId, req.ItemId)
			} else {
				UpdateCart(CartId, req)
			}
		} else {
			if req.Quantity != 0 {
				// データベースのハンドルを取得する
				db := ConnectSQL()

				// SQLの準備
				//UID,ItemId,Quantity

				ins, err := db.Prepare("INSERT INTO cart VALUES(?,?,?)")
				if err != nil {
					log.Fatal(err)
				}
				defer ins.Close()
				// SQLの実行
				_, err = ins.Exec(CartId, req.ItemId, req.Quantity)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		return GetCart(CartId)
	} else {
		return nil
	}

}

func InspectItems(ItemId string, Itemlist []models.Item) bool {
	for _, Item := range Itemlist {
		if ItemId == Item.ItemId {
			return true
		}
	}
	return false
}
func SearchCart(Carts []models.Cart, ItemId string) bool {
	for _, Cart := range Carts {
		if Cart.ItemId == ItemId {
			return true
		} else {
			return false
		}
	}
	return false
}

func GetCart(CartId string) (Carts []models.Cart) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	var Cart models.Cart
	// SQLの実行
	rows, err := db.Query("SELECT * FROM cart WHERE CartId = ?", CartId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Cart.CartId, &Cart.ItemId, &Cart.Quantity)
		if err != nil {
			panic(err.Error())
		}
		Carts = append(Carts, Cart)
	}

	return Carts
}

func LoginCart(SessionKey string, CartId string, UID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの準備
	//UID,ItemId,Quantity

	ins, err := db.Prepare("INSERT INTO cartlist VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(CartId, UID, SessionKey, GetDate(), true)
	if err != nil {
		log.Fatal(err)
	}
}
func SessionCart(SessionKey string, CartId string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの準備
	//UID,ItemId,Quantity

	ins, err := db.Prepare("INSERT INTO cartlist VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(CartId, nil, SessionKey, GetDate(), true)
	if err != nil {
		log.Fatal(err)
	}
}

func CartInvalid(SessionKey string) {
	log.Println("CartInvalid called")
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare("UPDATE cartlist SET Valid = false WHERE SessionKey = ?")
	if err != nil {
		log.Fatal(err)
	}
	// SQLの実行
	_, err = ins.Exec(SessionKey)
	defer ins.Close()

}
func VerifyCart(OldSessionKey string) bool {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT Valid FROM cartlist WHERE SessionKey = ?", OldSessionKey)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var valid bool
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&valid)

		if err != nil {
			panic(err.Error())
		}
	}
	return valid
}

func GetCartId(OldSessionKey string) (CartId string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT CartId FROM cartlist WHERE SessionKey = ?", OldSessionKey)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&CartId)

		if err != nil {
			panic(err.Error())
		}
	}
	return CartId
}

func SignUpCart(CartId string, UID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare("UPDATE cartlist SET UID = ? WHERE CartId = ?")
	if err != nil {
		log.Fatal(err)
	}
	// SQLの実行
	_, err = ins.Exec(UID, CartId)
	defer ins.Close()
}

func UpdateCart(CartId string, req models.CartRequestPayload) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare("UPDATE cart SET Quantity = ? WHERE CartId = ? AND ItemId = ?")
	if err != nil {
		log.Fatal(err)
	}
	// SQLの実行
	_, err = ins.Exec(req.Quantity, CartId, req.ItemId)
	defer ins.Close()
}

func DeleteCart(CartId string, ItemId string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare("DELETE FROM cart WHERE CartId = ? AND ItemId = ?")
	if err != nil {
		log.Fatal(err)
	}
	// SQLの実行
	_, err = ins.Exec(CartId, ItemId)
	defer ins.Close()
}
func DeleteItemFromCart(ItemId string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare("DELETE FROM cart WHERE ItemId = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(ItemId)
	if err != nil {
		log.Fatal(err)
	}
}
