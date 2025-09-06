# sesami2mackerel

[SESAMI](https://jp.candyhouse.co/) のデバイスを [Mackerel](https://mackerel.io/) で監視するためのツールです。

## 取得する情報

- SESAMIデバイスのバッテリー残量（%）
- ~~SESAMIの操作履歴~~（将来的に対応予定）

## Usage

### 以下の環境変数を設定した `.env` ファイルを作成してください

```sh
touch .env
```

```env
SESAMI_API_KEY=[your Sesame RESTful webAPI Key]
SESAMI_DEVICE_UUID=[your Sesame Device UUID]
MACKEREL_API_KEY=[your Mackerel API Key]
```

### 実行します

```sh
docker compose up -d
```

## LICENSE

[MIT](./LICENSE)
