package database

import "log"

type Cart struct {
	CartID     string `json:"CartID"`
	SessionKey string `json:"SessionKey"`
}

// CartSessionListからCartIDを取得する
func (c *Cart) CartSessionListGetCartID() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		CartID 
		
	FROM 
		CartSessionList 
	
	WHERE 
		SessionKey = ?`, c.SessionKey)

	if err != nil {
		return err
	}
	defer rows.Close()

	// SQLの実行
	for rows.Next() {
		err = rows.Scan(&c.CartID)
		if err != nil {
			return err
		}
	}
	return nil
}

// CartSessionListにCartIDを登録する
func (c *Cart) CartSessionListCreate() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの準備
	ins, err := db.Prepare(`
	INSERT 
	INTO 
		CartSessionList 
	(CartID,SessionKey)
	VALUES
	(?,?)
	`)
	if err != nil {
		return err
	}
	defer ins.Close()
	res, err := ins.Exec(c.CartID, c.SessionKey)
	log.Print("res:", res, "err:", err)
	if err != nil {
		return err
	}
	return nil
}

// CartSessionListの削除　セッションキーの更新に使う
func CartSessionListDelete(CartID string) error {
	db := ConnectSQL()
	defer db.Close()
	// SQLの準備
	ins, err := db.Prepare(`
	DELETE FROM 
		CartSessionList 
	WHERE 
		CartID = ?
	`)
	if err != nil {
		return err
	}
	defer ins.Close()
	_, err = ins.Exec(CartID)
	if err != nil {
		return err
	}
	return nil

}
