package errors

import (
	"errors"
)

// ユーザードメインのエラー定義
var (
	// ErrInvalidEmail はメールアドレスの形式が不正な場合のエラーです
	ErrInvalidEmail = errors.New("メールアドレスの形式が不正です")

	// ErrWeakPassword はパスワードの強度が不足している場合のエラーです
	ErrWeakPassword = errors.New("パスワードは8文字以上で、英字・数字・記号を含む必要があります")

	// ErrDuplicateEmail はメールアドレスが既に登録済みの場合のエラーです
	ErrDuplicateEmail = errors.New("このメールアドレスは既に登録されています")

	// ErrUserNotFound はユーザーが見つからない場合のエラーです
	ErrUserNotFound = errors.New("ユーザーが見つかりません")

	// ErrFirebaseAuthFailed はFirebase Authenticationでエラーが発生した場合のエラーです
	ErrFirebaseAuthFailed = errors.New("Firebase Authenticationでエラーが発生しました")

	// ErrDatabaseError はデータベースエラーが発生した場合のエラーです
	ErrDatabaseError = errors.New("データベースエラーが発生しました")

	// ErrJWTGenerationFailed はJWT生成に失敗した場合のエラーです
	ErrJWTGenerationFailed = errors.New("トークン生成に失敗しました")
)

// user_domain_errors はユーザードメインのエラー一覧です
var user_domain_errors = []error{
	ErrInvalidEmail,
	ErrWeakPassword,
	ErrDuplicateEmail,
	ErrUserNotFound,
	ErrFirebaseAuthFailed,
	ErrDatabaseError,
	ErrJWTGenerationFailed,
}

// IsUserDomainError はエラーがユーザードメインのエラーかどうかを判定します
func IsUserDomainError(err error) bool {
	for _, domain_err := range user_domain_errors {
		if errors.Is(err, domain_err) {
			return true
		}
	}
	return false
}
