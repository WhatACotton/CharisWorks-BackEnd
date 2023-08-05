package database

import (
	"log"
	"unify/internal/models"
)

func Storesession(sessionId string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,ItemId,Quantity

	ins, err := db.Prepare("INSERT INTO session VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(sessionId, GetDate(), true)
	if err != nil {
		log.Fatal(err)
	}
}

func Getsession(sessionId string) (SessionDate string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM session WHERE SessionKey = ?", sessionId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var Session models.CartSession
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Session.SessionId, &Session.Date)
		if err != nil {
			panic(err.Error())
		}
	}
	return Session.Date
}
