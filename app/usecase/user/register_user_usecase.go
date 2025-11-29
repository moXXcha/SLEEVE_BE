package user

import (
	"context"
	"fmt"

	domain_errors "sleeve/domain/errors"
	"sleeve/domain/models"
	"sleeve/usecase/utils"
)

// FirebaseUserRepositoryInterface はFirebaseユーザーリポジトリのインターフェースです
type FirebaseUserRepositoryInterface interface {
	CreateUser(ctx context.Context, email models.Email, password models.Password) (string, error)
	DeleteUser(ctx context.Context, firebase_uid string) error
}

// UserDAOInterface はUserDAOのインターフェースです
type UserDAOInterface interface {
	Save(ctx context.Context, user *models.User) error
}

// RegisterUserResult はユーザー登録の結果を表します
type RegisterUserResult struct {
	User         *models.User
	AccessToken  string
	RefreshToken string
}

// RegisterUserUseCase はユーザー登録のユースケースです
type RegisterUserUseCase struct {
	firebase_repo FirebaseUserRepositoryInterface
	user_dao      UserDAOInterface
	jwt_service   *utils.JWTService
}

// NewRegisterUserUseCase は新しいRegisterUserUseCaseを作成します
func NewRegisterUserUseCase(
	firebase_repo FirebaseUserRepositoryInterface,
	user_dao UserDAOInterface,
	jwt_service *utils.JWTService,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		firebase_repo: firebase_repo,
		user_dao:      user_dao,
		jwt_service:   jwt_service,
	}
}

// Execute はユーザー登録を実行します
func (uc *RegisterUserUseCase) Execute(ctx context.Context, email_str, password_str string) (*RegisterUserResult, error) {
	var email models.Email
	var password models.Password
	var firebase_uid string
	var user *models.User
	var token_pair *utils.TokenPair
	var err error

	// メールアドレスのバリデーション
	email, err = models.NewEmail(email_str)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain_errors.ErrInvalidEmail, err)
	}

	// パスワードのバリデーション
	password, err = models.NewPassword(password_str)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain_errors.ErrWeakPassword, err)
	}

	// Firebase Authenticationにユーザーを作成
	firebase_uid, err = uc.firebase_repo.CreateUser(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// ドメインモデルを作成
	user, err = models.NewUser(firebase_uid, email)
	if err != nil {
		// Firebaseユーザーをロールバック
		_ = uc.firebase_repo.DeleteUser(ctx, firebase_uid)
		return nil, fmt.Errorf("%w: %w", domain_errors.ErrDatabaseError, err)
	}

	// 自社DBに保存
	err = uc.user_dao.Save(ctx, user)
	if err != nil {
		// Firebaseユーザーをロールバック
		_ = uc.firebase_repo.DeleteUser(ctx, firebase_uid)
		return nil, fmt.Errorf("%w", err)
	}

	// JWTを発行
	token_pair, err = uc.jwt_service.GenerateTokenPair(user.PublicID().String(), firebase_uid)
	if err != nil {
		// Firebaseユーザーをロールバック（JWT発行失敗時もロールバック）
		_ = uc.firebase_repo.DeleteUser(ctx, firebase_uid)
		return nil, fmt.Errorf("%w", err)
	}

	return &RegisterUserResult{
		User:         user,
		AccessToken:  token_pair.AccessToken,
		RefreshToken: token_pair.RefreshToken,
	}, nil
}
