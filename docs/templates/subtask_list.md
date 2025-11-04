# Jira Subtask作成依頼

親タスク: [JiraのID] - [タスク名]

以下のSubtaskをJiraで作成してください。

---

## Subtask 1

**タイトル**: [Subtask名]

**説明**:
[このSubtaskで実装する内容を詳細に記述]

**依存関係**:
- なし / [依存先のSubtask番号]

**受け入れ条件**:
- [ ] [受け入れ条件1]
- [ ] [受け入れ条件2]
- [ ] [受け入れ条件3]

**対応ブランチ**: `feature/[SubtaskID]_[機能名]`（IDは作成後に確定）

---

## Subtask 2

**タイトル**: [Subtask名]

**説明**:
[このSubtaskで実装する内容を詳細に記述]

**依存関係**:
- なし / [依存先のSubtask番号]

**受け入れ条件**:
- [ ] [受け入れ条件1]
- [ ] [受け入れ条件2]

**対応ブランチ**: `feature/[SubtaskID]_[機能名]`（IDは作成後に確定）

---

（以下、必要なSubtaskを追加）

## 使用例

```markdown
# Jira Subtask作成依頼

親タスク: SLEEVE-100 - ユーザー認証機能の実装

以下のSubtaskをJiraで作成してください。

---

## Subtask 1

**タイトル**: Firebase Authenticationのセットアップ

**説明**:
Firebase Authenticationを導入し、プロジェクトでFirebase認証が利用できるように設定する。
- Firebase Admin SDKのインストールと設定
- 環境変数の設定
- Firebase初期化処理の実装

**依存関係**:
- なし

**受け入れ条件**:
- [ ] Firebase Admin SDKがインストールされている
- [ ] `firebase.secret.json`が適切に配置されている
- [ ] Firebase初期化処理が実装されている
- [ ] 接続テストが通過している

**対応ブランチ**: `feature/SLEEVE-101_firebase_setup`（IDは作成後に確定）

---

## Subtask 2

**タイトル**: Userモデルの作成

**説明**:
ユーザー情報を管理するためのモデルを作成する。
- Domain層にUserエンティティを定義
- Repository層にUserRepositoryを実装
- entスキーマの定義

**依存関係**:
- なし

**受け入れ条件**:
- [ ] User構造体が定義されている
- [ ] UserRepositoryインターフェースが定義されている
- [ ] entスキーマが定義されている
- [ ] 単体テストが通過している

**対応ブランチ**: `feature/SLEEVE-102_user_model`（IDは作成後に確定）

---

## Subtask 3

**タイトル**: ログインAPIの実装

**説明**:
Firebase AuthenticationのIDトークンを検証し、ユーザーログインを行うGraphQL APIを実装する。
- UseCase層にログインロジックを実装
- GraphQL Resolverを実装
- エラーハンドリング

**依存関係**:
- Subtask 2（Userモデルの作成）に依存

**受け入れ条件**:
- [ ] ログインUseCaseが実装されている
- [ ] GraphQL Resolverが実装されている
- [ ] 単体テストが通過している
- [ ] 結合テストが通過している

**対応ブランチ**: `feature/SLEEVE-103_login_api`（IDは作成後に確定）
```
