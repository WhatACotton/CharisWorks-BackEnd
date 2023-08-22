package database

import (
	"unify/validation"
)

type Cart_List struct {
	Cart_ID     string `json:"Cart_ID"`
	Session_Key string `json:"Session_Key"`
}

func (c *Cart_List) Refresh_Cart_List() {
	c.Get_Cart_List()
	c.Delete_Cart_List()
	c.Cart_ID = validation.GetUUID()
	c.Create_Cart_List()
}

func (c *Cart_List) Get_Cart_List() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT CartID FROM Cart_List WHERE Session_Key = ?", c.Session_Key)
	if err != nil {
		return err
	}

	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&c.Cart_ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cart_List) Create_Cart_List() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	ins, err := db.Prepare("INSERT INTO Cart_List VALUES(?,?)")
	if err != nil {
		return err
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(c.Session_Key, c.Cart_ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cart_List) Delete_Cart_List() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	del, err := db.Prepare("DELETE FROM Cart_List WHERE Cart_ID = ?")
	if err != nil {
		return err
	}
	defer del.Close()

	// SQLの実行
	_, err = del.Exec(c.Cart_ID)
	if err != nil {
		return err
	}
	return nil
}
