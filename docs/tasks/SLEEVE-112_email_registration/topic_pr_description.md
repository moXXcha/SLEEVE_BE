## 概要

Email認証（アカウント登録）機能を実装しました。GraphQL APIからEmail/Passwordでユーザー登録を行い、Firebase Authenticationと自社DBにユーザー情報を保存し、JWTを発行して返却します。

## 完了したSubtask

- [x] SLEEVE-112-1: usersテーブルのマイグレーション
- [x] SLEEVE-112-2: Firebase SDK初期化とセットアップ
- [x] SLEEVE-112-3: User Domain層の実装
- [x] SLEEVE-112-4: User Repository層の実装（Firebase + DB連携）
- [x] SLEEVE-112-5: JWT発行機能の実装
- [x] SLEEVE-112-6: RegisterUser UseCase層の実装
- [x] SLEEVE-112-7: RegisterUser GraphQL層の実装

## 変更内容

### データベース
- `users`テーブルを新規作成（public_id, firebase_uid, email, created_at, updated_at, deleted_at）

### Domain層
- `Email` ValueObject - メールアドレス形式バリデーション
- `Password` ValueObject - パスワード強度バリデーション（8文字以上、大文字・小文字・数字・記号必須）
- `User` Entity - ユーザードメインモデル
- ドメインエラー定義（ErrInvalidEmail, ErrWeakPassword, ErrDuplicateEmail等）

### Repository層
- Firebase User Repository - Firebase Authentication連携（CreateUser, CheckEmailExists, DeleteUser）
- User DAO - 自社DB連携（Save, FindByPublicID, FindByFirebaseUID, FindByEmail, ExistsByEmail）

### UseCase層
- JWT Service - アクセストークン（15分）/リフレッシュトークン（7日）の生成・検証
- RegisterUser UseCase - ユーザー登録ビジネスロジック（バリデーション→Firebase登録→DB保存→JWT発行）

### GraphQL層
- `registerUser` Mutation - ユーザー登録API

## API仕様

```graphql
mutation {
  registerUser(input: {
    email: "user@example.com"
    password: "Password123!"
  }) {
    user {
      id
      email
    }
    tokens {
      accessToken
      refreshToken
    }
  }
}
```

## テスト

- [x] 全Subtaskの単体テストが通過（計50+テスト）
- [x] 結合テストが通過（5テストケース）
  - 正常なユーザー登録フロー
  - Email重複時のエラーハンドリング
  - 不正なEmail形式のバリデーション
  - 弱いPasswordのバリデーション
  - JWTの有効性確認

```bash
# 全テスト実行
docker exec sleeve-be go test ./... -v
# 結果: 全てPASS
```

## 関連チケット

- Jira: SLEEVE-112

## チェックリスト

- [x] 全Subtaskがマージ済み
- [x] コーディング規約に準拠
- [x] Linterエラーがない
- [x] TDD（テスト駆動開発）の原則に従って実装
- [x] DB変更を行った場合は `docs/db_schema.md` を更新している
