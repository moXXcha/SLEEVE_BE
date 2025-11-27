package models

import (
	"fmt"
	"unicode"
)

// MinPasswordLength はパスワードの最小文字数です
const MinPasswordLength = 8

// Password はパスワードを表す値オブジェクトです
type Password struct {
	value string
}

// NewPassword は新しいPassword値オブジェクトを作成します
// パスワードは8文字以上で、大文字・小文字・数字・記号を含む必要があります
func NewPassword(value string) (Password, error) {
	var err error

	err = validate_password(value)
	if err != nil {
		return Password{}, err
	}
	return Password{value: value}, nil
}

// validate_password はパスワードの強度を検証します
func validate_password(value string) error {
	var has_upper bool
	var has_lower bool
	var has_digit bool
	var has_symbol bool

	has_upper = false
	has_lower = false
	has_digit = false
	has_symbol = false
	if len(value) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters", MinPasswordLength)
	}
	for _, char := range value {
		if unicode.IsUpper(char) {
			has_upper = true
		}
		if unicode.IsLower(char) {
			has_lower = true
		}
		if unicode.IsDigit(char) {
			has_digit = true
		}
		if is_symbol(char) {
			has_symbol = true
		}
	}
	if !has_upper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !has_lower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !has_digit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !has_symbol {
		return fmt.Errorf("password must contain at least one symbol")
	}
	return nil
}

// is_symbol は文字が記号かどうかを判定します
func is_symbol(char rune) bool {
	var symbols string

	symbols = "!@#$%^&*()_+-=[]{}|;':\",./<>?`~"
	for _, s := range symbols {
		if char == s {
			return true
		}
	}
	return false
}

// Value はパスワードの値を返します
func (p Password) Value() string {
	return p.value
}

// String はパスワードをマスクして返します（セキュリティ上の理由から）
func (Password) String() string {
	return "********"
}

// Equals は他のPasswordと等しいかどうかを判定します
func (p Password) Equals(other Password) bool {
	return p.value == other.value
}
