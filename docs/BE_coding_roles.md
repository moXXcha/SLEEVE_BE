# 目的
- AIでのvibe codingを前提としているため、プロンプトやエージェント、ツールによって生成されるコードの質が左右されない。
- AIや新規プロジェクトメンバーにプロジェクトのアーキテクチャやディレクトリ構造を理解してもらう。
- 自身でcodingをする際に一貫性を保たせる。

# アーキテクチャ
三層とクリーンアーキテクチャを組み合わせた構成。

1. **Resolver**
    - GraphQL I/O担当。
    - 入力を受け取り → UseCase呼び出し → 出力返すだけ。

2. **UseCase (Service)**
    - 手続きの流れを書く場所。
    - ユースケース単位の操作をまとめる。
    - AIに「自然言語でユースケースを説明」したら、ここを生成させる。

3. **Domain**
    - modelの定義
    - ルールや制約の定義
4. **Repository / Infrastructure**
    - 永続化・外部サービス接続。
    - interface と実装を分ける。

# ディレクトリ構造
```
be/
  app/
    server.go              # サーバー起動時の処理
    error_presenter.go     # ErrorPresenterでカスタムエラーを定義
    graph/                 # gqlgenで生成されたファイル
      schema.resolver.go   # プロジェクトのI/O
      schema.graphqls      # スキーマ定義
      resolver.go          # 依存注入
      models/              # 生成モデル、domainモデルは置かない
    domain/                # model, customErrorの定義
    usecase/               # domainモデルのinterface定義、処理呼び出し、resolverへの返信（外部APIはrepository層で行うこと）
      utils/               # 共通処理、AI補助、error判定関数、parseや計算処理など
    repository/            # 実際のDBアクセス
      external_api/        # 外部API用の記述
        [service_name]_api # それぞれのサービスのディレクトリ
    gqlgen.yml             # gqlgenコード生成設定
    middlewares/           # ミドルウェアを定義
  docs/                    # プロジェクトdocディレクトリ
    messages/
      error_messages.go    # ユーザーへのエラーメッセージの定義(日本語)
      log_messages.go      # ログメッセージの定義（日本語）
    tasks/
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
- 複数のgoroutineの完了を待つ場合は `sync.WaitGroup` を使用する。
- goroutine内でpanicが発生してもプログラム全体がクラッシュしないように `recover` を使用して適切に処理する。
- goroutineを起動した側が、そのgoroutineが確実に終了することを保証する設計にする。リソースリークを防ぐ。
- 原則として、goroutineは **UseCase層** で起動すること。ビジネスロジックの一部として非同期処理を開始・管理するのが責務であるため。
- Resolver層やRepository層で安易にgoroutineを起動しないこと。

# 命名規則
## 変数
- 名詞スネークケース
- ブール値: is_, has_, can_ で始める
- 数値カウンタ: xxx_count
- リスト/配列: 複数形

例: `user_id`, `todo_list`, `is_active`, `has_error`, `retry_count`

## 関数
### Public関数
- 動詞_目的語
- 動詞は get, list, create, update, delete, find, save
- ブール判定: is_, has_, can_

例: `create_todo`, `get_user_by_id`, `list_orders`, `is_valid_email`

### Private関数
- 動詞_対象
- 動詞は load, save, parse, build, exec

例: `parse_json`, `save_to_db`, `build_query`

## ファイル
- 対象_役割.go
- 役割: model, repository, usecase, handler, service, controller

例: `todo_repository.go`, `user_usecase.go`, `auth_handler.go`, `db_service.go`

## データベース
### テーブル
- 複数形スネークケース
- 中間テーブル: 単数形を_で結合

例: `users`, `todos`, `order_items`, `user_roles`, `todo_tags`

### カラム
- id は id 固定
- 外部キー: 対象_id
- 日時: *_at
- 論理削除: deleted_at
- フラグ: is_ で始める

例: `id`, `user_id`, `created_at`, `updated_at`, `deleted_at`, `is_active`, `todo_text`

# まとめ表（完全ルール）
| 区分       | 命名形式 | 例 |
|------------|-----------|--------------------------------|
| 変数       | 名詞スネークケース、真偽値は is/has/can | `user_id`, `is_active`, `todo_list` |
| Public関数 | 動詞_目的語（動詞は定義済みセットから選ぶ） | `create_todo`, `get_user_by_id` |
| Private関数| 動詞_対象 | `parse_json`, `save_to_db` |
| ファイル   | 対象_役割.go（役割は固定語） | `user_repository.go`, `auth_handler.go` |
| テーブル   | 複数形スネークケース | `users`, `order_items` |
| カラム     | id固定, *_id, *_at, is_* | `id`, `user_id`, `created_at`, `is_active` |

- 英語を使用
- スネークケースで統一

# 禁止事項
- ディレクトリ内でアーキテクチャでの役割以上のことをさせる
- 致命的なエラー以外のpanic
- エラーを schema.resolver.go でカスタムエラーに変換して FE に返さず関数内で処理すること
- 外部に export するものの最初の文字を小文字にする
- 循環 import
- マジックナンバー
- ハードコーディング
- UseCase外でのDBトランザクションの開始/コミット禁止。

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
