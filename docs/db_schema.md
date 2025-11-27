# データベーススキーマ図

本ドキュメントは、SLEEVE プロジェクトのデータベース全体の構造を [dbdiagram.io](https://dbdiagram.io) 形式で管理します。

## 📌 重要

**このファイルはDB変更時に必ず更新してください！**

詳細な手順については、`docs/DB_manual.md` を参照してください。

---

## スキーマ定義

以下のスキーマは、[dbdiagram.io](https://dbdiagram.io) にコピーして可視化できます。

### 現在のテーブル

```
Table users {
  id int [primary key, increment, note: '内部ID（auto increment、外部には公開しない）']
  public_id uuid [not null, unique, note: '公開用ユーザーID（UUID、外部APIで使用）']
  firebase_uid varchar [not null, unique, note: 'Firebase Authentication UID']
  email varchar [not null, unique, note: 'メールアドレス']
  created_at timestamptz [not null, note: '作成日時']
  updated_at timestamptz [not null, note: '更新日時']
  deleted_at timestamptz [null, note: '削除日時（論理削除）']

  indexes {
    public_id [unique, name: 'user_public_id']
    firebase_uid [unique, name: 'user_firebase_uid']
    email [unique, name: 'user_email']
    deleted_at [name: 'user_deleted_at']
  }
}
```

### ID設計方針

外部から参照されるテーブル（URLで直接アクセスされるもの、外部APIのレスポンスとして使用されるものなど）では、以下のID設計を採用します：

| カラム名 | 型 | 用途 | 公開 |
|---------|---|------|-----|
| id | int (auto increment) | 内部での参照・外部キー結合に使用 | 非公開 |
| public_id | uuid | 外部API・URLで使用（データ作成時にDBで自動生成） | 公開 |

**理由:**
- **セキュリティ**: auto incrementのIDは連番のため、ユーザー数やデータ量が推測されやすい。UUIDを公開することでこれを防ぐ
- **パフォーマンス**: 内部結合にはintのIDを使用することで、JOINの効率を維持
- **拡張性**: 将来的にシャーディングが必要になった場合、UUIDの方が分散に適している

---

## 変更履歴

| 日付 | 作成者 | 変更内容 | 関連Jira |
|------|--------|---------|---------|
| 2025-01-28 | Claude | usersテーブルのID設計を変更（id: uuid -> int auto increment, public_id: uuid追加） | SLEEVE-112 |
| 2025-01-16 | Claude | usersテーブルの作成 | SLEEVE-112-1 |

---

## 参考

- [dbdiagram.io ドキュメント](https://dbdiagram.io/docs)
- [dbdiagram.io エディタ](https://dbdiagram.io/d)
