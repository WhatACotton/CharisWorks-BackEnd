package database

import (
	"log"
)

// Item関連
type Item_List struct {
	Item_ID string `json:"id"`
	Info_ID string `json:"infoid"`
	Status  string `json:"status"`
}

func GetItemList() (Itemlist []Item_List) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM Item_List ")
	if err != nil {

		log.Fatal(err)
	}
	defer rows.Close()
	var resultItem Item_List
	var resultItemList []Item_List
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&resultItem.Item_ID, &resultItem.Info_ID)
		if err != nil {
			panic(err.Error())
		}
		resultItemList = append(resultItemList, resultItem)
	}
	return resultItemList
}

func GetItem(id string) (returnmodels Item_List) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM Item_List WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var resultItem Item_List
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&resultItem.Item_ID, &resultItem.Info_ID)
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
	rows, err := db.Query("SELECT price FROM Item_List WHERE ItemId = ?", id)
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
