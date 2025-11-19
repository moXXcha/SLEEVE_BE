# ユーザー機能 - ログメッセージ

このドキュメントは、ユーザー機能で使用されるログメッセージの定義と、それらがいつ・どのように出力されるかを説明します。

---

## LogUserRegistrationStarted

- **メッセージ**: "ユーザー登録処理を開始: email=%s"
- **出力タイミング**: RegisterUser UseCaseの開始時
- **ログレベル**: Info
- **関連関数**: `RegisterUser` (app/usecase/user/register_user_usecase.go)
- **含まれる情報**:
  - メールアドレス
- **用途**:
  - ユーザー登録処理の開始を記録
  - トラブルシューティング時の処理フロー追跡

---

## LogUserRegistrationSucceeded

- **メッセージ**: "ユーザー登録が成功: user_id=%s, firebase_uid=%s"
- **出力タイミング**: ユーザー登録が成功した時
- **ログレベル**: Info
- **関連関数**: `RegisterUser` (app/usecase/user/register_user_usecase.go)
- **含まれる情報**:
  - 自社ユーザーID
  - Firebase UID
- **用途**:
  - ユーザー登録成功の記録
  - ユーザー数の統計分析
  - 監査ログ

---

## LogUserRegistrationFailed

- **メッセージ**: "ユーザー登録が失敗: email=%s, error=%v"
- **出力タイミング**: ユーザー登録が失敗した時
- **ログレベル**: Error
- **関連関数**: `RegisterUser` (app/usecase/user/register_user_usecase.go)
- **含まれる情報**:
  - メールアドレス
  - エラー内容
- **用途**:
  - エラー発生の記録
  - トラブルシューティング
  - エラー傾向の分析

---

## LogFirebaseUserCreated

- **メッセージ**: "Firebase Authenticationにユーザーを作成: firebase_uid=%s"
- **出力タイミング**: Firebase Authenticationへのユーザー登録が成功した時
- **ログレベル**: Info
- **関連関数**: `CreateUserWithEmailPassword` (app/repository/external/firebase/user_repository.go)
- **含まれる情報**:
  - Firebase UID
- **用途**:
  - Firebase側のユーザー作成を記録
  - Firebase AuthenticationとDB間の整合性確認

---

## LogDatabaseUserSaved

- **メッセージ**: "DBにユーザー情報を保存: user_id=%s, firebase_uid=%s"
- **出力タイミング**: 自社DBへのユーザー情報保存が成功した時
- **ログレベル**: Info
- **関連関数**: `SaveUser` (app/repository/internal/user_dao.go)
- **含まれる情報**:
  - 自社ユーザーID
  - Firebase UID
- **用途**:
  - DB保存の記録
  - Firebase AuthenticationとDB間の整合性確認

---

## LogJWTGenerated

- **メッセージ**: "JWTを生成: user_id=%s, token_type=%s"
- **出力タイミング**: JWT生成が成功した時
- **ログレベル**: Info
- **関連関数**:
  - `GenerateAccessToken` (app/usecase/utils/jwt.go)
  - `GenerateRefreshToken` (app/usecase/utils/jwt.go)
- **含まれる情報**:
  - ユーザーID
  - トークン種別（access / refresh）
- **用途**:
  - トークン生成の記録
  - セキュリティ監査

---

## LogEmailValidationFailed

- **メッセージ**: "メールアドレスのバリデーションが失敗: email=%s"
- **出力タイミング**: メールアドレスのバリデーションが失敗した時
- **ログレベル**: Warn
- **関連関数**: `new_email` (app/domain/models/code_models/email.go)
- **含まれる情報**:
  - 不正なメールアドレス
- **用途**:
  - 不正な入力の記録
  - 攻撃パターンの検出

---

## LogPasswordValidationFailed

- **メッセージ**: "パスワードのバリデーションが失敗: reason=%s"
- **出力タイミング**: パスワードのバリデーションが失敗した時
- **ログレベル**: Warn
- **関連関数**: `new_password` (app/domain/models/code_models/password.go)
- **含まれる情報**:
  - 失敗理由（文字数不足、記号なし、など）
- **用途**:
  - 不正な入力の記録
  - ユーザー体験改善のためのデータ収集
- **注意**: パスワード自体はログに出力しない

---

## LogEmailDuplicateCheckStarted

- **メッセージ**: "Email重複チェックを開始: email=%s"
- **出力タイミング**: Email重複チェックの開始時
- **ログレベル**: Debug
- **関連関数**: `CheckEmailExists` (app/repository/external/firebase/user_repository.go)
- **含まれる情報**:
  - メールアドレス
- **用途**:
  - デバッグ時の処理フロー追跡

---

## LogEmailDuplicateDetected

- **メッセージ**: "Email重複を検出: email=%s"
- **出力タイミング**: Email重複が検出された時
- **ログレベル**: Warn
- **関連関数**: `CheckEmailExists` (app/repository/external/firebase/user_repository.go)
- **含まれる情報**:
  - 重複しているメールアドレス
- **用途**:
  - 重複登録の試行を記録
  - 不正行為の検出

---

## LogFirebaseAuthError

- **メッセージ**: "Firebase Authenticationエラー: error=%v"
- **出力タイミング**: Firebase Authenticationとの連携でエラーが発生した時
- **ログレベル**: Error
- **関連関数**: `CreateUserWithEmailPassword` (app/repository/external/firebase/user_repository.go)
- **含まれる情報**:
  - エラー内容
- **用途**:
  - Firebase連携エラーの記録
  - システム障害の検出

---

## LogDatabaseError

- **メッセージ**: "データベースエラー: operation=%s, error=%v"
- **出力タイミング**: DB操作でエラーが発生した時
- **ログレベル**: Error
- **関連関数**: `SaveUser` (app/repository/internal/user_dao.go)
- **含まれる情報**:
  - 操作種別（insert, update, delete, select）
  - エラー内容
- **用途**:
  - DB操作エラーの記録
  - システム障害の検出

---

## LogJWTGenerationError

- **メッセージ**: "JWT生成エラー: token_type=%s, error=%v"
- **出力タイミング**: JWT生成処理でエラーが発生した時
- **ログレベル**: Error
- **関連関数**:
  - `GenerateAccessToken` (app/usecase/utils/jwt.go)
  - `GenerateRefreshToken` (app/usecase/utils/jwt.go)
- **含まれる情報**:
  - トークン種別（access / refresh）
  - エラー内容
- **用途**:
  - JWT生成エラーの記録
  - セキュリティ問題の検出

---

## ログレベルの使い分け

### Debug
- 開発環境でのデバッグ用
- 詳細な処理フローの追跡
- 本番環境では通常出力しない

### Info
- 正常な処理の記録
- ユーザー登録成功、JWT生成成功など
- 監査ログとして使用

### Warn
- エラーではないが注意が必要な事象
- バリデーション失敗、重複登録の試行など
- 不正行為の検出に使用

### Error
- エラー発生時
- システム障害の可能性がある事象
- 即座に対応が必要な問題

---

## セキュリティ考慮事項

### ログに出力してはいけない情報

- ❌ パスワード（平文）
- ❌ JWT（トークン本体）
- ❌ Firebase秘密鍵
- ❌ データベース接続文字列

### ログに出力してもよい情報

- ✅ メールアドレス（個人情報だが、識別に必要）
- ✅ ユーザーID
- ✅ Firebase UID
- ✅ エラーメッセージ
- ✅ 処理開始・終了のタイムスタンプ

---

## 実装例

### UseCase層でのログメッセージ使用（app/messages/user/logs.go）

```go
package user

const (
	LogUserRegistrationStarted     = "ユーザー登録処理を開始: email=%s"
	LogUserRegistrationSucceeded   = "ユーザー登録が成功: user_id=%s, firebase_uid=%s"
	LogUserRegistrationFailed      = "ユーザー登録が失敗: email=%s, error=%v"
	LogFirebaseUserCreated         = "Firebase Authenticationにユーザーを作成: firebase_uid=%s"
	LogDatabaseUserSaved           = "DBにユーザー情報を保存: user_id=%s, firebase_uid=%s"
	LogJWTGenerated                = "JWTを生成: user_id=%s, token_type=%s"
	LogEmailValidationFailed       = "メールアドレスのバリデーションが失敗: email=%s"
	LogPasswordValidationFailed    = "パスワードのバリデーションが失敗: reason=%s"
	LogEmailDuplicateCheckStarted  = "Email重複チェックを開始: email=%s"
	LogEmailDuplicateDetected      = "Email重複を検出: email=%s"
	LogFirebaseAuthError           = "Firebase Authenticationエラー: error=%v"
	LogDatabaseError               = "データベースエラー: operation=%s, error=%v"
	LogJWTGenerationError          = "JWT生成エラー: token_type=%s, error=%v"
)
```

### 使用例

```go
// 登録開始
log.Infof(user.LogUserRegistrationStarted, email)

// 登録成功
log.Infof(user.LogUserRegistrationSucceeded, user_id, firebase_uid)

// 登録失敗
log.Errorf(user.LogUserRegistrationFailed, email, err)
```

---

## 関連ドキュメント

- **エラーメッセージ**: `docs/messages/user/errors.md`
- **UseCase層ログメッセージ定義**: `app/messages/user/logs.go`
