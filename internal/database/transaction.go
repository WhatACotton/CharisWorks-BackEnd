package database

import (
	"log"
	"unify/internal/models"
)

func PostTransaction(req models.TransactionRequestPayload) (res models.Transaction) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
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
	_, err = ins.Exec(req.UID, TransactionId, req.ItemId, TransactionDate, req.Count, false)
	if err != nil {
		log.Fatal(err)
	}
	return GetTransaction(req.UID)
}

func GetTransaction(transactionId string) (res models.Transaction) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM transaction WHERE transactionid = ?", transactionId)
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

func GetTransactionList(uid string) (res []models.Transaction) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM transaction WHERE uid=" + uid)
	if err != nil {

		log.Fatal(err)
	}
	defer rows.Close()
	var resultTransaction models.Transaction
	var result []models.Transaction
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(
			&resultTransaction.UID,
			&resultTransaction.ItemId,
			&resultTransaction.TransactionDate,
			&resultTransaction.Count,
			&resultTransaction.IsFinished)
		if err != nil {
			panic(err.Error())
		}
		result = append(result, resultTransaction)
	}
	return result
}

func ChangeTransactionState(isFinished bool, transactionid string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	upd, err := db.Prepare("UPDATE transaction SET isFinished = ? WHERE transactionid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer upd.Close()

	// SQLの実行
	_, err = upd.Exec(isFinished, transactionid)
	if err != nil {
		log.Fatal(err)
	}
}
