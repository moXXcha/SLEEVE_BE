# SLEEVE バックエンド ドキュメント

このディレクトリには、SLEEVE プロジェクトのバックエンド開発に関する全てのドキュメントが含まれています。

## 📚 目次

- [クイックスタート](#クイックスタート)
- [開発フロー・手順書](#開発フロー手順書)
- [コーディング規約・スタイル](#コーディング規約スタイル)
- [データベース](#データベース)
- [Git/GitHub 運用](#gitgithub運用)
- [テンプレート](#テンプレート)
- [プロダクトドキュメント](#プロダクトドキュメント)

---

## 🚀 クイックスタート

### 新しいタスクを始める前に

1. **[general.md](./general.md)** - 開発タスクの全体フローを理解する
2. **[BE_coding_roles.md](./BE_coding_roles.md)** - アーキテクチャと命名規則を確認する
3. **[coding_guide.md](./coding_guide.md)** - コーディングスタイルを確認する

### AI エージェント向け

AI エージェントは、タスク開始前に必ず **[general.md](./general.md)** を読み、プロジェクトの規約を完全に理解してください。

---

## 📖 開発フロー・手順書

### [general.md](./general.md)

**開発タスクの全体フロー**

- **対象**: AI エージェント、開発者全員
- **内容**:
  - Level 1: タスク全体の作業手順（親タスク）
  - Level 2: Subtask 内の作業手順
  - 作業計画の策定から PR 作成までの完全なフロー
- **いつ読む**: タスク開始前（必須）

### [task_procedure.md](./task_procedure.md)

**Subtask 内の詳細な実装手順**

- **対象**: AI エージェント、開発者
- **内容**:
  - 実装前準備
  - メッセージ定義の作成
  - 実装実行の詳細手順
  - Git Commit 手順
  - PR 作成とレビューフロー
  - ドキュメントテンプレートの参照
- **いつ読む**: 実装中（参照）

---

## 🎨 コーディング規約・スタイル

### [BE_coding_roles.md](./BE_coding_roles.md)

**アーキテクチャとコーディング規約**

- **対象**: 開発者全員（必読）
- **内容**:
  - プロジェクトアーキテクチャ（レイヤー構造）
  - 命名規則（関数、変数、ファイル、パッケージ）
  - ディレクトリ構造
  - 各層（Domain, Repository, UseCase, GraphQL）の役割と責務
- **いつ読む**: タスク開始前、命名に迷った時

### [coding_guide.md](./coding_guide.md)

**コーディングスタイルガイド**

- **対象**: 開発者全員（必読）
- **内容**:
  - 変数の定義ルール（スコープ、初期化）
  - 改行ルール
  - エラーハンドリング
  - コメント記述
  - 関数設計
  - パッケージ設計
  - チェックリスト
- **いつ読む**: コーディング中（参照）

### [lint_format_manual.md](./lint_format_manual.md)

**Linter と Formatter の手引き**

- **対象**: 開発者全員
- **内容**:
  - `golangci-lint` の設定内容
  - 実行方法（`--fix` オプション）
  - 開発フローでの使用タイミング
  - 各 Linter の詳細説明
  - エラー対応
  - トラブルシューティング
- **いつ読む**: Linter エラーが発生した時、コミット前

---

## 🗄️ データベース

### [DB_manual.md](./DB_manual.md)

**データベース取り扱いマニュアル**

- **対象**: DB 変更を行う開発者、AI エージェント
- **内容**:
  - DB 作成の方針
  - スキーマの位置（`ent/schema/`, `docs/db_schema.md`）
  - 実装手順（4 ステップ）
  - 変更計画書テンプレート
  - **重要**: `docs/db_schema.md` の更新が必須
  - トランザクション管理（UseCase 層と Repository 層の役割）
- **いつ読む**: テーブルの作成・編集・削除を行う前（必須）

### [db_schema.md](./db_schema.md)

**データベーススキーマ図**

- **対象**: 開発者全員
- **内容**:
  - 全テーブルの構造（dbdiagram.io 形式）
  - リレーション（外部キー）
  - インデックス
- **いつ読む**: DB 変更前、テーブル構造を確認したい時
- **更新**: テーブルの作成・編集・削除時に必ず更新

---

## 🔀 Git/GitHub 運用

### [github.md](./github.md)

**Git Flow 運用とコミット規約**

- **対象**: 開発者全員（必読）
- **内容**:
  - ブランチ戦略（Topic Branch, Feature Branch）
  - ブランチ命名規則
  - マージ方法（Squash and merge / Merge commit）
  - コミットメッセージ規約
  - PR 作成ルール
  - レビュー待機フロー（現状の運用）
- **いつ読む**: ブランチ作成前、コミット前、PR 作成前

---

## 📝 テンプレート

テンプレートは `templates/` ディレクトリに配置されています。

### [templates/task_plan.md](./templates/task_plan.md)

**作業計画書テンプレート**

- **用途**: `docs/tasks/[JiraのID]_[タスク名]/task_plan.md` として作成
- **内容**: ブランチ戦略、テスト通過基準、メッセージ定義、使用技術、作成/変更ファイル

### [templates/implementation_details.md](./templates/implementation_details.md)

**実装詳細書テンプレート**

- **用途**: `docs/tasks/[JiraのID]_[タスク名]/implementation_details.md` として作成
- **内容**: 概要、変更内容、テスト、考慮事項

### [templates/subtask_list.md](./templates/subtask_list.md)

**Jira Subtask 内容テンプレート**

- **用途**: `docs/tasks/[JiraのID]_[タスク名]/subtask_list.md` として作成
- **内容**: Subtask のタイトル、説明、依存関係、受け入れ条件

### [templates/topic_pr_description.md](./templates/topic_pr_description.md)

**Topic PR ディスクリプションテンプレート**

- **用途**: Topic Branch から main への PR 作成時
- **内容**: 完了した Subtask 一覧、変更内容、テスト、チェックリスト

### [templates/debug.md](./templates/debug.md)

**デバッグ調査書テンプレート**

- **用途**: `docs/tasks/[JiraのID]_[タスク名]/debug.md` として作成
- **内容**: 不具合内容、エラー文、エラー箇所、修正方針、影響範囲

### PR テンプレート

- **場所**: `../.github/pull_request_template.md`
- **用途**: GitHub 上で PR 作成時に自動適用
- **内容**: 概要、変更内容、依存関係、テスト、関連チケット、チェックリスト

---

## 📊 プロダクトドキュメント

プロダクトドキュメントは `prod_docs/` ディレクトリに配置されています。

### [prod_docs/SLEEVE 外部用サービス概要書.md](./prod_docs/SLEEVE%20外部用サービス概要書.md)

**サービス概要書**

- プロダクトの目的、ターゲットユーザー、提供価値

### [prod_docs/機能一覧.md](./prod_docs/機能一覧.md)

**機能一覧**

- 実装予定の全機能リスト

### [prod_docs/技術選定.md](./prod_docs/技術選定.md)

**技術選定**

- バックエンド・フロントエンドの技術スタック
- 選定理由

### [prod_docs/非機能要件書.md](./prod_docs/非機能要件書.md)

**非機能要件書**

- パフォーマンス、セキュリティ、可用性などの要件

### [prod_docs/開発運用フロー図.md](./prod_docs/開発運用フロー図.md)

**開発運用フロー図**

- 開発からデプロイまでのフロー

### [prod_docs/不正ユーザー対応フロー.md](./prod_docs/不正ユーザー対応フロー.md)

**不正ユーザー対応フロー**

- 不正行為を検知した場合の対応手順

### [prod_docs/SLEEVE 追加開発バックログ.md](./prod_docs/SLEEVE%20追加開発バックログ.md)

**追加開発バックログ**

- 今後の開発予定項目

### その他

- **画面遷移図**: `prod_docs/画面遷移図.canvas`
- **技術選定詳細**: `prod_docs/技術選定/` ディレクトリ

---

## 🎯 シチュエーション別ガイド

### 新しいタスクを始める時

1. **[general.md](./general.md)** - 全体フローを理解
2. **[BE_coding_roles.md](./BE_coding_roles.md)** - アーキテクチャを確認
3. **[coding_guide.md](./coding_guide.md)** - コーディングスタイルを確認
4. **[templates/task_plan.md](./templates/task_plan.md)** - 作業計画書を作成

### DB 変更を行う時

1. **[db_schema.md](./db_schema.md)** - 現在のスキーマを確認
2. **[DB_manual.md](./DB_manual.md)** - 変更手順を確認
3. 変更計画書を作成 → 承認
4. 実装
5. **[db_schema.md](./db_schema.md)** を更新（必須）

### コーディング中

1. **[coding_guide.md](./coding_guide.md)** - スタイルガイドを参照
2. **[BE_coding_roles.md](./BE_coding_roles.md)** - 命名規則を確認
3. こまめに `golangci-lint run --fix` を実行

### コミット前

1. `golangci-lint run --fix` を実行
2. 全テストを実行
3. **[github.md](./github.md)** - コミットメッセージ規約を確認
4. DB 変更を行った場合は **[db_schema.md](./db_schema.md)** の更新を確認

### PR 作成前

1. `golangci-lint run --fix` を実行
2. 全テストを実行
3. **[github.md](./github.md)** - PR テンプレートを確認
4. DB 変更を行った場合は **[db_schema.md](./db_schema.md)** の更新を確認

### Linter エラーが発生した時

1. **[lint_format_manual.md](./lint_format_manual.md)** - エラー対応を確認
2. **[coding_guide.md](./coding_guide.md)** - 該当ルールを確認
3. 解決できない場合は **[templates/debug.md](./templates/debug.md)** でデバッグ調査書を作成

---

## 🔄 ドキュメント間の関係

```
general.md （全体フロー）
    ↓
    ├─ task_procedure.md （詳細手順）
    │   ├─ BE_coding_roles.md （アーキテクチャ・命名）
    │   ├─ coding_guide.md （スタイル）
    │   ├─ lint_format_manual.md （Linter）
    │   ├─ DB_manual.md （DB操作）
    │   │   └─ db_schema.md （スキーマ図）
    │   └─ github.md （Git運用）
    │
    └─ templates/ （各種テンプレート）
        ├─ task_plan.md
        ├─ implementation_details.md
        ├─ subtask_list.md
        ├─ topic_pr_description.md
        └─ debug.md
```

---

## 📌 重要な注意事項

### DB 変更時の必須作業

テーブルの作成・編集・削除を行った場合は、**必ず `docs/db_schema.md` を更新してください**。これは必須事項です。

詳細は **[DB_manual.md](./DB_manual.md)** の「手順 4: スキーマの更新」を参照してください。

### コミット前の必須チェック

- [ ] `golangci-lint run --fix` でエラーが 0 件
- [ ] 全ての単体テストが通過
- [ ] コミットメッセージが規約に準拠
- [ ] DB 変更を行った場合は `docs/db_schema.md` を更新

### PR 作成前の必須チェック

- [ ] `golangci-lint run --fix` でエラーが 0 件
- [ ] 全ての単体テストが通過
- [ ] コーディング規約に準拠
- [ ] ドキュメントが更新されている
- [ ] DB 変更を行った場合は `docs/db_schema.md` を更新

---

## 🤖 AI エージェント向けガイドライン

AI エージェントがタスクを実施する際は、以下の順序でドキュメントを読んでください：

1. **[general.md](./general.md)** - 全体フローの理解（必須）
2. **[BE_coding_roles.md](./BE_coding_roles.md)** - アーキテクチャと命名規則の確認
3. **[coding_guide.md](./coding_guide.md)** - コーディングスタイルの確認
4. **[task_procedure.md](./task_procedure.md)** - 実装中の参照
5. DB 変更がある場合: **[DB_manual.md](./DB_manual.md)** と **[db_schema.md](./db_schema.md)**

### 必ず守るべきルール

- タスク開始前に必ず `docs/` 配下の全ドキュメントを読む
- DB 変更を行った場合は必ず `docs/db_schema.md` を更新
- コミット前に必ず `golangci-lint run --fix` を実行
- テンプレートを活用して作業計画書・実装詳細書を作成

---

## 📞 問い合わせ・フィードバック

ドキュメントに不明点や改善提案がある場合は、プロジェクトオーナーに相談してください。
github: moXXcha

---

**Last Updated**: 2025 年 1 月

**Maintained by**: SLEEVE Development Team
