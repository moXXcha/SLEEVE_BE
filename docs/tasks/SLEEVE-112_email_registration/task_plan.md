# SLEEVE-112: Email認証（アカウント登録）機能の実装

## 作業内容確認

新規ユーザーが Email と Password を入力してアカウントを登録できる機能をGraphQL API（gqlgen）側に実装する。登録後はJWTを発行し、FEに返却する。

### 目標（Goal）
- GraphQL Mutation 経由でアカウント登録を受け付ける
- Firebase Authentication にユーザーを登録
- 自社DBにもユーザー情報を保存（Firebase UIDと紐付け）
- 成功時に JWT + User 情報を返却
- 不正入力や登録失敗時に適切なエラーメッセージを返却

### 非目標（Non-Goal）
- パスワードリセット機能
- OAuth 認証
- Login 機能（既存アカウントのログイン）
- User プロフィール編集機能
- sign-in機能の実装

---

## ブランチ戦略

このタスクは複数のサブタスクに分割して進めます。

### タスクの分割判断

- [x] **複数サブタスクに分割**: Topic Branchを作成し、6つのSubtaskに分ける

### 親タスク（Story/Task）

- **Jira**: `SLEEVE-112: Email認証（アカウント登録）機能の実装`
- **対応ブランチ**: `topic/SLEEVE-100_email_registration`
- **分岐元**: `main`
- **マージ先**: `main`
- **マージ方法**: Merge commit

### サブタスク一覧

| Subtask ID | タスク名 | 依存関係 | 分岐元 | マージ先 | マージ方法 |
|-----------|---------|---------|--------|---------|-----------|
| SLEEVE-112-1 | usersテーブルのマイグレーション | なし | topic/SLEEVE-100 | topic/SLEEVE-100 | Squash and merge |
| SLEEVE-112-2 | Firebase SDK初期化とセットアップ | なし | topic/SLEEVE-100 | topic/SLEEVE-100 | Squash and merge |
| SLEEVE-112-3 | User Domain層の実装 | SLEEVE-101に依存 | topic/SLEEVE-100 | topic/SLEEVE-100 | Squash and merge |
| SLEEVE-112-4 | User Repository層の実装（Firebase + DB連携） | SLEEVE-101, 102, 103に依存 | topic/SLEEVE-100 | topic/SLEEVE-100 | Squash and merge |
| SLEEVE-112-5 | JWT発行機能の実装 | なし | topic/SLEEVE-100 | topic/SLEEVE-100 | Squash and merge |
| SLEEVE-112-6 | RegisterUser UseCase層の実装 | SLEEVE-103, 104, 105に依存 | topic/SLEEVE-100 | topic/SLEEVE-100 | Squash and merge |
| SLEEVE-112-7 | RegisterUser GraphQL層の実装 | SLEEVE-106に依存 | topic/SLEEVE-100 | topic/SLEEVE-100 | Squash and merge |

### ブランチ構造図

```
main
 └─ topic/SLEEVE-112_email_registration
      ├─ feature/SLEEVE-112-1_users_migration
      ├─ feature/SLEEVE-112-2_firebase_setup
      ├─ feature/SLEEVE-112-3_user_domain
      ├─ feature/SLEEVE-112-4_user_repository
      ├─ feature/SLEEVE-112-5_jwt_service
      ├─ feature/SLEEVE-112-6_register_usecase
      └─ feature/SLEEVE-112-7_register_graphql
```

### マージ順序

1. `feature/SLEEVE-112-1_users_migration` → `topic/SLEEVE-112`（最優先）
2. `feature/SLEEVE-112-2_firebase_setup` → `topic/SLEEVE-112`
3. `feature/SLEEVE-112-3_user_domain` → `topic/SLEEVE-112`
4. `feature/SLEEVE-112-4_user_repository` → `topic/SLEEVE-112`
5. `feature/SLEEVE-112-5_jwt_service` → `topic/SLEEVE-112`
6. `feature/SLEEVE-112-6_register_usecase` → `topic/SLEEVE-112`
7. `feature/SLEEVE-112-7_register_graphql` → `topic/SLEEVE-112`
8. 全Subtask完了後: `topic/SLEEVE-112` → `main`

---

## 実装方針

### TDD（テスト駆動開発）の徹底

**全てのSubtaskにおいて、TDD（テスト駆動開発）の原則に従って実装を進めます。**

#### TDDサイクル（Red → Green → Refactor）

**SLEEVE-101（Migration）を除く全てのSubtask**で、以下のサイクルを厳守してください：

1. 🔴 **Red（テストを先に書く）**
   - 実装コードを書く前に、まず単体テストを作成する
   - この時点ではテストは失敗する（コンパイルエラーまたはテスト失敗）
   - テストが失敗することを確認してからコミット

2. 🟢 **Green（最小限の実装でテストを通す）**
   - テストを通過させるための最小限の実装を行う
   - テストが全て通過することを確認
   - 実装コードをコミット

3. 🔵 **Refactor（リファクタリング）**
   - 必要に応じてコードを改善
   - テストが引き続き通過することを確認
   - リファクタリングした場合はコミット

#### TDD対象外のSubtask

- **SLEEVE-101（usersテーブルのマイグレーション）**: Migration処理のため、TDD対象外

#### TDDの重要性

- **品質保証**: テストファーストにより、要件を満たすコードを確実に実装
- **設計改善**: テストを先に書くことで、APIの使いやすさを考慮した設計になる
- **リファクタリング安全性**: テストがあることで、安心してリファクタリングできる
- **ドキュメント**: テストコードが実装の使い方を示すドキュメントになる

**重要**: TDD以外の方法で実装を進めないでください。必ずテストを先に書いてから実装してください。

---

## Subtask完了時の必須手順

**各Subtaskの実装が完了したら、必ず以下の手順を実行してください。この手順を飛ばすと作業手順が狂う重要なステップです。**

### 1. PR作成前の最終確認

```bash
# コミット漏れがないことを確認
git status
# 期待される出力: "nothing to commit, working tree clean"

# 全テストが通過することを確認
cd app
go test ./... -v
cd ..

# Linterエラーがないことを確認
cd app
golangci-lint run --fix
cd ..
```

### 2. リモートへプッシュ

```bash
git push origin feature/SLEEVE-XXX-X_[機能名]
```

**プッシュ前の確認事項**:
- [ ] `git status` で未コミットの変更がない
- [ ] 全ての単体テストが通過している
- [ ] `golangci-lint run --fix` でLinterエラーがない
- [ ] DB変更を行った場合は `docs/db_schema.md` を更新している

### 3. PRディスクリプションファイルの作成

`docs/tasks/SLEEVE-112_email_registration/pr_description_SLEEVE-XXX-X.md` を作成

### 4. GitHub CLIでPR作成

```bash
gh pr create \
  --base topic/SLEEVE-112_email_registration \
  --title "SLEEVE-XXX-X: [タスク名]" \
  --body-file docs/tasks/SLEEVE-112_email_registration/pr_description_SLEEVE-XXX-X.md
```

**失敗する場合**: PR descriptionファイルを作成済みなので、作業者に手動でのPR作成を依頼してください。

### 5. レビュー待機

- ✅ **PR作成後は、オーナーのレビュー・承認・マージを完全に待つ**
- ❌ **レビュー待ちの間、他のSubtaskに進まない**
- ✅ **レビュー中に他の作業を進めない**

### 6. マージ後の処理

PRがマージされたら：

```bash
# Feature Branchを削除
git branch -d feature/SLEEVE-XXX-X_[機能名]
git push origin --delete feature/SLEEVE-XXX-X_[機能名]

# Topic Branchを最新化
git checkout topic/SLEEVE-112_email_registration
git pull origin topic/SLEEVE-112_email_registration
```

**次のSubtaskへの準備確認**:
- [ ] 前のSubtaskのFeature Branchが削除されている
- [ ] Topic Branchが最新化されている
- [ ] 作業ディレクトリがクリーンな状態になっている（`git status`）

---

## テストの通過基準

### 単体テスト（各Subtaskで実施）

#### SLEEVE-112-1: usersテーブルのマイグレーション
- マイグレーション実行が成功すること（`task migrate-up`）
- usersテーブルが作成されていること
- インデックスが正しく設定されていること

#### SLEEVE-112-2: Firebase SDK初期化
- Firebase Admin SDKが正常に初期化できること
- Firebase Authクライアントが取得できること
- 環境変数未設定時にエラーが返されること

#### SLEEVE-112-3: User Domain層
- Email ValueObject: 正しい形式のメールアドレスがバリデーションを通過すること
- Email ValueObject: 不正な形式のメールアドレスがエラーとなること
- Password ValueObject: 強度要件を満たすパスワードがバリデーションを通過すること
- Password ValueObject: 強度要件を満たさないパスワードがエラーとなること
- User Entity: FirebaseUIDとEmailを持つUserが正常に生成されること

#### SLEEVE-112-4: User Repository層
- Firebase Authenticationにユーザーを登録できること
- 自社DBにユーザー情報を保存できること（Firebase UIDと紐付け）
- Email重複チェックが正常に動作すること
- Firebase Authenticationのエラーがドメインエラーに変換されること

#### SLEEVE-112-5: JWT発行機能
- アクセストークンが正常に生成されること
- リフレッシュトークンが正常に生成されること
- トークンのバリデーションが正常に動作すること
- 有効期限が正しく設定されていること

#### SLEEVE-112-6: RegisterUser UseCase層
- 正しいEmail/Passwordで登録が成功すること
- Email重複時にエラーが返されること
- 不正なEmailでエラーが返されること
- 弱いPasswordでエラーが返されること
- 登録成功時にJWTが発行されること

#### SLEEVE-112-7: RegisterUser GraphQL層
- registerUserミューテーションが正常に動作すること
- 入力バリデーションが正常に動作すること
- エラーレスポンスが適切に返されること

### 結合テスト（全Subtask完了後に実施）

- **テストケース 1: 正常なユーザー登録フロー**
  - 通過条件:
    - 正しいEmail/Passwordを渡した場合、Firebase Authenticationに新規ユーザーが登録される
    - 自社DBにもユーザー情報が保存される（Firebase UIDと紐付け）
    - JWT（アクセストークン + リフレッシュトークン）が返される
    - User情報が返される

- **テストケース 2: Email重複時のエラーハンドリング**
  - 通過条件:
    - 既に登録済みのEmailで登録を試みた場合、適切なエラーメッセージが返される
    - Firebase Authenticationにユーザーが登録されない
    - 自社DBにユーザー情報が保存されない

- **テストケース 3: 不正なEmail形式のバリデーション**
  - 通過条件:
    - 不正な形式のEmailで登録を試みた場合、適切なエラーメッセージが返される
    - Firebase Authenticationにユーザーが登録されない

- **テストケース 4: 弱いPasswordのバリデーション**
  - 通過条件:
    - 強度要件を満たさないPasswordで登録を試みた場合、適切なエラーメッセージが返される
    - Firebase Authenticationにユーザーが登録されない

- **テストケース 5: JWTの有効性確認**
  - 通過条件:
    - 発行されたアクセストークンが正しく検証できること
    - 発行されたリフレッシュトークンが正しく検証できること
    - トークンにFirebase UIDが含まれていること

---

## メッセージ定義

この機能で使用するエラーメッセージとログメッセージを `docs/messages/user/` に記述します。

### エラーメッセージ

`docs/messages/user/errors.md` に以下の形式で記述：

#### ErrInvalidEmail
- **メッセージ**: "メールアドレスの形式が不正です"
- **出力タイミング**: メールアドレスのバリデーションが失敗した場合
- **関連関数**: `new_email`, `RegisterUser`
- **HTTPステータス**: 400

#### ErrWeakPassword
- **メッセージ**: "パスワードは8文字以上で、英字・数字・記号を含む必要があります"
- **出力タイミング**: パスワードの強度チェックが失敗した場合
- **関連関数**: `new_password`, `RegisterUser`
- **HTTPステータス**: 400

#### ErrDuplicateEmail
- **メッセージ**: "このメールアドレスは既に登録されています"
- **出力タイミング**: Email重複チェックで既に登録済みのEmailが検出された場合
- **関連関数**: `CheckEmailExists`, `RegisterUser`
- **HTTPステータス**: 409

#### ErrUserNotFound
- **メッセージ**: "ユーザーが見つかりません"
- **出力タイミング**: ユーザーIDで検索した際に該当するユーザーが存在しない場合
- **関連関数**: `FindUserByID`
- **HTTPステータス**: 404

#### ErrFirebaseAuthFailed
- **メッセージ**: "Firebase Authenticationでエラーが発生しました"
- **出力タイミング**: Firebase Authenticationとの連携でエラーが発生した場合
- **関連関数**: `CreateUserWithEmailPassword`
- **HTTPステータス**: 500

#### ErrDatabaseError
- **メッセージ**: "データベースエラーが発生しました"
- **出力タイミング**: DB操作でエラーが発生した場合
- **関連関数**: `SaveUser`
- **HTTPステータス**: 500

#### ErrJWTGenerationFailed
- **メッセージ**: "トークン生成に失敗しました"
- **出力タイミング**: JWT生成処理でエラーが発生した場合
- **関連関数**: `GenerateAccessToken`, `GenerateRefreshToken`
- **HTTPステータス**: 500

### ログメッセージ

`docs/messages/user/logs.md` に以下の形式で記述：

#### LogUserRegistrationStarted
- **メッセージ**: "ユーザー登録処理を開始: email=%s"
- **出力タイミング**: RegisterUser UseCaseの開始時
- **ログレベル**: Info
- **関連関数**: `RegisterUser`

#### LogUserRegistrationSucceeded
- **メッセージ**: "ユーザー登録が成功: user_id=%s, firebase_uid=%s"
- **出力タイミング**: ユーザー登録が成功した時
- **ログレベル**: Info
- **関連関数**: `RegisterUser`

#### LogUserRegistrationFailed
- **メッセージ**: "ユーザー登録が失敗: email=%s, error=%v"
- **出力タイミング**: ユーザー登録が失敗した時
- **ログレベル**: Error
- **関連関数**: `RegisterUser`

#### LogFirebaseUserCreated
- **メッセージ**: "Firebase Authenticationにユーザーを作成: firebase_uid=%s"
- **出力タイミング**: Firebase Authenticationへのユーザー登録が成功した時
- **ログレベル**: Info
- **関連関数**: `CreateUserWithEmailPassword`

#### LogJWTGenerated
- **メッセージ**: "JWTを生成: user_id=%s"
- **出力タイミング**: JWT生成が成功した時
- **ログレベル**: Info
- **関連関数**: `GenerateAccessToken`, `GenerateRefreshToken`

---

## 使用技術（ライブラリなど、version まで）

### 新規追加
- `firebase.google.com/go/v4` - Firebase Admin SDK（最新版）
- `firebase.google.com/go/v4/auth` - Firebase Authentication（最新版）
- `github.com/golang-jwt/jwt/v5` - JWT生成・検証（v5.2.0以降）

### 既存（確認済み）
- `entgo.io/ent v0.14.5` - ORM
- `github.com/99designs/gqlgen v0.17.81` - GraphQL server
- `github.com/lib/pq v1.10.9` - PostgreSQL driver
- `github.com/google/uuid v1.6.0` - UUID生成
- `github.com/vektah/gqlparser/v2 v2.5.30` - GraphQL parser

---

## 作成ファイル

### SLEEVE-112-1: usersテーブルのマイグレーション
- `app/ent/schema/user.go` - entスキーマ定義
- `app/ent/user.go` - entが自動生成するファイル（複数）
- `app/migrations/YYYYMMDDHHMMSS_create_users.sql` - マイグレーションSQL
- `docs/db_schema.md` - スキーマドキュメント（更新）

### SLEEVE-112-2: Firebase SDK初期化
- `app/repository/external/firebase/init_firebase.go` - Firebase初期化処理
- `app/repository/external/firebase/init_firebase_test.go` - 初期化テスト

### SLEEVE-112-3: User Domain層
- `app/domain/models/code_models/email.go` - Email ValueObject
- `app/domain/models/code_models/email_test.go` - Emailテスト
- `app/domain/models/code_models/password.go` - Password ValueObject
- `app/domain/models/code_models/password_test.go` - Passwordテスト
- `app/domain/models/code_models/user.go` - User Entity
- `app/domain/models/code_models/user_test.go` - Userテスト
- `app/domain/errors/user_errors.go` - ドメインエラー定義
- `app/domain/errors/user_errors_test.go` - エラーテスト

### SLEEVE-112-4: User Repository層
- `app/repository/external/firebase/user_repository.go` - Firebase連携Repository
- `app/repository/external/firebase/user_repository_test.go` - Repositoryテスト
- `app/repository/internal/user_dao.go` - DB操作DAO
- `app/repository/internal/user_dao_test.go` - DAOテスト

### SLEEVE-112-5: JWT発行機能
- `app/usecase/utils/jwt.go` - JWT生成・検証処理
- `app/usecase/utils/jwt_test.go` - JWTテスト

### SLEEVE-112-6: RegisterUser UseCase層
- `app/usecase/user/register_user_usecase.go` - RegisterUser UseCase
- `app/usecase/user/register_user_usecase_test.go` - UseCaseテスト
- `app/messages/user/errors.go` - エラーメッセージ定義（実装ファイル）
- `app/messages/user/logs.go` - ログメッセージ定義（実装ファイル）
- `docs/messages/user/errors.md` - エラーメッセージドキュメント
- `docs/messages/user/logs.md` - ログメッセージドキュメント

### SLEEVE-112-7: RegisterUser GraphQL層
- `app/graph/schema.graphqls` - GraphQLスキーマ（registerUserミューテーション追加）
- `app/graph/schema.resolvers.go` - リゾルバー実装（registerUser追加）
- `app/graph/schema.resolvers_test.go` - リゾルバーテスト

---

## 変更ファイル

- `app/go.mod` - 依存関係追加（Firebase SDK, JWT）
- `app/go.sum` - 依存関係チェックサム
- `app/graph/resolver.go` - 依存注入の追加（RegisterUserUseCase）

---

## 影響範囲

### 新規機能のため、既存機能への影響は最小限

- **GraphQL Schema**: 新しいミューテーション（`registerUser`）を追加するが、既存のクエリ・ミューテーションには影響なし
- **データベース**: 新しいテーブル（`users`）を追加するが、既存テーブルには影響なし
- **依存関係**: Firebase SDK、JWTライブラリを追加するが、既存の依存関係とは競合しない

### 今後の拡張への影響

- **ログイン機能**: 本機能で作成したUser Entityとjwt.goが基盤となる
- **プロフィール編集**: usersテーブルにカラム追加が必要
- **OAuth認証**: Firebase Authenticationの他のプロバイダを追加することで対応可能
