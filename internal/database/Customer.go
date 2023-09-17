package database

import (
	"html"
	"log"
	"unify/validation"
)

type Customer struct {
	UserID          string `json:"UserID"`
	Name            string `json:"Name"`
	ZipCode         string `json:"ZipCode"`
	Address         string `json:"Address"`
	Email           string `json:"Contact"`
	PhoneNumber     string `json:"PhoneNumber"`
	IsRegistered    bool   `json:"IsRegistered"`
	CreatedDate     string `json:"CreatedDate"`
	IsEmailVerified bool   `json:"IsEmailVerified"`
	CartID          string `json:"CartID"`
	StripeAccountID string `json:"StripeAccountID,omitempty"`
}

func CustomerSignUp(req validation.CustomerReqPayload, NewSessionKey string, CartID string) {
	log.Printf("SignUpCustomer Called")
	log.Print("UserID : ", req.UserID)
	log.Print("SessionKey : ", NewSessionKey)
	db := ConnectSQL()
	tx, _ := db.Begin()
	_, err := tx.Exec(`
	INSERT INTO 
		Customer 
		(UserID,Email,LastSessionKey,CartID)
		VALUES
		(?,?,?,?)`, req.UserID, req.Email, NewSessionKey, CartID)
	if err != nil {
		tx.Rollback()
	}
	_, err = tx.Exec(`
	INSERT INTO
		LogInLog
		(UserID,SessionKey)
		VALUES
		(?,?)
	`, req.UserID, NewSessionKey)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()

}
func CustomerRegister(UserID string, customer validation.CustomerRegisterPayload) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UserID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,LastLogInDate
	ins, _ := db.Prepare(`
	UPDATE 
		Customer 
	SET 
		Name = ?,
		ZipCode = ?,
		Address = ?,
		IsRegistered = true

	WHERE 
		UserID = ?`)
	defer ins.Close()

	// SQLの実行
	ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.ZipCode), html.EscapeString(customer.Address), UserID)
}
func CustomerLogIn(UserID string, NewSessionKey string) {
	db := ConnectSQL()
	tx, _ := db.Begin()
	_, err := tx.Exec(`
	INSERT INTO 
		LogIn Log
		(UserID , SessionKey)
		VALUES
		(?,?)`, UserID, NewSessionKey)
	if err != nil {
		tx.Rollback()
	}
	_, err = tx.Exec(`
	UPDATE 
		Customer 

	SET 
		LastSessionKey = ? 

	WHERE 
		UserID = ?`, NewSessionKey, UserID)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
}
func CustomerEmailVerified(verify int, userid string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, _ := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		IsEmailVerified = ? 
	
	WHERE 
		UserID = ?`)

	// SQLの実行
	ins.Exec(verify, userid)
	defer ins.Close()
}

func CustomerDelete(userid string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, _ := db.Prepare(`
	DELETE FROM 
		Customer 
	WHERE 
		UserID = ?`)

	// SQLの実行
	ins.Exec(userid)
}
func CustomerDeleteSession(userid string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, _ := db.Prepare(`
	DELETE FROM 
		LogInLog 
	
	WHERE 
		UserID = ?`)
	// SQLの実行
	ins.Exec(userid)
	defer ins.Close()
}
func CustomerChangeEmail(userid string, email string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, _ := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		Email = ? 
	
	WHERE 
		UserID = ?`)
	// SQLの実行
	ins.Exec(email, userid)
	defer ins.Close()
}
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
func CustomerCreateStripeAccount(UserId string, StripeAccountID string) {
	db := ConnectSQL()
	tx, _ := db.Begin()

	_, err := tx.Exec(`
	UPDATE 
		Customer 
	
	SET 
		StripeAccountID = ? 
	
	WHERE 
		UserID = ?`, StripeAccountID, UserId)
	if err != nil {
		tx.Rollback()
	}
	_, err = tx.Exec(`
	INSERT
	INTO
		MakersDetails
		(MadeBy,Description,StripeAccountID)
		Value
		(?,?,?)`, UserId, "default", StripeAccountID)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
}
