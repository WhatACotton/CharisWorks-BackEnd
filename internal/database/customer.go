package database

import (
	"log"
	"unify/internal/models"
)

func PostCustomer(req models.CustomerRequestPayload) (res models.Customer) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	ins, err := db.Prepare("INSERT INTO user VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(req.UID, req.CreatedDate, "NULL", req.Contact, "NULL")
	if err != nil {
		log.Fatal(err)
	}
	return GetCustomer(req.UID)
}

func GetCustomer(uid string) (res models.Customer) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM user WHERE uid = ?", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var Customer models.Customer
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Customer.UID, &Customer.CreatedDate, &Customer.Name, &Customer.Contact, &Customer.Address)
		if err != nil {
			panic(err.Error())
		}
	}
	return Customer
}
