# SLEEVE-112-2: Firebase SDK初期化

## 概要

Email認証（アカウント登録）機能の実装の第2ステップとして、Firebase Admin SDK v4の初期化機能を実装しました。TDD（Test-Driven Development）アプローチに従い、テストファーストで開発を進めました。

## 変更内容

### テスト実装（Red Phase）
- `app/repository/external/firebase/init_firebase_test.go` を作成
- 以下の4つのテストケースを実装:
  - `TestInitializeFirebase_Success`: Firebase初期化成功のテスト
  - `TestInitializeFirebase_MissingCredentials`: 認証情報欠落時のエラーハンドリングテスト
  - `TestGetAuthClient_Success`: Auth Client取得成功のテスト
  - `TestGetAuthClient_NotInitialized`: 未初期化時のエラーハンドリングテスト
- グローバル変数の状態管理のため、各テストに適切なクリーンアップ処理を追加

### 実装（Green Phase）
- `app/repository/external/firebase/init_firebase.go` を作成
- `initialize_firebase()`: 環境変数から認証情報を読み込み、Firebase Appとauth.Clientを初期化
- `get_auth_client()`: 初期化済みのauth.Clientを取得、未初期化時はエラーを返す
- グローバル変数 `firebase_app` と `auth_client` でシングルトンパターンを実装

### 依存関係追加
- Firebase Admin SDK v4 (`firebase.google.com/go/v4 v4.18.0`) を追加
- `app/go.mod` と `app/go.sum` を更新

### テスト結果
- 全4テストケースが合格
- `golangci-lint run --fix` で0 issues

## 技術的な詳細

### 環境変数
- `GOOGLE_APPLICATION_CREDENTIALS`: Firebase service accountのJSONファイルパス

### エラーハンドリング
- 環境変数未設定時: 明確なエラーメッセージを返す
- Firebase App初期化失敗時: エラーをラップして返す
- Auth Client取得失敗時: エラーをラップして返す
- 未初期化状態でのAuth Client取得: 適切なエラーを返す

## 依存関係

- Depends on: SLEEVE-112-1（usersテーブルのマイグレーション）
- Blocks: SLEEVE-112-3（User Domain層）

## テスト

### 単体テスト
- [x] `TestInitializeFirebase_Success` が成功すること
- [x] `TestInitializeFirebase_MissingCredentials` が成功すること
- [x] `TestGetAuthClient_Success` が成功すること
- [x] `TestGetAuthClient_NotInitialized` が成功すること
- [x] すべてのテストが独立して実行可能であること

### 静的解析
- [x] `golangci-lint run --fix` で0 issues

## 関連チケット

- Jira: SLEEVE-112-2

## チェックリスト

- [x] TDDアプローチ（Red-Green-Refactor）に従った
- [x] テストファーストで実装
- [x] すべてのテストが合格
- [x] コーディング規約に準拠（変数宣言、`:=`の使用制限など）
- [x] Linterエラーがない（`golangci-lint run --fix`で0 issues）
- [x] エラーハンドリングが適切
- [x] 全ての変更がコミット済み
- [x] コミットメッセージが規約に従っている

## レビュー観点

- Firebase SDK初期化ロジックが正しいか
- エラーハンドリングが適切か
- テストケースが網羅的か
- グローバル変数の管理が適切か
- テストの独立性が保たれているか
- コーディング規約に準拠しているか
