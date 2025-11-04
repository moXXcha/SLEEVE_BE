# データベース取り扱いマニュアル

このドキュメントは、AI エージェントがデータベースのスキーマ変更やテーブル作成を行う際の指針と手順を定めたものです。

## 関連ドキュメント

- **全体フロー**: `docs/general.md` - 開発タスクの全体フロー
- **実装手順**: `docs/task_procedure.md` - 詳細な実装手順とテンプレート
- **コーディング規約**: `docs/BE_coding_roles.md` - アーキテクチャとコーディング規約
- **Git 運用**: `docs/github.md` - Git Flow 運用とコミット規約
- **Linter/Formatter**: `docs/lint_format_manual.md` - コード品質管理

## 1. DB 作成の方針

- **単一責任の原則:** テーブルは単一の責務を持つように設計してください。
- **正規化:** 原則として第三正規形を目指しますが、パフォーマンスとのトレードオフを考慮し、意図的に非正規化する場合はその理由をドキュメントに残してください。
- **命名規則:** `docs/BE_coding_roles.md` に記載されている命名規則を厳守してください。
- **汎用性:** 将来的な拡張性を考慮し、過度に特定のユースケースに最適化しすぎないように設計してください。
- **認証, user 情報について:** 認証、ユーザーの基本情報は firebase で管理しています。
- **orm:** orm は ent を使用しています

## 2. Schema の位置

- **Entity Schema:** Go のコードベースのスキーマ定義は `ent/schema/` ディレクトリ配下に配置されます。実際のテーブル構造はこれらのファイルによって定義されます。
- **全体スキーマ図:** データベース全体の構造を示すスキーマ図は `docs/db_schema.md` にて、[dbdiagram.io](https://dbdiagram.io) の形式で管理されます。

## 3. 実装手順

データベースのスキーマ変更は、以下の手順を厳密に守って進めてください。

### 手順 1: 現状の把握と不足の特定

- まず `docs/db_schema.md` と `ent/schema/` の内容を確認し、現在のデータベース構成を完全に理解します。
- タスクの要件と照らし合わせ、既存のテーブルに不足しているカラムや、新たに必要なテーブルを特定します。

### 手順 2: 変更方針の策定とドキュメント化

- 追加または変更するテーブルの方針を固めます。
- 以下のテンプレートに従い、変更計画書を `docs/tasks/[JiraのID]_[タスク名]/table_edit_plan.md` というファイル名で作成します。

### 手順 3: 作業者によるレビュー

- 作成した変更計画書 (`table_edit_plan.md`) を作業者（人間）に提示し、レビューを依頼します。
- 承認が得られたら、次のステップに進みます。承認なしにスキーマ変更に着手してはいけません。

### 手順 4: スキーマの更新（必須）

**重要**: テーブルの作成、編集、削除を行った場合は、**必ず `docs/db_schema.md` を更新してください**。これは必須事項です。

#### 4-1. 全体スキーマ図の更新（必須）

承認された変更計画に基づき、`docs/db_schema.md` のスキーマ図を最新の状態に更新します。

**更新内容**:
- テーブル作成: 新しいテーブル定義を追加
- テーブル編集: カラムの追加・変更・削除を反映
- テーブル削除: 該当テーブルの定義を削除
- リレーション変更: 外部キー関係を更新

**形式**: [dbdiagram.io](https://dbdiagram.io) の形式で記述

**例**:
```
Table users {
  id integer [primary key]
  username varchar
  email varchar [unique]
  created_at timestamp
}

Table posts {
  id integer [primary key]
  user_id integer [ref: > users.id]
  title varchar
  content text
  created_at timestamp
}
```

#### 4-2. Entity Schema の更新

`ent/schema/` ディレクトリ内の該当するスキーマファイルを変更、または新規に作成します。

#### 4-3. 更新確認チェックリスト

スキーマ更新後、以下を確認してください：

- [ ] `docs/db_schema.md` が最新の状態に更新されている
- [ ] `ent/schema/` のスキーマファイルが更新されている
- [ ] 両者の内容が一致している（矛盾がない）
- [ ] リレーション（外部キー）が正しく定義されている
- [ ] インデックスが適切に定義されている

---

## 4. 変更計画書テンプレート (`table_edit_plan.md`)

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

---

---

## 5. スキーマ更新の重要性

### 5.1 なぜ `docs/db_schema.md` の更新が必須なのか

`docs/db_schema.md` は、データベース全体の構造を可視化し、以下の目的で使用されます：

1. **全体像の把握**: プロジェクト全体のテーブル構造とリレーションを一目で理解できる
2. **設計レビュー**: 新しいテーブル追加時に、既存構造との整合性を確認できる
3. **ドキュメント**: 新しいメンバーがプロジェクトに参加した際の重要な資料
4. **AI エージェントの参照**: AI エージェントがタスク実施時に現状を把握するための情報源

### 5.2 更新を忘れた場合の問題

`docs/db_schema.md` の更新を忘れると、以下の問題が発生します：

- **設計ミス**: 既存のテーブル構造を把握できず、重複や矛盾が発生
- **バグの原因**: 最新の構造を知らずに実装し、不整合なデータが混入
- **レビュー困難**: レビュアーが変更の影響範囲を正確に把握できない
- **AI エージェントの誤動作**: 古い情報を元に実装し、エラーが発生

### 5.3 更新タイミング

以下のタイミングで `docs/db_schema.md` を更新してください：

1. **テーブル作成時**: 新しいテーブル定義を追加
2. **カラム追加/変更/削除時**: 該当テーブルの定義を更新
3. **テーブル削除時**: 該当テーブルの定義を削除
4. **リレーション変更時**: 外部キー関係を更新
5. **インデックス追加/削除時**: インデックス定義を更新

**重要**: 変更を commit する前に、必ず `docs/db_schema.md` を更新してください。

---

## 6. トランザクション管理

このセクションでは、entを使用したトランザクション管理の実装方法を説明します。

### 6.1 基本方針

- **トランザクションの開始・コミット・ロールバックは必ずUseCase層で行う**
- **Repository層はトランザクションオブジェクト（`*ent.Tx`）を受け取り、操作を実行する**
- **Repository層でトランザクションを開始してはいけない**

### 6.2 共通ヘルパー関数

トランザクションの定型処理（開始・コミット・ロールバック・panic対応）を共通化するため、`usecase/utils/transaction.go`にヘルパー関数を配置します。

#### ヘルパー関数の実装例

```go
// usecase/utils/transaction.go
package utils

import (
    "context"
    "fmt"

    "your-project/ent"
)

// with_transaction は、トランザクション処理を共通化するヘルパー関数です。
// UseCase層でトランザクションが必要な処理を実行する際に使用します。
//
// 使用例:
//   err := with_transaction(ctx, client, func(tx *ent.Tx) error {
//       // トランザクション内で実行したい処理
//       return nil
//   })
func with_transaction(ctx context.Context, client *ent.Client, fn func(*ent.Tx) error) error {
    var tx *ent.Tx
    var err error

    // トランザクション開始
    tx, err = client.Tx(ctx)
    if err != nil {
        return fmt.Errorf("failed to start transaction: %w", err)
    }

    // panic対応：panicが発生した場合も確実にロールバック
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()

    // 関数実行
    err = fn(tx)
    if err != nil {
        // エラー発生時はロールバック
        tx.Rollback()
        return err
    }

    // 成功時はコミット
    err = tx.Commit()
    if err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}
```

### 6.3 UseCase層での実装例

```go
// usecase/feature/user/create_user_usecase.go
package user

import (
    "context"
    "fmt"

    "your-project/ent"
    "your-project/usecase/utils"
)

type UserUsecase struct {
    client      *ent.Client
    userRepo    UserRepository
    profileRepo ProfileRepository
}

func (uc *UserUsecase) create_user(ctx context.Context, name string, email string) (*ent.User, error) {
    var user *ent.User
    var err error

    // ヘルパー関数を使用してトランザクション処理を実行
    err = utils.with_transaction(ctx, uc.client, func(tx *ent.Tx) error {
        var create_err error

        // Repository層を呼び出し（txを渡す）
        user, create_err = uc.userRepo.Create(ctx, tx, name, email)
        if create_err != nil {
            return fmt.Errorf("failed to create user: %w", create_err)
        }

        // 複数の操作を1つのトランザクション内で実行
        _, create_err = uc.profileRepo.CreateDefault(ctx, tx, user.ID)
        if create_err != nil {
            return fmt.Errorf("failed to create profile: %w", create_err)
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    return user, nil
}
```

### 6.4 Repository層での実装例

Repository層ではトランザクションオブジェクト（`*ent.Tx`）を受け取り、操作を実行します。

```go
// repository/user_repository.go
package repository

import (
    "context"
    "fmt"

    "your-project/ent"
)

type UserRepository interface {
    Create(ctx context.Context, tx *ent.Tx, name string, email string) (*ent.User, error)
    FindByID(ctx context.Context, client *ent.Client, id string) (*ent.User, error)
}

type userRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
    return &userRepositoryImpl{}
}

// Create はトランザクション内でユーザーを作成します
func (r *userRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, name string, email string) (*ent.User, error) {
    var user *ent.User
    var err error

    // txを使用してトランザクション内で操作
    user, err = tx.User.Create().
        SetName(name).
        SetEmail(email).
        Save(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}

// FindByID は通常のクエリでユーザーを取得します（トランザクション不要）
func (r *userRepositoryImpl) FindByID(ctx context.Context, client *ent.Client, id string) (*ent.User, error) {
    var user *ent.User
    var err error

    // 読み取り専用操作なのでclientを直接使用
    user, err = client.User.Get(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }

    return user, nil
}
```

### 6.5 実装時の注意点

1. **defer + recover**: panicが発生してもロールバックされるように必ず実装する
2. **エラー時は必ずロールバック**: エラーが発生したら即座にロールバックする
3. **Repository層でトランザクション開始禁止**: Repository層では`client.Tx()`を呼び出さない
4. **読み取り専用操作**: トランザクションが不要な単純な読み取り操作では`*ent.Client`を直接使用する
5. **トランザクションの粒度**: 必要最小限の範囲でトランザクションを使用し、長時間のロックを避ける

### 6.6 トランザクションが必要なケース

以下のような場合にトランザクションを使用します：

- 複数のテーブルへの書き込みが関連している場合
- データの整合性を保証する必要がある場合
- ロールバックが必要な可能性がある処理

### 6.7 トランザクションが不要なケース

以下のような場合はトランザクションを使用せず、`*ent.Client`を直接使用します：

- 単一のテーブルへの単純な読み取り
- 単一のテーブルへの単純な書き込み（他のテーブルへの影響がない場合）
- データの整合性が問題にならない場合
