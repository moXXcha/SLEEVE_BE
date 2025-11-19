# ユーザー機能 - エラーメッセージ

このドキュメントは、ユーザー機能で使用されるエラーメッセージの定義と、それらがいつ・どのように出力されるかを説明します。

---

## ErrInvalidEmail

- **メッセージ**: "メールアドレスの形式が不正です"
- **出力タイミング**: メールアドレスのバリデーションが失敗した場合
- **関連関数**:
  - `new_email` (app/domain/models/code_models/email.go)
  - `RegisterUser` (app/usecase/user/register_user_usecase.go)
- **HTTPステータス**: 400 Bad Request
- **エラーコード**: `INVALID_EMAIL`
- **想定されるケース**:
  - Email形式が不正（例: `invalid`, `@example.com`, `user@`）
  - 空文字列
  - 特殊文字が含まれている

---

## ErrWeakPassword

- **メッセージ**: "パスワードは8文字以上で、英字・数字・記号を含む必要があります"
- **出力タイミング**: パスワードの強度チェックが失敗した場合
- **関連関数**:
  - `new_password` (app/domain/models/code_models/password.go)
  - `RegisterUser` (app/usecase/user/register_user_usecase.go)
- **HTTPステータス**: 400 Bad Request
- **エラーコード**: `WEAK_PASSWORD`
- **想定されるケース**:
  - 8文字未満
  - 英字のみ
  - 数字のみ
  - 記号が含まれていない

---

## ErrDuplicateEmail

- **メッセージ**: "このメールアドレスは既に登録されています"
- **出力タイミング**: Email重複チェックで既に登録済みのEmailが検出された場合
- **関連関数**:
  - `CheckEmailExists` (app/repository/external/firebase/user_repository.go)
  - `RegisterUser` (app/usecase/user/register_user_usecase.go)
- **HTTPステータス**: 409 Conflict
- **エラーコード**: `DUPLICATE_EMAIL`
- **想定されるケース**:
  - Firebase Authenticationに既に同じEmailが登録されている
  - 自社DBに既に同じEmailが登録されている

---

## ErrUserNotFound

- **メッセージ**: "ユーザーが見つかりません"
- **出力タイミング**: ユーザーIDで検索した際に該当するユーザーが存在しない場合
- **関連関数**:
  - `FindUserByID` (app/repository/internal/user_dao.go)
  - 今後のログイン機能・ユーザー取得機能で使用予定
- **HTTPステータス**: 404 Not Found
- **エラーコード**: `USER_NOT_FOUND`
- **想定されるケース**:
  - 指定されたユーザーIDがDBに存在しない
  - 論理削除済みのユーザー

---

## ErrFirebaseAuthFailed

- **メッセージ**: "Firebase Authenticationでエラーが発生しました"
- **出力タイミング**: Firebase Authenticationとの連携でエラーが発生した場合
- **関連関数**:
  - `CreateUserWithEmailPassword` (app/repository/external/firebase/user_repository.go)
- **HTTPステータス**: 500 Internal Server Error
- **エラーコード**: `FIREBASE_AUTH_FAILED`
- **想定されるケース**:
  - Firebase APIとの通信エラー
  - Firebase側のサーバーエラー
  - サービスアカウントキーの問題
  - Firebase プロジェクトの設定問題

---

## ErrDatabaseError

- **メッセージ**: "データベースエラーが発生しました"
- **出力タイミング**: DB操作でエラーが発生した場合
- **関連関数**:
  - `SaveUser` (app/repository/internal/user_dao.go)
- **HTTPステータス**: 500 Internal Server Error
- **エラーコード**: `DATABASE_ERROR`
- **想定されるケース**:
  - DB接続エラー
  - SQL実行エラー
  - トランザクションエラー
  - 制約違反（ユニーク制約など）

---

## ErrJWTGenerationFailed

- **メッセージ**: "トークン生成に失敗しました"
- **出力タイミング**: JWT生成処理でエラーが発生した場合
- **関連関数**:
  - `GenerateAccessToken` (app/usecase/utils/jwt.go)
  - `GenerateRefreshToken` (app/usecase/utils/jwt.go)
- **HTTPステータス**: 500 Internal Server Error
- **エラーコード**: `JWT_GENERATION_FAILED`
- **想定されるケース**:
  - JWT秘密鍵が未設定
  - JWT秘密鍵の形式が不正
  - 署名アルゴリズムのエラー

---

## エラーハンドリングのガイドライン

### クライアント側のエラー（4xx）

- **ErrInvalidEmail**: ユーザーに入力フォームの修正を促す
- **ErrWeakPassword**: パスワード要件を明示して再入力を促す
- **ErrDuplicateEmail**: ログイン画面への誘導、またはパスワードリセットの提案
- **ErrUserNotFound**: 入力内容の確認を促す

### サーバー側のエラー（5xx）

- **ErrFirebaseAuthFailed**: システム管理者に通知、ユーザーには一時的なエラーメッセージを表示
- **ErrDatabaseError**: システム管理者に通知、ユーザーには一時的なエラーメッセージを表示
- **ErrJWTGenerationFailed**: システム管理者に通知、ユーザーには一時的なエラーメッセージを表示

---

## 実装例

### Domain層でのエラー定義（app/domain/errors/user_errors.go）

```go
package errors

import "fmt"

var (
	ErrInvalidEmail        = fmt.Errorf("invalid email format")
	ErrWeakPassword        = fmt.Errorf("password does not meet security requirements")
	ErrDuplicateEmail      = fmt.Errorf("email already exists")
	ErrUserNotFound        = fmt.Errorf("user not found")
	ErrFirebaseAuthFailed  = fmt.Errorf("firebase authentication error")
	ErrDatabaseError       = fmt.Errorf("database error")
	ErrJWTGenerationFailed = fmt.Errorf("jwt generation failed")
)

func WrapUserError(domain_err error, wrapped_err error) error {
	return fmt.Errorf("%w: %v", domain_err, wrapped_err)
}
```

### UseCase層でのエラーメッセージ使用（app/messages/user/errors.go）

```go
package user

const (
	MsgInvalidEmail        = "メールアドレスの形式が不正です"
	MsgWeakPassword        = "パスワードは8文字以上で、英字・数字・記号を含む必要があります"
	MsgDuplicateEmail      = "このメールアドレスは既に登録されています"
	MsgUserNotFound        = "ユーザーが見つかりません"
	MsgFirebaseAuthFailed  = "Firebase Authenticationでエラーが発生しました"
	MsgDatabaseError       = "データベースエラーが発生しました"
	MsgJWTGenerationFailed = "トークン生成に失敗しました"
)
```

---

## 関連ドキュメント

- **ログメッセージ**: `docs/messages/user/logs.md`
- **Domain層エラー定義**: `app/domain/errors/user_errors.go`
- **UseCase層メッセージ定義**: `app/messages/user/errors.go`
