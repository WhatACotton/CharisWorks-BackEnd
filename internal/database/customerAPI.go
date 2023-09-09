package database

import (
	"log"

	"github.com/pkg/errors"
)

func (c *Customer) Get_Customer(UID string) error {
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		UID,
		Name,
		Address,
		Email,
		Phone_Number,
		Register,
		Created_Date,
		Last_Session_Date 

	FROM 
		Customer 

	WHERE 
		UID= ?`, UID)
	if err != nil {
		log.Fatal(err)
		return errors.Wrap(err, "error in getting Customer /LogIn_Customer_1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		err := rows.Scan(&c.UID, &c.Name, &c.Address, &c.Email, &c.Phone_Number, &c.Register, &c.CreatedDate, &c.LastSessionDate)
		if err != nil {
			log.Fatal(err)
			return errors.Wrap(err, "error in scanning Customer /LogIn_Customer_2")
		}
	}

	return nil
}
func Get_Transaction_Lists(UID string) (Transaction_Lists []Transaction_List, err error) {
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		Transaction_ID,
		Stripe_ID,
		Total_Amount,
		Status,
		Transaction_Time,
		Address
	FROM 
		Transaction_List 
	WHERE 
		UID= ?
	AND
		Status != "決済前"	
		`, UID)
	if err != nil {
		log.Fatal(err)
		return Transaction_Lists, errors.Wrap(err, "error in getting Customer /LogIn_Customer_1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		Transaction_List := new(Transaction_List)
		//err := rows.Scan(&Customer)
		err := rows.Scan(&Transaction_List.Transaction_ID, &Transaction_List.Stripe_ID, &Transaction_List.Total_Amount, &Transaction_List.Status, &Transaction_List.Transaction_Time, &Transaction_List.Address)
		if err != nil {
			log.Fatal(err)
			return Transaction_Lists, errors.Wrap(err, "error in scanning Customer /LogIn_Customer_2")
		}
		Transaction_Lists = append(Transaction_Lists, *Transaction_List)
	}
	return Transaction_Lists, nil
}
func (t *Transaction_List) Get_Transactions() (Transacitons []Transaction, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	Transaciton := new(Transaction)
	// SQLの実行
	rows, err := db.Query(`
		SELECT * FROM Transaction WHERE Transaction_ID = ?`, t.Transaction_ID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting prepare Cart_ID /Get_Cart_Info_1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		Transaciton = new(Transaction)
		err := rows.Scan(&Transaciton.Order, &Transaciton.Info_ID, &Transaciton.Transaction_ID, &Transaciton.Quantity)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning Cart_ID /Get_Cart_Info_2")
		}
		Transacitons = append(Transacitons, *Transaciton)
	}
	return Transacitons, nil
}
