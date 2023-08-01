package database

import (
	"log"
	"unify/internal/models"
)

func PostCart(req models.CartRequestPayload, uid string) (cart models.Cart) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,ItemId,Quantity

	ins, err := db.Prepare("INSERT INTO user VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()
	CartId := GetUUID()
	// SQLの実行
	_, err = ins.Exec(uid, req.ItemId, req.Quantity, CartId, GetDate())
	if err != nil {
		log.Fatal(err)
	}
	return GetCart(CartId)
}

func GetCart(CartId string) (Cart models.Cart) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM cartlist WHERE id = ?", CartId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Cart)
		if err != nil {
			panic(err.Error())
		}
	}
	return Cart
}
