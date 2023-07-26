package handler

import (
	"database/sql"
	"log"
	"os"
	"unify/internal/models"
)

func PostTransaction(req models.TransactionRequestPayload) (res models.Transaction) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()
	TransactionDate := GetDate()
	TransactionId := GettransactionId()
	// SQLの準備
	ins, err := db.Prepare("INSERT INTO transaction VALUES(?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(req.UID, TransactionId, req.ItemId, TransactionDate, req.Count, req.IsFinished)
	if err != nil {
		log.Fatal(err)
	}
	return GetTransaction(req.UID)
}

func GetTransaction(uid string) (res models.Transaction) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM transaction WHERE uid = ?", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var Customer models.Transaction
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Customer.UID, &Customer.ItemId, &Customer.TransactionDate, &Customer.Count, &Customer.IsFinished)
		if err != nil {
			panic(err.Error())
		}
	}
	return Customer
}
func UpdateTransaction(req models.Transaction) (res models.Transaction) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

	// SQLの準備
	upd, err := db.Prepare("UPDATE transaction SET id = ?, transactionDate = ?, count = ?, isFinished = ? WHERE uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer upd.Close()
	// SQLの実行
	_, err = upd.Exec(req.ItemId, req.TransactionDate, req.Count, req.IsFinished, req.UID)
	if err != nil {
		log.Fatal(err)
	}
	return GetTransaction(req.UID)
}
func DeleteTransaction(uid string) (s string) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

	// SQLの実行
	del, err := db.Prepare("DELETE FROM transaction WHERE uid = ?")
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
