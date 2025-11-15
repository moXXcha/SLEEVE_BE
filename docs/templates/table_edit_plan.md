```markdown
# 変更の名前

（例: ユーザープロフィールテーブルの追加）

# 背景

（なぜこの変更が必要なのか、Jira チケットの要件などを基に記述）

# 変更内容

（どのテーブルに何を追加/変更するのか、具体的に記述）

- **テーブル名:** `users`
- **追加カラム:** `profile_image_url (string, nullable)`
- **変更カラム:** `age (int)` -> `birthday (datetime)`

# 影響範囲

（このスキーマ変更によって影響を受ける既存の機能や API、テーブルなどを記述）

# 最終 schema (dbdiagram.io 形式)

（変更後の **テーブル単体** または **関連するテーブル群** のスキーマを dbdiagram.io 形式で記述）

Table users {
id integer [primary key]
username varchar
email varchar
// ... (既存のカラム)
profile_image_url varchar
birthday datetime
}
```
