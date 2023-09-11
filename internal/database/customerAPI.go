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
func (t *Transaction) GetTransactionContents() (TransacitonContents TransactionContents, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	TransactionContent := new(TransactionContent)
	// SQLの実行
	rows, err := db.Query(`
		SELECT * FROM TransactionContent WHERE TransactionID = ?`, t.TransactionID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting prepare CartID /GetCartInfo1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&TransactionContent.Order, &TransactionContent.InfoID, &TransactionContent.TransactionID, &TransactionContent.Quantity)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetCartInfo2")
		}
		TransacitonContents = append(TransacitonContents, *TransactionContent)
	}
	return TransacitonContents, nil
}
