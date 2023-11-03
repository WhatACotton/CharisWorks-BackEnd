package database

import "log"

type Maker struct {
	UserID          string `json:"UserID"`
	MakerName       string `json:"MakerName"`
	Description     string `json:"MakerDescription"`
	StripeAccountID string `json:"StripeAccountID"`
}

func CustomerGetStripeAccountID(UserID string) (role string) {
	db := ConnectSQL()
	// SQLの実行
	rows, _ := db.Query(`
	SELECT 
		role
	FROM 
		Customer 
	WHERE 
		UserID= ?`, UserID)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		rows.Scan(&role)
	}
	return role
}

// Stripeのアカウント作成
func MakerAccountCreate(MakerName string, StripeAccountID string) {
	db := ConnectSQL()
	ins, _ := db.Prepare(`
	INSERT INTO
		Maker
		(UserID,StripeAccountID)
	VALUES
		(?,?)`)

	ins.Exec(MakerName, StripeAccountID)
	defer ins.Close()
}

// Stripeのアカウント情報の修正
func (m *Maker) MakerAccountModyfy() {
	//アカウント情報の修正
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	res, err := db.Exec(`
	UPDATE
		Maker
	SET
		MakerName = ?,
		Description = ?
	WHERE
		StripeAccountID = ?`, m.MakerName, m.Description, m.StripeAccountID)
	log.Print("res:", res, "err:", err)
}

// MakerNameからStripeIDを取得
func MakerGetStripeID(UserID string) (StripeAccountID string) {
	db := ConnectSQL()
	defer db.Close()
	rows, _ := db.Query(`
	SELECT 
		StripeAccountID
	FROM 
		Maker
	WHERE 
		UserID = ?`, UserID)
	for rows.Next() {
		rows.Scan(&StripeAccountID)
	}
	log.Print("MakerName : " + UserID + " StripeAccountID : " + StripeAccountID)
	defer rows.Close()
	return StripeAccountID
}

// MakerNameをStripeIDから取得
func MakerStripeAccountIDGet(StripeAccountID string) (MakerName string) {
	db := ConnectSQL()
	defer db.Close()
	rows, _ := db.Query(`
	SELECT 
		MakerName
	FROM 
		Maker 
	WHERE 
		StripeAccountID = ?`, StripeAccountID)
	for rows.Next() {
		err := rows.Scan(&MakerName)
		if err != nil {
			log.Print(err)
		}
	}
	defer rows.Close()
	return MakerName
}

func (Maker *Maker) MakerDetailsGet() {
	db := ConnectSQL()
	defer db.Close()
	rows, _ := db.Query(`
	SELECT 
		MakerName,
		Description
	FROM 
		Maker 
	WHERE 
		StripeAccountID = ?`, Maker.StripeAccountID)
	for rows.Next() {
		rows.Scan(&Maker.MakerName, &Maker.Description)
	}
	defer rows.Close()
}
func MakerNameGet(UserID string) (StribeAccountID string) {
	db := ConnectSQL()
	defer db.Close()
	rows, _ := db.Query(`
	SELECT 
		StripeAccountID
	FROM 
		Maker 
	WHERE 
		UserID = ?`, UserID)
	for rows.Next() {
		rows.Scan(&StribeAccountID)
	}
	defer rows.Close()
	return StribeAccountID
}
