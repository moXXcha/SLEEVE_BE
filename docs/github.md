
# GitHub 運用ルール

## 概要

このドキュメントは、SLEEVE プロジェクトの GitHub 運用に関するルールとガイドラインを定義します。
AI エージェントが一貫性のあるコミットとプルリクエスト管理を行うための規約です。

## 関連ドキュメント

- **全体フロー**: `docs/general.md` - 開発タスクの全体フロー
- **実装手順**: `docs/task_procedure.md` - 詳細な実装手順とテンプレート
- **コーディング規約**: `docs/BE_coding_roles.md` - アーキテクチャとコーディング規約
- **コーディングスタイル**: `docs/coding_guide.md` - コーディングスタイルガイド
- **DB 操作**: `docs/DB_manual.md` - データベース操作手順
- **Linter/Formatter**: `docs/lint_format_manual.md` - コード品質管理

## ブランチ戦略

### 基本方針

- **Topic Branch = 1つの機能単位**: Topic Branchは1つの完結した機能（例: ユーザー登録機能）
- **各PR = レイヤー × 操作単位**: 1つのレイヤーの1操作実装、単体テストで動作確認可能
- **Topic Branch運用**: 大きなタスクは複数のサブタスクに分割し、Topic Branchを親として管理する
- **依存関係の許容**: PRの依存関係を許容し、依存先を親としてブランチを切る
- **Jira連携**: ブランチごとにJiraでSubtaskを作成し、進捗管理を明確にする

**重要**: PR粒度の詳細は「PR粒度の詳細定義」セクションを参照してください。

### ブランチの種類

#### main ブランチ
- 本番環境にデプロイされる安定版（リリース済みコード）
- 直接コミット禁止
- 保護設定あり（レビュー必須、CI通過必須）

#### topic/[JiraのID]_[機能名]
- **用途**: 1つの完結した機能の親ブランチ（例: ユーザー登録機能、ユーザー取得機能）
- **分岐元**: main
- **マージ先**: main
- **作成タイミング**:
  - 1つの機能をレイヤー単位で分割して実装する場合
  - 通常3〜4個のSubtaskに分かれる（Domain、Repository、UseCase、GraphQL）
- **対応Jira**: Story/Task
- **例**: `topic/SLEEVE-100_user_registration`、`topic/SLEEVE-200_user_retrieval`

#### feature/[JiraのID]_[機能名]
- **用途**: 実際の機能実装・修正を行うブランチ
- **分岐元**:
  - 依存関係がない場合: Topic Branch
  - 依存関係がある場合: 依存先のFeature Branch
- **マージ先**:
  - 依存関係がない場合: Topic Branch
  - 依存関係がある場合: 依存先のFeature Branch
- **対応Jira**: Subtask
- **例**: `feature/SLEEVE-101_firebase_setup`

#### hotfix/[JiraのID]_[修正内容]
- **用途**: 緊急修正用ブランチ
- **分岐元**: main
- **マージ先**: main
- **例**: `hotfix/SLEEVE-999_fix_login_bug`

### ブランチ構造の例

「ユーザー登録機能」を実装する場合（細粒度パターン）：

```
main
 └─ topic/SLEEVE-100_user_registration (1つの機能: ユーザー登録)
      ├─ feature/SLEEVE-101_user_model (Domain層)
      ├─ feature/SLEEVE-102_user_repository_create (Repository層 - Create操作)
      ├─ feature/SLEEVE-103_create_user_usecase (UseCase層)
      └─ feature/SLEEVE-104_user_registration_api (GraphQL層)
```

**マージフロー**:
1. `feature/SLEEVE-101` → `topic/SLEEVE-100` (Domain層、依存なし)
2. `feature/SLEEVE-102` → `topic/SLEEVE-100` (Repository層、101に依存)
3. `feature/SLEEVE-103` → `topic/SLEEVE-100` (UseCase層、101, 102に依存)
4. `feature/SLEEVE-104` → `topic/SLEEVE-100` (GraphQL層、103に依存)
5. 全てのサブタスク完了後: `topic/SLEEVE-100` → `main`

詳細な例は「PR粒度の詳細定義」セクションを参照してください。

### マージ戦略

#### Feature Branch → Topic Branch (または依存先Feature Branch)
- **マージ方法**: Squash and merge
- **理由**: コミット履歴を整理し、Topic Branch上で見やすくする

#### Feature Branch → Feature Branch (依存関係がある場合)
- **マージ方法**: Merge commit
- **理由**: 依存先の変更を保持したまま、自分の変更を追加する

#### Topic Branch → main
- **マージ方法**: Merge commit
- **理由**: Topic Branchの履歴を残し、機能単位での変更を追跡可能にする

#### hotfix → main
- **マージ方法**: Squash and merge
- **理由**: 緊急修正は単一のコミットとして記録

---

## PR粒度の詳細定義

### 基本方針

```
Topic Branch = 1つの機能単位（例: ユーザー登録機能、ユーザー取得機能）
各PR（Subtask） = 1つのレイヤーの1操作実装 + 単体テスト
                 単体テストで動作確認可能
```

**重要原則**:
- 各PRは単体テストで動作確認できる状態で作成する
- レイヤー × CRUD操作単位で分割することで、レビューを高速化する
- 機能として意味のある単位で分割する
- 行数での強制的な分割は行わない（参考程度）

### 1 PRの粒度

**定義**: 1 PR = 1つのレイヤーの1操作実装 + 単体テスト

**動作確認**: そのレイヤーの単体テストが通れば動作確認完了

### 完結性の基準（必須条件）

各PRは以下を**必ず**満たす必要があります：

- [ ] 単独でビルドが通る
- [ ] **そのレイヤーの単体テストが全て通過する**
- [ ] 他の機能に影響を与えない（または依存が明確に定義されている）
- [ ] 機能として説明できる（「UserRepositoryのCreate操作を実装」など）
- [ ] 単独でロールバック可能

### サイズ目安（参考のみ）

以下は参考値です。**無理に守る必要はありません**。機能の完結性を優先してください。

| レイヤー | 変更行数（参考） | ファイル数（参考） | レビュー時間（参考） |
|---------|----------------|------------------|-------------------|
| Domain層 | 50〜150行 | 1〜3個 | 10〜15分 |
| Repository層（1操作） | 50〜100行 | 1〜2個 | 5〜10分 |
| Repository層（全操作） | 150〜300行 | 2〜4個 | 15〜25分 |
| UseCase層（1操作） | 80〜150行 | 1〜2個 | 10〜15分 |
| GraphQL層（1操作） | 60〜120行 | 2〜3個 | 10〜15分 |

**注意事項**:
- 機能の完結性を優先してください
- レビュー時間が30分を大きく超える場合は分割を検討してください

### 推奨する分割パターン（レイヤー × CRUD操作）

基本的に以下のパターンで分割します：

```
Topic: [機能名]（例: ユーザー登録機能）
│
├─ PR1: [機能名]モデル定義（Domain層）
│   - エンティティ定義、バリデーション、エラー定義
│   - 単体テスト
│   ✅ 動作確認: `go test ./domain/models/code_models`
│
├─ PR2: [機能名]Repository - [操作]操作実装（Repository層）
│   - インターフェース定義、実装
│   - 単体テスト
│   ✅ 動作確認: `go test ./repository`
│
├─ PR3: [操作][機能名]UseCase実装（UseCase層）
│   - UseCase実装、ビジネスロジック
│   - 単体テスト（Repositoryモック）
│   ✅ 動作確認: `go test ./usecase/[feature]`
│
└─ PR4: [機能名][操作]GraphQL API実装（GraphQL層）
    - スキーマ定義、リゾルバー実装
    - 単体テスト（UseCaseモック）
    ✅ 動作確認: `go test ./graph`
```

### 各レイヤーのPR内容

#### Domain層のPR

```
含むもの:
- エンティティ定義（app/domain/models/code_models/[entity].go）
- 値オブジェクト定義
- バリデーションロジック
- ドメインエラー定義（app/domain/errors/[entity]_errors.go）
- 単体テスト（*_test.go）

動作確認:
- `go test ./domain/models/code_models`

サイズ感（参考）:
- 50〜150行、1〜3ファイル
- レビュー時間: 10〜15分
```

#### Repository層のPR（1操作）

```
含むもの:
- DBモデル定義（app/domain/models/db_models/[entity].go）
- Repositoryインターフェース定義（その操作部分）
- 実装（app/repository/[entity]_repository.go）
- エラーハンドリング
- 単体テスト（app/repository/[entity]_repository_test.go）

動作確認:
- `go test ./repository`

サイズ感（参考）:
- 50〜150行、1〜3ファイル
- レビュー時間: 5〜15分
```

#### UseCase層のPR（1操作）

```
含むもの:
- UseCase実装（app/usecase/[feature]/[operation]_[entity]_usecase.go）
- ビジネスロジック
- Repository呼び出し
- エラーハンドリング
- メッセージ定義（app/messages/[feature]/errors.go, docs/messages/[feature]/errors.md）
- 単体テスト（*_test.go、Repositoryモック使用）

動作確認:
- `go test ./usecase/[feature]`

サイズ感（参考）:
- 80〜180行、2〜4ファイル
- レビュー時間: 10〜18分
```

#### GraphQL層のPR（1操作）

```
含むもの:
- スキーマ定義（app/graph/schema.graphqls - MutationまたはQuery追加）
- Input型定義、Payload型定義
- リゾルバー実装（app/graph/schema.resolvers.go）
- 入力バリデーション
- エラーレスポンス整形
- 単体テスト（app/graph/schema.resolvers_test.go、UseCaseモック使用）

動作確認:
- `go test ./graph`

サイズ感（参考）:
- 60〜130行、2〜3ファイル
- レビュー時間: 10〜13分
```

### 実例: ユーザー登録機能

```
Topic Branch: topic/SLEEVE-100_user_registration
完結した機能: ユーザーがメールアドレスとパスワードで新規登録できる

Subtask分割:

├─ PR1: SLEEVE-101: Userモデル定義
│   📁 変更ファイル:
│   - app/domain/models/code_models/user.go (新規)
│   - app/domain/models/code_models/user_test.go (新規)
│   - app/domain/errors/user_errors.go (新規)
│
│   📝 実装内容:
│   - User構造体定義（ID, Email, Password, CreatedAt等）
│   - EmailアドレスバリデーションValueObject
│   - パスワードバリデーションロジック
│   - ドメインエラー定義（ErrInvalidEmail, ErrWeakPassword等）
│
│   📊 規模感: 約120行、3ファイル
│   ⏱️ レビュー時間: 12分
│   ✅ 動作確認: `go test ./domain/models/code_models`
│   🔗 依存: なし
│
├─ PR2: SLEEVE-102: UserRepository - Create操作実装
│   📁 変更ファイル:
│   - app/domain/models/db_models/user.go (新規)
│   - app/repository/user_repository.go (新規)
│   - app/repository/user_repository_test.go (新規)
│
│   📝 実装内容:
│   - DBモデル定義（ent/gorm用）
│   - UserRepositoryインターフェース定義（CreateUserメソッド）
│   - CreateUser実装（DBへの挿入処理）
│   - エラーハンドリング
│   - 単体テスト（モックDB使用）
│
│   📊 規模感: 約150行、3ファイル
│   ⏱️ レビュー時間: 15分
│   ✅ 動作確認: `go test ./repository`
│   🔗 依存: PR1（Userモデルを使用）
│
├─ PR3: SLEEVE-103: CreateUserUseCase実装
│   📁 変更ファイル:
│   - app/usecase/user/create_user_usecase.go (新規)
│   - app/usecase/user/create_user_usecase_test.go (新規)
│   - app/messages/user/errors.go (新規)
│   - docs/messages/user/errors.md (新規)
│
│   📝 実装内容:
│   - CreateUserUseCase構造体定義
│   - メールアドレス重複チェックロジック
│   - パスワードハッシュ化
│   - Repository呼び出し
│   - エラーメッセージ定義
│   - 単体テスト（Repositoryモック使用）
│
│   📊 規模感: 約180行、4ファイル
│   ⏱️ レビュー時間: 18分
│   ✅ 動作確認: `go test ./usecase/user`
│   🔗 依存: PR1, PR2
│
└─ PR4: SLEEVE-104: ユーザー登録GraphQL API実装
    📁 変更ファイル:
    - app/graph/schema.graphqls (変更 - createUserミューテーション追加)
    - app/graph/schema.resolvers.go (変更 - createUserリゾルバー追加)
    - app/graph/schema.resolvers_test.go (変更)

    📝 実装内容:
    - createUserミューテーション定義（スキーマ）
    - CreateUserInput型定義
    - CreateUserPayload型定義
    - createUserリゾルバー実装（UseCaseを呼び出し）
    - 入力バリデーション
    - エラーレスポンス整形
    - 単体テスト（UseCaseモック使用）

    📊 規模感: 約130行、3ファイル
    ⏱️ レビュー時間: 13分
    ✅ 動作確認: `go test ./graph`
    🔗 依存: PR3

Topic全体:
📊 合計: 約580行、13ファイル、4 PR
⏱️ 合計レビュー時間: 58分（1 PRあたり約14分）
✅ 各PRで単体テストによる動作確認可能
✅ Topic全体で「ユーザー登録機能」が完結
```

### 実例: ユーザー取得機能

```
Topic Branch: topic/SLEEVE-200_user_retrieval
完結した機能: ユーザーIDまたはメールアドレスでユーザー情報を取得できる

Subtask分割:

├─ PR1: SLEEVE-201: UserRepository - Read操作実装
│   📁 変更ファイル:
│   - app/repository/user_repository.go (変更 - Read関連メソッド追加)
│   - app/repository/user_repository_test.go (変更)
│
│   📝 実装内容:
│   - UserRepositoryインターフェースにFindByID, FindByEmailメソッド追加
│   - FindByID実装（IDでユーザー検索）
│   - FindByEmail実装（メールアドレスでユーザー検索）
│   - エラーハンドリング（ユーザーが見つからない場合等）
│   - 単体テスト
│
│   📊 規模感: 約100行、2ファイル
│   ⏱️ レビュー時間: 10分
│   ✅ 動作確認: `go test ./repository`
│   🔗 依存: SLEEVE-101（Userモデル）, SLEEVE-102（UserRepository）
│
├─ PR2: SLEEVE-202: GetUserUseCase実装
│   📁 変更ファイル:
│   - app/usecase/user/get_user_usecase.go (新規)
│   - app/usecase/user/get_user_usecase_test.go (新規)
│   - app/messages/user/errors.go (変更)
│   - docs/messages/user/errors.md (変更)
│
│   📝 実装内容:
│   - GetUserUseCase構造体定義
│   - IDまたはEmailでユーザー取得ロジック
│   - ユーザーが見つからない場合のエラーハンドリング
│   - エラーメッセージ追加（ErrUserNotFound）
│   - 単体テスト（Repositoryモック使用）
│
│   📊 規模感: 約120行、4ファイル
│   ⏱️ レビュー時間: 12分
│   ✅ 動作確認: `go test ./usecase/user`
│   🔗 依存: PR1
│
└─ PR3: SLEEVE-203: ユーザー取得GraphQL API実装
    📁 変更ファイル:
    - app/graph/schema.graphqls (変更 - getUserクエリ追加)
    - app/graph/schema.resolvers.go (変更 - getUserリゾルバー追加)
    - app/graph/schema.resolvers_test.go (変更)

    📝 実装内容:
    - getUserクエリ定義（スキーマ）
    - GetUserInput型定義（idまたはemail）
    - getUserリゾルバー実装（UseCaseを呼び出し）
    - エラーレスポンス整形
    - 単体テスト（UseCaseモック使用）

    📊 規模感: 約100行、3ファイル
    ⏱️ レビュー時間: 10分
    ✅ 動作確認: `go test ./graph`
    🔗 依存: PR2

Topic全体:
📊 合計: 約320行、9ファイル、3 PR
⏱️ 合計レビュー時間: 32分（1 PRあたり約11分）
✅ 各PRで単体テストによる動作確認可能
✅ Topic全体で「ユーザー取得機能」が完結
```

### 実例: ユーザー更新機能

```
Topic Branch: topic/SLEEVE-300_user_update
完結した機能: ユーザー情報（メールアドレス、パスワード）を更新できる

Subtask分割:

├─ PR1: SLEEVE-301: UserRepository - Update操作実装
│   📁 変更ファイル:
│   - app/repository/user_repository.go (変更 - Update関連メソッド追加)
│   - app/repository/user_repository_test.go (変更)
│
│   📝 実装内容:
│   - UserRepositoryインターフェースにUpdateUserメソッド追加
│   - UpdateUser実装（ユーザー情報更新）
│   - エラーハンドリング
│   - 単体テスト
│
│   📊 規模感: 約90行、2ファイル
│   ⏱️ レビュー時間: 9分
│   ✅ 動作確認: `go test ./repository`
│   🔗 依存: SLEEVE-101（Userモデル）, SLEEVE-102（UserRepository）
│
├─ PR2: SLEEVE-302: UpdateUserUseCase実装
│   📁 変更ファイル:
│   - app/usecase/user/update_user_usecase.go (新規)
│   - app/usecase/user/update_user_usecase_test.go (新規)
│   - app/messages/user/errors.go (変更)
│   - docs/messages/user/errors.md (変更)
│
│   📝 実装内容:
│   - UpdateUserUseCase構造体定義
│   - ユーザー存在チェック
│   - メールアドレス重複チェック（変更時）
│   - パスワードハッシュ化（変更時）
│   - Repository呼び出し
│   - エラーメッセージ追加
│   - 単体テスト（Repositoryモック使用）
│
│   📊 規模感: 約160行、4ファイル
│   ⏱️ レビュー時間: 16分
│   ✅ 動作確認: `go test ./usecase/user`
│   🔗 依存: PR1, SLEEVE-201（FindByID）
│
└─ PR3: SLEEVE-303: ユーザー更新GraphQL API実装
    📁 変更ファイル:
    - app/graph/schema.graphqls (変更 - updateUserミューテーション追加)
    - app/graph/schema.resolvers.go (変更 - updateUserリゾルバー追加)
    - app/graph/schema.resolvers_test.go (変更)

    📝 実装内容:
    - updateUserミューテーション定義（スキーマ）
    - UpdateUserInput型定義
    - UpdateUserPayload型定義
    - updateUserリゾルバー実装（UseCaseを呼び出し）
    - 入力バリデーション
    - エラーレスポンス整形
    - 単体テスト（UseCaseモック使用）

    📊 規模感: 約120行、3ファイル
    ⏱️ レビュー時間: 12分
    ✅ 動作確認: `go test ./graph`
    🔗 依存: PR2

Topic全体:
📊 合計: 約370行、9ファイル、3 PR
⏱️ 合計レビュー時間: 37分（1 PRあたり約12分）
✅ 各PRで単体テストによる動作確認可能
✅ Topic全体で「ユーザー更新機能」が完結
```

### 実例: ユーザー削除機能

```
Topic Branch: topic/SLEEVE-400_user_delete
完結した機能: ユーザーを論理削除できる

Subtask分割:

├─ PR1: SLEEVE-401: UserRepository - Delete操作実装
│   📁 変更ファイル:
│   - app/repository/user_repository.go (変更 - Delete関連メソッド追加)
│   - app/repository/user_repository_test.go (変更)
│
│   📝 実装内容:
│   - UserRepositoryインターフェースにDeleteUserメソッド追加
│   - DeleteUser実装（論理削除: deleted_atを更新）
│   - エラーハンドリング
│   - 単体テスト
│
│   📊 規模感: 約70行、2ファイル
│   ⏱️ レビュー時間: 7分
│   ✅ 動作確認: `go test ./repository`
│   🔗 依存: SLEEVE-101（Userモデル）, SLEEVE-102（UserRepository）
│
├─ PR2: SLEEVE-402: DeleteUserUseCase実装
│   📁 変更ファイル:
│   - app/usecase/user/delete_user_usecase.go (新規)
│   - app/usecase/user/delete_user_usecase_test.go (新規)
│   - app/messages/user/errors.go (変更)
│   - docs/messages/user/errors.md (変更)
│
│   📝 実装内容:
│   - DeleteUserUseCase構造体定義
│   - ユーザー存在チェック
│   - Repository呼び出し（論理削除）
│   - エラーメッセージ追加
│   - 単体テスト（Repositoryモック使用）
│
│   📊 規模感: 約100行、4ファイル
│   ⏱️ レビュー時間: 10分
│   ✅ 動作確認: `go test ./usecase/user`
│   🔗 依存: PR1, SLEEVE-201（FindByID）
│
└─ PR3: SLEEVE-403: ユーザー削除GraphQL API実装
    📁 変更ファイル:
    - app/graph/schema.graphqls (変更 - deleteUserミューテーション追加)
    - app/graph/schema.resolvers.go (変更 - deleteUserリゾルバー追加)
    - app/graph/schema.resolvers_test.go (変更)

    📝 実装内容:
    - deleteUserミューテーション定義（スキーマ）
    - DeleteUserInput型定義
    - DeleteUserPayload型定義
    - deleteUserリゾルバー実装（UseCaseを呼び出し）
    - エラーレスポンス整形
    - 単体テスト（UseCaseモック使用）

    📊 規模感: 約90行、3ファイル
    ⏱️ レビュー時間: 9分
    ✅ 動作確認: `go test ./graph`
    🔗 依存: PR2

Topic全体:
📊 合計: 約260行、9ファイル、3 PR
⏱️ 合計レビュー時間: 26分（1 PRあたり約9分）
✅ 各PRで単体テストによる動作確認可能
✅ Topic全体で「ユーザー削除機能」が完結
```

### 分割の判断フローチャート

```
機能実装を開始
    ↓
Topic Branchを作成（1つの機能単位: ユーザー登録、ユーザー取得等）
    ↓
レイヤー × CRUD操作単位でSubtaskを分割
    ├─ Domain層 → PR1（最初の機能のみ）
    ├─ Repository層 - [操作]操作 → PR
    ├─ UseCase層 - [操作]UseCase → PR
    └─ GraphQL層 - [操作]API → PR
    ↓
各PRが単体テストで動作確認できるか？
    ├─ Yes → そのまま実装
    └─ No → 単体テストを含むように調整
    ↓
レビュー時間が30分を大きく超えるか？
    ├─ No → そのまま実装
    └─ Yes → さらに細かく分割を検討
```

### 依存関係の管理

レイヤー × CRUD操作分割の場合、典型的な依存関係：

```
PR1: Domain層（最初の機能のみ）
    ↓
PR2: Repository層 - Create操作（PR1に依存）
    ↓
PR3: UseCase層 - CreateUseCase（PR1, PR2に依存）
    ↓
PR4: GraphQL層 - CreateAPI（PR3に依存）
```

**別の機能（例: ユーザー取得）の場合**:

```
PR1: Repository層 - Read操作（SLEEVE-101 Userモデル、SLEEVE-102 UserRepositoryに依存）
    ↓
PR2: UseCase層 - GetUserUseCase（PR1に依存）
    ↓
PR3: GraphQL層 - GetUserAPI（PR2に依存）
```

**ブランチ戦略**:
- 前のPRがTopicにマージ済み → 次のPRはTopicから分岐
- 前のPRがまだマージされていない → 次のPRは前のPRから分岐

詳細は「依存関係があるSubtaskの詳細マージ手順」を参照。

---

### 依存関係があるSubtaskの詳細マージ手順

**ケース: topic → task1 → task2 のようにブランチが切られている場合**

依存関係がある場合のマージは、**変更差分を見やすくする**ために以下の順序で行います。

#### 前提条件
```
- topic/SLEEVE-100_user_authentication (親ブランチ)
- feature/SLEEVE-102_user_model (task1)
- feature/SLEEVE-103_login_api (task2: 102に依存)
```

#### ステップ1: task1を実装・PR作成

```bash
# task1をtopicから分岐して実装
git checkout topic/SLEEVE-100_user_authentication
git checkout -b feature/SLEEVE-102_user_model

# 実装・テスト・コミット...

# PR作成（base: topic）
gh pr create --base topic/SLEEVE-100_user_authentication --title "SLEEVE-102: Userモデルの作成"
```

**レビュー・マージ:**
- オーナーがレビュー・承認
- Squash and mergeでtopicにマージ
- task1のFeature Branchを削除

#### ステップ2: topicブランチを最新化

```bash
git checkout topic/SLEEVE-100_user_authentication
git pull origin topic/SLEEVE-100_user_authentication
```

#### ステップ3: task2を実装・PR作成

**重要な判断基準:**
- **task1が既にtopicにマージ済みの場合**: task2をtopicから分岐
- **task1がまだマージされていない場合**: task2をtask1から分岐

**パターンA: task1が既にtopicにマージ済みの場合**

```bash
# task2をtopicから分岐して実装（task1の変更を含むため）
git checkout topic/SLEEVE-100_user_authentication
git pull origin topic/SLEEVE-100_user_authentication
git checkout -b feature/SLEEVE-103_login_api

# 実装・テスト・コミット...
# （この時点でtask2にはtask1の変更も含まれている）

# PR作成（base: topic）
gh pr create --base topic/SLEEVE-100_user_authentication --title "SLEEVE-103: ログインAPIの実装"
```

**レビュー・マージ:**
- オーナーがレビュー・承認
- Squash and mergeでtopicにマージ（**task2の変更差分だけがマージされる**）
- task2のFeature Branchを削除

**パターンB: task1がまだマージされていない場合**

```bash
# task2をtask1から分岐して実装
git checkout feature/SLEEVE-102_user_model
git pull origin feature/SLEEVE-102_user_model
git checkout -b feature/SLEEVE-103_login_api

# 実装・テスト・コミット...

# PR作成（base: task1）
gh pr create --base feature/SLEEVE-102_user_model --title "SLEEVE-103: ログインAPIの実装"

# task1がtopicにマージされた後、task2もtopicにPR作成
gh pr create --base topic/SLEEVE-100_user_authentication --title "SLEEVE-103: ログインAPIの実装"
```

**レビュー・マージ:**
- オーナーがレビュー・承認
- task1がtopicにマージ後、task2もtopicにマージ
- task2のFeature Branchを削除

#### 理由
- task1がtopicにマージ済みであれば、task2もtopicから分岐することで、task2のPRではtask2の変更だけが見やすくなる
- task1がまだマージされていない場合、task2をtask1から分岐することで、依存関係を明確にできる
- 最終的にはtopicに両方マージされるが、レビュー時に変更差分が明確になる
- Squash and mergeによって、task1の変更は重複せず、task2の変更差分だけがコミットされる

### ブランチの最新化ルール

**重要:** 各Subtaskのマージ後、必ずtopicまたはmainブランチを最新化してください。

```bash
# 複数Subtaskの場合: Topic Branchを最新化
git checkout topic/[親ID]_[機能名]
git pull origin topic/[親ID]_[機能名]

# 単一タスクの場合: Main Branchを最新化
git checkout main
git pull origin main
```

**最新化のタイミング:**
- 他のSubtaskがマージされた時
- 新しいSubtaskの実装を開始する前
- PR作成前

### Jira連携

#### Epic
- プロジェクト全体の大きな目標
- 例: 「ユーザー管理機能の実装」

#### Story/Task
- Topic Branchに対応
- 1つの完結した機能（例: ユーザー登録機能、ユーザー取得機能）
- 通常3〜4個のSubtaskを持つ（Domain、Repository、UseCase、GraphQL）
- 例:
  - `SLEEVE-100: ユーザー登録機能の実装`
  - `SLEEVE-200: ユーザー取得機能の実装`
  - `SLEEVE-300: ユーザー更新機能の実装`
  - `SLEEVE-400: ユーザー削除機能の実装`

#### Subtask
- Feature Branchに対応
- 1つのレイヤーの実装（または1つのレイヤーの1操作）
- 例（ユーザー登録機能の場合）:
  - `SLEEVE-101: Userモデル定義`
  - `SLEEVE-102: UserRepository - Create操作実装`
  - `SLEEVE-103: CreateUserUseCase実装`
  - `SLEEVE-104: ユーザー登録GraphQL API実装`

### プルリクエスト運用

#### PR作成ルール

- **レビュー必須**: オーナーのレビューと承認が必要
- **CI通過必須**: 全てのテストが通過している必要がある
- **Linterチェック**: golangci-lintエラーがないこと
- **依存関係の明記**: PR descriptionに依存関係を記載

#### PR descriptionテンプレート

```markdown
## 概要
[このPRで実装する内容を簡潔に記述]

## 変更内容
- [変更内容1]
- [変更内容2]

## 依存関係
- Depends on: #123, #124 (このPRがマージされる前に必要なPR)
- Blocks: #126 (このPRに依存しているPR)

## テスト
- [ ] 単体テストが通過
- [ ] 結合テストが通過
- [ ] 手動テストが完了

## 関連チケット
- Jira: SLEEVE-XXX

## チェックリスト
- [ ] コーディング規約に準拠
- [ ] Linterエラーがない
- [ ] ドキュメントが更新されている
- [ ] DB変更を行った場合は `docs/db_schema.md` を更新している
```

#### レビュー待機フロー（現状の運用）

**⚠️ 注意：以下は現状の運用です。メンバーが増えた際は変更される可能性があります。**

##### レビュアー
- **オーナーがレビューを行います**

##### レビュー待ちの動き
- ❌ **レビュー待ちの間、他のSubtaskに進まない**
- ✅ **PR作成後は、オーナーのレビュー・承認・マージを完全に待つ**
- ✅ **レビュー中に他の作業を進めない**

**理由:**
- 現状は少人数開発のため、並行作業によるコンフリクトを避ける
- レビューでの修正指示を即座に反映できるようにする
- ブランチ管理をシンプルに保つ

**将来の変更予定:**
- メンバーが増えた際は、依存関係のないSubtaskを並行開発できるように変更する可能性があります

##### レビュー中の対応

1. **PR作成後**
   - オーナーに通知
   - レビュー待ちの状態で待機

2. **修正依頼があった場合**
   - 指摘された箇所を修正
   - 追加のコミットを作成
   - プッシュして再レビュー依頼

3. **承認された場合**
   - オーナーがマージを実行
   - マージ完了を確認
   - Feature Branchを削除
   - Topic/Main Branchを最新化

### CI/CD

#### CI実行タイミング
- **全てのブランチ**: PR作成時およびpush時にCIを実行
- **main**: マージ時に必ずCIを実行
- **Topic Branch**: PR作成時およびpush時にCIを実行

#### CI内容
- 単体テスト
- 結合テスト
- Linterチェック
- ビルド確認

## 使用可能コマンド

### ブランチ操作

#### Topic Branchの作成
```bash
# mainから新しいTopic Branchを作成
git checkout main
git pull origin main
git checkout -b topic/[JiraのID]_[機能名]
git push -u origin topic/[JiraのID]_[機能名]

# 例（ユーザー登録機能）
git checkout -b topic/SLEEVE-100_user_registration
```

#### Feature Branchの作成（依存関係なし）
```bash
# Topic Branchから新しいFeature Branchを作成
git checkout topic/[JiraのID]_[機能名]
git pull origin topic/[JiraのID]_[機能名]
git checkout -b feature/[JiraのSubtaskID]_[機能名]
git push -u origin feature/[JiraのSubtaskID]_[機能名]

# 例（Domain層のPR作成）
git checkout topic/SLEEVE-100_user_registration
git checkout -b feature/SLEEVE-101_user_model
```

#### Feature Branchの作成（依存関係あり）
```bash
# 依存先のFeature Branchから新しいFeature Branchを作成
git checkout feature/[依存先のJiraID]_[機能名]
git pull origin feature/[依存先のJiraID]_[機能名]
git checkout -b feature/[JiraのSubtaskID]_[機能名]
git push -u origin feature/[JiraのSubtaskID]_[機能名]

# 例: Repository層のPRがDomain層のPRに依存する場合
git checkout feature/SLEEVE-101_user_model
git checkout -b feature/SLEEVE-102_user_repository_create
```

#### ブランチの切り替え
```bash
git checkout [ブランチ名]
```

#### ブランチ一覧の表示
```bash
# ローカルブランチ一覧
git branch

# リモートを含む全ブランチ一覧
git branch -a
```

#### ブランチの削除
```bash
# ローカルブランチを削除
git branch -d [ブランチ名]

# リモートブランチを削除（マージ後）
git push origin --delete [ブランチ名]
```

### コミット操作

```bash
# 変更をステージング
git add [ファイルパス]
git add .  # 全ての変更をステージング

# コミット実行
git commit -m "[コミットメッセージ]"

# 直前のコミットを修正
git commit --amend -m "[修正後のメッセージ]"
```

### プッシュ・プル操作

```bash
# リモートにプッシュ
git push origin [ブランチ名]

# 初回プッシュ時
git push -u origin [ブランチ名]

# リモートから最新を取得
git pull origin [ブランチ名]
```

## 使用禁止コマンド

### 環境破壊の危険があるコマンド（絶対禁止）

```bash
# 絶対に実行してはいけないコマンド
git reset --hard HEAD~[数]  # コミット履歴を強制削除
git push --force            # 強制プッシュ（履歴上書き）
git push --force-with-lease # 強制プッシュ（条件付き）
git rebase --onto [base] [upstream] [branch]  # 複雑なリベース
git filter-branch          # 履歴書き換え
git gc --aggressive --prune=now  # 強制ガベージコレクション
git clean -fd              # 追跡されていないファイルを強制削除
git checkout -f            # 強制チェックアウト
git merge --no-ff -X theirs  # 強制マージ
```

### 注意が必要なコマンド（使用前に確認必須）

```bash
# 使用前に必ず確認が必要
git rebase -i HEAD~[数]    # インタラクティブリベース
git cherry-pick [commit]   # コミットの移植
git revert [commit]        # コミットの取り消し
git stash drop             # スタッシュの削除
git branch -D [ブランチ名]  # 強制ブランチ削除
git tag -d [タグ名]        # タグの削除
```

### 安全な代替手段

```bash
# 危険なコマンドの代わりに使用する安全なコマンド

# git reset --hard の代わり
git checkout [ブランチ名]  # ブランチを切り替え
git restore .              # 変更を元に戻す

# git push --force の代わり
git push origin [ブランチ名]  # 通常のプッシュ
git pull origin [ブランチ名]  # 最新を取得してからプッシュ

# git clean -fd の代わり
git status                 # 状態確認
git add [ファイル]          # 必要なファイルを追加
git restore [ファイル]      # 特定ファイルを元に戻す
```

## コミットメッセージ規約

***コミットメッセージはPrefix以外日本語で行うこと!!!***

### 基本形式

```
[prefix]: [簡潔な説明]

- 詳細な変更内容1
- 詳細な変更内容2
- 関連するJiraチケットID: [JiraのID]
```

### Prefix 一覧

#### 機能関連

- `feat:` - 新機能の追加
- `feature:` - 大きな機能の追加（feat:より大規模）
- `enhance:` - 既存機能の改善・拡張

#### 修正関連

- `fix:` - バグ修正
- `hotfix:` - 緊急バグ修正
- `patch:` - 小さな修正

#### 削除関連

- `delete:` - ファイル・機能の削除
- `remove:` - コード・設定の削除
- `cleanup:` - 不要なコードの整理

#### リファクタリング関連

- `refactor:` - リファクタリング（機能変更なし）
- `reorganize:` - ファイル・ディレクトリ構造の整理
- `rename:` - ファイル・関数・変数の名前変更

#### テスト関連

- `test:` - テストの追加・修正
- `test:unit` - 単体テストの追加・修正
- `test:integration` - 統合テストの追加・修正

#### ドキュメント関連

- `docs:` - ドキュメントの更新
- `readme:` - README の更新
- `comment:` - コメントの追加・修正

#### 設定・環境関連

- `config:` - 設定ファイルの変更
- `env:` - 環境変数の変更
- `deps:` - 依存関係の追加・更新・削除
- `build:` - ビルド設定の変更
- `ci:` - CI/CD 設定の変更

#### パフォーマンス関連

- `perf:` - パフォーマンスの改善
- `optimize:` - 最適化

#### セキュリティ関連

- `security:` - セキュリティ関連の修正
- `auth:` - 認証・認可関連の修正

#### データベース関連

- `db:` - データベース関連の変更
- `migration:` - データベースマイグレーション
- `schema:` - スキーマの変更

#### その他

- `style:` - コードスタイルの修正（機能に影響なし）
- `format:` - コードフォーマットの適用
- `lint:` - Linter エラーの修正
- `chore:` - その他の雑務

### コミットメッセージ例

#### 良い例

```
feat: ユーザー認証機能を追加

- JWT認証の実装
- ログイン・ログアウトAPIの追加
- パスワードハッシュ化の実装
- 関連するJiraチケットID: SLEEVE-123
```

```
fix: ユーザー登録時のバリデーションエラーを修正

- メールアドレス形式チェックの修正
- パスワード強度チェックの改善
- エラーメッセージの日本語化
- 関連するJiraチケットID: SLEEVE-456
```

```
refactor: Repository層の構造を整理

- インターフェースの分離
- エラーハンドリングの統一
- テストカバレッジの向上
- 関連するJiraチケットID: SLEEVE-789
```

#### 悪い例

```
修正
```

```
update
```

```
いろいろ修正した
```

### コミット頻度のガイドライン

#### 推奨されるコミット頻度

- **機能単位**: 1 つの機能が完成したらコミット
- **テスト単位**: テストが通るようになったらコミット
- **修正単位**: 1 つのバグが修正されたらコミット
- **リファクタリング単位**: リファクタリングが完了したらコミット

#### 避けるべきコミット

- 動作しないコードのコミット
- 複数の機能を一度にコミット
- テストが通らない状態でのコミット
- 意味のないコミット（「作業中」など）


## 実践例：ユーザー登録機能の開発フロー

この例では、細粒度のブランチ戦略（レイヤー × CRUD操作単位）を使用して「ユーザー登録機能」を開発するフローを説明します。

### 1. Jiraでのタスク管理

**Epic**: `EPIC-50: ユーザー管理機能の実装`

**Story/Task**: `SLEEVE-100: ユーザー登録機能の実装`

**Subtask**:
- `SLEEVE-101: Userモデル定義`（Domain層）
- `SLEEVE-102: UserRepository - Create操作実装`（Repository層）
- `SLEEVE-103: CreateUserUseCase実装`（UseCase層、101と102に依存）
- `SLEEVE-104: ユーザー登録GraphQL API実装`（GraphQL層、103に依存）

### 2. Topic Branchの作成

```bash
git checkout main
git pull origin main
git checkout -b topic/SLEEVE-100_user_registration
git push -u origin topic/SLEEVE-100_user_registration
```

### 3. 各Subtaskの実装

#### Subtask 1: Userモデル定義（依存なし）

```bash
# Feature Branchを作成
git checkout topic/SLEEVE-100_user_registration
git checkout -b feature/SLEEVE-101_user_model

# 実装内容:
# - app/domain/models/code_models/user.go
# - app/domain/models/code_models/user_test.go
# - app/domain/errors/user_errors.go

# テスト実行
cd app
go test ./domain/models/code_models

# Linter実行
golangci-lint run --fix

# コミット & プッシュ
git add .
git commit -m "feat: Userモデル定義を追加

- User構造体定義
- EmailバリデーションValueObject
- パスワードバリデーションロジック
- ドメインエラー定義
- 関連するJiraチケットID: SLEEVE-101"

git push origin feature/SLEEVE-101_user_model

# PR作成
gh pr create --base topic/SLEEVE-100_user_registration \
  --title "SLEEVE-101: Userモデル定義" \
  --body-file docs/tasks/SLEEVE-100_user_registration/pr_description_SLEEVE-101.md

# レビュー・承認・マージ待ち（オーナーが実施）
# マージ方法: Squash and merge
```

#### Subtask 2: UserRepository - Create操作実装（101に依存）

```bash
# Topic Branchを最新化（101がマージ済み）
git checkout topic/SLEEVE-100_user_registration
git pull origin topic/SLEEVE-100_user_registration

# Feature Branchを作成
git checkout -b feature/SLEEVE-102_user_repository_create

# 実装内容:
# - app/domain/models/db_models/user.go
# - app/repository/user_repository.go
# - app/repository/user_repository_test.go

# テスト実行
cd app
go test ./repository

# Linter実行
golangci-lint run --fix

# コミット & プッシュ
git add .
git commit -m "feat: UserRepository - Create操作を実装

- DBモデル定義
- UserRepositoryインターフェース定義（CreateUserメソッド）
- CreateUser実装
- エラーハンドリング
- 関連するJiraチケットID: SLEEVE-102"

git push origin feature/SLEEVE-102_user_repository_create

# PR作成
gh pr create --base topic/SLEEVE-100_user_registration \
  --title "SLEEVE-102: UserRepository - Create操作実装" \
  --body-file docs/tasks/SLEEVE-100_user_registration/pr_description_SLEEVE-102.md

# レビュー・承認・マージ待ち
# マージ方法: Squash and merge
```

#### Subtask 3: CreateUserUseCase実装（101と102に依存）

```bash
# Topic Branchを最新化（102がマージ済み）
git checkout topic/SLEEVE-100_user_registration
git pull origin topic/SLEEVE-100_user_registration

# Feature Branchを作成
git checkout -b feature/SLEEVE-103_create_user_usecase

# 実装内容:
# - app/usecase/user/create_user_usecase.go
# - app/usecase/user/create_user_usecase_test.go
# - app/messages/user/errors.go
# - docs/messages/user/errors.md

# テスト実行
cd app
go test ./usecase/user

# Linter実行
golangci-lint run --fix

# コミット & プッシュ
git add .
git commit -m "feat: CreateUserUseCaseを実装

- CreateUserUseCase構造体定義
- メールアドレス重複チェックロジック
- パスワードハッシュ化
- エラーメッセージ定義
- 関連するJiraチケットID: SLEEVE-103"

git push origin feature/SLEEVE-103_create_user_usecase

# PR作成
gh pr create --base topic/SLEEVE-100_user_registration \
  --title "SLEEVE-103: CreateUserUseCase実装" \
  --body-file docs/tasks/SLEEVE-100_user_registration/pr_description_SLEEVE-103.md

# レビュー・承認・マージ待ち
# マージ方法: Squash and merge
```

#### Subtask 4: ユーザー登録GraphQL API実装（103に依存）

```bash
# Topic Branchを最新化（103がマージ済み）
git checkout topic/SLEEVE-100_user_registration
git pull origin topic/SLEEVE-100_user_registration

# Feature Branchを作成
git checkout -b feature/SLEEVE-104_user_registration_api

# 実装内容:
# - app/graph/schema.graphqls (createUserミューテーション追加)
# - app/graph/schema.resolvers.go (createUserリゾルバー追加)
# - app/graph/schema.resolvers_test.go

# テスト実行
cd app
go test ./graph

# Linter実行
golangci-lint run --fix

# コミット & プッシュ
git add .
git commit -m "feat: ユーザー登録GraphQL APIを実装

- createUserミューテーション定義
- CreateUserInput型、CreateUserPayload型定義
- createUserリゾルバー実装
- 入力バリデーション
- 関連するJiraチケットID: SLEEVE-104"

git push origin feature/SLEEVE-104_user_registration_api

# PR作成
gh pr create --base topic/SLEEVE-100_user_registration \
  --title "SLEEVE-104: ユーザー登録GraphQL API実装" \
  --body-file docs/tasks/SLEEVE-100_user_registration/pr_description_SLEEVE-104.md

# レビュー・承認・マージ待ち
# マージ方法: Squash and merge
```

### 4. Topic Branchをmainにマージ

全てのSubtaskが完了したら、Topic Branchをmainにマージします。

```bash
# PR作成
gh pr create --base main \
  --title "SLEEVE-100: ユーザー登録機能の実装" \
  --body-file docs/tasks/SLEEVE-100_user_registration/topic_pr_description.md

# レビュー・承認・マージ（オーナーが実施）
# マージ方法: Merge commit
```

### 5. ブランチのクリーンアップ

```bash
# ローカルブランチを削除
git branch -d feature/SLEEVE-101_user_model
git branch -d feature/SLEEVE-102_user_repository_create
git branch -d feature/SLEEVE-103_create_user_usecase
git branch -d feature/SLEEVE-104_user_registration_api
git branch -d topic/SLEEVE-100_user_registration

# リモートブランチを削除
git push origin --delete feature/SLEEVE-101_user_model
git push origin --delete feature/SLEEVE-102_user_repository_create
git push origin --delete feature/SLEEVE-103_create_user_usecase
git push origin --delete feature/SLEEVE-104_user_registration_api
git push origin --delete topic/SLEEVE-100_user_registration
```

## トラブルシューティング

### よくある問題

#### コミットメッセージの修正

```bash
# 直前のコミットメッセージを修正
git commit --amend -m "[修正後のメッセージ]"

# 複数のコミットを修正（インタラクティブリベース）
git rebase -i HEAD~3
```

#### プッシュエラーの解決

```bash
# リモートの最新を取得してマージ
git pull origin [ブランチ名]

# 注意: 強制プッシュは使用禁止
# git push --force-with-lease は絶対に使用しない
```

#### ブランチの復旧

```bash
# 削除されたブランチを復旧
git checkout -b [ブランチ名] [コミットハッシュ]

# リモートブランチを復旧
git push origin [ブランチ名]
```
