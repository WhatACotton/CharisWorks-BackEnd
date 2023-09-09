package database

import (
	"unify/cashing"

	"github.com/pkg/errors"
)

type Transaction struct {
	UID             string `json:"UID"`
	TransactionID   string `json:"TransactionID"`
	Name            string `json:"Name"`
	TotalAmount     int    `json:"TotalAmount"`
	Address         string `json:"Address"`
	PhoneNumber     string `json:"PhoneNumber"`
	TransactionTime string `json:"TransactionTime"`
	StripeID        string `json:"StripeID"`
	Status          string `json:"status"`
}
type Transactions []Transaction
type TransactionContent struct {
	Order         int    `json:"Order"`
	TransactionID string `json:"TransactionID"`
	InfoID        string `json:"InfoID"`
	Quantity      int    `json:"Quantity"`
}
type TransactionContents []TransactionContent

func PostTransaction(Cart Cart, Customer Customer, StripeInfo cashing.StripeInfo, TransactionID string, CartContents CartContents) {
	t := new(Transaction)
	t.UID = Customer.UID
	t.TransactionID = TransactionID
	t.Name = Customer.Name
	t.TotalAmount = int(StripeInfo.AmountTotal)
	t.Address = Customer.Address
	t.PhoneNumber = Customer.PhoneNumber
	t.TransactionTime = GetDate()
	t.StripeID = StripeInfo.ID
	t.Status = "決済前"
	t.postTransaction()

	TransactionContents := new(TransactionContent)
	for _, CartContent := range CartContents {
		TransactionContents.ConstructTransaction(CartContent, Cart, TransactionID)
		TransactionContents.PostTransactionContent()
	}
}
func (t *TransactionContent) ConstructTransaction(CartContent CartContent, Cart Cart, TransactionID string) {
	t.TransactionID = TransactionID
	t.InfoID = CartContent.InfoID
	t.Quantity = CartContent.Quantity
}

func (t *Transaction) postTransaction() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare(`
	INSERT 
		INTO 
			Transaction
			(UID,TransactionID,Name,TotalAmount,Address,PhoneNumber,TransactionTime,StripeID,Status)
			VALUES
			(?,?,?,?,?,?,?,?,?)
	`)
	if err != nil {
		return err
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(
		t.UID,
		t.TransactionID,
		t.Name,
		t.TotalAmount,
		t.Address,
		t.PhoneNumber,
		t.TransactionTime,
		t.StripeID,
		t.Status,
	)
	if err != nil {
		return err
	}
	return nil
}

func ChangeTransactionStatus(status string, stripeID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare(`UPDATE Transaction SET Status = ? WHERE StripeID = ?`)
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
func GetUIDfromStripeID(ID string) (string, error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの実行

	rows, err := db.Query(`
	SELECT 
		UID 
	FROM 
		Transaction
	WHERE 
		StripeID = ?`, ID)
	if err != nil {
		return "", errors.Wrap(err, "error in getting UID /Get_UID_1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {

		err := rows.Scan(&ID)

		if err != nil {
			return "", errors.Wrap(err, "error in scanning Email /Get_UID_2")
		}
	}
	return ID, nil
}
func (t *TransactionContent) PostTransactionContent() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの準備
	ins, err := db.Prepare(`
	INSERT 
		INTO 
			TransactionContent
			(InfoID,TransactionID,Quantity)
			VALUES
			(?,?,?)
	`)
	if err != nil {
		return err

	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(
		t.InfoID,
		t.TransactionID,
		t.Quantity,
	)
	if err != nil {
		return err

	}
	return nil
}
func Purchased(TransactionContent TransactionContent) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの準備
	ins, err := db.Prepare(`UPDATE ItemDetails SET Stock = Stock-? WHERE InfoID = ?`)
	if err != nil {
		panic(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(
		TransactionContent.Quantity,
		TransactionContent.InfoID,
	)
	if err != nil {
		panic(err)
	}
}
