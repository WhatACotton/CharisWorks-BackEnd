package database

import (
	"log"

	"github.com/WhatACotton/go-backend-test/cashing"
)

type CartContent struct {

	//From CartRequestPayload
	ItemID   string `json:"ItemID"`
	Quantity int    `json:"Quantity"`
	//From Item
	Status string `json:"Status"`
	Price  int    `json:"Price"`
	Name   string `json:"Name"`
	Stock  int    `json:"Stock"`
}
type CartContents []CartContent
type Transaction struct {
	TransactionID   string `json:"TransactionID"`
	UserID          string `json:"UserID"`
	CustomerName    string `json:"Name"`
	TotalAmount     int    `json:"TotalAmount"`
	ZipCode         string `json:"ZipCode"`
	Address1        string `json:"Address1"`
	Address2        string `json:"Address2"`
	Address3        string `json:"Address3"`
	PhoneNumber     string `json:"PhoneNumber"`
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
type CartRequestPayload struct {
	ItemID   string `json:"ItemID"`
	Quantity int    `json:"Quantity"`
}
type CartRequestPayloads []CartRequestPayload

// 取引履歴の作成
func TransactionPost(Customer Customer, StripeInfo cashing.StripeInfo, TransactionID string, CartRequestPayloads CartRequestPayloads) {
	t := new(Transaction)
	t.UserID = Customer.UserID
	t.TransactionID = TransactionID
	t.CustomerName = Customer.CustomerName
	t.TotalAmount = int(StripeInfo.AmountTotal)
	t.Address1 = Customer.Address1
	t.Address2 = Customer.Address2
	t.Address3 = Customer.Address3
	t.PhoneNumber = Customer.PhoneNumber
	t.TransactionTime = GetDate()
	t.StripeID = StripeInfo.ID
	t.Status = "決済前"
	t.ZipCode = Customer.ZipCode
	t.transactionPost()
	for _, CartRequestPayload := range CartRequestPayloads {
		CartRequestPayload.transactionDetailsPost(TransactionID)
	}
	PurchasedCart(Customer.UserID)
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
			(UserID,TransactionID,CustomerName,TotalAmount,ZipCode,Address1,Address2,Address3,PhoneNumber,TransactionTime,StripeID,Status)
			VALUES
			(?,?,?,?,?,?,?,?,?,?,?,?)
	`)

	defer ins.Close()

	// SQLの実行
	ins.Exec(
		t.UserID,
		t.TransactionID,
		t.CustomerName,
		t.TotalAmount,
		t.ZipCode,
		t.Address1,
		t.Address2,
		t.Address3,
		t.PhoneNumber,
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
func (t *CartRequestPayload) transactionDetailsPost(TransactionID string) error {
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
		TransactionID,
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
func PurchasedCart(UserID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの準備
	ins, err := db.Prepare(`UPDATE Customer SET Cart='Purchased' WHERE UserID = ?`)
	if err != nil {
		panic(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(
		UserID,
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
		rows.Scan(&TransactionDetail.TransactionID, &TransactionDetail.Quantity, &TransactionDetail.ItemOrder, &TransactionDetail.ItemID)
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
		CustomerName,
		TotalAmount,
		ZipCode,
		Address1,
		Address2,
		Address3,
		PhoneNumber,
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
		rows.Scan(&Transaction.TransactionID, &Transaction.CustomerName, &Transaction.TotalAmount, &Transaction.ZipCode, &Transaction.Address1, &Transaction.Address2, &Transaction.Address3, &Transaction.PhoneNumber, &Transaction.TransactionTime, &Transaction.StripeID, &Transaction.Status)
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

// カートの中身を確認し、購入可能かどうかを判定する。
func (carts *CartRequestPayloads) InspectCart() int {
	price := 0
	if len(*carts) == 0 {
		return 0
	}
	for _, Cart := range *carts {
		Item := new(Item)
		Item.ItemGet(Cart.ItemID)
		if Item.Status != "Available" {
			return 0
		}
		flag, stock := IsItemExist(Cart.ItemID)
		if !flag {
			return 0
		}
		if Cart.Quantity <= 0 {
			return 0
		}
		if stock < Cart.Quantity {
			return 0
		}
	}
	if hasDuplicates(*carts) {
		for _, Cart := range *carts {
			Item := new(Item)
			Item.ItemGet(Cart.ItemID)
			price += Item.Price * Cart.Quantity
		}
		return price
	}
	return 0
}
func hasDuplicates(slice CartRequestPayloads) bool {
	encountered := make(map[string]bool)
	for _, item := range slice {
		if encountered[item.ItemID] {
			return false
		}
		encountered[item.ItemID] = true
	}
	return true
}
