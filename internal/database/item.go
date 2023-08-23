package database

// Item関連
type Item_List struct {
	Item_ID string `json:"id"`
	Info_ID string `json:"infoid"`
	Status  string `json:"status"`
}

func Get_Item_List() (Itemlist []Item_List, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM Item_List ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resultItem Item_List
	var resultItemList []Item_List
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&resultItem.Item_ID, &resultItem.Info_ID)
		if err != nil {
			return nil, err
		}
		resultItemList = append(resultItemList, resultItem)
	}
	return resultItemList, nil
}

func (returnmodels Item_List) Get_Item_List(id string) (err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM Item_List WHERE id = ?", id)
	if err != nil {
		return err
	}
	defer rows.Close()
	var resultItem Item_List
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&resultItem.Item_ID, &resultItem.Info_ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func Get_Price(Info_ID string) (price int, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT price FROM Item WHERE Info_ID = ?", Info_ID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var resultPrice int
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&resultPrice)
		if err != nil {
			return 0, err
		}
	}
	return resultPrice, nil
}
