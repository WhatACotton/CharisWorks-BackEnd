package database

import (
	"html"
	"log"

	"github.com/WhatACotton/go-backend-test/validation"
)

type Customer struct {
	UserID           string `json:"UserID"`
	Name             string `json:"Name"`
	ZipCode          string `json:"ZipCode"`
	Address          string `json:"Address"`
	Email            string `json:"Contact"`
	IsRegistered     bool   `json:"IsRegistered"`
	CreatedDate      string `json:"CreatedDate"`
	IsEmailVerified  bool   `json:"IsEmailVerified"`
	CartID           string `json:"CartID"`
	StripeAccountID  string `json:"StripeAccountID,omitempty"`
	LastAccessedDate string `json:"LastAccessedDate"`
}

// サインアップ処理　LoginLog,Customerにデータを追加
func CustomerSignUp(req validation.CustomerReqPayload, NewSessionKey string, CartID string) error {
	log.Printf("SignUpCustomer Called")
	log.Print("UserID : ", req.UserID)
	log.Print("SessionKey : ", NewSessionKey)
	log.Print("Email : ", req.Email)
	log.Print("CartID : ", CartID)
	db := ConnectSQL()
	defer db.Close()
	tx, _ := db.Begin()
	//Customerに追加
	_, err := tx.Exec(`
	INSERT INTO 
		Customer 
		(UserID,Email,CartID) 
		VALUES 
		(?,?,?)`, req.UserID, req.Email, CartID)

	if err != nil {
		tx.Rollback()
		return err
	}

	//LoginLogに追加
	_, err = tx.Exec(`
	INSERT INTO
		LogInLog
		(UserID,SessionKey)
		VALUES
		(?,?)
	`, req.UserID, NewSessionKey)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer tx.Rollback()
	return nil
}

// 顧客情報の登録・変更
func CustomerRegister(UserID string, customer validation.CustomerRegisterPayload) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UserID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,LastLogInDate
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	SET 
		Name = ?,
		ZipCode = ?,
		Address = ?,
		IsRegistered = true

	WHERE 
		UserID = ?`)
	if err != nil {
		return err
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.ZipCode), html.EscapeString(customer.Address), UserID)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// ログイン処理　LoginLog,Customerの最終セッションを更新
func CustomerLogIn(UserID string, NewSessionKey string) error {
	db := ConnectSQL()
	defer db.Close()
	tx, _ := db.Begin()

	_, err := tx.Exec(`
	DELETE FROM 
		LogInLog
	WHERE
		UserID = ?`, UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
	INSERT INTO 
		LogInLog
		(UserID , SessionKey)
		VALUES
		(?,?)`, UserID, NewSessionKey)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	defer tx.Rollback()
	return nil
}

// Email認証情報の更新　メールが変更されたとき、メール認証がリセットされてしまうので、未認証への変更も対応
func CustomerEmailVerified(verify int, userid string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		IsEmailVerified = ? 
	
	WHERE 
		UserID = ?`)
	if err != nil {
		return err
	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(verify, userid)
	if err != nil {
		return err
	}
	defer ins.Close()
	return nil
}

// 顧客情報の削除　LoginLog,Customerからデータを削除
func CustomerDelete(UserID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	tx, _ := db.Begin()
	_, err := tx.Exec(`
	DELETE FROM 
		Customer 
	WHERE 
		UserID = ?`, UserID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`
	DELETE FROM 
		LogInLog 
	WHERE 
		UserID = ?`, UserID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Emailの変更
func CustomerChangeEmail(userid string, email string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		Email = ? 
	
	WHERE 
		UserID = ?`)
	if err != nil {
		return err
	}
	// SQLの実行
	_, err = ins.Exec(email, userid)
	if err != nil {
		return err
	}
	defer ins.Close()
	return nil

}

// CartIDをCustomerに登録　カートがからのときで、セッションカートリストのカートに商品がある場合は、そのカートIDを登録するので、内部関数化できない
func CustomerSetCartID(userid string, CartID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, _ := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		CartID = ? 
	
	WHERE 
		UserID = ?`)

	// SQLの実行
	ins.Exec(CartID, userid)
}

// 顧客の基本情報を取得
func (c *Customer) CustomerGet(UserID string) {
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
		LastAccessedDate,
		IsEmailVerified,
		CartID,
		StripeAccountID
	FROM 
		Customer 

	WHERE 
		UserID= ?`, UserID)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		rows.Scan(&c.Name, &c.ZipCode, &c.Address, &c.Email, &c.IsRegistered, &c.CreatedDate, &c.LastAccessedDate, &c.IsEmailVerified, &c.CartID, &c.StripeAccountID)
	}
}

// UserIDの取得　GetDatafromSessionKeyで使用し、直接呼び出さない
func GetUserID(SessionKey string) (UserID string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの実行
	defer db.Close()
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

// Emailの取得
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

// CartIDを取得
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
		UserID = ?`, UID)

	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&CartID)
		if err != nil {
			log.Print(err)
		}
	}
	return CartID
}
