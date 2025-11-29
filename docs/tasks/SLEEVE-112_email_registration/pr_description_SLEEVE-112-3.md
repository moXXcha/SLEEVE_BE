## 概要

User Domain層の実装を行いました。Email/Password ValueObjectとUser Entity、ドメインエラーを定義しています。

## 変更内容

### 新規作成ファイル

- `app/domain/models/code_models/email.go` - Email ValueObject
- `app/domain/models/code_models/email_test.go` - Emailテスト
- `app/domain/models/code_models/password.go` - Password ValueObject
- `app/domain/models/code_models/password_test.go` - Passwordテスト
- `app/domain/models/code_models/user.go` - User Entity
- `app/domain/models/code_models/user_test.go` - Userテスト
- `app/domain/errors/user_errors.go` - ドメインエラー定義
- `app/domain/errors/user_errors_test.go` - エラーテスト

### 実装詳細

#### Email ValueObject
- メールアドレスの形式バリデーション（正規表現）
- 空文字・スペース含むメールアドレスの拒否
- `Value()`, `String()`, `Equals()` メソッド

#### Password ValueObject
- 8文字以上の長さチェック
- 大文字・小文字・数字・記号の含有チェック
- `Value()`, `String()`（マスク表示）, `Equals()` メソッド

#### User Entity
- `id` (UUID), `firebase_uid`, `email`, `created_at`, `updated_at`, `deleted_at` フィールド
- `NewUser()` - 新規ユーザー作成
- `NewUserWithID()` - DBからの復元用
- `IsDeleted()` - 論理削除判定

#### ドメインエラー
- `ErrInvalidEmail` - メールアドレス形式不正
- `ErrWeakPassword` - パスワード強度不足
- `ErrDuplicateEmail` - メールアドレス重複
- `ErrUserNotFound` - ユーザー未検出
- `ErrFirebaseAuthFailed` - Firebase認証エラー
- `ErrDatabaseError` - DBエラー
- `ErrJWTGenerationFailed` - JWT生成失敗

## 依存関係

- Depends on: なし（SLEEVE-112-1, SLEEVE-112-2はマージ済み）
- Blocks: SLEEVE-112-4（User Repository層の実装）

## テスト

- [x] Email ValueObject: 正しい形式のメールアドレスがバリデーションを通過すること
- [x] Email ValueObject: 不正な形式のメールアドレスがエラーとなること
- [x] Password ValueObject: 強度要件を満たすパスワードがバリデーションを通過すること
- [x] Password ValueObject: 強度要件を満たさないパスワードがエラーとなること
- [x] User Entity: FirebaseUIDとEmailを持つUserが正常に生成されること
- [x] Domain Errors: 各エラーが正しいメッセージを返すこと
- [x] Domain Errors: `IsUserDomainError()` が正しく判定すること

```bash
# テスト実行結果
cd app
go test ./domain/models/code_models/... -v
go test ./domain/errors/... -v
# 全テストPASS
```

## 関連チケット

- Jira: SLEEVE-112-3

## チェックリスト

- [x] コーディング規約に準拠
- [x] 単体テストが通過
- [x] TDD（テスト駆動開発）の原則に従って実装
