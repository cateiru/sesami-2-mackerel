# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

SESAMIスマートロックデバイスの状態（バッテリー残量、ロック状態）をMackerelに送信して監視するためのGoアプリケーションです。24時間間隔で定期実行され、メトリクスを自動収集します。

## 開発・実行コマンド

### ローカル実行
```bash
# 環境変数を設定して実行
SESAMI_API_KEY=test_sesami_key SESAMI_DEVICE_UUID=test_device_uuid MACKEREL_API_KEY=test_mackerel_key go run main.go

# または直接実行
go run main.go
```

### Docker Compose実行
```bash
# .envファイルを作成後
docker compose up -d
```

### ビルド
```bash
go build -o sesami-2-mackerel .
```

## 必須環境変数

以下の環境変数を`.env`ファイルまたは環境に設定する必要があります：

- `SESAMI_API_KEY`: SESAMI RESTful webAPI Key
- `SESAMI_DEVICE_UUID`: SESAMI Device UUID  
- `MACKEREL_API_KEY`: Mackerel API Key

## アーキテクチャ

- `main.go`: エントリーポイント、設定読み込みとスケジューラー開始
- `internal/config/`: 環境変数からの設定読み込みと検証
- `internal/scheduler/`: cronライブラリを使用した24時間間隔の定期実行
- `internal/sesami/`: SESAMI APIクライアント（デバイス状態取得）
- `internal/mackerel/`: Mackerel APIクライアント（メトリクス送信）

## 現在の実装状況

- **SESAMI APIクライアント**: 現在はモック実装。実際のHTTP APIコールが必要
- **Mackerel APIクライアント**: 現在はモック実装。実際のHTTP APIコールが必要
- **スケジューラーとconfig**: 実装完了

## 依存関係

- Go 1.24.6
- `github.com/robfig/cron/v3`: cronスケジューリング