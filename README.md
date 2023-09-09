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
Customer||--o{Cart:""
Cart}o--o|ItemList:"一つのカートに同じItem_Idが複数存在することはない"
ItemList|o--o|ItemInfo:""
Transaction}o--||ItemInfo:"ここでInfo_Idにしているので変更があっても旧Idに遡れる"
TransactionList }o--||Customer:""
TransactionList ||--|{Transaction:""

Customer{
 string UID
 string Address
 string Name
 string PhoneNumber
 string CartId
 string LastSession_Key
}

TransactionList{
 string CartID
 timestamp TransactionTime
 string UID
 int TotalPrice
 string Address
 string Name
 string PhoneNumber
}

Transaction{
 string CartId
 string InfoId
 int quantity
}



Cart{
 int order "Auto Inctriment: 商品の登録した順番を管理"
 string Cart_Id
 string Item_Id
 int quantity
}

ItemList{
 string ItemId
 string InfoId "マイナーチェンジはItem_Idに紐付ける"
 string status "購入可能かどうか"
}
ItemInfo{
 string InfoId
 int Price
 string Name
 int stock
 string Color
 string Key_Words
 string Description
 bool Top "Topに表示するかどうか"
}
```
