# SLEEVE-112-1: usersテーブルのマイグレーション

## 概要

Email認証（アカウント登録）機能の実装の第1ステップとして、usersテーブルのマイグレーションを実装しました。

## 変更内容

### entスキーマ定義
- `app/ent/schema/user.go` を作成
- User Entityの定義（id, firebase_uid, email, created_at, updated_at, deleted_at）
- firebase_uidとemailにユニークインデックスを設定
- 論理削除用のdeleted_atフィールドを追加

### entコード生成
- `task ent-gen` により、User EntityのCRUD操作コードを自動生成
- User用のクエリビルダーを生成
- entクライアントへのUser統合

### マイグレーションSQL生成
- `task sql-gen` により、usersテーブル作成のマイグレーションSQLを生成
- `app/migrations/20251116144252.sql` を作成
- firebase_uidとemailのユニークインデックスを作成
- deleted_atのインデックスを作成

### マイグレーション実行
- `task migrate-up` により、usersテーブルをデータベースに作成
- テーブルとインデックスが正常に作成されたことを確認

### スキーマドキュメント更新
- `docs/db_schema.md` を更新
- usersテーブルの定義を追加（dbdiagram.io形式）
- インデックス情報を記載
- 変更履歴を記録

## 依存関係

- Depends on: なし
- Blocks: SLEEVE-112-3（User Domain層）

## テスト

### マイグレーション実行確認
- [x] `task migrate-up` が成功すること
- [x] usersテーブルが作成されていること
- [x] インデックスが正しく設定されていること

## 関連チケット

- Jira: SLEEVE-112-1

## チェックリスト

- [x] コーディング規約に準拠
- [x] Linterエラーがない（`golangci-lint run --fix`で0 issues）
- [x] ドキュメントが更新されている（`docs/db_schema.md`）
- [x] DB変更を行ったため `docs/db_schema.md` を更新している
- [x] 全ての変更がリモートにプッシュされている
- [x] コミットメッセージが規約に従っている

## レビュー観点

- entスキーマ定義が正しいか
- usersテーブルの構造が要件を満たしているか
- インデックスが適切に設定されているか
- docs/db_schema.mdが正確に更新されているか
