# 空目bot
## 作りたいやつ

形態素解析で「AをBに空目」というツイートを適当に判定
UserStreamとか使ってリアルタイム反応
ユーザー毎の統計を個人ページで見れたり、とか

## 使うライブラリ
* Goji - Sinatra風のやつ
* kagome - Pure Go 形態素解析器
* go-twitter - Twitter API
* gorm - ORマッパ
* go-oauth - TwitterOAuth
* ~~Configor - 設定ファイル取り扱い~~
* yaml - 設定ファイル取り扱い

## 今のところ
* GojiでのHTTPサーバ(server.go)
* UserStream受信(bot.go)
