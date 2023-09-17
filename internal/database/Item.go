package database

// Item関連
type Item struct {
	ItemID    string `json:"id"`
	DetailsID string `json:"infoid"`
	Status    string `json:"status"`
}

type ItemDetailsImage struct {
	DetailsID string `json:"DetailsID"`
	Image     string `json:"Image"`
	Main      bool   `json:"Main"`
}
type ItemForCart struct {
	ItemID    string `json:"ItemID"`
	DetailsID string `json:"DetailsID"`
	Status    string `json:"Status"`
	Stock     int    `json:"Stock"`
}

func ItemCreate() {
	// Itemの作成
}
func ItemChange() {
	//引数は変更前,変更後
	//紐付いているItemの変更
}
func ItemChangeStatus() {
	// Itemの状態の変更
}

func (i *ItemForCart) ItemForCartGet(ItemID string) {
	//ItemForCartの取得
}
func ItemDetailsImageCreate() {
	//画像を取得して、特定の場所に保存
	//画像のインスペクト・サイズの確認・
}
