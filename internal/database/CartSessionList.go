package database

type Cart struct {
	CartID     string `json:"CartID"`
	SessionKey string `json:"SessionKey"`
}

func (c *Cart) CartSessionListGetCartID() {
	// データベースのハンドルを取得する
	db := ConnectSQL()

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

func (c *Cart) CartSessionListCreate() {
	// データベースのハンドルを取得する
	db := ConnectSQL()

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
