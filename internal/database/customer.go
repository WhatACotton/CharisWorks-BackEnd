package database

import (
	"html"
	"log"
	"unify/internal/models"
	"unify/validation"

	"github.com/pkg/errors"
) // Customer関連
type Customer struct {
	UID            string `json:"UID"`
	Name           string `json:"Name"`
	Address        string `json:"address"`
	Email          string `json:"Contact"`
	PhoneNumber    string `json:"PhoneNumber"`
	Register       bool
	CreatedDate    string
	ModifiedDate   string
	RegisteredDate string
	LastSessionId  string
}

func SignUp_Customer(req models.CustomerRequestPayload, SessionID string) error {
	log.Printf("SignUpCustomer Called")
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate

	ins, err := db.Prepare("INSERT INTO Customer VALUES(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /SignUp_Customer_1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(req.UID, "default", "default", req.Email, "00000000000", false, GetDate(), 20000101, 20000101, SessionID)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /SignUp_Customer_2")
	}
	log.Printf(req.Email)

	return nil
}

func Register_Customer(usr validation.User, customer models.CustomerRegisterPayload) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate
	ins, err := db.Prepare(
		`UPDATE Customer SET 
		Name = ?,
		Address = ?,
		PhoneNumber = ?,
		Register = ?,
		RegisteredDate = ?,
		WHERE UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /Register_Customer_1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.Address), customer.PhoneNumber, true, GetDate(), usr.Userdata.UID)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /Register_Customer_2")
	}
	return nil
}

func Modify_Customer(usr validation.User, customer models.CustomerRegisterPayload) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate
	ins, err := db.Prepare(
		`UPDATE Customer SET 
		Name = ?,
		Address = ?,
		PhoneNumber = ?,
		Register = ?,
		ModifiedDate = ?,
		WHERE UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /Modify_Customer_1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.Address), customer.PhoneNumber, true, GetDate(), usr.Userdata.UID)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /Modify_Customer_2")
	}
	return nil
}

func Verify_Customer(uid string, OldSessionKey string) bool {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT Available FROM LogIn WHERE UID = ? AND Session_Key = ?", uid, OldSessionKey)
	if err != nil {
		log.Fatal(err)
		return false
	}

	defer rows.Close()
	var LogIn_Log int
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&LogIn_Log)

		if err != nil {
			return false
		}
	}
	if LogIn_Log == 1 {
		return true
	} else {
		return false
	}
}

func (c *Customer) LogIn_Customer(uid string, NewSessionKey string) error {
	LogIn_Log(uid, NewSessionKey)
	Update_Session_ID(uid, NewSessionKey)
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query("SELECT * FROM Customer WHERE UID = ?", uid)
	if err != nil {
		return errors.Wrap(err, "error in getting Customer /LogIn_Customer_1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		err := rows.Scan(&c.UID, &c.Name, &c.Address, &c.Email, &c.PhoneNumber, &c.Register, &c.CreatedDate, &c.ModifiedDate, &c.RegisteredDate, &c.LastSessionId)
		if err != nil {
			return errors.Wrap(err, "error in scanning Customer /LogIn_Customer_2")
		}
	}
	return nil
}

func Update_Session_ID(uid string, NewSessionKey string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare("UPDATE Customer SET Last_Session_ID = ? WHERE UID = ?")
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /Update_Session_ID_1")
	}
	// SQLの実行
	_, err = ins.Exec(NewSessionKey, uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /Update_Session_ID_2")
	}
	defer ins.Close()
	return nil
}

func LogIn_Log(uid string, NewSessionKey string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID SessionKey LoginedDate Available
	ins, err := db.Prepare("INSERT INTO LogIn (UID , Session_Key,LogIn_Date,Available)VALUES(?,?,?,1)")
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /LogIn_Log_1")
	}

	// SQLの実行
	_, err = ins.Exec(uid, NewSessionKey, GetDate())
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /LogIn_Log_2")
	}
	defer ins.Close()
	return nil
}

func Invalid(SessionKey string) error {
	log.Println("Invalid called")
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare("UPDATE LogIn SET Available = 0 WHERE Session_Key = ?")
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /Invalid_1")
	}
	// SQLの実行
	_, err = ins.Exec(SessionKey)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /Invalid_2")
	}
	defer ins.Close()
	return nil
}

func Get_UID(SessionKey string) (uid string, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの実行

	rows, err := db.Query("SELECT UID FROM LogIn WHERE Session_Key = ?", SessionKey)
	if err != nil {
		return "error", errors.Wrap(err, "error in getting UID /Get_UID_1")
	}
	defer rows.Close()
	var UID string
	// SQLの実行
	for rows.Next() {

		err := rows.Scan(&UID)

		if err != nil {
			return "error", errors.Wrap(err, "error in scanning UID /Get_UID_2")
		}
	}
	return UID, nil
}

func Delete_Customer(uid string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare("DELETE FROM Customer WHERE UID = ?")
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /Delete_Customer_1")
	}
	// SQLの実行
	_, err = ins.Exec(uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /Delete_Customer_2")
	}
	defer ins.Close()
	return nil
}
