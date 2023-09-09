package database

// Item関連
type ItemList struct {
	ItemID string `json:"id"`
	InfoID string `json:"infoid"`
	Status string `json:"status"`
	Stock  int    `json:"stock"`
}

func (ItemList *ItemList) GetItemList(ItemID string) (err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		Item.ItemID,
		Item.InfoID,
		Item.Status,
		ItemDetails.Stock 
	
	FROM 
		Item
	
	JOIN 
		ItemDetails
		
	ON 
		Item.InfoID = ItemDetails.InfoID 
	
	WHERE 
		ItemID = ?`, ItemID)
	if err != nil {
		return err
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&ItemList.ItemID, &ItemList.InfoID, &ItemList.Status, &ItemList.Stock)
		if err != nil {
			return err
		}
	}
	return nil
}
