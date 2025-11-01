
# GitHub 運用ルール

## 概要

このドキュメントは、SLEEVE プロジェクトの GitHub 運用に関するルールとガイドラインを定義します。
AI エージェントが一貫性のあるコミットとプルリクエスト管理を行うための規約です。

## 関連ドキュメント

- **全体フロー**: `docs/general.md` - 開発タスクの全体フロー
- **実装手順**: `docs/task_procedure.md` - 詳細な実装手順とテンプレート
- **コーディング規約**: `docs/BE_coding_roles.md` - アーキテクチャとコーディング規約
- **DB 操作**: `docs/DB_manual.md` - データベース操作手順
- **Linter/Formatter**: `docs/lint_format_manual.md` - コード品質管理

## 制約

### ブランチ運用（Git Flow）

- **main ブランチ**: 本番環境にデプロイされる安定版（リリース済みコード）
- **feature/[Jira の ID]\_[機能名]**: 機能開発用ブランチ（main から分岐）
- **release/[Jira の ID]\_[バージョン]**: リリース準備用ブランチ（main から分岐）
- **hotfix/[Jira の ID]\_[修正内容]**: 緊急修正用ブランチ（main から分岐）

### プルリクエスト制約

- プルリクエストは必ず main ブランチに対して作成
- レビュー必須（最低 1 名の承認）
- 全てのテストが通過している必要がある
- Linter エラーがない必要がある
- リリースブランチは main にマージ後、タグを作成

## 使用可能コマンド

### ブランチ操作

```bash
# 新しいブランチを作成
git checkout -b feature/[JiraのID]_[機能名]

# ブランチを切り替え
git checkout [ブランチ名]

# ブランチ一覧を表示
git branch -a

# ブランチを削除（ローカル）
git branch -d [ブランチ名]
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
- `test:e2e` - E2E テストの追加・修正
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

## プルリクエスト規約

### プルリクエストタイトル

```
[JiraのID]: [機能名・修正内容]
```

例:

- `SLEEVE-123: ユーザー認証機能の実装`
- `SLEEVE-456: ログイン画面のバグ修正`

### プルリクエスト説明

```markdown
## 概要

[変更内容の概要]

## 変更内容

- [具体的な変更内容 1]
- [具体的な変更内容 2]

## テスト

- [ ] 単体テストが通過
- [ ] E2E テストが通過
- [ ] 手動テストが完了

## 関連チケット

- Jira チケット ID: [Jira の ID]

## チェックリスト

- [ ] コーディング規約に準拠
- [ ] Linter エラーがない
- [ ] ドキュメントが更新されている
```

## ブランチ命名規約

### 機能開発ブランチ

```
feature/[JiraのID]_[機能名]
```

例: `feature/SLEEVE-123_user_auth`

### リリースブランチ

```
release/[JiraのID]_[バージョン]
```

例: `release/SLEEVE-200_v1.2.0`

### 緊急修正ブランチ

```
hotfix/[JiraのID]_[修正内容]
```

例: `hotfix/SLEEVE-789_security_fix`

### リファクタリングブランチ

```
refactor/[JiraのID]_[リファクタリング内容]
```

例: `refactor/SLEEVE-101_repository_structure`

## マージ規約

### マージ方法（Git Flow）

- **Squash and merge**: 機能ブランチのマージ時（feature → main）
- **Merge commit**: リリースブランチのマージ時（release → main）
- **Merge commit**: 緊急修正のマージ時（hotfix → main）
- **Rebase and merge**: リファクタリングのマージ時（refactor → main）

### リリースフロー

1. **機能開発**: feature ブランチで開発
2. **リリース準備**: release ブランチで最終調整
3. **リリース**: main ブランチにマージ + タグ作成
4. **緊急修正**: hotfix ブランチで修正

### マージ後の処理

1. ブランチの削除（リモート）
2. ローカルブランチの削除
3. 最新の main ブランチを取得
4. 次のタスクブランチの作成

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
