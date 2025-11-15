# SLEEVE バックエンド

## 📖 プロジェクト概要

SLEEVEは、コーデ投稿とアイテム売買を融合させたCtoC型ファッションリユースプラットフォームです。

**スローガン**: 「あの人らしいが自分らしいになるところ。」

### 主な機能
- コーデ投稿機能（自己表現＆収益化）
- アイテムのCtoC売買
- マネキン買い機能（気に入ったコーデをそのまま購入）
- リユース品・類似品による低コストなトレンド再現

詳細は [docs/prod_docs/SLEEVE 外部用サービス概要書.md](./docs/prod_docs/SLEEVE%20外部用サービス概要書.md) を参照してください。

---

## 🛠️ 使用技術

### バックエンド
| 項目 | 技術 | バージョン |
|------|------|-----------|
| 言語 | Go | 1.25.1 |
| フレームワーク | gqlgen | 0.17.81 |
| API | GraphQL | - |
| ORM | ent | - |
| Linter | golangci-lint | - |
| Formatter | goimports | - |

### データベース
| 項目 | 技術 |
|------|------|
| メインDB | PostgreSQL 14 |
| キャッシュ・ベクトル検索 | Redis |
| 検索エンジン | Managed Elasticsearch on GCP |

### 認証・セキュリティ
| 項目 | 技術 |
|------|------|
| 認証基盤 | Firebase Authentication |
| 認可方針 | JWT + リフレッシュToken |
| パスワード暗号化 | Argon2 |
| OAuth連携 | Google / Apple / X(Twitter) |
| MFA | SMS |
| 秘密情報管理 | GCP Secret Manager |

### インフラ
| 項目 | 技術 |
|------|------|
| クラウド | GCP |
| コンテナ | Docker / Kubernetes |
| CI/CD | GitHub Actions |
| IaC | Terraform |

### アーキテクチャ
- **設計**: 3層アーキテクチャ（クリーンアーキテクチャの要素を含む）
  - **Resolver**: GraphQL I/O担当
  - **UseCase (Service)**: ビジネスロジックの流れ
  - **Domain**: エンティティ、ビジネスルール、バリデーション
  - **Repository / Infrastructure**: 永続化・外部サービス接続

詳細は [docs/prod_docs/技術選定.md](./docs/prod_docs/技術選定.md) および [docs/BE_coding_roles.md](./docs/BE_coding_roles.md) を参照してください。

---

## 🚀 環境構築方法

本プロジェクトではできるだけローカルの環境への依存が少ないように基本ツールはコンテナに入れ、コンテナを通して使用する形にしています

### 必要なツール
- Docker & Docker Compose
- Go 1.25.1以上
- Task (タスクランナー) - [インストール方法](https://taskfile.dev/installation/)
- golangci-lint - [インストール方法](https://golangci-lint.run/usage/install/)

### セットアップ手順

#### 1. リポジトリのクローン
```bash
git clone <repository-url>
cd be
```

#### 2. 環境変数の設定
```bash
# config/.env ファイルを作成
cp config/.env.example config/.env

# .env ファイルを編集（必要に応じて値を設定）
# - POSTGRES_HOST
# - POSTGRES_PORT
# - POSTGRES_USER
# - POSTGRES_PASSWORD
# - POSTGRES_DB
```

#### 3. Dockerコンテナの起動
```bash
# 初回起動（全てクリーンアップして起動）
task init

# 通常起動
task up
```

#### 4. 動作確認
```bash
# GraphQL Playgroundにアクセス
open http://localhost:8080/playground
```

#### 5. 開発環境のセットアップ（ローカル開発時）
```bash
cd app

# 依存関係のインストール
go mod download

# GraphQLスキーマからコード生成
task generate

# Linter & Formatterの実行
task lint-fmt
```

---

## 📋 Taskコマンド一覧

このプロジェクトでは、[Task](https://taskfile.dev/)をタスクランナーとして使用しています。

### 主要コマンド

| コマンド | 説明 |
|---------|------|
| `task` | デフォルトタスク（Hello, World!を表示） |
| `task up` | Dockerコンテナを起動 |
| `task down` | Dockerコンテナを停止 |
| `task init` | 環境を完全にクリーンアップして再起動 |
| `task build` | Dockerイメージをビルド |
| `task generate` | GraphQLスキーマからGoコードを生成（gqlgen） |
| `task lint-fmt` | Linterとフォーマッターを実行 |
| `task sql-gen` | migrationファイルを生成 |
| `task ent-gen` | entのコードを生成 |
| `task migrate-up` | migrateの実行 |
| `task migrate-rollback` | 1件migrateをrollback |
| `task move-app-container` | appコンテナに移動 |
| `task move-db-container` | dbコンテナに移動 |

### コマンド詳細

#### `task up`
```bash
task up
```
- Dockerコンテナをバックグラウンドで起動します
- PostgreSQLとバックエンドサーバーが起動します

#### `task down`
```bash
task down
```
- 実行中のDockerコンテナを停止します

#### `task init`
```bash
task init
```
- **注意**: 全てのDockerイメージ、ボリューム、孤立したコンテナを削除します
- データベースのデータも削除されます
- 開発環境を完全にクリーンな状態から再構築したい場合に使用します

#### `task build`
```bash
task build
```
- Dockerイメージを再ビルドします
- Dockerfileに変更を加えた場合に実行します

#### `task generate`
```bash
task generate
```
- `app/graph/schema.graphqls` からGoコードを自動生成します
- GraphQLスキーマを変更した場合は必ず実行してください

#### `task lint-fmt`
```bash
task lint-fmt
```
- `golangci-lint`を実行してコード品質をチェックします
- `--exclude-dirs ./graph` オプションで生成されたコードを除外しています
- コミット前に必ず実行してください

etc...

---

## 🤖 使用可能な生成AI

### ✅ 推奨AI
このプロジェクトでの使用を推奨するAIツール:

- **GitHub Copilot** - コード補完・生成
- **Claude** / **Claude Code** - コードレビュー、設計相談、実装支援
- **Codex** - コード生成
- **Cursor** - AI統合エディタ

### ❌ 使用禁止AI
以下のAIツールは使用を禁止します:

- **DeepShark** - 使用禁止
- **Gemini CLI** - 使用禁止

**理由**: プロジェクトのセキュリティポリシーおよびコード品質管理の観点から、上記AIツールの使用は認められていません。

---

## 📂 ディレクトリ構造

```
be/
  app/                    # アプリケーションコード
    config/               # シークレット情報（.gitignore対象）
    graph/                # GraphQLスキーマ・リゾルバー
      schema.graphqls     # スキーマ定義
      schema.resolver.go  # リゾルバー実装
    domain/               # ドメインモデル・エラー定義
      models/
        db_models/        # DBテーブル構造体
        code_models/      # コード上のモデル構造体
      services/           # ドメインサービス
      errors/             # ドメインエラー
    usecase/              # ビジネスロジック
    repository/           # データアクセス層
    messages/             # メッセージ定義（エラー・ログ）
    middlewares/          # ミドルウェア
    tests/                # 結合テスト
  config/                 # 環境設定ファイル
  docs/                   # プロジェクトドキュメント
    prod_docs/            # プロダクト仕様書
    templates/            # ドキュメントテンプレート
    tasks/                # タスクドキュメント
  migrations/             # マイグレーションファイル
  docker-compose.yml      # Docker構成
  Taskfile.yml            # Taskコマンド定義
```

---

## 📚 ドキュメント

### 開発者向けドキュメント
プロジェクトの開発に関する詳細なドキュメントは `docs/` ディレクトリに配置されています。

#### クイックスタート
- **[docs/readme.md](./docs/readme.md)** - ドキュメント全体の目次
- **[docs/general.md](./docs/general.md)** - 開発タスクの全体フロー（必読）
- **[docs/BE_coding_roles.md](./docs/BE_coding_roles.md)** - アーキテクチャとコーディング規約（必読）

#### 詳細ドキュメント
- **コーディング規約**:
  - [docs/coding_guide.md](./docs/coding_guide.md) - コーディングスタイルガイド
  - [docs/lint_format_manual.md](./docs/lint_format_manual.md) - Linter・Formatter手引き
- **開発フロー**:
  - [docs/task_procedure.md](./docs/task_procedure.md) - Subtask実装手順
  - [docs/github.md](./docs/github.md) - Git Flow運用とコミット規約
- **データベース**:
  - [docs/DB_manual.md](./docs/DB_manual.md) - DB操作マニュアル
  - [docs/db_schema.md](./docs/db_schema.md) - DBスキーマ図

#### プロダクトドキュメント
- [docs/prod_docs/SLEEVE 外部用サービス概要書.md](./docs/prod_docs/SLEEVE%20外部用サービス概要書.md) - サービス概要
- [docs/prod_docs/技術選定.md](./docs/prod_docs/技術選定.md) - 技術選定理由
- [docs/prod_docs/機能一覧.md](./docs/prod_docs/機能一覧.md) - 実装予定機能
- [docs/prod_docs/非機能要件書.md](./docs/prod_docs/非機能要件書.md) - 非機能要件

### AIエージェント向け
AIエージェントがタスクを実施する際は、以下のドキュメントを必ず読んでください:

1. **[docs/general.md](./docs/general.md)** - 全体フローの理解（必須）
2. **[docs/BE_coding_roles.md](./docs/BE_coding_roles.md)** - アーキテクチャと命名規則
3. **[docs/coding_guide.md](./docs/coding_guide.md)** - コーディングスタイル
4. **[docs/task_procedure.md](./docs/task_procedure.md)** - 実装手順

---

## 🧪 テスト

### テストの実行
```bash
# 全テストを実行
cd app
go test ./...

# カバレッジ付きで実行
go test -cover ./...

# 結合テストのみ実行
cd tests
go test ./...
```

### テスト作成ガイドライン
- **単体テスト**: 各関数の隣に `*_test.go` として配置
- **結合テスト**: `tests/integration/[feature]/` に配置
- **テストカバレッジ**: 80%以上を目標

---

## 🔧 開発フロー

### 1. ブランチ戦略
- **main**: 本番環境（直接コミット禁止）
- **topic/[ID]_[機能名]**: 機能単位の親ブランチ
- **feature/[ID]_[機能名]**: 実際の実装ブランチ

詳細は [docs/github.md](./docs/github.md) を参照してください。

### 2. コミット前のチェックリスト
```bash
# 1. Linterを実行
cd app
golangci-lint run --fix

# 2. テストを実行
go test ./...

# 3. コミット（規約に従ったメッセージ）
git add .
git commit -m "feat: 新機能を追加

- 詳細な変更内容
- 関連するJiraチケットID: SLEEVE-XXX"
```

### 3. PR作成
```bash
# GitHub CLIを使用
gh pr create --base topic/[親ID]_[機能名] --title "[ID]: [タスク名]" --body-file docs/tasks/[ID]_[タスク名]/pr_description.md
```

---

## 🚨 トラブルシューティング

### Dockerコンテナが起動しない
```bash
# ログを確認
docker compose logs

# 環境をクリーンアップして再起動
task init
```

### ポート競合エラー
```bash
# 既存のコンテナを停止
task down

# ポートを使用しているプロセスを確認
lsof -i :8080
lsof -i :5432
```

### GraphQLコード生成エラー
```bash
# スキーマファイルの構文を確認
cat app/graph/schema.graphqls

# 再生成
task generate
```

### Linterエラー
詳細は [docs/lint_format_manual.md](./docs/lint_format_manual.md) を参照してください。

---

## 👥 コントリビューション

### 新規メンバーへ
1. このREADMEを読む
2. [docs/readme.md](./docs/readme.md) でドキュメント全体を把握
3. [docs/general.md](./docs/general.md) で開発フローを理解
4. [docs/BE_coding_roles.md](./docs/BE_coding_roles.md) でコーディング規約を確認

### 開発タスクの開始前に
1. **[docs/general.md](./docs/general.md)** を読み、全体フローを理解
2. 作業計画書を作成し、オーナーの承認を得る
3. ブランチを作成し、実装を開始

---

## 📞 お問い合わせ

プロジェクトオーナー: [@moXXcha](https://github.com/moXXcha)

---

## 📄 ライセンス

このプロジェクトは内部プロジェクトです。外部への公開や再配布は禁止されています。

---

**Last Updated**: 2025年1月
**Maintained by**: SLEEVE Development Team
