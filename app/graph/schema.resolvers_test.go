package graph

import (
	"context"
	"testing"

	domain_errors "sleeve/domain/errors"
	"sleeve/domain/models"
	"sleeve/graph/model"
	"sleeve/usecase/user"
	"sleeve/usecase/utils"
)

// テスト用定数
const (
	testEmail       = "test@example.com"
	testPassword    = "Password123!"
	testFirebaseUID = "firebase_uid_123"
	testSecretKey   = "test_secret_key_for_testing_1234567890"
)

// TestRegisterUser_Success は正常なユーザー登録をテストします
func TestRegisterUser_Success(t *testing.T) {
	var ctx context.Context
	var resolver *mutationResolver
	var input model.RegisterUserInput
	var result *model.RegisterUserPayload
	var err error

	ctx = context.Background()
	resolver = createTestMutationResolver(
		NewMockFirebaseUserRepository(),
		NewMockUserDAO(),
	)
	input = model.RegisterUserInput{
		Email:    testEmail,
		Password: testPassword,
	}
	result, err = resolver.RegisterUser(ctx, input)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result to be non-nil")
	}
	if result.User == nil {
		t.Fatal("expected user to be non-nil")
	}
	if result.User.Email != testEmail {
		t.Errorf("expected email %s, got %s", testEmail, result.User.Email)
	}
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

// TestRegisterUser_InvalidEmail は不正なメールでエラーを返すケースをテストします
func TestRegisterUser_InvalidEmail(t *testing.T) {
	var ctx context.Context
	var resolver *mutationResolver
	var input model.RegisterUserInput
	var err error

	ctx = context.Background()
	resolver = createTestMutationResolver(
		NewMockFirebaseUserRepository(),
		NewMockUserDAO(),
	)
	input = model.RegisterUserInput{
		Email:    "invalid-email",
		Password: testPassword,
	}
	_, err = resolver.RegisterUser(ctx, input)
	if err == nil {
		t.Error("expected error for invalid email, got nil")
	}
}

// TestRegisterUser_WeakPassword は弱いパスワードでエラーを返すケースをテストします
func TestRegisterUser_WeakPassword(t *testing.T) {
	var ctx context.Context
	var resolver *mutationResolver
	var input model.RegisterUserInput
	var err error

	ctx = context.Background()
	resolver = createTestMutationResolver(
		NewMockFirebaseUserRepository(),
		NewMockUserDAO(),
	)
	input = model.RegisterUserInput{
		Email:    testEmail,
		Password: "weak",
	}
	_, err = resolver.RegisterUser(ctx, input)
	if err == nil {
		t.Error("expected error for weak password, got nil")
	}
}

// TestRegisterUser_DuplicateEmail はメール重複時にエラーを返すケースをテストします
func TestRegisterUser_DuplicateEmail(t *testing.T) {
	var ctx context.Context
	var resolver *mutationResolver
	var input model.RegisterUserInput
	var err error

	ctx = context.Background()
	resolver = createTestMutationResolver(
		NewMockFirebaseUserRepositoryWithDuplicateEmail(),
		NewMockUserDAO(),
	)
	input = model.RegisterUserInput{
		Email:    testEmail,
		Password: testPassword,
	}
	_, err = resolver.RegisterUser(ctx, input)
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
}

// createTestMutationResolver はテスト用のmutationResolverを作成します
func createTestMutationResolver(
	firebase_repo user.FirebaseUserRepositoryInterface,
	user_dao user.UserDAOInterface,
) *mutationResolver {
	var resolver *Resolver
	var use_case *user.RegisterUserUseCase

	use_case = user.NewRegisterUserUseCase(
		firebase_repo,
		user_dao,
		utils.NewJWTService(testSecretKey),
	)
	resolver = &Resolver{
		RegisterUserUseCase: use_case,
	}
	return &mutationResolver{resolver}
}

// MockFirebaseUserRepository はテスト用のFirebaseリポジトリモックです
type MockFirebaseUserRepository struct {
	should_return_duplicate_error bool
	should_return_error           bool
}

// NewMockFirebaseUserRepository は新しいMockFirebaseUserRepositoryを作成します
func NewMockFirebaseUserRepository() *MockFirebaseUserRepository {
	return &MockFirebaseUserRepository{
		should_return_duplicate_error: false,
		should_return_error:           false,
	}
}

// NewMockFirebaseUserRepositoryWithDuplicateEmail はメール重複エラーを返すモックを作成します
func NewMockFirebaseUserRepositoryWithDuplicateEmail() *MockFirebaseUserRepository {
	return &MockFirebaseUserRepository{
		should_return_duplicate_error: true,
		should_return_error:           false,
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

// Save はモックのユーザー保存を行います
func (m *MockUserDAO) Save(_ context.Context, _ *models.User) error {
	if m.should_return_error {
		return domain_errors.ErrDatabaseError
	}
	return nil
}
