package database

type ItemDetailsImage struct {
	DetailsID string `json:"DetailsID"`
	Image     string `json:"Image"`
	Main      bool   `json:"Main"`
	Stock     int    `json:"Stock"`
	Price     int    `json:"Price"`
	ItemName  string `json:"ItemName"`
}
type ItemDetails struct {
	DetailsID   string `json:"DetailsID"`
	Description string `json:"Description"`
	Color       string `json:"Color"`
	Series      string `json:"Series"`
}

func ItemDetailsCreate() {
	// ItemDetailsの作成
}
func ItemDetailsModyfy() {
	// ItemDetailsの編集
}
func ItemDetailsGet(DetailsID string) ItemDetails {
	ItemDetails := new(ItemDetails)
	ItemDetails.DetailsID = DetailsID
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, _ := db.Query(
		`SELECT 
			Description,
			Color,
			Series
		From 
			ItemDetails

		WHERE 
			DetailsID = ?`,
		DetailsID)
	for rows.Next() {
		rows.Scan(&ItemDetails.Description, &ItemDetails.Color, &ItemDetails.Series)
	}
	return *ItemDetails
}
func ItemDetailsIDGet(ItemID string) (DetailsID string, Status string) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, _ := db.Query(
		`SELECT 
			DetailsID,
			status
		FROM 
			Item 
		
		WHERE 
			ItemID = ?`,
		ItemID)

	for rows.Next() {
		rows.Scan(&DetailsID, &Status)
	}
	return DetailsID, Status
}

func ItemDetailsUploadImage() {
	//画像を取得して、特定の場所に保存
	//画像のインスペクト・サイズの確認・
}
