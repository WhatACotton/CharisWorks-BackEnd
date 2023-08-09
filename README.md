# gin-server

**これは golang ベースで作られた開発中のバックエンドサーバです。**

今まで商品管理・アカウント・取引を別々に開発していましたが、今回はそれをすべて統合した形となります。

データ構造は次のようなものを考えています。

## カート・取引　データベース

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
 stirng UID
 string SessionKey
 bool Valid
}

Cart{
 timestamp CartTime
 string Cart_Id
 string Item_Id
 int quantity
}

ItemList{
 string Item_Id
 string Info_Id "マイナーチェンジはItem_Idに紐付ける"
 string status "購入可能かどうか"
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

## セッション

```mermaid
sequenceDiagram

participant Server
participant Client

Client ->> Server:SessionRequest
Note left of Server:validation
Server ->> Server:issue New SessionKey
Server ->> Server:invalidation of Requested SessionKey
Server ->> Client:New SessionKey


```

## ログイン

```mermaid
sequenceDiagram

participant Server
participant firebase
participant Client

Client ->> firebase:email,password
firebase ->> Client: uid,context
Client ->> Server: uid,context
Server ->> firebase:　
firebase ->> Server:　
Note left of firebase:Verify

```
