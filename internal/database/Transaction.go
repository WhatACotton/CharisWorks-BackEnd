package database

import (
	"log"

	"github.com/WhatACotton/go-backend-test/cashing"
)

type Transaction struct {
	TransactionID   string `json:"TransactionID"`
	UserID          string `json:"UserID"`
	Name            string `json:"Name"`
	TotalAmount     int    `json:"TotalAmount"`
	ZipCode         string `json:"Zipcode"`
	Address         string `json:"Address"`
	TransactionTime string `json:"TransactionTime"`
	StripeID        string `json:"StripeID"`
	Status          string `json:"status"`
}
type Transactions []Transaction
type TransactionDetail struct {
	ItemOrder     int    `json:"ItemOrder"`
	TransactionID string `json:"TransactionID"`
	ItemID        string `json:"ItemID"`
	Quantity      int    `json:"Quantity"`
}
type TransactionDetails []TransactionDetail

// 取引履歴の作成
func TransactionPost(Cart Cart, Customer Customer, StripeInfo cashing.StripeInfo, TransactionID string, CartContents CartContents) {
	t := new(Transaction)
	t.UserID = Customer.UserID
	t.TransactionID = TransactionID
	t.Name = Customer.Name
	t.TotalAmount = int(StripeInfo.AmountTotal)
	t.Address = Customer.Address
	t.TransactionTime = GetDate()
	t.StripeID = StripeInfo.ID
	t.Status = "決済前"
	t.ZipCode = Customer.ZipCode
	t.transactionPost()
	TransactionContents := new(TransactionDetail)
	for _, CartContent := range CartContents {
		TransactionContents.transactionConstruct(CartContent, Cart, TransactionID)
		TransactionContents.transactionDetailsPost()
	}
}
func (t *TransactionDetail) transactionConstruct(CartContent CartContent, Cart Cart, TransactionID string) {
	t.TransactionID = TransactionID
	t.ItemID = CartContent.ItemID
	t.Quantity = CartContent.Quantity
}
func (t *Transaction) transactionPost() {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, _ := db.Prepare(`
	INSERT 
		INTO 
			Transactions
			(UserID,TransactionID,Name,TotalAmount,ZipCode,Address,TransactionTime,StripeID,Status)
			VALUES
			(?,?,?,?,?,?,?,?,?)
	`)

	defer ins.Close()

	// SQLの実行
	ins.Exec(
		t.UserID,
		t.TransactionID,
		t.Name,
		t.TotalAmount,
		t.ZipCode,
		t.Address,
		t.TransactionTime,
		t.StripeID,
		t.Status,
	)
}

// 取引履歴のステータスを取得
func TransactionGetStatus(stripeID string) (status string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		Status 
	FROM 
		Transactions
	WHERE 
		StripeID = ?`, stripeID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&status)
		if err != nil {
			panic(err)
		}
	}
	return status
}

// 取引履歴のステータスを更新
func TransactionSetStatus(status string, stripeID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare(`UPDATE Transactions SET Status = ? WHERE StripeID = ?`)
	if err != nil {
		panic(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(
		status,
		stripeID,
	)
	if err != nil {
		panic(err)
	}
}

// UserIDからStripeIDを取得
func TransactionGetUserIDfromStripeID(ID string) (StripeID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの実行

	rows, _ := db.Query(`
	SELECT 
		UserID 
	FROM 
		Transactions
	WHERE 
		StripeID = ?`, ID)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {

		rows.Scan(&StripeID)

	}
	return StripeID
}

// 取引履歴の登録
func (t *TransactionDetail) transactionDetailsPost() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの準備
	ins, err := db.Prepare(`
	INSERT 
		INTO 
			TransactionDetails
			(ItemID,TransactionID,Quantity)
			VALUES
			(?,?,?)
	`)
	if err != nil {
		return err

	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(
		t.ItemID,
		t.TransactionID,
		t.Quantity,
	)
	if err != nil {
		return err

	}
	return nil
}

// 購入処理
func Purchased(TransactionDetail TransactionDetail) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの準備
	ins, err := db.Prepare(`UPDATE Item SET Stock = Stock - ? WHERE ItemID = ?`)
	if err != nil {
		panic(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(
		TransactionDetail.Quantity,
		TransactionDetail.ItemID,
	)
	if err != nil {
		panic(err)
	}
}

// 取引履歴の詳細の取得
func (t *Transaction) TransactionDetailsGet() (TransacitonDetails TransactionDetails) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	TransactionDetail := new(TransactionDetail)
	// SQLの実行
	rows, _ := db.Query(`
		SELECT TransactionID,Quantity,ItemOrder,ItemID FROM TransactionDetails WHERE TransactionID = ?`, t.TransactionID)

	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&TransactionDetail.TransactionID, &TransactionDetail.Quantity, &TransactionDetail.ItemOrder, &TransactionDetail.ItemID)
		log.Print("TransactionContent:", TransactionDetail, "err:", err)
		TransacitonDetails = append(TransacitonDetails, *TransactionDetail)
	}
	return TransacitonDetails
}

// 取引履歴の取得
func TransactionGet(UserID string) (Transactions Transactions) {
	db := ConnectSQL()
	// SQLの実行
	rows, _ := db.Query(`
	SELECT 
		TransactionID,
		Name,
		TotalAmount,
		ZipCode,
		Address,
		TransactionTime,
		StripeID,
		Status
	FROM 
		Transactions 
	WHERE 
		UserID= ?
	AND
		Status != "決済前"	
		`, UserID)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		Transaction := new(Transaction)
		//err := rows.Scan(&Customer)
		rows.Scan(&Transaction.TransactionID, &Transaction.Name, &Transaction.TotalAmount, &Transaction.ZipCode, &Transaction.Address, &Transaction.TransactionTime, &Transaction.StripeID, &Transaction.Status)
		Transactions = append(Transactions, *Transaction)
	}
	return Transactions
}

// 　StripeIDからTransactionIDを取得
func TransactionGetID(StripeID string) (TransactionID string) {
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		TransactionID
	FROM 
		Transactions 
	WHERE 
		StripeID = ?
		`, StripeID)
	log.Print("rows:", rows, "err:", err)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		rows.Scan(&TransactionID)
	}
	log.Print("TransactionID:", TransactionID)
	return TransactionID
}
