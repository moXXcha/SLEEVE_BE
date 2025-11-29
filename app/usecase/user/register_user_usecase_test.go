package user

import (
	"context"
	"testing"

	domain_errors "sleeve/domain/errors"
	"sleeve/domain/models"
	"sleeve/usecase/utils"
)

// テスト用定数
const (
	testEmail       = "test@example.com"
	testPassword    = "Password123!"
	testFirebaseUID = "firebase_uid_123"
	testSecretKey   = "test_secret_key_for_testing_1234567890"
)

// TestRegisterUserUseCase_Execute_Success は正常なユーザー登録をテストします
func TestRegisterUserUseCase_Execute_Success(t *testing.T) {
	var ctx context.Context
	var use_case *RegisterUserUseCase
	var result *RegisterUserResult
	var err error

	ctx = context.Background()
	use_case = NewRegisterUserUseCase(
		NewMockFirebaseUserRepository(),
		NewMockUserDAO(),
		utils.NewJWTService(testSecretKey),
	)
	result, err = use_case.Execute(ctx, testEmail, testPassword)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result to be non-nil")
	}
	if result.User == nil {
		t.Fatal("expected user to be non-nil")
	}
	if result.User.Email().Value() != testEmail {
		t.Errorf("expected email %s, got %s", testEmail, result.User.Email().Value())
	}
	if result.AccessToken == "" {
		t.Error("expected access_token to be non-empty")
	}
	if result.RefreshToken == "" {
		t.Error("expected refresh_token to be non-empty")
	}
}

// TestRegisterUserUseCase_Execute_InvalidEmail は不正なメールでエラーを返すケースをテストします
func TestRegisterUserUseCase_Execute_InvalidEmail(t *testing.T) {
	var ctx context.Context
	var use_case *RegisterUserUseCase
	var err error

	ctx = context.Background()
	use_case = NewRegisterUserUseCase(
		NewMockFirebaseUserRepository(),
		NewMockUserDAO(),
		utils.NewJWTService(testSecretKey),
	)
	_, err = use_case.Execute(ctx, "invalid-email", testPassword)
	if err == nil {
		t.Error("expected error for invalid email, got nil")
	}
	if !domain_errors.IsUserDomainError(err) {
		t.Errorf("expected user domain error, got %v", err)
	}
}

// TestRegisterUserUseCase_Execute_WeakPassword は弱いパスワードでエラーを返すケースをテストします
func TestRegisterUserUseCase_Execute_WeakPassword(t *testing.T) {
	var ctx context.Context
	var use_case *RegisterUserUseCase
	var err error

	ctx = context.Background()
	use_case = NewRegisterUserUseCase(
		NewMockFirebaseUserRepository(),
		NewMockUserDAO(),
		utils.NewJWTService(testSecretKey),
	)
	_, err = use_case.Execute(ctx, testEmail, "weak")
	if err == nil {
		t.Error("expected error for weak password, got nil")
	}
	if !domain_errors.IsUserDomainError(err) {
		t.Errorf("expected user domain error, got %v", err)
	}
}

// TestRegisterUserUseCase_Execute_DuplicateEmail はメール重複時にエラーを返すケースをテストします
func TestRegisterUserUseCase_Execute_DuplicateEmail(t *testing.T) {
	var ctx context.Context
	var use_case *RegisterUserUseCase
	var err error

	ctx = context.Background()
	use_case = NewRegisterUserUseCase(
		NewMockFirebaseUserRepositoryWithDuplicateEmail(),
		NewMockUserDAO(),
		utils.NewJWTService(testSecretKey),
	)
	_, err = use_case.Execute(ctx, testEmail, testPassword)
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
	if !domain_errors.IsUserDomainError(err) {
		t.Errorf("expected user domain error, got %v", err)
	}
}

// TestRegisterUserUseCase_Execute_FirebaseError はFirebaseエラー時のロールバックをテストします
func TestRegisterUserUseCase_Execute_FirebaseError(t *testing.T) {
	var ctx context.Context
	var use_case *RegisterUserUseCase
	var err error

	ctx = context.Background()
	use_case = NewRegisterUserUseCase(
		NewMockFirebaseUserRepositoryWithError(),
		NewMockUserDAO(),
		utils.NewJWTService(testSecretKey),
	)
	_, err = use_case.Execute(ctx, testEmail, testPassword)
	if err == nil {
		t.Error("expected error for firebase failure, got nil")
	}
	if !domain_errors.IsUserDomainError(err) {
		t.Errorf("expected user domain error, got %v", err)
	}
}

// TestRegisterUserUseCase_Execute_DatabaseError はDBエラー時のロールバックをテストします
func TestRegisterUserUseCase_Execute_DatabaseError(t *testing.T) {
	var ctx context.Context
	var use_case *RegisterUserUseCase
	var mock_firebase *MockFirebaseUserRepository
	var err error

	ctx = context.Background()
	mock_firebase = NewMockFirebaseUserRepository()
	use_case = NewRegisterUserUseCase(
		mock_firebase,
		NewMockUserDAOWithError(),
		utils.NewJWTService(testSecretKey),
	)
	_, err = use_case.Execute(ctx, testEmail, testPassword)
	if err == nil {
		t.Error("expected error for database failure, got nil")
	}
	if !domain_errors.IsUserDomainError(err) {
		t.Errorf("expected user domain error, got %v", err)
	}
	// Firebaseユーザーがロールバックされていることを確認
	if !mock_firebase.DeleteUserCalled {
		t.Error("expected DeleteUser to be called for rollback")
	}
}

// TestRegisterUserUseCase_Execute_JWTGenerated は登録成功時にJWTが発行されることをテストします
func TestRegisterUserUseCase_Execute_JWTGenerated(t *testing.T) {
	var ctx context.Context
	var use_case *RegisterUserUseCase
	var jwt_service *utils.JWTService
	var result *RegisterUserResult
	var access_claims *utils.JWTClaims
	var refresh_claims *utils.JWTClaims
	var err error

	ctx = context.Background()
	jwt_service = utils.NewJWTService(testSecretKey)
	use_case = NewRegisterUserUseCase(
		NewMockFirebaseUserRepository(),
		NewMockUserDAO(),
		jwt_service,
	)
	result, err = use_case.Execute(ctx, testEmail, testPassword)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// アクセストークンの検証
	access_claims, err = jwt_service.ValidateToken(result.AccessToken)
	if err != nil {
		t.Errorf("expected valid access token, got error: %v", err)
	}
	if access_claims.TokenType != utils.TokenTypeAccess {
		t.Errorf("expected token type %s, got %s", utils.TokenTypeAccess, access_claims.TokenType)
	}
	// リフレッシュトークンの検証
	refresh_claims, err = jwt_service.ValidateToken(result.RefreshToken)
	if err != nil {
		t.Errorf("expected valid refresh token, got error: %v", err)
	}
	if refresh_claims.TokenType != utils.TokenTypeRefresh {
		t.Errorf("expected token type %s, got %s", utils.TokenTypeRefresh, refresh_claims.TokenType)
	}
}

// MockFirebaseUserRepository はテスト用のFirebaseリポジトリモックです
type MockFirebaseUserRepository struct {
	should_return_duplicate_error bool
	should_return_error           bool
	DeleteUserCalled              bool
}

// NewMockFirebaseUserRepository は新しいMockFirebaseUserRepositoryを作成します
func NewMockFirebaseUserRepository() *MockFirebaseUserRepository {
	return &MockFirebaseUserRepository{
		should_return_duplicate_error: false,
		should_return_error:           false,
		DeleteUserCalled:              false,
	}
}

// NewMockFirebaseUserRepositoryWithDuplicateEmail はメール重複エラーを返すモックを作成します
func NewMockFirebaseUserRepositoryWithDuplicateEmail() *MockFirebaseUserRepository {
	return &MockFirebaseUserRepository{
		should_return_duplicate_error: true,
		should_return_error:           false,
		DeleteUserCalled:              false,
	}
}

// NewMockFirebaseUserRepositoryWithError はエラーを返すモックを作成します
func NewMockFirebaseUserRepositoryWithError() *MockFirebaseUserRepository {
	return &MockFirebaseUserRepository{
		should_return_duplicate_error: false,
		should_return_error:           true,
		DeleteUserCalled:              false,
	}
}

// CreateUser はモックのユーザー作成を行います
func (m *MockFirebaseUserRepository) CreateUser(_ context.Context, _ models.Email, _ models.Password) (string, error) {
	if m.should_return_duplicate_error {
		return "", domain_errors.ErrDuplicateEmail
	}
	if m.should_return_error {
		return "", domain_errors.ErrFirebaseAuthFailed
	}
	return testFirebaseUID, nil
}

// DeleteUser はモックのユーザー削除を行います
func (m *MockFirebaseUserRepository) DeleteUser(_ context.Context, _ string) error {
	m.DeleteUserCalled = true
	return nil
}

// MockUserDAO はテスト用のUserDAOモックです
type MockUserDAO struct {
	should_return_error bool
}

// NewMockUserDAO は新しいMockUserDAOを作成します
func NewMockUserDAO() *MockUserDAO {
	return &MockUserDAO{
		should_return_error: false,
	}
}

// NewMockUserDAOWithError はエラーを返すモックを作成します
func NewMockUserDAOWithError() *MockUserDAO {
	return &MockUserDAO{
		should_return_error: true,
	}
}

// Save はモックのユーザー保存を行います
func (m *MockUserDAO) Save(_ context.Context, _ *models.User) error {
	if m.should_return_error {
		return domain_errors.ErrDatabaseError
	}
	return nil
}
