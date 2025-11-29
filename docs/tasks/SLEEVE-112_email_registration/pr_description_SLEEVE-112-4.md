## 概要

User Repository層（Firebase + DB連携）の実装を行いました。Firebase AuthenticationとのユーザーCRUD機能と、自社DBへのユーザー情報永続化機能を実装しています。

## 変更内容

### 新規作成ファイル

- `app/repository/external/firebase/user_repository.go` - Firebase連携Repository
- `app/repository/external/firebase/user_repository_test.go` - Repositoryテスト
- `app/repository/external/firebase/mock_firebase_client.go` - テスト用モック
- `app/repository/internal/user_dao.go` - DB操作DAO
- `app/repository/internal/user_dao_test.go` - DAOテスト
- `app/repository/internal/mock_ent_client.go` - テスト用モック

### 実装詳細

#### Firebase User Repository

Firebase Authenticationとの連携機能:

- `CreateUser` - Firebase Authenticationにユーザーを作成し、Firebase UIDを返却
- `CheckEmailExists` - メールアドレスがFirebaseに登録済みかを確認
- `DeleteUser` - Firebase Authenticationからユーザーを削除
- Firebaseエラー（EMAIL_EXISTS, USER_NOT_FOUND等）のドメインエラーへの変換

#### User DAO

自社DBへのユーザー情報永続化機能:

- `Save` - ユーザーをDBに保存（Firebase UIDと紐付け）
- `FindByPublicID` - 公開ID（UUID）でユーザーを検索
- `FindByFirebaseUID` - Firebase UIDでユーザーを検索
- `FindByEmail` - メールアドレスでユーザーを検索
- `ExistsByEmail` - メールアドレスがDBに存在するかを確認
- DBエラー（UNIQUE制約違反、NOT FOUND等）のドメインエラーへの変換
- entのUserエンティティからドメインモデルへの変換処理

## 依存関係

- Depends on: SLEEVE-112-1（usersテーブル）, SLEEVE-112-2（Firebase SDK初期化）, SLEEVE-112-3（User Domain層）
- Blocks: SLEEVE-112-6（RegisterUser UseCase層の実装）

## テスト

- [x] Firebase Authenticationにユーザーを登録できること
- [x] 自社DBにユーザー情報を保存できること（Firebase UIDと紐付け）
- [x] Email重複チェックが正常に動作すること（Firebase）
- [x] Email重複チェックが正常に動作すること（DB）
- [x] Firebase Authenticationのエラーがドメインエラーに変換されること
- [x] DBのエラーがドメインエラーに変換されること
- [x] 各種検索（PublicID, FirebaseUID, Email）が正常に動作すること

```bash
# テスト実行結果
cd app
go test ./repository/... -v
# 全テストPASS
```

## 関連チケット

- Jira: SLEEVE-112-4

## チェックリスト

- [x] コーディング規約に準拠
- [x] 単体テストが通過
- [x] TDD（テスト駆動開発）の原則に従って実装
- [x] Linterエラーなし
