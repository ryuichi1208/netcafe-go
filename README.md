# netcafe-go

東京都内のネットカフェ営業時間検索ツール

## インストール

```bash
go get github.com/ryuichi1208/netcafe-go
```

## 使い方

```bash
# 全店舗表示
./netcafe

# キーワード検索
./netcafe 新宿
./netcafe 渋谷

# Web最新情報取得
./netcafe -scrape

# ヘルプ
./netcafe -help
```

## 機能

- 店舗情報表示（名前、住所、営業時間、電話番号、URL）
- キーワード検索（店舗名・住所）
- Webスクレイピングによる最新情報取得
  - 快活CLUB
  - 自遊空間
  - マンボー

## 開発

```bash
# テスト実行
go test -v

# カバレッジ確認
go test -cover

# ビルド
go build -o netcafe
```

## 依存関係

- github.com/PuerkitoBio/goquery（スクレイピング用）