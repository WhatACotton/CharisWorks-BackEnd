package handler

import (
	"database/sql"
	"log"
	"os"
	"unify/internal/models"
)

func PostCustomer(req models.CustomerRequestPayload) (res models.Customer) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

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
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

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
func UpdateCustomer(req models.Customer) (res models.Customer) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

	// SQLの準備
	upd, err := db.Prepare("UPDATE user SET Name = ?, Address = ? WHERE uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer upd.Close()
	// SQLの実行
	_, err = upd.Exec(req.Name, req.Address, req.UID)
	if err != nil {
		log.Fatal(err)
	}
	return GetCustomer(req.UID)
}
func DeleteCustomer(uid string) (s string) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

	// SQLの実行
	del, err := db.Prepare("DELETE FROM user WHERE uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer del.Close()

	// SQLの実行
	_, err = del.Exec(uid)
	if err != nil {
		log.Fatal(err)
	}

	return
}
