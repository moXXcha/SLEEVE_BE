package internal

import (
	"context"
	"testing"

	domain_errors "sleeve/domain/errors"
	"sleeve/domain/models"

	"github.com/google/uuid"
)

// TestUserDAO_Save_Success はユーザー保存が成功するケースをテストします
func TestUserDAO_Save_Success(t *testing.T) {
	var ctx context.Context
	var email models.Email
	var user *models.User
	var dao *UserDAO
	var err error

	ctx = context.Background()
	email, err = models.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	user, err = models.NewUser("firebase_uid_123", email)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	dao = NewUserDAO(NewMockEntClient())
	err = dao.Save(ctx, user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestUserDAO_Save_DuplicateEmail は重複メール時にエラーを返すケースをテストします
func TestUserDAO_Save_DuplicateEmail(t *testing.T) {
	var ctx context.Context
	var email models.Email
	var user *models.User
	var dao *UserDAO
	var err error

	ctx = context.Background()
	email, err = models.NewEmail("duplicate@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	user, err = models.NewUser("firebase_uid_123", email)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	dao = NewUserDAO(NewMockEntClientWithDuplicateError())
	err = dao.Save(ctx, user)
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
	if !domain_errors.IsUserDomainError(err) {
		t.Errorf("expected user domain error, got %v", err)
	}
}

// TestUserDAO_FindByPublicID_Success は公開IDでユーザーを取得するケースをテストします
func TestUserDAO_FindByPublicID_Success(t *testing.T) {
	var ctx context.Context
	var dao *UserDAO
	var found_user *models.User
	var public_id uuid.UUID
	var err error

	ctx = context.Background()
	public_id = uuid.New()
	dao = NewUserDAO(NewMockEntClientWithUser(public_id))
	found_user, err = dao.FindByPublicID(ctx, public_id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if found_user == nil {
		t.Error("expected user to be found")
	}
	if found_user != nil && found_user.PublicID() != public_id {
		t.Errorf("expected public_id %s, got %s", public_id, found_user.PublicID())
	}
}

// TestUserDAO_FindByPublicID_NotFound はユーザーが見つからない場合をテストします
func TestUserDAO_FindByPublicID_NotFound(t *testing.T) {
	var ctx context.Context
	var dao *UserDAO
	var public_id uuid.UUID
	var err error

	ctx = context.Background()
	public_id = uuid.New()
	dao = NewUserDAO(NewMockEntClient())
	_, err = dao.FindByPublicID(ctx, public_id)
	if err == nil {
		t.Error("expected error for not found user, got nil")
	}
	if !domain_errors.IsUserDomainError(err) {
		t.Errorf("expected user domain error, got %v", err)
	}
}

// TestUserDAO_FindByFirebaseUID_Success はFirebase UIDでユーザーを取得するケースをテストします
func TestUserDAO_FindByFirebaseUID_Success(t *testing.T) {
	var ctx context.Context
	var dao *UserDAO
	var found_user *models.User
	var firebase_uid string
	var err error

	ctx = context.Background()
	firebase_uid = "firebase_uid_123"
	dao = NewUserDAO(NewMockEntClientWithUserByFirebaseUID(firebase_uid))
	found_user, err = dao.FindByFirebaseUID(ctx, firebase_uid)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if found_user == nil {
		t.Error("expected user to be found")
	}
	if found_user != nil && found_user.FirebaseUID() != firebase_uid {
		t.Errorf("expected firebase_uid %s, got %s", firebase_uid, found_user.FirebaseUID())
	}
}

// TestUserDAO_FindByEmail_Success はEmailでユーザーを取得するケースをテストします
func TestUserDAO_FindByEmail_Success(t *testing.T) {
	var ctx context.Context
	var dao *UserDAO
	var found_user *models.User
	var email models.Email
	var err error

	ctx = context.Background()
	email, err = models.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	dao = NewUserDAO(NewMockEntClientWithUserByEmail(email.Value()))
	found_user, err = dao.FindByEmail(ctx, email)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if found_user == nil {
		t.Error("expected user to be found")
	}
	if found_user != nil && found_user.Email().Value() != email.Value() {
		t.Errorf("expected email %s, got %s", email.Value(), found_user.Email().Value())
	}
}

// TestUserDAO_ExistsByEmail_Exists はメールが存在する場合をテストします
func TestUserDAO_ExistsByEmail_Exists(t *testing.T) {
	var ctx context.Context
	var email models.Email
	var dao *UserDAO
	var exists bool
	var err error

	ctx = context.Background()
	email, err = models.NewEmail("existing@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	dao = NewUserDAO(NewMockEntClientWithUserByEmail(email.Value()))
	exists, err = dao.ExistsByEmail(ctx, email)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !exists {
		t.Error("expected email to exist")
	}
}

// TestUserDAO_ExistsByEmail_NotExists はメールが存在しない場合をテストします
func TestUserDAO_ExistsByEmail_NotExists(t *testing.T) {
	var ctx context.Context
	var email models.Email
	var dao *UserDAO
	var exists bool
	var err error

	ctx = context.Background()
	email, err = models.NewEmail("notexisting@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	dao = NewUserDAO(NewMockEntClient())
	exists, err = dao.ExistsByEmail(ctx, email)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if exists {
		t.Error("expected email to not exist")
	}
}
