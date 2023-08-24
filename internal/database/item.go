package database

import "log"

// Item関連
type Item_List struct {
	Item_ID string `json:"id"`
	Info_ID string `json:"infoid"`
	Status  string `json:"status"`
	Stock   int    `json:"stock"`
}

func (Item_List *Item_List) Get_Item_List(Item_ID string) (err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query("SELECT Item_List.Item_ID,Item_List.Info_ID,Item_List.Status,Item_Info.Stock FROM Item_List JOIN Item_Info ON Item_List.Info_ID = Item_Info.Info_ID WHERE Item_ID = ?", Item_ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Item_List.Item_ID, &Item_List.Info_ID, &Item_List.Status, &Item_List.Stock)
		if err != nil {
			return err
		}
	}
	log.Print(Item_List)
	return nil
}
