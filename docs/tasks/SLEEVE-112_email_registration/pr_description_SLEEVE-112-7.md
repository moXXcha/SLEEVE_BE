## 概要

RegisterUser GraphQL層の実装を行いました。registerUserミューテーションを追加し、Email/Passwordでのユーザー登録APIを公開しています。

## 変更内容

### 変更ファイル

- `app/graph/schema.graphqls` - GraphQLスキーマ
- `app/graph/resolver.go` - Resolver依存注入
- `app/graph/schema.resolvers.go` - Resolver実装
- `app/graph/generated.go` - gqlgen生成コード
- `app/graph/model/models_gen.go` - gqlgen生成モデル

### 新規作成ファイル

- `app/graph/schema.resolvers_test.go` - Resolverテスト

### 実装詳細

#### GraphQLスキーマ

```graphql
# 入力型
input RegisterUserInput {
  email: String!
  password: String!
}

# 出力型
type AuthTokens {
  accessToken: String!
  refreshToken: String!
}

type RegisteredUser {
  id: ID!
  email: String!
}

type RegisterUserPayload {
  user: RegisteredUser!
  tokens: AuthTokens!
}

# ミューテーション
type Mutation {
  registerUser(input: RegisterUserInput!): RegisterUserPayload!
}
```

#### Resolver実装

- `RegisterUserUseCase`を呼び出し
- UseCaseの結果をGraphQLレスポンス型に変換
- エラー時はそのまま返却（GraphQLエラーとして処理）

#### 依存注入

`Resolver`構造体に`RegisterUserUseCase`を追加し、外部から注入可能に

## 依存関係

- Depends on: SLEEVE-112-6（RegisterUser UseCase層）
- Blocks: なし（SLEEVE-112の最終Subtask）

## テスト

- [x] registerUserミューテーションが正常に動作すること
- [x] 不正なEmail形式でエラーが返されること
- [x] 弱いPasswordでエラーが返されること
- [x] Email重複時にエラーが返されること

```bash
# テスト実行結果
docker exec sleeve-be go test ./graph/... -v
# 全4テストPASS
```

## 使用例

```graphql
mutation {
  registerUser(input: {
    email: "user@example.com"
    password: "Password123!"
  }) {
    user {
      id
      email
    }
    tokens {
      accessToken
      refreshToken
    }
  }
}
```

## 関連チケット

- Jira: SLEEVE-112-7

## チェックリスト

- [x] コーディング規約に準拠
- [x] 単体テストが通過
- [x] TDD（テスト駆動開発）の原則に従って実装
- [x] Linterエラーなし
