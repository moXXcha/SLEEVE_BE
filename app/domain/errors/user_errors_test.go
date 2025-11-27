package errors

import (
	"errors"
	"testing"
)

func TestErrInvalidEmail(t *testing.T) {
	// Arrange & Act
	var err error

	err = ErrInvalidEmail
	// Assert
	if err == nil {
		t.Error("expected ErrInvalidEmail to be not nil")
	}
	if err.Error() != "メールアドレスの形式が不正です" {
		t.Errorf("expected error message to be 'メールアドレスの形式が不正です', got '%s'", err.Error())
	}
}

func TestErrWeakPassword(t *testing.T) {
	// Arrange & Act
	var err error

	err = ErrWeakPassword
	// Assert
	if err == nil {
		t.Error("expected ErrWeakPassword to be not nil")
	}
	if err.Error() != "パスワードは8文字以上で、英字・数字・記号を含む必要があります" {
		t.Errorf("expected error message to be 'パスワードは8文字以上で、英字・数字・記号を含む必要があります', got '%s'", err.Error())
	}
}

func TestErrDuplicateEmail(t *testing.T) {
	// Arrange & Act
	var err error

	err = ErrDuplicateEmail
	// Assert
	if err == nil {
		t.Error("expected ErrDuplicateEmail to be not nil")
	}
	if err.Error() != "このメールアドレスは既に登録されています" {
		t.Errorf("expected error message to be 'このメールアドレスは既に登録されています', got '%s'", err.Error())
	}
}

func TestErrUserNotFound(t *testing.T) {
	// Arrange & Act
	var err error

	err = ErrUserNotFound
	// Assert
	if err == nil {
		t.Error("expected ErrUserNotFound to be not nil")
	}
	if err.Error() != "ユーザーが見つかりません" {
		t.Errorf("expected error message to be 'ユーザーが見つかりません', got '%s'", err.Error())
	}
}

func TestErrFirebaseAuthFailed(t *testing.T) {
	// Arrange & Act
	var err error

	err = ErrFirebaseAuthFailed
	// Assert
	if err == nil {
		t.Error("expected ErrFirebaseAuthFailed to be not nil")
	}
	if err.Error() != "Firebase Authenticationでエラーが発生しました" {
		t.Errorf("expected error message to be 'Firebase Authenticationでエラーが発生しました', got '%s'", err.Error())
	}
}

func TestErrDatabaseError(t *testing.T) {
	// Arrange & Act
	var err error

	err = ErrDatabaseError
	// Assert
	if err == nil {
		t.Error("expected ErrDatabaseError to be not nil")
	}
	if err.Error() != "データベースエラーが発生しました" {
		t.Errorf("expected error message to be 'データベースエラーが発生しました', got '%s'", err.Error())
	}
}

func TestErrJWTGenerationFailed(t *testing.T) {
	// Arrange & Act
	var err error

	err = ErrJWTGenerationFailed
	// Assert
	if err == nil {
		t.Error("expected ErrJWTGenerationFailed to be not nil")
	}
	if err.Error() != "トークン生成に失敗しました" {
		t.Errorf("expected error message to be 'トークン生成に失敗しました', got '%s'", err.Error())
	}
}

func TestIsUserDomainError(t *testing.T) {
	// Arrange
	var user_errors []error
	var other_error error

	user_errors = []error{
		ErrInvalidEmail,
		ErrWeakPassword,
		ErrDuplicateEmail,
		ErrUserNotFound,
		ErrFirebaseAuthFailed,
		ErrDatabaseError,
		ErrJWTGenerationFailed,
	}
	other_error = errors.New("some other error")
	// Act & Assert
	for _, err := range user_errors {
		if !IsUserDomainError(err) {
			t.Errorf("expected %v to be a user domain error", err)
		}
	}
	if IsUserDomainError(other_error) {
		t.Error("expected other_error to not be a user domain error")
	}
}
