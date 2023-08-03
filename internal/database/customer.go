package database

import (
	"html"
	"log"
	"unify/internal/models"
	"unify/validation"
)

func SignUpCustomer(req models.CustomerRequestPayload, SessionID string) (res models.Customer) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate

	ins, err := db.Prepare("INSERT INTO user VALUES(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(req.UID, nil, nil, req.Email, nil, false, GetDate(), nil, nil, nil, SessionID)
	if err != nil {
		log.Fatal(err)
	}
	return GetCustomer(req.UID)
}

func RegisterCustomer(usr validation.User, customer models.CustomerRegisterPayload) (res models.Customer) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate
	ins, err := db.Prepare(
		`UPDATE user SET 
		Name = ?,
		Address = ?,
		PhoneNumber = ?,
		Register = ?,
		RegisteredDate = ?,
		WHERE uid = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.Address), customer.PhoneNumber, true, GetDate(), usr.Userdata.UID)
	if err != nil {
		log.Fatal(err)
	}
	return GetCustomer(usr.Userdata.UID)
}

func ModifyCustomer(usr validation.User, customer models.CustomerRegisterPayload) (res models.Customer) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate
	ins, err := db.Prepare(
		`UPDATE user SET 
		Name = ?,
		Address = ?,
		PhoneNumber = ?,
		Register = ?,
		ModifiedDate = ?,
		WHERE uid = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.Address), customer.PhoneNumber, true, GetDate(), usr.Userdata.UID)
	if err != nil {
		log.Fatal(err)
	}
	return GetCustomer(usr.Userdata.UID)
}

func LogInCustomer(uid string, SessionId string) (res models.Customer) {
	LogInTimeStamp(uid)
	LogInLog(uid)
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query("SELECT * FROM user WHERE uid = ?", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var Customer models.Customer
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		err := rows.Scan(&Customer.UID, &Customer.CreatedDate, &Customer.Name, &Customer.Email, &Customer.Address, &Customer.PhoneNumber, &Customer.LastSessionId)
		if err != nil {
			panic(err.Error())
		}
	}
	return Customer
}

func LogInTimeStamp(uid string) {
	// SQLの準備
	db := ConnectSQL()

	upd, err := db.Prepare("UPDATE user SET LastLogInDate = ? WHERE ? = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer upd.Close()
	// SQLの実行
	_, err = upd.Exec(GetDate(), uid)
}

func LogInLog(uid string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate

	ins, err := db.Prepare("INSERT INTO loginlog VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(uid, GetUUID(), GetDate())
	if err != nil {
		log.Fatal(err)
	}
}

func GetCustomer(uid string) (res models.Customer) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM user WHERE uid = ?", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var Customer models.Customer
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Customer.UID, &Customer.CreatedDate, &Customer.Name, &Customer.Email, &Customer.Address)
		if err != nil {
			panic(err.Error())
		}
	}
	return Customer
}
