
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

- **1PR = 1コンテキスト**: 1つのPR（1つのブランチ）では、修正・作成内容が1つのコンテキストにまとまるようにする
- **Topic Branch運用**: 大きなタスクは複数のサブタスクに分割し、Topic Branchを親として管理する
- **依存関係の許容**: PRの依存関係を許容し、依存先を親としてブランチを切る
- **Jira連携**: ブランチごとにJiraでSubtaskを作成し、進捗管理を明確にする

### ブランチの種類

#### main ブランチ
- 本番環境にデプロイされる安定版（リリース済みコード）
- 直接コミット禁止
- 保護設定あり（レビュー必須、CI通過必須）

#### topic/[JiraのID]_[機能名]
- **用途**: 大きな機能追加や、複数のサブタスクに分かれる作業の親ブランチ
- **分岐元**: main
- **マージ先**: main
- **作成タイミング**:
  - 3つ以上のサブタスクに分かれる場合
  - 大きな機能追加の場合
  - 複数人で並行開発する場合
- **対応Jira**: Story/Task
- **例**: `topic/SLEEVE-100_user_authentication`

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

大きな機能「ユーザー認証機能」を実装する場合：

```
main
 └─ topic/SLEEVE-100_user_authentication (Story/Task)
      ├─ feature/SLEEVE-101_firebase_setup (Subtask 1: 依存なし)
      ├─ feature/SLEEVE-102_user_model (Subtask 2: 依存なし)
      ├─ feature/SLEEVE-103_login_api (Subtask 3: 102に依存)
      │    └─ feature/SLEEVE-104_login_validation (Subtask 4: 103に依存)
      └─ feature/SLEEVE-105_logout_api (Subtask 5: 依存なし)
```

**マージフロー**:
1. `feature/SLEEVE-101` → `topic/SLEEVE-100` (依存なし)
2. `feature/SLEEVE-102` → `topic/SLEEVE-100` (依存なし)
3. `feature/SLEEVE-103` → `feature/SLEEVE-102` (102に依存) → 102がtopicにマージ後、103もtopicにマージ
4. `feature/SLEEVE-104` → `feature/SLEEVE-103` (103に依存) → 103がtopicにマージ後、104もtopicにマージ
5. `feature/SLEEVE-105` → `topic/SLEEVE-100` (依存なし)
6. 全てのサブタスク完了後: `topic/SLEEVE-100` → `main`

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
- 複数のSubtaskを持つ
- 例: `SLEEVE-100: ユーザー認証機能の実装`

#### Subtask
- Feature Branchに対応
- 1つの具体的な実装・修正内容
- 例:
  - `SLEEVE-101: Firebase Authenticationのセットアップ`
  - `SLEEVE-102: Userモデルの作成`
  - `SLEEVE-103: ログインAPIの実装`

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

# 例
git checkout -b topic/SLEEVE-100_user_authentication
```

#### Feature Branchの作成（依存関係なし）
```bash
# Topic Branchから新しいFeature Branchを作成
git checkout topic/[JiraのID]_[機能名]
git pull origin topic/[JiraのID]_[機能名]
git checkout -b feature/[JiraのSubtaskID]_[機能名]
git push -u origin feature/[JiraのSubtaskID]_[機能名]

# 例
git checkout topic/SLEEVE-100_user_authentication
git checkout -b feature/SLEEVE-101_firebase_setup
```

#### Feature Branchの作成（依存関係あり）
```bash
# 依存先のFeature Branchから新しいFeature Branchを作成
git checkout feature/[依存先のJiraID]_[機能名]
git pull origin feature/[依存先のJiraID]_[機能名]
git checkout -b feature/[JiraのSubtaskID]_[機能名]
git push -u origin feature/[JiraのSubtaskID]_[機能名]

# 例: SLEEVE-103がSLEEVE-102に依存する場合
git checkout feature/SLEEVE-102_user_model
git checkout -b feature/SLEEVE-103_login_api
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


## 実践例：ユーザー認証機能の開発フロー

この例では、新しいブランチ戦略を使用して「ユーザー認証機能」を開発するフローを説明します。

### 1. Jiraでのタスク管理

**Epic**: `EPIC-50: ユーザー管理機能の実装`

**Story/Task**: `SLEEVE-100: ユーザー認証機能の実装`

**Subtask**:
- `SLEEVE-101: Firebase Authenticationのセットアップ`
- `SLEEVE-102: Userモデルの作成`
- `SLEEVE-103: ログインAPIの実装` (102に依存)
- `SLEEVE-104: ログアウトAPIの実装`

### 2. Topic Branchの作成

```bash
git checkout main
git pull origin main
git checkout -b topic/SLEEVE-100_user_authentication
git push -u origin topic/SLEEVE-100_user_authentication
```

### 3. 各Subtaskの実装

#### Subtask 1: Firebase Authenticationのセットアップ（依存なし）

```bash
# Feature Branchを作成
git checkout topic/SLEEVE-100_user_authentication
git checkout -b feature/SLEEVE-101_firebase_setup

# 実装...

# PR作成: feature/SLEEVE-101 → topic/SLEEVE-100
# マージ方法: Squash and merge
```

#### Subtask 2: Userモデルの作成（依存なし）

```bash
# Feature Branchを作成
git checkout topic/SLEEVE-100_user_authentication
git checkout -b feature/SLEEVE-102_user_model

# 実装...

# PR作成: feature/SLEEVE-102 → topic/SLEEVE-100
# マージ方法: Squash and merge
```

#### Subtask 3: ログインAPIの実装（102に依存）

```bash
# 依存先のFeature Branchから作成
git checkout feature/SLEEVE-102_user_model
git checkout -b feature/SLEEVE-103_login_api

# 実装...

# PR作成: feature/SLEEVE-103 → feature/SLEEVE-102
# マージ方法: Merge commit

# 102がtopic branchにマージされた後、103もtopic branchにマージ
# PR作成: feature/SLEEVE-103 → topic/SLEEVE-100
# マージ方法: Squash and merge
```

#### Subtask 4: ログアウトAPIの実装（依存なし）

```bash
# Feature Branchを作成
git checkout topic/SLEEVE-100_user_authentication
git checkout -b feature/SLEEVE-104_logout_api

# 実装...

# PR作成: feature/SLEEVE-104 → topic/SLEEVE-100
# マージ方法: Squash and merge
```

### 4. Topic Branchをmainにマージ

全てのSubtaskが完了したら、Topic Branchをmainにマージします。

```bash
# PR作成: topic/SLEEVE-100_user_authentication → main
# マージ方法: Merge commit
```

### 5. ブランチのクリーンアップ

```bash
# ローカルブランチを削除
git branch -d feature/SLEEVE-101_firebase_setup
git branch -d feature/SLEEVE-102_user_model
git branch -d feature/SLEEVE-103_login_api
git branch -d feature/SLEEVE-104_logout_api
git branch -d topic/SLEEVE-100_user_authentication

# リモートブランチを削除
git push origin --delete feature/SLEEVE-101_firebase_setup
git push origin --delete feature/SLEEVE-102_user_model
git push origin --delete feature/SLEEVE-103_login_api
git push origin --delete feature/SLEEVE-104_logout_api
git push origin --delete topic/SLEEVE-100_user_authentication
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
