# Users テーブル作成計画書

## 変更内容

- **種類**: 新規テーブル作成
- **テーブル名**: users
- **対象Subtask**: SLEEVE-101

## テーブル定義

### カラム一覧

| カラム名 | 型 | NULL許可 | デフォルト値 | 制約 | 説明 |
|---------|---|---------|------------|-----|-----|
| id | varchar | NO | - | PRIMARY KEY | 自社独自のユーザーID（UUID） |
| firebase_uid | varchar | NO | - | UNIQUE | Firebase Authentication UID |
| email | varchar | YES | NULL | - | メールアドレス（任意） |
| created_at | timestamp | NO | now() | - | 作成日時 |
| updated_at | timestamp | NO | now() | - | 更新日時 |
| deleted_at | timestamp | YES | NULL | - | 論理削除日時 |

### インデックス

| インデックス名 | 種類 | カラム | 目的 |
|-------------|------|-------|------|
| users_pkey | PRIMARY KEY | id | 主キー |
| user_firebase_uid | UNIQUE INDEX | firebase_uid | Firebase UIDでの高速検索と一意性保証 |
| user_email | INDEX | email | Emailでの検索性能向上 |

### リレーション

現時点ではリレーションなし。今後、以下のテーブルとのリレーションが予想されます：
- `profiles` テーブル（ユーザープロフィール）
- `items` テーブル（出品アイテム）
- `coordinates` テーブル（コーディネート投稿）

---

## 変更理由

### ビジネス要件
- Firebase Authenticationで管理するユーザー情報と、自社DBのユーザー情報を紐付けるため
- Firebase UIDを保存し、自社独自のユーザーIDも発行することで、将来的にFirebaseからの移行も可能にする
- Email認証機能の基盤となるテーブル

### 技術的理由
- 認証技術選定ドキュメント（`docs/prod_docs/技術選定/認証技術選定.md`）に従い、Firebase UIDと自社DBを連携
- Firebase Authenticationは認証のみを担当し、ユーザーに関連するビジネスデータは自社DBで管理
- 論理削除（soft delete）を採用し、ユーザーデータの削除履歴を保持

---

## 影響範囲

### 既存機能への影響
- **影響なし**: 新規テーブルのため、既存機能への影響はありません

### 今後の機能への影響
- **ログイン機能**: このusersテーブルを使用してユーザー認証を行う
- **プロフィール機能**: usersテーブルと1対1のリレーションでprofilesテーブルを作成予定
- **コンテンツ投稿機能**: usersテーブルのidを外部キーとして参照
- **フォロー機能**: usersテーブルを自己結合して実装予定

---

## 実装手順

### 手順 1: ent/schema/user.go の作成

```go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			Comment("自社独自のユーザーID"),
		field.String("firebase_uid").
			Unique().
			Immutable().
			Comment("Firebase Authentication UID"),
		field.String("email").
			Optional().
			Comment("メールアドレス（任意）"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("作成日時"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新日時"),
		field.Time("deleted_at").
			Optional().
			Nillable().
			Comment("論理削除日時"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("firebase_uid").Unique(),
		index.Fields("email"),
	}
}
```

### 手順 2: entコード生成

```bash
cd app
task ent-gen
cd ..
```

生成されるファイル:
- `app/ent/user.go`
- `app/ent/user_create.go`
- `app/ent/user_update.go`
- `app/ent/user_delete.go`
- `app/ent/user_query.go`
- その他entが自動生成するファイル

### 手順 3: migration SQL生成

```bash
cd app
task sql-gen
cd ..
```

生成されるファイル:
- `app/migrations/YYYYMMDDHHMMSS_create_users.sql`

期待されるSQL内容:
```sql
-- create "users" table
CREATE TABLE "users" (
  "id" character varying NOT NULL,
  "firebase_uid" character varying NOT NULL,
  "email" character varying NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "user_firebase_uid" to table: "users"
CREATE UNIQUE INDEX "user_firebase_uid" ON "users" ("firebase_uid");
-- create index "user_email" to table: "users"
CREATE INDEX "user_email" ON "users" ("email");
```

### 手順 4: docs/db_schema.md の更新（必須）

`docs/db_schema.md` に以下を追加:

```markdown
Table users {
  id varchar [primary key, note: '自社独自のユーザーID（UUID）']
  firebase_uid varchar [unique, not null, note: 'Firebase Authentication UID']
  email varchar [note: 'メールアドレス（任意）']
  created_at timestamp [not null, default: `now()`, note: '作成日時']
  updated_at timestamp [not null, default: `now()`, note: '更新日時']
  deleted_at timestamp [note: '論理削除日時（論理削除）']

  indexes {
    firebase_uid [unique, name: 'user_firebase_uid']
    email [name: 'user_email']
  }
}
```

更新履歴に追加:
```markdown
| 更新日 | 作業者 | 変更内容 | 関連Jira |
|------|--------|---------|---------|
| 2025-01-16 | Claude | usersテーブルを追加 | SLEEVE-101 |
```

### 手順 5: マイグレーション実行（開発環境）

```bash
cd app
task migrate-up
cd ..
```

### 手順 6: マイグレーション確認

PostgreSQLに接続してテーブルを確認:
```bash
psql -U youruser -d yourdb -c "\d users"
```

期待される出力:
```
                Table "public.users"
    Column     |           Type           | Nullable
---------------+--------------------------+----------
 id            | character varying        | not null
 firebase_uid  | character varying        | not null
 email         | character varying        |
 created_at    | timestamp with time zone | not null
 updated_at    | timestamp with time zone | not null
 deleted_at    | timestamp with time zone |
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
    "user_firebase_uid" UNIQUE, btree (firebase_uid)
    "user_email" btree (email)
```

---

## ロールバック手順

マイグレーション実行後に問題が発生した場合:

```bash
cd app
task migrate-rollback
cd ..
```

**注意**: 最新の1件のみロールバック可能です。複数件戻す場合は、`ent/schema/user.go`を編集して新しいマイグレーションを生成してください。

---

## セキュリティ考慮事項

### データ保護
- **Firebase UID**: 不変（Immutable）として定義し、変更を防止
- **Email**: 個人情報のため、適切なアクセス制御が必要
- **論理削除**: 物理削除ではなく論理削除を採用し、削除履歴を保持

### インデックス設計
- **firebase_uid**: ユニークインデックスを設定し、重複を防止
- **email**: 検索性能向上のためインデックスを設定（ユニークではない：Emailなしでも登録可能なため）

---

## パフォーマンス考慮事項

### インデックス戦略
- `firebase_uid`にユニークインデックス: Firebase UIDでの高速検索
- `email`にインデックス: Emailでの検索性能向上

### 将来的な最適化
- ユーザー数が増加した場合、以下を検討:
  - パーティショニング（created_atベース）
  - レプリケーション（読み取り専用レプリカ）
  - キャッシュ層の追加（Redis）

---

## テスト計画

### マイグレーションテスト
- [ ] マイグレーション実行が成功すること
- [ ] usersテーブルが作成されていること
- [ ] 全てのカラムが正しく定義されていること
- [ ] インデックスが正しく設定されていること
- [ ] ロールバックが正常に動作すること

### データ整合性テスト
- [ ] firebase_uidにユニーク制約が機能すること
- [ ] emailがNULL許可であること
- [ ] created_atにデフォルト値（now()）が設定されること
- [ ] 論理削除が正常に動作すること

---

## 承認

この変更計画書を作成したら、作業者に内容の確認を促してください。承認が得られたら、実装を開始してください。

**承認者**: （作業者が記入）
**承認日**: （作業者が記入）
**備考**: （作業者が記入）
