package database

import (
	"github.com/pkg/errors"
)

type Cart_List struct {
	Cart_ID     string `json:"Cart_ID"`
	Session_Key string `json:"Session_Key"`
}

func (c *Cart_List) Get_Cart_ID_from_SessionKey() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		Cart_ID 
		
	FROM 
		Cart_List 
	
	WHERE 
		Session_Key = ?`, c.Session_Key)
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Get_Cart_List_1")
	}

	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&c.Cart_ID)
		if err != nil {
			return errors.Wrap(err, "error in scanning Cart_ID /Get_Cart_List_2")
		}
	}

	return nil
}

func (c *Cart_List) Create_Cart_List() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	ins, err := db.Prepare(`
	INSERT 
	INTO 
		Cart_List

	(Cart_ID,Session_Key)
	VALUES
	(?,?)
	`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Cart_List /Create_Cart_List_1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(c.Cart_ID, c.Session_Key)
	if err != nil {
		return errors.Wrap(err, "error in inserting Cart_List /Create_Cart_List_2")
	}
	return nil
}

func Delete_Cart_List(Cart_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	del, err := db.Prepare(`
	DELETE 
	FROM 
		Cart_List 
	WHERE 
		Cart_ID = ?
	`)
	if err != nil {
		return errors.Wrap(err, "error in preparing to delete Cart_List /Delete_Cart_List_1")
	}
	defer del.Close()

	// SQLの実行
	_, err = del.Exec(Cart_ID)
	if err != nil {
		return errors.Wrap(err, "error in deleting Cart_List /Delete_Cart_List_2")
	}
	return nil
}
