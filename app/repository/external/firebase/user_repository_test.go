package firebase

import (
	"context"
	"testing"

	domain_errors "sleeve/domain/errors"
	"sleeve/domain/models"
)

// TestFirebaseUserRepository_CreateUser_Success はFirebaseへのユーザー登録が成功するケースをテストします
func TestFirebaseUserRepository_CreateUser_Success(t *testing.T) {
	var ctx context.Context
	var email models.Email
	var password models.Password
	var repo *FirebaseUserRepository
	var firebase_uid string
	var err error

	ctx = context.Background()
	email, err = models.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	password, err = models.NewPassword("Password123!")
	if err != nil {
		t.Fatalf("failed to create password: %v", err)
	}
	repo = NewFirebaseUserRepository(NewMockFirebaseAuthClient())
	firebase_uid, err = repo.CreateUser(ctx, email, password)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if firebase_uid == "" {
		t.Error("expected firebase_uid to be non-empty")
	}
}

// TestFirebaseUserRepository_CreateUser_DuplicateEmail はメール重複時にエラーを返すケースをテストします
func TestFirebaseUserRepository_CreateUser_DuplicateEmail(t *testing.T) {
	var ctx context.Context
	var email models.Email
	var password models.Password
	var repo *FirebaseUserRepository
	var err error

	ctx = context.Background()
	email, err = models.NewEmail("duplicate@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	password, err = models.NewPassword("Password123!")
	if err != nil {
		t.Fatalf("failed to create password: %v", err)
	}
	repo = NewFirebaseUserRepository(NewMockFirebaseAuthClientWithDuplicateEmail())
	_, err = repo.CreateUser(ctx, email, password)
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
	if !domain_errors.IsUserDomainError(err) {
		t.Errorf("expected user domain error, got %v", err)
	}
}

// TestFirebaseUserRepository_CheckEmailExists_Exists はメールが存在する場合をテストします
func TestFirebaseUserRepository_CheckEmailExists_Exists(t *testing.T) {
	var ctx context.Context
	var email models.Email
	var repo *FirebaseUserRepository
	var exists bool
	var err error

	ctx = context.Background()
	email, err = models.NewEmail("existing@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	repo = NewFirebaseUserRepository(NewMockFirebaseAuthClientWithExistingEmail())
	exists, err = repo.CheckEmailExists(ctx, email)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !exists {
		t.Error("expected email to exist")
	}
}

// TestFirebaseUserRepository_CheckEmailExists_NotExists はメールが存在しない場合をテストします
func TestFirebaseUserRepository_CheckEmailExists_NotExists(t *testing.T) {
	var ctx context.Context
	var email models.Email
	var repo *FirebaseUserRepository
	var exists bool
	var err error

	ctx = context.Background()
	email, err = models.NewEmail("notexisting@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	repo = NewFirebaseUserRepository(NewMockFirebaseAuthClient())
	exists, err = repo.CheckEmailExists(ctx, email)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if exists {
		t.Error("expected email to not exist")
	}
}

// TestFirebaseUserRepository_DeleteUser_Success はFirebaseからのユーザー削除が成功するケースをテストします
func TestFirebaseUserRepository_DeleteUser_Success(t *testing.T) {
	var ctx context.Context
	var repo *FirebaseUserRepository
	var err error

	ctx = context.Background()
	repo = NewFirebaseUserRepository(NewMockFirebaseAuthClient())
	err = repo.DeleteUser(ctx, "firebase_uid_123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
