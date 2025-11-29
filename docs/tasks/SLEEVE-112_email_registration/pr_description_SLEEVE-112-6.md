## 概要

RegisterUser UseCase層の実装を行いました。Email/Passwordでのユーザー登録ビジネスロジックを実装し、Firebase AuthenticationとDB保存、JWT発行を統合しています。

## 変更内容

### 新規作成ファイル

- `app/usecase/user/register_user_usecase.go` - RegisterUser UseCase
- `app/usecase/user/register_user_usecase_test.go` - UseCaseテスト + モック

### 実装詳細

#### RegisterUserUseCase

以下の処理フローを実装:

1. **入力バリデーション**
   - メールアドレス形式の検証（ドメインエラーへのラップ）
   - パスワード強度の検証（ドメインエラーへのラップ）

2. **Firebase Authentication登録**
   - `FirebaseUserRepositoryInterface`経由でFirebaseにユーザー作成
   - Firebase UIDを取得

3. **ドメインモデル作成**
   - `models.NewUser()`でUserエンティティを生成

4. **自社DB保存**
   - `UserDAOInterface`経由でDBにユーザー情報を永続化

5. **JWT発行**
   - `JWTService`でアクセストークン・リフレッシュトークンを発行

6. **エラーハンドリング・ロールバック**
   - DB保存失敗時にFirebaseユーザーを削除
   - JWT発行失敗時にもFirebaseユーザーを削除

#### インターフェース定義

依存性逆転の原則に従い、以下のインターフェースを定義:
- `FirebaseUserRepositoryInterface` - Firebase連携
- `UserDAOInterface` - DB操作

#### RegisterUserResult

結果構造体として以下を返却:
- `User` - 登録されたユーザードメインモデル
- `AccessToken` - アクセストークン（15分有効期限）
- `RefreshToken` - リフレッシュトークン（7日有効期限）

## 依存関係

- Depends on: SLEEVE-112-3（User Domain層）, SLEEVE-112-4（User Repository層）, SLEEVE-112-5（JWT発行機能）
- Blocks: SLEEVE-112-7（RegisterUser GraphQL層の実装）

## テスト

- [x] 正しいEmail/Passwordで登録が成功すること
- [x] 不正なEmailでエラーが返されること
- [x] 弱いPasswordでエラーが返されること
- [x] Email重複時にエラーが返されること
- [x] Firebaseエラー時にエラーが返されること
- [x] DBエラー時にFirebaseユーザーがロールバックされること
- [x] 登録成功時にJWTが発行されること

```bash
# テスト実行結果
docker exec sleeve-be go test ./usecase/user/... -v
# 全7テストPASS
```

## 関連チケット

- Jira: SLEEVE-112-6

## チェックリスト

- [x] コーディング規約に準拠
- [x] 単体テストが通過
- [x] TDD（テスト駆動開発）の原則に従って実装
- [x] Linterエラーなし
