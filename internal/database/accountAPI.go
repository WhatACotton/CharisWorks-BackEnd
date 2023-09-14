package database

import (
	"log"

	"github.com/pkg/errors"
)

func (c *Customer) GetCustomer(UID string) error {
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		UID,
		Name,
		ZipCode,
		Address,
		Email,
		PhoneNumber,
		Register,
		CreatedDate,
		LastSessionDate,
		EmailVerified,
		CartID,
		LastSessionKey
	FROM 
		Customer 

	WHERE 
		UID= ?`, UID)
	if err != nil {
		log.Fatal(err)
		return errors.Wrap(err, "error in getting Customer /LogInCustomer1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		err := rows.Scan(&c.UID, &c.Name, &c.ZipCode, &c.Address, &c.Email, &c.PhoneNumber, &c.Register, &c.CreatedDate, &c.LastSessionDate, &c.EmailVerified, &c.CartID, &c.LastSessionKey)
		if err != nil {
			log.Fatal(err)
			return errors.Wrap(err, "error in scanning Customer /LogInCustomer2")
		}
	}
	return nil
}
func GetTransactions(UID string) (Transactions Transactions, err error) {
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		UID,
		Name,
		TransactionID,
		StripeID,
		TotalAmount,
		Status,
		TransactionTime,
		Address,
		PhoneNumber
	FROM 
		Transaction 
	WHERE 
		UID= ?
	AND
		Status != "決済前"	
		`, UID)
	if err != nil {
		log.Fatal(err)
		return Transactions, errors.Wrap(err, "error in getting Customer /LogInCustomer1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		Transaction := new(Transaction)
		//err := rows.Scan(&Customer)
		err := rows.Scan(&Transaction.UID, &Transaction.Name, &Transaction.TransactionID, &Transaction.StripeID, &Transaction.TotalAmount, &Transaction.Status, &Transaction.TransactionTime, &Transaction.Address, &Transaction.PhoneNumber)
		if err != nil {
			log.Fatal(err)
			return Transactions, errors.Wrap(err, "error in scanning Customer /LogInCustomer2")
		}
		Transactions = append(Transactions, *Transaction)
	}
	return Transactions, nil
}
func GetTransactionID(StripeID string) (TransactionID string, err error) {
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		TransactionID
	FROM 
		Transaction 
	WHERE 
		StripeID = ?
		`, StripeID)
	if err != nil {
		log.Fatal(err)
		return TransactionID, errors.Wrap(err, "error in getting Customer /LogInCustomer1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		err := rows.Scan(&TransactionID)
		if err != nil {
			log.Fatal(err)
			return TransactionID, errors.Wrap(err, "error in scanning Customer /LogInCustomer2")
		}
	}
	return TransactionID, nil
}
func GetStripeAccountID(UID string) (StripeAccountID string, srr error) {
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		StripeAccountID
	FROM 
		Maker 
	WHERE 
		UID= ?`, UID)
	if err != nil {
		log.Fatal(err)
		return StripeAccountID, errors.Wrap(err, "error in getting Customer /AllowCreateMaker1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		var UID string
		err := rows.Scan(&StripeAccountID)
		if err != nil {
			log.Fatal(err)
			return StripeAccountID, errors.Wrap(err, "error in scanning Customer /AllowCreateMaker2")
		}
		if UID == "" {
			return StripeAccountID, nil
		}
	}
	return StripeAccountID, nil
}
func GetUID(SessionKey string) (UID string, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの実行

	rows, err := db.Query(`
	SELECT 
		UID

	FROM 
		LogIn

	WHERE 
		SessionKey = ?`, SessionKey)
	if err != nil {
		return "error", errors.Wrap(err, "error in getting UID /GetUID1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&UID)

		if err != nil {
			return "error", errors.Wrap(err, "error in scanning UID /GetUID2")
		}
	}
	return UID, nil
}
func GetEmail(UID string) (Email string, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの実行

	rows, err := db.Query(`
	SELECT 
		Email 
	
	FROM 
		Customer 
	
	WHERE 
		UID = ?`, UID)
	if err != nil {
		return "error", errors.Wrap(err, "error in getting Email /GetEmail1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {

		err := rows.Scan(&Email)

		if err != nil {
			return "error", errors.Wrap(err, "error in scanning Email /GetEmail2")
		}
	}
	return Email, nil
}
