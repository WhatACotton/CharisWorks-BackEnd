package database

import (
	"unify/cashing"
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
	DetailsID     string `json:"DetailsID"`
	Quantity      int    `json:"Quantity"`
}
type TransactionDetails []TransactionDetail

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
	t.transactionPost()

	TransactionContents := new(TransactionDetail)
	for _, CartContent := range CartContents {
		TransactionContents.TransactionConstruct(CartContent, Cart, TransactionID)
		TransactionContents.transactionPostDetails()
	}
}
func (t *TransactionDetail) TransactionConstruct(CartContent CartContent, Cart Cart, TransactionID string) {
	t.TransactionID = TransactionID
	t.DetailsID = CartContent.DetailsID
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
			Transaction
			(UserID,TransactionID,Name,TotalAmount,ZipCode,Address,TransactionTime,StripeID,Status)
			VALUES
			(?,?,?,?,?,?,?,?,?,?)
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

func (t *TransactionDetail) transactionPostDetails() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの準備
	ins, err := db.Prepare(`
	INSERT 
		INTO 
			TransactionDetails
			(DetailsID,TransactionID,Quantity)
			VALUES
			(?,?,?)
	`)
	if err != nil {
		return err

	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(
		t.DetailsID,
		t.TransactionID,
		t.Quantity,
	)
	if err != nil {
		return err

	}
	return nil
}
func Purchased(TransactionDetail TransactionDetail) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの準備
	ins, err := db.Prepare(`UPDATE ItemDetails SET Stock = Stock - ? WHERE DetailsID = ?`)
	if err != nil {
		panic(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(
		TransactionDetail.Quantity,
		TransactionDetail.DetailsID,
	)
	if err != nil {
		panic(err)
	}
}
func (t *Transaction) TransactionGetContents() (TransacitonContents TransactionDetails) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	TransactionContent := new(TransactionDetail)
	// SQLの実行
	rows, _ := db.Query(`
		SELECT * FROM TransactionContent WHERE TransactionID = ?`, t.TransactionID)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		rows.Scan(&TransactionContent.ItemOrder, &TransactionContent.DetailsID, &TransactionContent.TransactionID, &TransactionContent.Quantity)
		TransacitonContents = append(TransacitonContents, *TransactionContent)
	}
	return TransacitonContents
}
