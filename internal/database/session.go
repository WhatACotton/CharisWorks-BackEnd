package database

import (
	"log"
	"time"
	"unify/internal/models"
)

func Storesession(sessionId string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,ItemId,Quantity

	ins, err := db.Prepare("INSERT INTO session VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(sessionId, GetDate())
	if err != nil {
		log.Fatal(err)
	}
}
func Getsession(sessionId string) (SessionDate time.Time) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM user WHERE sessionId = ?", sessionId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var Customer models.Session
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Customer)
		if err != nil {
			panic(err.Error())
		}
	}
	return Customer.Date
}
