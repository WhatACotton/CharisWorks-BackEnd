# gin-server

**これは golang ベースで作られた開発中のバックエンドサーバです。**

今まで商品管理・アカウント・取引を別々に開発していましたが、今回はそれをすべて統合した形となります。

## 構成

```mermaid
graph TD
client --- Nginx
subgraph  [Server]
Nginx --- FrontEndServer
FrontEndServer --- APIServer
APIServer --- DB
FrontEndServer --- BackEndServer
FireBaseAuth --- BackEndServer
BackEndServer --- DB
BackEndServer --- CashServer
CashServer --- Stripe
FrontEndServer --- FireBaseAuth
end

```

## フロー

### セッション

```mermaid
sequenceDiagram

participant Client
participant Server
participant DB

Client ->> Server:SessionKey
Server ->> DB:Session_Key
DB ->> Server:UID
Note over Server: newSession_Key
Server ->> DB:NewSession_Key UID
Server ->> Client:newSessionKey


```

### ログイン

```mermaid
sequenceDiagram
participant Client
participant firebase

participant Server


Client ->> firebase:email password
firebase ->> Client: userCredential
Client ->> firebase:userCredential.user
firebase ->> Client:IdToken(JWT)
Client ->> Server: IdToken(JWT)
Server ->> firebase:IdToken(JWT)
firebase ->>Server:Token
Note over Server:issue Session_Key
Server ->> Client:Sesison_Key

```

### カート管理

#### カート ID 取得フロー

```mermaid
graph TB
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

Client ->> Server:POST/Session_Key Item_ID Quantity
Note left of DB:カートID取得フロー
DB ->> Server:Cart_ID　

Note over Server:issue newSession_Key
Server ->> Client:newSession_Key
Note right of DB:Cart
Server ->> DB:Cart_ID Item_ID Quantity
DB ->> Server:Carts from Cart_ID
Server ->> Client:Carts
```

## データ構造

### カート・取引　データベース

```mermaid
erDiagram
Transactionlist ||--||Cartlist:"取引とカートは1対1"
Transactionlist ||--|{Transaction:""
Cartlist}o--|{Cart:""
Cart}o--||ItemList:"一つのカートに同じItem_Idが複数存在することはない"
ItemList||--o{ItemInfo:""
Transaction}o--||ItemInfo:"ここでInfo_Idにしているので変更があっても旧Idに遡れる"

Transactionlist{
 string Cart_Id
 timestamp TransactionTime
 string UID
}
Transaction{
 string Cart_Id
 string Info_Id
 int quantity
}

Cartlist{
 string Cart_Id
 string SessionKey
}

Cart{
 int order "Auto Inctriment: 商品の登録した順番を管理"
 string Cart_Id
 string Item_Id
 int quantity
}

ItemList{
 string Item_Id
 string Info_Id "マイナーチェンジはItem_Idに紐付ける"
 string status "購入可能かどうか"
 int stock
}
ItemInfo{
 string Info_Id
 int Price
 string Name
 int Stone_Size
 int Min_length
 int Max_length
 string Description
 string Keyword
}
```
