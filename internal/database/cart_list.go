package database

import (
	"log"

	"github.com/pkg/errors"
)

type Cart struct {
	CartID     string `json:"CartID"`
	SessionKey string `json:"SessionKey"`
}

func (c *Cart) GetCartIDfromSessionKey() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		CartID 
		
	FROM 
		Cart 
	
	WHERE 
		SessionKey = ?`, c.SessionKey)
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Get_Cart_List_1")
	}

	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&c.CartID)
		if err != nil {
			return errors.Wrap(err, "error in scanning Cart_ID /Get_Cart_List_2")
		}
	}

	return nil
}

func (c *Cart) CreateCartList() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	ins, err := db.Prepare(`
	INSERT 
	INTO 
		Cart

	(CartID,SessionKey)
	VALUES
	(?,?)
	`)
	if err != nil {
		return errors.Wrap(err, "error in preparing Cart_List /Create_Cart_List_1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(c.CartID, c.SessionKey)
	if err != nil {
		return errors.Wrap(err, "error in inserting Cart_List /Create_Cart_List_2")
	}
	return nil
}

func DeleteCartList(CartID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	del, err := db.Prepare(`
	DELETE 
	FROM 
		Cart 
	WHERE 
		CartID = ?
	`)
	if err != nil {
		return errors.Wrap(err, "error in preparing to delete Cart_List /Delete_Cart_List_1")
	}
	defer del.Close()

	// SQLの実行
	_, err = del.Exec(CartID)
	if err != nil {
		return errors.Wrap(err, "error in deleting Cart_List /Delete_Cart_List_2")
	}
	return nil
}
func (c *Cart) SessionGet() bool {
	if c.SessionKey == "new" {
		log.Print("don't have sessionKey")
		return false
	} else {
		err := c.GetCartIDfromSessionKey()
		if err != nil {
			log.Fatal(err)
		}
		DeleteCartList(c.CartID)
		return true
	}
}
