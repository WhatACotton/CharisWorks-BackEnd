package database

type Maker struct {
	MakerName       string `json:"MakerName"`
	Description     string `json:"MakerDescription"`
	StripeAccountID string `json:"StripeAccountID"`
}

// Stripeのアカウントを登録
func CustomerCreateStripeAccount(UserId string, StripeAccountID string) {
	db := ConnectSQL()
	tx, _ := db.Begin()

	_, err := tx.Exec(`
	UPDATE 
		Customer 
	
	SET 
		StripeAccountID = ? 
	
	WHERE 
		UserID = ?`, StripeAccountID, UserId)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
}

// StripeAccountIDを取得
func CustomerGetStripeAccountID(UserID string) (StripeAccountID string) {
	db := ConnectSQL()
	// SQLの実行
	rows, _ := db.Query(`
	SELECT 
		StripeAccountID
	FROM 
		Customer 
	WHERE 
		UserID= ?`, UserID)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		rows.Scan(&StripeAccountID)
	}
	return StripeAccountID
}

// Stripeのアカウント作成
func (m Maker) MakerAccountCreate() {
	db := ConnectSQL()
	ins, _ := db.Prepare(`
	UPDATE
		Customer
	SET
		StripeAccountID = ?,
		MakerName = ?,
		MakerDescription = ?
	WHERE
		UserID = ?
	`)
	ins.Exec(m.MakerName, m.Description, m.StripeAccountID)
	defer ins.Close()
}

// Stripeのアカウント情報の修正
func (m *Maker) MakerAccountModyfy() {
	//アカウント情報の修正
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	db.Exec(`
	UPDATE
		Customer
	SET
		MakerName = ?,
		MakerDescription = ?
	WHERE
		StripeAccountID = ?`, m.MakerName, m.Description, m.StripeAccountID)
}

// MakerNameからStripeIDを取得
func MakerGetStripeID(MakerName string) (StripeAccountID string) {
	db := ConnectSQL()
	defer db.Close()
	rows, _ := db.Query(`
	SELECT 
		StripeAccountID
	FROM 
		Cutomer
	WHERE 
		MakerName = ?`, MakerName)
	rows.Scan(&StripeAccountID)
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
		Customer 
	WHERE 
		StripeAccountID = ?`, StripeAccountID)
	rows.Scan(&MakerName)
	defer rows.Close()
	return MakerName
}

func (Maker *Maker) MakerDetailsGet() {
	db := ConnectSQL()
	defer db.Close()
	rows, _ := db.Query(`
	SELECT 
		MakerName,
		MakerDescription
	FROM 
		Customer 
	WHERE 
		StripeAccountID = ?`, Maker.StripeAccountID)
	for rows.Next() {
		rows.Scan(&Maker.MakerName, &Maker.Description)
	}
	defer rows.Close()
}
