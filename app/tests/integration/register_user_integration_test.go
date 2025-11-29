package integration

import (
	"context"
	"testing"

	domain_errors "sleeve/domain/errors"
	"sleeve/domain/models"
	"sleeve/graph"
	"sleeve/graph/model"
	"sleeve/usecase/user"
	"sleeve/usecase/utils"
)

// テスト用定数
const (
	testSecretKey = "test_secret_key_for_integration_testing_1234567890"
)

// 結合テスト用のセットアップ
func setupIntegrationTest() (*graph.Resolver, *MockFirebaseUserRepository, *MockUserDAO, *utils.JWTService) {
	var firebase_repo *MockFirebaseUserRepository
	var user_dao *MockUserDAO
	var jwt_service *utils.JWTService
	var use_case *user.RegisterUserUseCase
	var resolver *graph.Resolver

	firebase_repo = NewMockFirebaseUserRepository()
	user_dao = NewMockUserDAO()
	jwt_service = utils.NewJWTService(testSecretKey)
	use_case = user.NewRegisterUserUseCase(firebase_repo, user_dao, jwt_service)
	resolver = &graph.Resolver{
		RegisterUserUseCase: use_case,
	}
	return resolver, firebase_repo, user_dao, jwt_service
}

// TestIntegration_RegisterUser_NormalFlow は正常なユーザー登録フローをテストします
// 通過条件:
// - 正しいEmail/Passwordを渡した場合、Firebase Authenticationに新規ユーザーが登録される
// - 自社DBにもユーザー情報が保存される（Firebase UIDと紐付け）
// - JWT（アクセストークン + リフレッシュトークン）が返される
// - User情報が返される
func TestIntegration_RegisterUser_NormalFlow(t *testing.T) {
	var ctx context.Context
	var resolver *graph.Resolver
	var firebase_repo *MockFirebaseUserRepository
	var user_dao *MockUserDAO
	var mutation_resolver graph.MutationResolver
	var input model.RegisterUserInput
	var result *model.RegisterUserPayload
	var err error

	ctx = context.Background()
	resolver, firebase_repo, user_dao, _ = setupIntegrationTest()
	mutation_resolver = resolver.Mutation()
	input = model.RegisterUserInput{
		Email:    "integration_test@example.com",
		Password: "SecurePassword123!",
	}

	// ユーザー登録を実行
	result, err = mutation_resolver.RegisterUser(ctx, input)

	// エラーがないことを確認
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 結果がnilでないことを確認
	if result == nil {
		t.Fatal("expected result to be non-nil")
	}

	// User情報が返されることを確認
	verifyUserResult(t, result, input.Email)

	// JWTが返されることを確認
	verifyTokensExist(t, result)

	// Firebase Authenticationにユーザーが登録されたことを確認
	if !firebase_repo.CreateUserCalled {
		t.Error("expected Firebase CreateUser to be called")
	}

	// 自社DBにユーザー情報が保存されたことを確認
	if !user_dao.SaveCalled {
		t.Error("expected UserDAO Save to be called")
	}
}

// verifyUserResult はユーザー結果を検証するヘルパー関数です
func verifyUserResult(t *testing.T, result *model.RegisterUserPayload, expected_email string) {
	t.Helper()
	if result.User == nil {
		t.Fatal("expected user to be non-nil")
	}
	if result.User.ID == "" {
		t.Error("expected user ID to be non-empty")
	}
	if result.User.Email != expected_email {
		t.Errorf("expected email %s, got %s", expected_email, result.User.Email)
	}
}

// verifyTokensExist はトークンの存在を検証するヘルパー関数です
func verifyTokensExist(t *testing.T, result *model.RegisterUserPayload) {
	t.Helper()
	if result.Tokens == nil {
		t.Fatal("expected tokens to be non-nil")
	}
	if result.Tokens.AccessToken == "" {
		t.Error("expected access_token to be non-empty")
	}
	if result.Tokens.RefreshToken == "" {
		t.Error("expected refresh_token to be non-empty")
	}
}

// TestIntegration_RegisterUser_DuplicateEmail はEmail重複時のエラーハンドリングをテストします
// 通過条件:
// - 既に登録済みのEmailで登録を試みた場合、適切なエラーメッセージが返される
// - Firebase Authenticationにユーザーが登録されない
// - 自社DBにユーザー情報が保存されない
func TestIntegration_RegisterUser_DuplicateEmail(t *testing.T) {
	var ctx context.Context
	var resolver *graph.Resolver
	var firebase_repo *MockFirebaseUserRepository
	var user_dao *MockUserDAO
	var mutation_resolver graph.MutationResolver
	var input model.RegisterUserInput
	var err error

	ctx = context.Background()
	resolver, firebase_repo, user_dao, _ = setupIntegrationTest()
	firebase_repo.ShouldReturnDuplicateError = true
	mutation_resolver = resolver.Mutation()
	input = model.RegisterUserInput{
		Email:    "duplicate@example.com",
		Password: "SecurePassword123!",
	}

	// ユーザー登録を実行
	_, err = mutation_resolver.RegisterUser(ctx, input)

	// エラーが返されることを確認
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}

	// ドメインエラーであることを確認
	if !domain_errors.IsUserDomainError(err) {
		t.Errorf("expected user domain error, got %v", err)
	}

	// 自社DBにユーザー情報が保存されていないことを確認
	if user_dao.SaveCalled {
		t.Error("expected UserDAO Save NOT to be called for duplicate email")
	}
}

// TestIntegration_RegisterUser_InvalidEmail は不正なEmail形式のバリデーションをテストします
// 通過条件:
// - 不正な形式のEmailで登録を試みた場合、適切なエラーメッセージが返される
// - Firebase Authenticationにユーザーが登録されない
func TestIntegration_RegisterUser_InvalidEmail(t *testing.T) {
	var ctx context.Context
	var resolver *graph.Resolver
	var firebase_repo *MockFirebaseUserRepository
	var user_dao *MockUserDAO
	var mutation_resolver graph.MutationResolver
	var input model.RegisterUserInput
	var err error

	ctx = context.Background()
	resolver, firebase_repo, user_dao, _ = setupIntegrationTest()
	mutation_resolver = resolver.Mutation()
	input = model.RegisterUserInput{
		Email:    "invalid-email-format",
		Password: "SecurePassword123!",
	}

	// ユーザー登録を実行
	_, err = mutation_resolver.RegisterUser(ctx, input)

	// エラーが返されることを確認
	if err == nil {
		t.Fatal("expected error for invalid email, got nil")
	}

	// ドメインエラーであることを確認
	if !domain_errors.IsUserDomainError(err) {
		t.Errorf("expected user domain error, got %v", err)
	}

	// Firebase Authenticationにユーザーが登録されていないことを確認
	if firebase_repo.CreateUserCalled {
		t.Error("expected Firebase CreateUser NOT to be called for invalid email")
	}

	// 自社DBにユーザー情報が保存されていないことを確認
	if user_dao.SaveCalled {
		t.Error("expected UserDAO Save NOT to be called for invalid email")
	}
}

// TestIntegration_RegisterUser_WeakPassword は弱いPasswordのバリデーションをテストします
// 通過条件:
// - 強度要件を満たさないPasswordで登録を試みた場合、適切なエラーメッセージが返される
// - Firebase Authenticationにユーザーが登録されない
func TestIntegration_RegisterUser_WeakPassword(t *testing.T) {
	var ctx context.Context
	var resolver *graph.Resolver
	var firebase_repo *MockFirebaseUserRepository
	var user_dao *MockUserDAO
	var mutation_resolver graph.MutationResolver
	var input model.RegisterUserInput
	var err error

	ctx = context.Background()
	resolver, firebase_repo, user_dao, _ = setupIntegrationTest()
	mutation_resolver = resolver.Mutation()
	input = model.RegisterUserInput{
		Email:    "test@example.com",
		Password: "weak",
	}

	// ユーザー登録を実行
	_, err = mutation_resolver.RegisterUser(ctx, input)

	// エラーが返されることを確認
	if err == nil {
		t.Fatal("expected error for weak password, got nil")
	}

	// ドメインエラーであることを確認
	if !domain_errors.IsUserDomainError(err) {
		t.Errorf("expected user domain error, got %v", err)
	}

	// Firebase Authenticationにユーザーが登録されていないことを確認
	if firebase_repo.CreateUserCalled {
		t.Error("expected Firebase CreateUser NOT to be called for weak password")
	}

	// 自社DBにユーザー情報が保存されていないことを確認
	if user_dao.SaveCalled {
		t.Error("expected UserDAO Save NOT to be called for weak password")
	}
}

// TestIntegration_RegisterUser_JWTValidity はJWTの有効性を確認するテストです
// 通過条件:
// - 発行されたアクセストークンが正しく検証できること
// - 発行されたリフレッシュトークンが正しく検証できること
// - トークンにFirebase UIDが含まれていること
func TestIntegration_RegisterUser_JWTValidity(t *testing.T) {
	var ctx context.Context
	var resolver *graph.Resolver
	var jwt_service *utils.JWTService
	var mutation_resolver graph.MutationResolver
	var input model.RegisterUserInput
	var result *model.RegisterUserPayload
	var access_claims *utils.JWTClaims
	var refresh_claims *utils.JWTClaims
	var err error

	ctx = context.Background()
	resolver, _, _, jwt_service = setupIntegrationTest()
	mutation_resolver = resolver.Mutation()
	input = model.RegisterUserInput{
		Email:    "jwt_test@example.com",
		Password: "SecurePassword123!",
	}

	// ユーザー登録を実行
	result, err = mutation_resolver.RegisterUser(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// アクセストークンの検証
	access_claims, err = jwt_service.ValidateToken(result.Tokens.AccessToken)
	if err != nil {
		t.Fatalf("failed to validate access token: %v", err)
	}

	// アクセストークンのクレームを確認
	if access_claims.UserID == "" {
		t.Error("expected user_id in access token claims")
	}
	if access_claims.FirebaseUID == "" {
		t.Error("expected firebase_uid in access token claims")
	}
	if access_claims.TokenType != utils.TokenTypeAccess {
		t.Errorf("expected token type %s, got %s", utils.TokenTypeAccess, access_claims.TokenType)
	}

	// リフレッシュトークンの検証
	refresh_claims, err = jwt_service.ValidateToken(result.Tokens.RefreshToken)
	if err != nil {
		t.Fatalf("failed to validate refresh token: %v", err)
	}

	// リフレッシュトークンのクレームを確認
	if refresh_claims.UserID == "" {
		t.Error("expected user_id in refresh token claims")
	}
	if refresh_claims.FirebaseUID == "" {
		t.Error("expected firebase_uid in refresh token claims")
	}
	if refresh_claims.TokenType != utils.TokenTypeRefresh {
		t.Errorf("expected token type %s, got %s", utils.TokenTypeRefresh, refresh_claims.TokenType)
	}

	// アクセストークンとリフレッシュトークンのuser_idが一致することを確認
	if access_claims.UserID != refresh_claims.UserID {
		t.Errorf("user_id mismatch: access=%s, refresh=%s", access_claims.UserID, refresh_claims.UserID)
	}

	// アクセストークンとリフレッシュトークンのfirebase_uidが一致することを確認
	if access_claims.FirebaseUID != refresh_claims.FirebaseUID {
		t.Errorf("firebase_uid mismatch: access=%s, refresh=%s", access_claims.FirebaseUID, refresh_claims.FirebaseUID)
	}
}

// MockFirebaseUserRepository は結合テスト用のFirebaseリポジトリモックです
type MockFirebaseUserRepository struct {
	ShouldReturnDuplicateError bool
	ShouldReturnError          bool
	CreateUserCalled           bool
	DeleteUserCalled           bool
	last_firebase_uid          string
}

// NewMockFirebaseUserRepository は新しいMockFirebaseUserRepositoryを作成します
func NewMockFirebaseUserRepository() *MockFirebaseUserRepository {
	return &MockFirebaseUserRepository{
		ShouldReturnDuplicateError: false,
		ShouldReturnError:          false,
		CreateUserCalled:           false,
		DeleteUserCalled:           false,
	}
}

// CreateUser はモックのユーザー作成を行います
func (m *MockFirebaseUserRepository) CreateUser(_ context.Context, _ models.Email, _ models.Password) (string, error) {
	m.CreateUserCalled = true
	if m.ShouldReturnDuplicateError {
		return "", domain_errors.ErrDuplicateEmail
	}
	if m.ShouldReturnError {
		return "", domain_errors.ErrFirebaseAuthFailed
	}
	m.last_firebase_uid = "integration_test_firebase_uid_123"
	return m.last_firebase_uid, nil
}

// DeleteUser はモックのユーザー削除を行います
func (m *MockFirebaseUserRepository) DeleteUser(_ context.Context, _ string) error {
	m.DeleteUserCalled = true
	return nil
}

// MockUserDAO は結合テスト用のUserDAOモックです
type MockUserDAO struct {
	ShouldReturnError bool
	SaveCalled        bool
}

// NewMockUserDAO は新しいMockUserDAOを作成します
func NewMockUserDAO() *MockUserDAO {
	return &MockUserDAO{
		ShouldReturnError: false,
		SaveCalled:        false,
	}
}

// Save はモックのユーザー保存を行います
func (m *MockUserDAO) Save(_ context.Context, _ *models.User) error {
	m.SaveCalled = true
	if m.ShouldReturnError {
		return domain_errors.ErrDatabaseError
	}
	return nil
}
