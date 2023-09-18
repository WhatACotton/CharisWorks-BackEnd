package database

type Cart struct {
	CartID     string `json:"CartID"`
	SessionKey string `json:"SessionKey"`
}

// CartSessionListからCartIDを取得する
func (c *Cart) CartSessionListGetCartID() {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, _ := db.Query(`
	SELECT 
		CartID 
		
	FROM 
		CartSessionList 
	
	WHERE 
		SessionKey = ?`, c.SessionKey)

	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		rows.Scan(&c.CartID)
	}
}

// CartSessionListにCartIDを登録する
func (c *Cart) CartSessionListCreate() {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの準備
	ins, _ := db.Prepare(`
	INSERT 
	INTO 
		CartSessionList 
	(CartID,SessionKey)
	VALUES
	(?,?)
	`)
	defer ins.Close()
	ins.Exec(c.CartID, c.SessionKey)
}

// CartSessionListの削除　セッションキーの更新に使う
func CartSessionListDelete(CartID string) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの準備
	ins, _ := db.Prepare(`
	DELETE FROM 
		CartSessionList 
	WHERE 
		CartID = ?
	`)
	defer ins.Close()
	ins.Exec(CartID)
}
