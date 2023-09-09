package database

import (
	"log"
	"unify/cashing"

	"github.com/pkg/errors"
)

type Transaction_List struct {
	UID              string `json:"UID"`
	Transaction_ID   string `json:"Transaction_ID"`
	Name             string `json:"Name"`
	Total_Amount     int    `json:"Total_Amount"`
	Address          string `json:"Address"`
	Phone_Number     string `json:"Phone_Number"`
	Transaction_Time string `json:"Transaction_Time"`
	Stripe_ID        string `json:"Stripe_ID"`
	Status           string `json:"status"`
}

type Transaction struct {
	Order          int    `json:"Order"`
	Transaction_ID string `json:"Transaction_ID"`
	Info_ID        string `json:"Info_ID"`
	Quantity       int    `json:"Quantity"`
}

func (t *Transaction_List) Construct_Transaction_List(Cart_List Cart_List, Customer Customer, stripe_info cashing.Stripe_info, Transaction_ID string) {
	t.UID = Customer.UID
	t.Transaction_ID = Transaction_ID
	t.Name = Customer.Name
	t.Total_Amount = int(stripe_info.AmountTotal)
	t.Address = Customer.Address
	t.Phone_Number = Customer.Phone_Number
	t.Transaction_Time = GetDate()
	t.Stripe_ID = stripe_info.ID
	t.Status = "決済前"
}
func (t *Transaction) Construct_Transaction(Cart Cart, Cart_List Cart_List, Transaction_ID string) {
	t.Transaction_ID = Transaction_ID
	t.Info_ID = Cart.Info_ID
	t.Quantity = Cart.Quantity
}

func (t *Transaction_List) Post_Transaction_List() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare(`
	INSERT 
		INTO 
			Transaction_List
			(UID,Transaction_ID,Name,Total_Amount,Address,Phone_Number,Transaction_Time,Stripe_ID,Status)
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
		t.Transaction_ID,
		t.Name,
		t.Total_Amount,
		t.Address,
		t.Phone_Number,
		t.Transaction_Time,
		t.Stripe_ID,
		t.Status,
	)
	if err != nil {
		return err
	}
	return nil
}

func (t *Transaction) Post_Transaction() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	log.Print(t)
	// SQLの準備
	ins, err := db.Prepare(`
	INSERT 
		INTO 
			Transaction
			(Info_ID,Transaction_ID,Quantity)
			VALUES
			(?,?,?)
	`)
	if err != nil {
		return err

	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(
		t.Info_ID,
		t.Transaction_ID,
		t.Quantity,
	)
	if err != nil {
		return err

	}
	return nil
}
func Change_Transaction_Status(status string, stripe_ID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare(`UPDATE Transaction_List SET Status = ? WHERE Stripe_ID = ?`)
	if err != nil {
		panic(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(
		status,
		stripe_ID,
	)
	if err != nil {
		panic(err)
	}
}
func Get_UID_from_Stripe_ID(ID string) (string, error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの実行

	rows, err := db.Query(`
	SELECT 
		UID 
	
	FROM 
		Transaction_List
	
	WHERE 
		Stripe_ID = ?`, ID)
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
