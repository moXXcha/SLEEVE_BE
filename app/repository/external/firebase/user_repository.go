package firebase

import (
	"context"
	"fmt"
	"strings"

	domain_errors "sleeve/domain/errors"
	"sleeve/domain/models"

	"firebase.google.com/go/v4/auth"
)

// FirebaseAuthClientInterface はFirebase Auth Clientのインターフェースです
type FirebaseAuthClientInterface interface {
	CreateUser(ctx context.Context, params *auth.UserToCreate) (*auth.UserRecord, error)
	GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error)
	DeleteUser(ctx context.Context, uid string) error
}

// FirebaseUserRepository はFirebase Authenticationを使用したユーザーリポジトリです
type FirebaseUserRepository struct {
	auth_client FirebaseAuthClientInterface
}

// NewFirebaseUserRepository は新しいFirebaseUserRepositoryを作成します
func NewFirebaseUserRepository(auth_client FirebaseAuthClientInterface) *FirebaseUserRepository {
	return &FirebaseUserRepository{
		auth_client: auth_client,
	}
}

// CreateUser はFirebase Authenticationに新しいユーザーを作成します
func (r *FirebaseUserRepository) CreateUser(ctx context.Context, email models.Email, password models.Password) (string, error) {
	var params *auth.UserToCreate
	var user_record *auth.UserRecord
	var err error

	params = (&auth.UserToCreate{}).
		Email(email.Value()).
		Password(password.Value()).
		EmailVerified(false)
	user_record, err = r.auth_client.CreateUser(ctx, params)
	if err != nil {
		if is_duplicate_email_error(err) {
			return "", fmt.Errorf("%w: %s", domain_errors.ErrDuplicateEmail, email.Value())
		}
		return "", fmt.Errorf("%w: %w", domain_errors.ErrFirebaseAuthFailed, err)
	}
	return user_record.UID, nil
}

// CheckEmailExists はメールアドレスが既に登録されているかをチェックします
func (r *FirebaseUserRepository) CheckEmailExists(ctx context.Context, email models.Email) (bool, error) {
	var err error

	_, err = r.auth_client.GetUserByEmail(ctx, email.Value())
	if err != nil {
		if is_user_not_found_error(err) {
			return false, nil
		}
		return false, fmt.Errorf("%w: %w", domain_errors.ErrFirebaseAuthFailed, err)
	}
	return true, nil
}

// DeleteUser はFirebase Authenticationからユーザーを削除します
func (r *FirebaseUserRepository) DeleteUser(ctx context.Context, firebase_uid string) error {
	var err error

	err = r.auth_client.DeleteUser(ctx, firebase_uid)
	if err != nil {
		return fmt.Errorf("%w: %w", domain_errors.ErrFirebaseAuthFailed, err)
	}
	return nil
}

// is_duplicate_email_error はFirebaseのメール重複エラーかどうかを判定します
func is_duplicate_email_error(err error) bool {
	var error_message string

	error_message = err.Error()
	return strings.Contains(error_message, "EMAIL_EXISTS") ||
		strings.Contains(error_message, "email-already-exists")
}

// is_user_not_found_error はFirebaseのユーザー未検出エラーかどうかを判定します
func is_user_not_found_error(err error) bool {
	var error_message string

	error_message = err.Error()
	return strings.Contains(error_message, "USER_NOT_FOUND") ||
		strings.Contains(error_message, "user-not-found")
}
