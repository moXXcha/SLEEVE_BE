# コーディングスタイルガイド

## 概要

このドキュメントは、SLEEVE プロジェクトにおけるコーディングスタイルの統一ルールを定義します。
可読性と保守性を高めるため、全ての開発者はこのガイドに従ってコードを記述してください。

## 関連ドキュメント

- **コーディング規約**: `docs/BE_coding_roles.md` - アーキテクチャと命名規則
- **全体フロー**: `docs/general.md` - 開発タスクの全体フロー
- **実装手順**: `docs/task_procedure.md` - 詳細な実装手順

---

## 1. 変数の定義

### 1.1 変数定義の位置と初期化

- 変数は**その変数が使われるスコープの最上位**で定義する
- **定義と初期化は必ず別で行う**
- **変数定義ブロックの後に改行を1つ入れて、その後に初期化を行う**
- **初期化と実装の間には改行を入れない**
- **例外：for文のrange句では`:=`の使用を許可する**（例：`for _, user := range users`）

#### 構造

```
1. 変数の宣言
2. 改行（1つ）
3. 初期化
4. 実装（改行なしで続く）
```

#### スコープの判断基準

1. **関数全体で使う変数** → 関数の最上位で宣言
2. **if文の中でしか使わない変数** → そのif文のブロックの最上位で宣言
3. **forループの中でしか使わない変数** → そのforループのブロックの最上位で宣言
4. **switch文の中でしか使わない変数** → そのswitch文のブロックの最上位で宣言

#### 良い例（スコープごとの宣言）

```go
func get_user_by_id(user_id string) (*User, error) {
    var user *User
    var err error
    var query string

    query = "SELECT * FROM users WHERE id = ?"
    user, err = db.QueryRow(query, user_id)
    if err != nil {
        return nil, err
    }
    // この変数はif文の中でしか使わないので、if文のスコープで宣言
    if user.Age < 18 {
        var minor_message string

        minor_message = "未成年のユーザーです"
        log.Info(minor_message)
    }
    return user, nil
}
```

#### 良い例（forループ内のスコープ）

```go
func process_users(users []*User) error {
    var err error

    for _, user := range users {
        var validation_err error
        var message string

        validation_err = validate_user(user)
        if validation_err != nil {
            return validation_err
        }
        message = fmt.Sprintf("Processed user: %s", user.Name)
        log.Info(message)
    }
    return nil
}
```

#### 悪い例

```go
func get_user_by_id(user_id string) (*User, error) {
    // 悪い例1: 定義と初期化が同時に行われている（:= を使用）
    // 注：for文のrange句以外では := を使用しない
    query := "SELECT * FROM users WHERE id = ?"

    // 悪い例2: 変数が途中で定義されている
    user, err := db.QueryRow(query, user_id)
    if err != nil {
        return nil, err
    }
    // 悪い例3: 初期化と実装の間に改行がある

    // 悪い例4: if文の中でしか使わないのに関数の最上位で宣言している
    var minor_message string
    if user.Age < 18 {
        minor_message = "未成年のユーザーです"
        log.Info(minor_message)
    }
    return user, nil
}
```

### 1.2 変数のグループ化

関連する変数はまとめて定義し、用途ごとに空行で区切らない（そのスコープ内で一括定義）。

```go
func process_order(order_id string) error {
    var order *Order
    var user *User
    var err error

    order, err = get_order(order_id)
    if err != nil {
        return err
    }
    user, err = get_user(order.UserID)
    if err != nil {
        return err
    }
    // このスコープでしか使わない変数
    if order.Status == "pending" {
        var payment *Payment
        var total_amount int
        var is_valid bool

        payment, err = get_payment(order.PaymentID)
        if err != nil {
            return err
        }
        total_amount = calculate_total(order)
        is_valid = validate_payment(payment, total_amount)
        if !is_valid {
            return errors.New("invalid payment")
        }
    }
    return nil
}
```

---

## 2. 改行ルール

### 2.1 必須の改行

以下の箇所では**改行を1つ**入れる：

1. **import文とpackage文の間**
2. **import文とコードの間**
3. **変数定義ブロックと初期化の間**
4. **関数と関数の間**
5. **構造体定義と次の要素の間**
6. **型定義と次の要素の間**

#### 例

```go
package usecase

import (
    "context"
    "errors"
)

type UserUsecase struct {
    repo UserRepository
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
    return &UserUsecase{
        repo: repo,
    }
}

func (u *UserUsecase) get_user(ctx context.Context, user_id string) (*User, error) {
    var user *User
    var err error

    user, err = u.repo.FindByID(ctx, user_id)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (u *UserUsecase) create_user(ctx context.Context, name string) (*User, error) {
    var user *User
    var err error

    user = &User{Name: name}
    err = u.repo.Save(ctx, user)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

### 2.2 禁止の改行

以下の箇所では**無駄な改行を入れない**：

1. **変数定義ブロック内**
2. **初期化と実装の間**（重要！）
3. **if文とその実装の間**
4. **for文とその実装の間**
5. **switch文とcaseの間**
6. **関数の戻り値とreturn文の間**
7. **連続する処理の間**（論理的なまとまりがある場合を除く）

#### 良い例

```go
func validate_user(user *User) error {
    var err error

    if user == nil {
        return errors.New("user is nil")
    }
    if user.Name == "" {
        return errors.New("name is empty")
    }
    if user.Email == "" {
        return errors.New("email is empty")
    }
    return nil
}
```

#### 悪い例

```go
func validate_user(user *User) error {
    var err error

    if user == nil {

        return errors.New("user is nil")
    }

    if user.Name == "" {

        return errors.New("name is empty")
    }

    if user.Email == "" {

        return errors.New("email is empty")
    }

    return nil
}
```

---

## 3. インデントとフォーマット

### 3.1 インデント

- **タブ文字**を使用（スペースは使用しない）
- ネストは1段階につき1タブ

### 3.2 行の長さ

- 1行は**120文字以内**を推奨
- 超える場合は適切に改行し、インデントを1段階深くする

```go
func create_user_with_profile(
    ctx context.Context,
    name string,
    email string,
    age int,
) (*User, error) {
    // 実装
}
```

---

## 4. 関数の記述スタイル

### 4.1 関数の構造

```go
func 関数名(引数) (戻り値) {
    // 1. 変数宣言（スコープ最上位）
    var 変数1 型
    var 変数2 型

    // 2. 初期化
    変数1 = 初期値1
    変数2 = 初期値2
    // 3. 実装部分（初期化と実装の間に改行なし）
    // ロジック
    return 結果, nil
}
```

### 4.2 早期リターン

エラーチェックは早期リターンを使用し、ネストを深くしない。

#### 良い例

```go
func get_user_by_id(user_id string) (*User, error) {
    var user *User
    var err error

    if user_id == "" {
        return nil, errors.New("user_id is empty")
    }
    user, err = db.FindByID(user_id)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

#### 悪い例

```go
func get_user_by_id(user_id string) (*User, error) {
    var user *User
    var err error

    if user_id != "" {
        user, err = db.FindByID(user_id)
        if err == nil {
            return user, nil
        } else {
            return nil, err
        }
    } else {
        return nil, errors.New("user_id is empty")
    }
}
```

---

## 5. エラーハンドリング

### 5.1 エラーチェック

- エラーは即座にチェックする
- エラーメッセージは `docs/messages/` で定義したものを使用

```go
func save_user(user *User) error {
    var err error

    err = validate_user(user)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    err = db.Save(user)
    if err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }
    return nil
}
```

### 5.2 エラーのラップ

- エラーは `fmt.Errorf` と `%w` を使ってラップする
- コンテキスト情報を追加する

---

## 6. 構造体の定義

### 6.1 構造体の記述

```go
type User struct {
    ID        string
    Name      string
    Email     string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}
```

### 6.2 構造体のフィールド

- フィールドは型ごとにグループ化しない
- 関連するフィールドは近くに配置
- タグは見やすく整列させる

```go
type User struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Email     string    `json:"email" db:"email"`
    Age       int       `json:"age" db:"age"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

---

## 7. 定数の定義

### 7.1 定数のグループ化

- 関連する定数は `const` ブロックでまとめる
- 定数名は大文字スネークケースではなく、キャメルケース（Go標準）

```go
const (
    MaxRetryCount = 3
    DefaultTimeout = 30 * time.Second
    BufferSize = 1024
)

const (
    StatusActive = "active"
    StatusInactive = "inactive"
    StatusPending = "pending"
)
```

---

## 8. インポート文の整理

### 8.1 インポートの順序

1. 標準ライブラリ
2. サードパーティライブラリ
3. プロジェクト内パッケージ

各グループ間に空行を入れる。

```go
import (
    "context"
    "errors"
    "fmt"

    "github.com/lib/pq"
    "github.com/pkg/errors"

    "project/app/domain/models"
    "project/app/repository"
)
```

---

## 9. コメントの記述

### 9.1 コメントスタイル

- 関数の前に説明コメントを記述（Go標準）
- コメントは日本語で記述
- コメントは完全な文章で記述する

```go
// get_user_by_id は指定されたIDのユーザーを取得します。
// ユーザーが存在しない場合はエラーを返します。
func get_user_by_id(user_id string) (*User, error) {
    var user *User
    var err error

    // バリデーション
    if user_id == "" {
        return nil, errors.New("user_id is empty")
    }

    user, err = db.FindByID(user_id)
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

### 9.2 TODO コメント

```go
// TODO: キャッシュ機能の実装
// FIXME: エラーハンドリングの改善が必要
// NOTE: この処理はパフォーマンスに影響する
```

---

## 10. 制御構文

### 10.1 if文

- 条件式は括弧で囲まない（Go標準）
- else は避け、早期リターンを使用

```go
func is_adult(age int) bool {
    if age < 0 {
        return false
    }
    if age < 18 {
        return false
    }
    return true
}
```

### 10.2 switch文

```go
func get_status_message(status string) string {
    var message string

    switch status {
    case StatusActive:
        message = "アクティブです"
    case StatusInactive:
        message = "非アクティブです"
    case StatusPending:
        message = "保留中です"
    default:
        message = "不明なステータスです"
    }
    return message
}
```

### 10.3 for文

```go
func process_users(users []*User) error {
    var err error

    for _, user := range users {
        err = validate_user(user)
        if err != nil {
            return err
        }
    }
    return nil
}
```

---

## 11. テストコードのスタイル

### 11.1 テスト関数の命名

```go
func TestGetUserByID_Success(t *testing.T) {
    // テストコード
}

func TestGetUserByID_UserNotFound(t *testing.T) {
    // テストコード
}

func TestGetUserByID_InvalidID(t *testing.T) {
    // テストコード
}
```

### 11.2 テストの構造

```go
func TestGetUserByID_Success(t *testing.T) {
    // Arrange（準備）
    var user_id string
    var expected_user *User
    var actual_user *User
    var err error

    user_id = "test_user_id"
    expected_user = &User{ID: user_id, Name: "Test User"}
    // Act（実行）
    actual_user, err = get_user_by_id(user_id)
    // Assert（検証）
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if actual_user.ID != expected_user.ID {
        t.Errorf("expected ID %s, got %s", expected_user.ID, actual_user.ID)
    }
}
```

---

## 12. その他のベストプラクティス

### 12.1 nil チェック

```go
func process_user(user *User) error {
    if user == nil {
        return errors.New("user is nil")
    }
    // 処理
    return nil
}
```

### 12.2 空スライス・マップの初期化

```go
// 良い例
func get_all_users() []*User {
    var users []*User

    users = make([]*User, 0)
    // 処理
    return users
}

func create_user_map() map[string]*User {
    var user_map map[string]*User

    user_map = make(map[string]*User)
    // 処理
    return user_map
}

// 悪い例（定義と初期化が同時）
func get_all_users() []*User {
    users := make([]*User, 0)  // := を使っている
    return users
}

// 悪い例（初期化していない）
func get_all_users() []*User {
    var users []*User  // nil スライスのまま
    return users
}
```

### 12.3 defer の使用

```go
func read_file(file_path string) (string, error) {
    var file *os.File
    var content []byte
    var err error

    file, err = os.Open(file_path)
    if err != nil {
        return "", err
    }
    defer file.Close()
    content, err = io.ReadAll(file)
    if err != nil {
        return "", err
    }
    return string(content), nil
}
```

---

## 13. チェックリスト

コードレビュー時に以下を確認してください：

- [ ] 変数はその変数が使われるスコープの最上位で定義されているか
- [ ] 定義と初期化は分離されているか（例外：for文のrange句は`:=`使用可）
- [ ] 関数全体で使わない変数を関数の最上位で宣言していないか
- [ ] 特定のスコープ（if、for、switch）でしか使わない変数はそのスコープ内で宣言されているか
- [ ] 変数定義と初期化の間に改行が1つあるか
- [ ] 初期化と実装の間に改行がないか（重要！）
- [ ] 関数間に改行が1つあるか
- [ ] import文の後に改行が1つあるか
- [ ] 無駄な改行がないか
- [ ] 早期リターンが使用されているか
- [ ] エラーハンドリングが適切か
- [ ] コメントが適切に記述されているか
- [ ] 命名規則に従っているか（`docs/BE_coding_roles.md` 参照）

---

## 14. Linter とフォーマッター

コーディングスタイルの遵守は以下のツールで自動チェックされます：

```bash
# app/ ディレクトリで実行
golangci-lint run --fix
```

詳細は `docs/lint_format_manual.md` を参照してください。
