package database

import (
	"log"
	"unify/internal/models"
	"unify/validation"
)

func PostCart(req models.CartRequestPayload, CartId string) (carts []models.Cart) {
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
	return GetCart(CartId)
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

func NewCartSession(SessionKey string) (CartId string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	CartId = validation.GetUUID()
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
	return CartId
}
func NewCartLogin(SessionKey string, UID string) (CartId string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	CartId = validation.GetUUID()
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
	return CartId
}

func StoredLoginCart(SessionKey string, CartId string, UID string) {
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
func StoredCartSession(SessionKey string, CartId string) {
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
