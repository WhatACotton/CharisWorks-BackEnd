package database

type Maker struct {
	MadeBy          string `json:"MadeBy"`
	Description     string `json:"Description"`
	Image           string `json:"Image"`
	StripeAccountID string `json:"StripeAccountID"`
}

func MakerAccountCreate(m Maker) {
	//stripeAccountを元にアカウントを作成する。
	//アカウントの画像はpathで指定して特定の場所に保存しておく必要がある。
}

func MakerAccountGet(MadeBy string) {
	//stripeAccountID以外を取得してくる。
}
func MakerAccountModyfy() {
	//アカウント情報の修正
}
func MakerAccountDelete() {
	//アカウント削除
	//関連しているテーブルに影響を与えざるを得ないので、あまり使いたくない。
}
func MakerItemDetailsModyfy() {
	//商品の価格変更など、大きな変更には使えない。
}
func MakerItemDetailsCreate() {

}
