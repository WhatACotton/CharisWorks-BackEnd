# go-backend-server

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

### カート管理

#### カート ID 取得フロー

```mermaid
graph LR
A([start])-->B{login}
    B--Yes-->C([get CartID from Customer])
    B--No-->D([get CartID from Session])
    D--failed-->E([issue CartID])
    E-->F
    D--SessionReflesh-->F
    C--failed-->D
    C-->F([Cart dealing])

```

```mermaid
sequenceDiagram
participant Client
participant Server
participant DB
Client ->> Server:Item_ID Quantity
alt 未ログイン
opt 2回目以降
Client ->> Server:Cart_Session_Key
Server ->> DB:Cart_Session_Key
DB ->> Server:Cart_ID　
Server ->> DB: delete Cart_List
end
Note over Server:issue newSession_Key
Server ->> Client:newSession_Key
Server ->> DB:newSession_Key
else ログイン済み
Client ->> Server: Session_Key
rect rgba(255, 0, 255, 0.2)
Note right of DB: loginlog
Server ->> DB: Session_Key
DB ->> Server: UID
end
Server ->> DB: UID
Note right of DB: Customer
DB ->> Server:Cart_ID
DB ->> Server:Carts from Cart_ID


opt Cartに商品が入っていなかったとき
Client ->> Server: Cart_Session_Key
Server ->> Client: Delete Cart_Session_Key
rect rgba(255, 0, 255, 0.2)
Note right of DB: Cart_List
Server ->> DB: Cart_Session_Key
DB ->> Server:Cart_ID

end
end
Note over Server: Session Reflesh
rect rgba(255, 0, 255, 0.2)
Note right of DB: loginlog,Customer
Server ->> DB:newSession_Key,Cart_ID
end
Server ->> Client :newSession_Key
end
Server ->> DB:Cart_ID Item_ID Quantity
DB ->> Server:Carts from Cart_ID
Server ->> Client:Carts
```

## データ構造

### カート・取引　データベース

```mermaid
erDiagram
Customer }|--|{Transactions:""
Transactions}|--|{TransactionDetails:""
TransactionDetails}|--|{Item:""
Item}|--|{CartContents:""
Item}|--|{ItemDetailsImage:"商品の画像"
CartContents}|--|{Customer:""
CartContents}|--|{CartSessionList:""
Customer}|--|{LogInLog:"ログインとセッションの管理"
Item}|--|{Customer:"出品者が商品を出品"
Customer{
    UserID string PK "FirebaseTokenから取得"
    Email string "FirebaseTokenから取得"
    CartID string UK
    IsEmailVerified bool "FirebaseTokenから取得"
    CreatedDate timestamp "作成日時"
    Name string "名前"
    ZipCode string "郵便番号"
    Address string "郵便番号以降の住所"
    IsRegistered bool "本登録"
    StripeAccountID string "出品者のみ アカウントが許可されたとき'Allow'となる"
    MakerName string "出品者のみ"
    MakerDescription string "出品者のみ"
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
    MakerName string FK "出品者の名前"
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
}
CartContents{
    CartOrder int PK "カートに入っている商品の順番"
    CartID string UK
    ItemID string FK
    Quantity int "数量"
}
CartSessionList{
    CartID string FK
    SessionKey string
}
LogInLog{
    UserID string PK
    SessionKey string
}
```
