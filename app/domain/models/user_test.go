package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

const test_firebase_uid = "firebase_uid_123"

func TestNewUser_Success(t *testing.T) {
	// Arrange
	var firebase_uid string
	var email Email
	var user *User
	var err error

	firebase_uid = test_firebase_uid
	email, err = NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	// Act
	user, err = NewUser(firebase_uid, email)
	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user to be not nil")
	}
	if user.PublicID() == uuid.Nil {
		t.Error("expected user PublicID to be not nil")
	}
	if user.FirebaseUID() != firebase_uid {
		t.Errorf("expected firebase_uid %s, got %s", firebase_uid, user.FirebaseUID())
	}
	if !user.Email().Equals(email) {
		t.Errorf("expected email %s, got %s", email.Value(), user.Email().Value())
	}
	if user.CreatedAt().IsZero() {
		t.Error("expected created_at to be set")
	}
	if user.UpdatedAt().IsZero() {
		t.Error("expected updated_at to be set")
	}
	if user.DeletedAt() != nil {
		t.Error("expected deleted_at to be nil")
	}
}

func TestNewUser_EmptyFirebaseUID(t *testing.T) {
	// Arrange
	var firebase_uid string
	var email Email
	var err error

	firebase_uid = ""
	email, err = NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	// Act
	_, err = NewUser(firebase_uid, email)
	// Assert
	if err == nil {
		t.Error("expected error for empty firebase_uid, got nil")
	}
}

func TestNewUserWithPublicID_Success(t *testing.T) {
	// Arrange
	var public_id uuid.UUID
	var firebase_uid string
	var email Email
	var created_at time.Time
	var updated_at time.Time
	var user *User
	var err error

	public_id = uuid.New()
	firebase_uid = test_firebase_uid
	email, err = NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	created_at = time.Now().Add(-time.Hour)
	updated_at = time.Now()
	// Act
	user, err = NewUserWithPublicID(public_id, firebase_uid, email, created_at, updated_at, nil)
	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user to be not nil")
	}
	if user.PublicID() != public_id {
		t.Errorf("expected public_id %s, got %s", public_id, user.PublicID())
	}
	if user.FirebaseUID() != firebase_uid {
		t.Errorf("expected firebase_uid %s, got %s", firebase_uid, user.FirebaseUID())
	}
	if !user.Email().Equals(email) {
		t.Errorf("expected email %s, got %s", email.Value(), user.Email().Value())
	}
	if !user.CreatedAt().Equal(created_at) {
		t.Errorf("expected created_at %v, got %v", created_at, user.CreatedAt())
	}
	if !user.UpdatedAt().Equal(updated_at) {
		t.Errorf("expected updated_at %v, got %v", updated_at, user.UpdatedAt())
	}
	if user.DeletedAt() != nil {
		t.Error("expected deleted_at to be nil")
	}
}

func TestUser_IsDeleted(t *testing.T) {
	// Arrange
	var public_id uuid.UUID
	var firebase_uid string
	var email Email
	var created_at time.Time
	var updated_at time.Time
	var deleted_at time.Time
	var user_not_deleted *User
	var user_deleted *User
	var err error

	public_id = uuid.New()
	firebase_uid = test_firebase_uid
	email, err = NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	created_at = time.Now().Add(-time.Hour)
	updated_at = time.Now()
	deleted_at = time.Now()
	user_not_deleted, err = NewUserWithPublicID(public_id, firebase_uid, email, created_at, updated_at, nil)
	if err != nil {
		t.Fatalf("failed to create user_not_deleted: %v", err)
	}
	user_deleted, err = NewUserWithPublicID(uuid.New(), firebase_uid, email, created_at, updated_at, &deleted_at)
	if err != nil {
		t.Fatalf("failed to create user_deleted: %v", err)
	}
	// Act & Assert
	if user_not_deleted.IsDeleted() {
		t.Error("expected user_not_deleted to not be deleted")
	}
	if !user_deleted.IsDeleted() {
		t.Error("expected user_deleted to be deleted")
	}
}
