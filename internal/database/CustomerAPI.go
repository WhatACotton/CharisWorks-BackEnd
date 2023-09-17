package database

func (c *Customer) GetCustomer(UserID string) {
	db := ConnectSQL()
	c.UserID = UserID
	// SQLの実行
	rows, _ := db.Query(`
	SELECT 
		Name,
		ZipCode,
		Address,
		Email,
		IsRegistered,
		CreatedDate,
		IsEmailVerified,
		CartID,
		LastSessionKey,
		StripeAccountID
	FROM 
		Customer 

	WHERE 
		UserID= ?`, UserID)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		rows.Scan(&c.Name, &c.ZipCode, &c.Address, &c.Email, &c.IsRegistered, &c.CreatedDate, &c.IsEmailVerified, &c.CartID)
	}
}
func GetTransactions(UserID string) (Transactions Transactions) {
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
func GetTransactionID(StripeID string) (TransactionID string) {
	db := ConnectSQL()
	// SQLの実行
	rows, _ := db.Query(`
	SELECT 
		TransactionID
	FROM 
		Transactions 
	WHERE 
		StripeID = ?
		`, StripeID)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		rows.Scan(&TransactionID)
	}
	return TransactionID
}
func GetStripeAccountID(UserID string) (StripeAccountID string) {
	db := ConnectSQL()
	// SQLの実行
	rows, _ := db.Query(`
	SELECT 
		StripeAccountID
	FROM 
		MakersDetails 
	WHERE 
		UserID= ?`, UserID)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		var UserID string
		rows.Scan(&StripeAccountID)
		if UserID == "" {
			return StripeAccountID
		}
	}
	return StripeAccountID
}
func GetUserID(SessionKey string) (UserID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの実行

	rows, _ := db.Query(`
	SELECT 
		UserID

	FROM 
		LogInLog

	WHERE 
		SessionKey = ?`, SessionKey)

	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		rows.Scan(&UserID)
	}
	return UserID
}
func GetEmail(UserID string) (Email string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの実行
	rows, _ := db.Query(`
	SELECT 
		Email 
	
	FROM 
		Customer 
	
	WHERE 
		UserID = ?`, UserID)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		rows.Scan(&Email)
	}
	return Email
}
func GetCartID(UID string) (CartID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, _ := db.Query(`
	SELECT 
		CartID
	FROM 
		Customer 
	WHERE 
		UID = ?`, UID)

	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		rows.Scan(&CartID)
	}
	return CartID
}
