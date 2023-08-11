package database

import (
	"log"
	"unify/internal/models"
)

func GetItemList() (Itemlist []models.Item) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
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
		err := rows.Scan(&resultItem.ItemId, &resultItem.InfoId)
		if err != nil {
			panic(err.Error())
		}
		resultItemList = append(resultItemList, resultItem)
	}
	return resultItemList
}

func GetItem(id string) (returnmodels models.Item) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
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
		err := rows.Scan(&resultItem.ItemId, &resultItem.InfoId)
		if err != nil {
			panic(err.Error())
		}
	}
	return resultItem
}

func GetPrice(id string) (price int) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT price FROM infolist WHERE ItemId = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var resultPrice int
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&resultPrice)
		if err != nil {
			panic(err.Error())
		}
	}
	return resultPrice
}
