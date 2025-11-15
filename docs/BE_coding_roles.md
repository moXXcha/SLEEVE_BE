# 目的

- AI での vibe coding を前提としているため、プロンプトやエージェント、ツールによって生成されるコードの質が左右されない。
- AI や新規プロジェクトメンバーにプロジェクトのアーキテクチャやディレクトリ構造を理解してもらう。
- 自身で coding をする際に一貫性を保たせる。

## 関連ドキュメント

- **全体フロー**: `docs/general.md` - 開発タスクの全体フロー
- **実装手順**: `docs/task_procedure.md` - 詳細な実装手順とテンプレート
- **コーディングスタイル**: `docs/coding_guide.md` - コーディングスタイルガイド
- **Git 運用**: `docs/github.md` - Git Flow 運用とコミット規約
- **DB 操作**: `docs/DB_manual.md` - データベース操作手順
- **Linter/Formatter**: `docs/lint_format_manual.md` - コード品質管理

# アーキテクチャ

三層とクリーンアーキテクチャを組み合わせた構成。

1. **Resolver**

   - GraphQL I/O 担当。
   - 入力を受け取り → UseCase 呼び出し → 出力返すだけ。

2. **UseCase (Service)**

   - 手続きの流れを書く場所。
   - ユースケース単位の操作をまとめる。
   - AI に「自然言語でユースケースを説明」したら、ここを生成させる。

3. **Domain**
   - **Entity（エンティティ）**: ビジネス上の実体を表現（例: User, Item, Coordinate）
   - **ValueObject（値オブジェクト）**: 不変の値を表現（例: Email, Price, ImageURL）
   - **DomainService**: 複数のEntityにまたがるビジネスロジック
   - **DomainError**: ドメイン固有のエラー定義
   - **配置場所**:
     - `domain/models/code_models/` - Entity、ValueObject
     - `domain/models/db_models/` - DBスキーマに対応する構造体
     - `domain/services/` - DomainService（複数Entity間のロジック）
     - `domain/errors/` - ドメイン固有エラー
   - **責務の境界**:
     - Domain層：ビジネスルールの検証（例: メールアドレス形式チェック）
     - UseCase層：ビジネスフローの制御（例: ユーザー登録の手続き）
4. **Repository / Infrastructure**
   - 永続化・外部サービス接続。
   - interface と実装を分ける。

# ディレクトリ構造

```
be/
  app/
    config/                # シークレット情報があります（エージェントは参照することを禁止します）
      .env
      .firebase.secret.json
      etc...
    server.go              # サーバー起動時の処理
    error_presenter.go     # ErrorPresenterでカスタムエラーを定義
    ent/                   # ent関連のファイル
    tests/                 # 結合テスト置き場
      integration/         # 結合テストのディレクトリ
        [feature]/         # 機能ごとのディレクトリ
        fixtures/          # テストデータ
        helpers/           # ヘルパー関数
    graph/                 # gqlgenで生成されたファイル
      schema.resolver.go   # プロジェクトのI/O
      schema.graphqls      # スキーマ定義
      resolver.go          # 依存注入
      models/              # 生成モデル、domainモデルは置かない
    domain/                # model, customErrorの定義
      modles/
        code_models/       # コード上で使用するmodelの構造体
      services/            # DomainService（複数Entity間のロジック）
      errors/
    usecase/               # domainモデルのinterface定義、処理呼び出し、resolverへの返信（外部APIはrepository層で行うこと）
      feature/             # 機能のディレクトリ
        [機能名]/           # 機能ごとのディレクトリ
          utils/           # 機能ごとのファイルAI補助、error判定関数、parseや計算処理など
        utils/             # featureを跨ぐ処理
      utils/               # 全共通処理
    repository/            # 実際のDBアクセス
      [entity名]_dao.go    # data access objectファイル
      external_api/        # 外部API用の記述
        [service_name]_api # それぞれのサービスのディレクトリ
    messages/              # メッセージ定義（実装ファイル）
      [feature]/           # 機能ごとのメッセージディレクトリ
        errors.go          # 機能固有のエラーメッセージ定義
        logs.go            # 機能固有のログメッセージ定義
      common/              # 共通メッセージディレクトリ
        errors.go          # 共通エラーメッセージ定義
        logs.go            # 共通ログメッセージ定義
    gqlgen.yml             # gqlgenコード生成設定
    middlewares/           # ミドルウェアを定義
  docs/                    # プロジェクトdocディレクトリ
    messages/              # メッセージドキュメントディレクトリ
      [feature]/           # 機能ごとのメッセージドキュメント
        errors.md          # エラーメッセージの説明（いつ・どのメッセージが出るか）
        logs.md            # ログメッセージの説明（いつ・どのログが出るか）
      common/              # 共通メッセージドキュメント
        errors.md          # 共通エラーメッセージの説明
        logs.md            # 共通ログメッセージの説明
    tasks/                  # タスクドキュメントディレクトリ
      [jiraのtaskID]_[task名]/ # taskのドキュメントディレクトリ（タスクのたびにこのディレクトリが生成される。AIが書き込む）
  config/                  # プロジェクトconfigディレクトリ
    .env
    firebase.secret.json
  migrations/              # migrationファイル置き場
  .env
  .gitignore
  docker-compose.yml
  README.md
  Taskfile.yml
```

# 非同期(gorutine)に関しての制約

- `context.Context` を適切に引き回し、タイムアウトやキャンセルを処理できるようにする。
- 複数の goroutine の完了を待つ場合は `sync.WaitGroup` を使用する。
- goroutine 内で panic が発生してもプログラム全体がクラッシュしないように `recover` を使用して適切に処理する。
- goroutine を起動した側が、その goroutine が確実に終了することを保証する設計にする。リソースリークを防ぐ。
- 原則として、goroutine は **UseCase 層** で起動すること。ビジネスロジックの一部として非同期処理を開始・管理するのが責務であるため。
- Resolver 層や Repository 層で安易に goroutine を起動しないこと。

# 命名規則

## 変数

- 名詞スネークケース
- ブール値: is*, has*, can\_ で始める
- 数値カウンタ: xxx_count
- リスト/配列: 複数形

例: `user_id`, `todo_list`, `is_active`, `has_error`, `retry_count`

## 関数

### Public 関数

- 動詞\_目的語
- 動詞は get, list, create, update, delete, find, save
- ブール判定: is*, has*, can\_

例: `create_todo`, `get_user_by_id`, `list_orders`, `is_valid_email`

### Private 関数

- 動詞\_対象
- 動詞は load, save, parse, build, exec

例: `parse_json`, `save_to_db`, `build_query`

## ファイル

- 対象\_役割.go
- 役割: model, repository, usecase, handler, service, controller

例: `todo_repository.go`, `user_usecase.go`, `auth_handler.go`, `db_service.go`

## データベース

### テーブル

- 複数形スネークケース
- 中間テーブル: 単数形を\_で結合

例: `users`, `todos`, `order_items`, `user_roles`, `todo_tags`

### カラム

- id は id 固定
- 外部キー: 対象\_id
- 日時: \*\_at
- 論理削除: deleted_at
- フラグ: is\_ で始める

例: `id`, `user_id`, `created_at`, `updated_at`, `deleted_at`, `is_active`, `todo_text`

# まとめ表（完全ルール）

| 区分         | 命名形式                                     | 例                                         |
| ------------ | -------------------------------------------- | ------------------------------------------ |
| 変数         | 名詞スネークケース、真偽値は is/has/can      | `user_id`, `is_active`, `todo_list`        |
| Public 関数  | 動詞\_目的語（動詞は定義済みセットから選ぶ） | `create_todo`, `get_user_by_id`            |
| Private 関数 | 動詞\_対象                                   | `parse_json`, `save_to_db`                 |
| ファイル     | 対象\_役割.go（役割は固定語）                | `user_repository.go`, `auth_handler.go`    |
| テーブル     | 複数形スネークケース                         | `users`, `order_items`                     |
| カラム       | id 固定, _\_id, \_\_at, is_\*                | `id`, `user_id`, `created_at`, `is_active` |

- 英語を使用
- スネークケースで統一

# 禁止事項

- ディレクトリ内でアーキテクチャでの役割以上のことをさせる
- 致命的なエラー以外の panic
- エラーを schema.resolver.go でカスタムエラーに変換して FE に返さず関数内で処理すること
- 外部に export するものの最初の文字を小文字にする
- 循環 import
- マジックナンバー
- ハードコーディング
- UseCase 外での DB トランザクションの開始/コミット禁止。

# 制限事項

- 大きすぎる interface の作成
- 無駄な型変換
- 不要なポインタ渡し
- 不必要なグローバル変数の定義

# 推奨事項

- 標準ライブラリの使用 (ライブラリの指定がなかった場合)
- コンテナ型の初期化
- コメント
- 単一責任の interface
