## 概要

JWT発行・検証サービスの実装を行いました。アクセストークン（15分有効期限）とリフレッシュトークン（7日有効期限）の生成・検証機能を提供します。

## 変更内容

### 新規作成ファイル

- `app/usecase/utils/jwt.go` - JWTサービス実装
- `app/usecase/utils/jwt_test.go` - JWTサービステスト

### 実装詳細

#### JWTService

- `NewJWTService` - シークレットキーを受け取りJWTServiceを初期化
- `GenerateAccessToken` - アクセストークンを生成（有効期限15分）
- `GenerateRefreshToken` - リフレッシュトークンを生成（有効期限7日）
- `GenerateTokenPair` - アクセストークンとリフレッシュトークンのペアを生成
- `ValidateToken` - トークンを検証しクレームを返却

#### JWTClaims

カスタムクレームとして以下を含む:
- `user_id` - アプリケーション内のユーザーID
- `firebase_uid` - Firebase Authentication UID
- `token_type` - トークン種別（access/refresh）
- 標準クレーム（exp, iat, nbf, iss, sub）

### 依存関係追加

- `github.com/golang-jwt/jwt/v5` - JWT生成・検証ライブラリ

## 依存関係

- Depends on: SLEEVE-112-3（User Domain層のエラー定義）
- Blocks: SLEEVE-112-6（RegisterUser UseCase層の実装）

## テスト

- [x] アクセストークン生成が成功すること
- [x] リフレッシュトークン生成が成功すること
- [x] 有効なアクセストークンの検証が成功すること
- [x] 有効なリフレッシュトークンの検証が成功すること
- [x] 無効なトークンの検証がエラーを返すこと
- [x] 異なるシークレットで署名されたトークンの検証がエラーを返すこと
- [x] アクセストークンの有効期限が正しいこと（15分）
- [x] リフレッシュトークンの有効期限が正しいこと（7日）
- [x] トークンペア生成が成功すること

```bash
# テスト実行結果
docker exec sleeve-be go test ./usecase/utils/... -v
# 全9テストPASS
```

## 関連チケット

- Jira: SLEEVE-112-5

## チェックリスト

- [x] コーディング規約に準拠
- [x] 単体テストが通過
- [x] TDD（テスト駆動開発）の原則に従って実装
- [x] Linterエラーなし
