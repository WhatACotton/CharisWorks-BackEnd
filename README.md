# CharisWorks-BackEnd

**これは golang ベースで作られた開発中のバックエンドサーバです。**

今まで商品管理・アカウント・取引を別々に開発していましたが、今回はそれをすべて統合した形となります。

## 構成

```mermaid
graph LR
client --- Nginx
subgraph  [Server]
Nginx --- FrontEndServer
FrontEndServer --- BackEndServer
FireBaseAuth --- BackEndServer
BackEndServer --- DB
BackEndServer --- Stripe
FrontEndServer --- FireBaseAuth
Stripe --> FrontEndServer
end

```

## フロー

### セッション

```mermaid
sequenceDiagram

participant Client
participant Server
participant DB

Client ->> Server:Session_Key
Note right of DB: loginlog
Server ->> DB:Session_Key
DB ->> Server:UID
Note over Server: newSession_Key
Note right of DB: loginlog,Customer
Server ->> DB:NewSession_Key UID
Server ->> Client:NewSession_Key


```

### ログイン

```mermaid
sequenceDiagram
participant Client
participant firebase
participant Server
participant DB

Client ->> firebase:email password
firebase ->> Client: userCredential
Client ->> firebase:userCredential.user
firebase ->> Client:IdToken(JWT)
Client ->> Server: IdToken(JWT)
Server ->> firebase:IdToken(JWT)
firebase ->>Server:Token
Note over Server:issue Session_Key
Server ->> Client:Sesison_Key
Note right of DB: loginlog,Customer
Server ->> DB: Token.UID,Session_Key
```

## 取引フロー

```mermaid
sequenceDiagram
participant Client
participant Server
participant Stripe
participant DB


Client ->> Server:request with cart
Note over Server:check Item,Stock,Price
Server ->> DB:Register Transaction
Server ->> Stripe:Create PayentIntent
Stripe ->> Server:PaymentIntent Info
Server ->> Client:URL
Client ->> Stripe:Stripe Checkout
Stripe ->> Server:Status of Charge
Server ->> DB:Update Transaction
Server ->> DB:Update Item
```

## データ構造

### カート・取引　データベース

```mermaid
erDiagram
Customer }|--|{Transactions:""
Transactions}|--|{TransactionDetails:""
TransactionDetails}|--|{Item:""
Item}|--|{ItemDetailsImage:"商品の画像"
Item}|--|{Customer:"出品者が商品を出品"
Customer}|--|{LogInLog:"ログインとセッションの管理"
Maker}|--|{Item:"出品者が商品を出品"
Maker}|--|{Customer:"出品者が商品を出品"
Customer{
    UserID string PK "FirebaseTokenから取得"
    Email string "FirebaseTokenから取得"
    IsEmailVerified bool "FirebaseTokenから取得"
    CreatedDate timestamp "作成日時"
    Name string "名前"
    ZipCode string "郵便番号"
    Address string "郵便番号以降の住所"
    IsRegistered bool "本登録"

}
ItemDetailsImage{
    ItemID string PK
    Image string "画像へのpath"
    Order int
}
Item{
    ItemOrder int PK "商品の順番"
    ItemID string
    Status string "商品の状態(買えるかどうかなど)"
    Price int "価格"
    Stock int "在庫"
    ItemName string "商品の名前"
    Description string "商品説明"
    Color stirng "色"
    Series string "シリーズ"
    ItemSize string "商品のサイズなど"
    StripeAccountID string FK
    Top int "Topに表示するかどうか"
}
TransactionDetails{
    ItemOrder int PK "取引商品の順番"
    TransactionID string UK
    Quantity int "数量"
    ItemID string FK
}
Transactions{
    UserID string UK
    TransactionID string PK
    Name string "顧客の名前"
    TotalAmount int "購入金額"
    ZipCode string "顧客の郵便番号"
    Address string "郵便番号の先の住所"
    TransactionTime timestamp "取引された時間"
    StripeID string "Stripeから振られたID"
    Status string "取引の状態"
    ShipID string "郵便局の追跡番号"
}
LogInLog{
    UserID string PK
    SessionKey string
}
Maker{
    UserID string PK
    StripeAccountID string UK
    MakerName string
    Description string
}
```
