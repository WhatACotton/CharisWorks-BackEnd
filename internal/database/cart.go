package database

import (
	"log"
	"unify/internal/models"
	"unify/validation"
)

func PostCart(req models.CartRequestPayload, uid string) (carts []models.Cart) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,ItemId,Quantity

	ins, err := db.Prepare("INSERT INTO cartlist VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()
	CartId := validation.GetUUID()
	// SQLの実行
	_, err = ins.Exec(uid, CartId, req.ItemId, req.Quantity, GetDate())
	if err != nil {
		log.Fatal(err)
	}
	return GetCart(uid)
}

func GetCart(uid string) (Carts []models.Cart) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	var Cart models.Cart
	// SQLの実行
	rows, err := db.Query("SELECT * FROM cartlist WHERE uid = ?", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Cart.UID, &Cart.CartId, &Cart.ItemId, &Cart.Quantity, &Cart.RegisteredDate)
		if err != nil {
			panic(err.Error())
		}
		Carts = append(Carts, Cart)
	}

	return Carts
}
