package database

import (
	"log"
	"unify/internal/models"
)

func SignUpCustomer(req models.CustomerRequestPayload) (res models.Customer) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	ins, err := db.Prepare("INSERT INTO user VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(req.UID, req.CreatedDate, "NULL", req.Email, "NULL")
	if err != nil {
		log.Fatal(err)
	}
	return GetCustomer(req.UID)
}

func RegisterCustomer(customer models.Customer) (res models.Customer) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	ins, err := db.Prepare("INSERT INTO user VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(customer.UID, customer.CreatedDate, "NULL", customer.Email, "NULL")
	if err != nil {
		log.Fatal(err)
	}
	return GetCustomer(customer.UID)

	// // SQLの準備
	// upd, err := db.Prepare("UPDATE ? SET ? = ? WHERE ? = ?")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer upd.Close()

	// if http.DetectContentType([]byte(patchItem.Value)) == "int" {
	// 	value, err := strconv.Atoi(patchItem.Value)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	// SQLの実行
	// 	_, err = upd.Exec(table, patchItem.Attribute, value, where, patchItem.ID)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// } else {
	// 	// SQLの実行
	// 	_, err = upd.Exec(table, patchItem.Attribute, patchItem.Attribute, where, patchItem.ID)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	//ココらへん使う
}
func ModifyCustomer(modify models.Customer) {

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
