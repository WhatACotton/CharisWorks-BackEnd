package database

import (
	"database/sql"
	"log"
	"os"
	"unify/internal/models"
)

func PostItem(newItem models.Item) (Itemlist []models.Item) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		//log.Fatal(err)
		panic(err.Error())

	}
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare("INSERT INTO itemlist VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		//log.Fatal(err)
		panic(err.Error())

	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(
		newItem.ID,
		newItem.Price,
		newItem.Name,
		newItem.Stonesize,
		newItem.Minlength,
		newItem.Maxlength,
		newItem.Decsription,
		newItem.Keyword,
	)
	if err != nil {
		//log.Fatal(err)
		panic(err.Error())

	}
	return GetItemList()
}

func GetItemList() (Itemlist []models.Item) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM itemlist ")
	if err != nil {

		log.Fatal(err)
	}
	defer rows.Close()
	var resultItem models.Item
	var resultItemList []models.Item
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(
			&resultItem.ID,
			&resultItem.Price,
			&resultItem.Name,
			&resultItem.Stonesize,
			&resultItem.Minlength,
			&resultItem.Maxlength,
			&resultItem.Decsription,
			&resultItem.Keyword)
		if err != nil {
			panic(err.Error())
		}
		resultItemList = append(resultItemList, resultItem)
	}
	return resultItemList
}

func GetItem(id string) (returnmodels models.Item) {
	// データベースのハンドルを取得する

	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")
	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM itemlist WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var resultItem models.Item
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(
			&resultItem.ID,
			&resultItem.Price,
			&resultItem.Name,
			&resultItem.Stonesize,
			&resultItem.Minlength,
			&resultItem.Maxlength,
			&resultItem.Decsription,
			&resultItem.Keyword)
		if err != nil {
			panic(err.Error())
		}
	}
	return resultItem
}
