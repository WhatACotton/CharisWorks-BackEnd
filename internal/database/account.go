package database

import (
	"html"
	"log"
	"unify/internal/models"
	"unify/validation"

	"github.com/pkg/errors"
) // Customer関連
type Customer struct {
	UID             string `json:"UID"`
	Name            string `json:"Name"`
	Address         string `json:"address"`
	Email           string `json:"Contact"`
	Phone_Number    string `json:"PhoneNumber"`
	Register        bool   `json:"Register"`
	CreatedDate     string `json:"CreatedDate"`
	RegisteredDate  string `json:"RegisteredDate"`
	LastSessionId   string
	LastSessionDate string `json:"LastSessionDate"`
	Email_Verified  bool   `json:"Email_Verified"`
	Cart_ID         string `json:"Cart_ID"`
}

func Get_Customer(UID string) (string, error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	Cart_ID := new(string)
	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		Cart_ID 
	FROM 
		Customer 
	WHERE 
		UID = ?`, UID)
	if err != nil {
		return "new", errors.Wrap(err, "error in getting Customer /Get_Customer_1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		err := rows.Scan(&Cart_ID)
		if err != nil {
			return "new", errors.Wrap(err, "error in scanning Customer /Get_Customer_2")
		}
	}
	return *Cart_ID, nil
}
func SignUp_Customer(req models.CustomerRequestPayload, SessionID string, Cart_ID string) error {
	log.Printf("SignUpCustomer Called")
	log.Print("UID : ", req.UID)
	log.Print("Session_ID : ", SessionID)
	err := LogIn_Log(req.UID, SessionID)
	if err != nil {
		return errors.Wrap(err, "error in LogIn_Log /SignUp_Customer_1")
	}
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate

	ins, err := db.Prepare(`
	INSERT 
	INTO 
		Customer 
		(UID,Name,Address,Email,Phone_Number,Register,Last_Session_ID,Email_Verified,Cart_ID)
		VALUES
		(?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /SignUp_Customer_1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(req.UID, "default", "default", req.Email, "00000000000", false, SessionID, false, Cart_ID)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /SignUp_Customer_2")
	}
	log.Printf(req.Email)

	return nil
}

func Register_Customer(usr validation.UserReqPayload, customer models.CustomerRegisterPayload) error {
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
		RegisteredDate = ?,

	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /Register_Customer_1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.Address), customer.PhoneNumber, true, GetDate(), usr.UID)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /Register_Customer_2")
	}
	return nil
}

func Modify_Customer(usr validation.UserReqPayload, customer models.CustomerRegisterPayload) error {
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
		return errors.Wrap(err, "error in preparing Customer /Modify_Customer_1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.Address), customer.PhoneNumber, true, GetDate(), usr.UID)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /Modify_Customer_2")
	}
	return nil
}

func (c *Customer) LogIn_Customer(UID string, NewSessionKey string) error {
	LogIn_Log(UID, NewSessionKey)
	Update_Session_ID(UID, NewSessionKey)
	err := c.Get_Customer(UID)
	if err != nil {
		return errors.Wrap(err, "error in getting Customer /LogIn_Customer_1")
	}
	return nil
}

func Email_Verified(uid string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		Email_Verified = 1 
	
	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /Email_Verified_1")
	}
	// SQLの実行
	_, err = ins.Exec(uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /Email_Verified_2")
	}
	defer ins.Close()
	return nil
}

func Update_Session_ID(uid string, NewSessionKey string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		Last_Session_ID = ? 
	
	WHERE 
		UID = ?`)
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
	ins, err := db.Prepare(`
	INSERT 
	INTO 
		LogIn 
		(UID , Session_Key)
		VALUES
		(?,?)`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /LogIn_Log_1")
	}

	// SQLの実行
	_, err = ins.Exec(uid, NewSessionKey)
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
	ins, err := db.Prepare(`
	UPDATE 
		LogIn 
	
	SET 
		Available = 0 
		
	WHERE 
		Session_Key = ?`)
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

	rows, err := db.Query(`
	SELECT 
		UID

	FROM 
		LogIn

	WHERE 
		Session_Key = ?`, SessionKey)
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
	ins, err := db.Prepare(`
	DELETE 
	FROM 
		Customer 
	WHERE 
		UID = ?`)
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
func Delete_Session(uid string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	DELETE 
	FROM 
		LogIn 
	
	WHERE 
		UID = ?`)
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
func Get_Email(UID string) (Email string, err error) {
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
		return "error", errors.Wrap(err, "error in getting Email /Get_Email_1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {

		err := rows.Scan(&Email)

		if err != nil {
			return "error", errors.Wrap(err, "error in scanning Email /Get_Email_2")
		}
	}
	return Email, nil
}
func Change_Email(uid string, email string) error {
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
		return errors.Wrap(err, "error in preparing Customer /Change_Email_1")
	}
	// SQLの実行
	_, err = ins.Exec(email, uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /Change_Email_2")
	}
	defer ins.Close()
	return nil
}
func Set_Cart_ID(uid string, Cart_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare(`
	UPDATE 
		Customer 
	
	SET 
		Cart_ID = ? 
	
	WHERE 
		UID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Customer /Change_Email_1")
	}
	// SQLの実行
	_, err = ins.Exec(Cart_ID, uid)
	if err != nil {
		return errors.Wrap(err, "error in inserting Customer /Change_Email_2")
	}
	defer ins.Close()
	return nil
}
