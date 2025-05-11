# 基本的な内容

## データ型
redisでサポートされるデータ型は以下の通り

|型名|説明|ユースケース|
|---|---|---|
|String|文字列や数値|セッション管理|
|List|リスト|タイムライン|
|Hash|ハッシュ||
|Set|集合|タグ管理|
|Sorted Set|順序付きの集合|ランキングシステム|

言語における文字列、バイナリ、整数値・浮動小数値など、いずれもredisにおいては文字列として集約されてる。

## コマンド

### 基本操作
```
// データ型毎にコマンドが定義されてる
// 文字列・・・SET / GET / DEL...
// ハッシュ・・・HSET / HGET / HDEL...
SET USER 100
> OK
GET USER
> "100"
DEL USER
> (integer) 1
```

### ユーティリティコマンド
```
// キーの一覧を取得(※本番での利用は非推奨)
KEYS <pattern>

// キーの存在確認
EXISTS <key>

// データ型の確認
TYPE <key>
```

### TTLの設定
キーには生存期間を設定できる
```
// SETEX <key> <ttl> <value>
SETEX user:7 300 "tarou tanaka"
> OK
TTL user:7
> (integer) 299 // 残りのttlが帰ってくる
```

# ユースケース
## String型：セッションキャッシュ
```
session:XXX: {"user_id": "AAA"}
session:XXY: {"user_id": "BBB"}
...
```

キーから高速にセッション情報を取得できる

## List型：人気コンテンツの表示（タイムライン）
RDBMSから都度都度上位10件取得。。。だと応答性が悪いのでredis上のキャッシュとして上位コンテンツを保持しておく方が適切
```
LPUSH timeline topic1 topic2 topic3
> (integer) 3
LRANGE timeline 0 -1
> 1) topic1
> 2) topic2
> 3) topic3
```

## Sorted Set型：リアルタイムランキング表示
ユニークなIDのリストを保持し、リアルタイム性を頻繁な更新を可能とする仕組みはRDBMSでは実現が難しく、redisの`Sorted Set型`が適切

```
ZADD key [NX(キーがないなら) | XX(キーがあるなら)] [GT(grater than) | LT(less than)] score membera
```

# その他機能
## pub/sub機能
pub/sub機能も提供しており、
チャットルームにおけるメッセージをpublisherが送信し、
各subscriberがそれを受け取るような仕組みが作れる。
※websocketとredisを組み合わせることで、異なるサーバーに接続したユーザーが同一のチャットルームにおけるメッセージを同期するような仕組みが作れる（redisのpub/sub機能による同期を図ることでAPIサーバ側をスケール可能な構成として作れる）

```
// SUBSCRIBERは購読対象のチャンネルを選択
SUBSCRIBE channel1 channel2 channel3
> Reading messages...

// PUBLISHERはチャンネルへメッセージを送る
PUBLISH channnel1 "hello world"
> (integer) 1
```

# 付録
## キー名の命名規則
キー名の自由度が高いためチーム内での命名規則を定義することが重要

必要に応じて`<特定スキーマ名>:<キー>`のような形でセミコロンで連結する必要もある。
※セミコロン以外にも"."や"/"など、スキーマとキーの区別がつくような形式になっていれば良い

```
例）user-group:group-id:users:user-id => SET user-group:100:users:100 {"session": "10000"}
```
