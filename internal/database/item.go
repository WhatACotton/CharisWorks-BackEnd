package database

import (
	"log"
)

// Item関連
type Item struct {
	ItemId string `json:"id"`
	InfoId string `json:"infoid"`
}

type ItemInfo struct {
	InfoId      string `json:"infoid"`
	Price       int    `json:"price"`
	Name        string `json:"Name"`
	Stonesize   int    `json:"Stonesize"`
	Minlength   int    `json:"Minlength"`
	Maxlength   int    `json:"Maxlength"`
	Decsription string `json:"Description"`
	Keyword     string `json:"Keyword"`
}

func GetItemList() (Itemlist []Item) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM itemlist ")
	if err != nil {

		log.Fatal(err)
	}
	defer rows.Close()
	var resultItem Item
	var resultItemList []Item
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

func GetItem(id string) (returnmodels Item) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM itemlist WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var resultItem Item
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
