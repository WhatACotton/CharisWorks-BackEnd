package database

import (
	"html"
	"log"
	"unify/validation"

	"github.com/pkg/errors"
) // Customer関連
type Customer struct {
	UID             string `json:"UID"`
	Name            string `json:"Name"`
	ZipCode         string `json:"ZipCode"`
	Address         string `json:"Address"`
	Email           string `json:"Contact"`
	PhoneNumber     string `json:"PhoneNumber"`
	Register        bool   `json:"Register"`
	CreatedDate     string `json:"CreatedDate"`
	LastSessionKey  string
	LastSessionDate string `json:"LastSessionDate"`
	EmailVerified   bool   `json:"EmailVerified"`
	CartID          string `json:"CartID"`
	StripeAccountID string `json:"StripeAccountID,omitempty"`
}

func SignUpCustomer(req validation.CustomerReqPayload, SessionID string, CartID string) error {
	log.Printf("SignUpCustomer Called")
	log.Print("UID : ", req.UID)
	log.Print("SessionKey : ", SessionID)
	err := LogInLog(req.UID, SessionID)
	if err != nil {
		return errors.Wrap(err, "error in LogInLog /SignUpCustomer1")
	}
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,LastLogInDate

	ins, err := db.Prepare(`
	INSERT 
	INTO 
		Customer 
		(UID,Name,Address,ZipCode,Email,PhoneNumber,Register,LastSessionKey,EmailVerified,CartID)
		VALUES
		(?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /SignUpCustomer1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(req.UID, "default", "default", "default", req.Email, "00000000000", false, SessionID, false, CartID)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /SignUpCustomer2")
	}
	log.Printf(req.Email)

	return nil
}
func RegisterCustomer(UID string, customer validation.CustomerRegisterPayload) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,LastLogInDate
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	SET 
		Name = ?,
		ZipCode = ?,
		Address = ?,
		PhoneNumber = ?,
		Register = ?

	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /RegisterCustomer1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.ZipCode), html.EscapeString(customer.Address), customer.PhoneNumber, true, UID)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /RegisterCustomer2")
	}
	return nil
}
func ModifyCustomer(usr validation.CustomerReqPayload, customer validation.CustomerRegisterPayload) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
		
	SET 
		Name = ?,
		Address = ?,
		PhoneNumber = ?,
		Register = ?,
		ModifiedDate = ?,
		
	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /ModifyCustomer1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.Address), customer.PhoneNumber, true, GetDate(), usr.UID)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /ModifyCustomer2")
	}
	return nil
}
func LogIn(UID string, NewSessionKey string) error {
	LogInLog(UID, NewSessionKey)
	UpdateSessionID(UID, NewSessionKey)
	return nil
}
func EmailVerified(verify int, uid string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		EmailVerified = ? 
	
	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /EmailVerified1")
	}
	// SQLの実行
	_, err = ins.Exec(verify, uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /EmailVerified2")
	}
	defer ins.Close()
	return nil
}
func UpdateSessionID(uid string, NewSessionKey string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		LastSessionKey = ? 
	
	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /UpdateSessionID1")
	}
	// SQLの実行
	_, err = ins.Exec(NewSessionKey, uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /UpdateSessionID2")
	}
	defer ins.Close()
	return nil
}
func LogInLog(uid string, NewSessionKey string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID SessionKey LoginedDate Available
	ins, err := db.Prepare(`
	INSERT 
	INTO 
		LogIn 
		(UID , SessionKey)
		VALUES
		(?,?)`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /LogInLog1")
	}

	// SQLの実行
	_, err = ins.Exec(uid, NewSessionKey)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /LogInLog2")
	}
	defer ins.Close()
	return nil
}
func Invalid(SessionKey string) error {
	log.Println("Invalid called")
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	UPDATE 
		LogIn 
	
	SET 
		Available = 0 
		
	WHERE 
		SessionKey = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /Invalid1")
	}
	// SQLの実行
	_, err = ins.Exec(SessionKey)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /Invalid2")
	}
	defer ins.Close()
	return nil
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
func DeleteCustomer(uid string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	DELETE 
	FROM 
		Customer 
	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /DeleteCustomer1")
	}
	// SQLの実行
	_, err = ins.Exec(uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /DeleteCustomer2")
	}
	defer ins.Close()
	return nil
}
func DeleteSession(uid string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	DELETE 
	FROM 
		LogIn 
	
	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /DeleteCustomer1")
	}
	// SQLの実行
	_, err = ins.Exec(uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /DeleteCustomer2")
	}
	defer ins.Close()
	return nil
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
func ChangeEmail(uid string, email string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		Email = ? 
	
	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /ChangeEmail1")
	}
	// SQLの実行
	_, err = ins.Exec(email, uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /ChangeEmail2")
	}
	defer ins.Close()
	return nil
}
func SetCartID(uid string, CartID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		CartID = ? 
	
	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /ChangeEmail1")
	}
	// SQLの実行
	_, err = ins.Exec(CartID, uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /ChangeEmail2")
	}
	defer ins.Close()
	return nil
}

func CreateStripeAccount(uid string, StripeAccountID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		StripeAccountID = ? 
	
	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /CreateStripeAccount1")
	}
	// SQLの実行
	_, err = ins.Exec(StripeAccountID, uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /CreateStripeAccount2")
	}
	defer ins.Close()
	return nil
}
