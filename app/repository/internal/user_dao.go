package internal

import (
	"context"
	"fmt"
	"strings"

	domain_errors "sleeve/domain/errors"
	"sleeve/domain/models"
	"sleeve/ent"

	"github.com/google/uuid"
)

// EntClientInterface はEnt Clientのインターフェースです
type EntClientInterface interface {
	GetUserClient() UserClientInterface
}

// UserClientInterface はEnt User Clientのインターフェースです
type UserClientInterface interface {
	Create() UserCreateInterface
	Query() UserQueryInterface
}

// UserCreateInterface はEnt User Create Builderのインターフェースです
type UserCreateInterface interface {
	SetPublicID(uuid.UUID) UserCreateInterface
	SetFirebaseUID(string) UserCreateInterface
	SetEmail(string) UserCreateInterface
	Save(ctx context.Context) (*ent.User, error)
}

// UserQueryInterface はEnt User Query Builderのインターフェースです
type UserQueryInterface interface {
	Where(predicates ...any) UserQueryInterface
	Only(ctx context.Context) (*ent.User, error)
	Exist(ctx context.Context) (bool, error)
}

// UserDAO はユーザーのデータアクセスオブジェクトです
type UserDAO struct {
	client EntClientInterface
}

// NewUserDAO は新しいUserDAOを作成します
func NewUserDAO(client EntClientInterface) *UserDAO {
	return &UserDAO{
		client: client,
	}
}

// Save はユーザーをDBに保存します
func (d *UserDAO) Save(ctx context.Context, user *models.User) error {
	var err error

	_, err = d.client.GetUserClient().
		Create().
		SetPublicID(user.PublicID()).
		SetFirebaseUID(user.FirebaseUID()).
		SetEmail(user.Email().Value()).
		Save(ctx)
	if err != nil {
		if is_unique_constraint_error(err) {
			return fmt.Errorf("%w: %s", domain_errors.ErrDuplicateEmail, user.Email().Value())
		}
		return fmt.Errorf("%w: %w", domain_errors.ErrDatabaseError, err)
	}
	return nil
}

// FindByPublicID は公開IDでユーザーを検索します
func (d *UserDAO) FindByPublicID(ctx context.Context, public_id uuid.UUID) (*models.User, error) {
	var ent_user *ent.User
	var err error

	ent_user, err = d.client.GetUserClient().
		Query().
		Where("public_id", public_id).
		Only(ctx)
	if err != nil {
		return nil, handle_query_error(err)
	}
	return convert_ent_user_to_domain(ent_user)
}

// FindByFirebaseUID はFirebase UIDでユーザーを検索します
func (d *UserDAO) FindByFirebaseUID(ctx context.Context, firebase_uid string) (*models.User, error) {
	var ent_user *ent.User
	var err error

	ent_user, err = d.client.GetUserClient().
		Query().
		Where("firebase_uid", firebase_uid).
		Only(ctx)
	if err != nil {
		return nil, handle_query_error(err)
	}
	return convert_ent_user_to_domain(ent_user)
}

// FindByEmail はメールアドレスでユーザーを検索します
func (d *UserDAO) FindByEmail(ctx context.Context, email models.Email) (*models.User, error) {
	var ent_user *ent.User
	var err error

	ent_user, err = d.client.GetUserClient().
		Query().
		Where("email", email.Value()).
		Only(ctx)
	if err != nil {
		return nil, handle_query_error(err)
	}
	return convert_ent_user_to_domain(ent_user)
}

// ExistsByEmail はメールアドレスが存在するかをチェックします
func (d *UserDAO) ExistsByEmail(ctx context.Context, email models.Email) (bool, error) {
	var exists bool
	var err error

	exists, err = d.client.GetUserClient().
		Query().
		Where("email", email.Value()).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("%w: %w", domain_errors.ErrDatabaseError, err)
	}
	return exists, nil
}

// convert_ent_user_to_domain はEntのUserエンティティをドメインモデルに変換します
func convert_ent_user_to_domain(ent_user *ent.User) (*models.User, error) {
	var domain_user *models.User
	var email models.Email
	var err error

	email, err = models.NewEmail(ent_user.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain_errors.ErrDatabaseError, err)
	}
	domain_user, err = models.NewUserWithPublicID(
		ent_user.PublicID,
		ent_user.FirebaseUID,
		email,
		ent_user.CreatedAt,
		ent_user.UpdatedAt,
		ent_user.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain_errors.ErrDatabaseError, err)
	}
	return domain_user, nil
}

// handle_query_error はクエリエラーをドメインエラーに変換します
func handle_query_error(err error) error {
	if is_not_found_error(err) {
		return domain_errors.ErrUserNotFound
	}
	return fmt.Errorf("%w: %w", domain_errors.ErrDatabaseError, err)
}

// is_unique_constraint_error はユニーク制約エラーかどうかを判定します
func is_unique_constraint_error(err error) bool {
	var error_message string

	error_message = err.Error()
	return strings.Contains(error_message, "UNIQUE") ||
		strings.Contains(error_message, "duplicate key") ||
		strings.Contains(error_message, "unique constraint")
}

// is_not_found_error はレコード未検出エラーかどうかを判定します
func is_not_found_error(err error) bool {
	var error_message string

	error_message = err.Error()
	return strings.Contains(error_message, "not found") ||
		strings.Contains(error_message, "no rows")
}
