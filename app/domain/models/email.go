package models

import (
	"fmt"
	"regexp"
	"strings"
)

// Email はメールアドレスを表す値オブジェクトです
type Email struct {
	value string
}

// email_regex はメールアドレスの正規表現パターンです
var email_regex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewEmail は新しいEmail値オブジェクトを作成します
// 不正な形式のメールアドレスの場合はエラーを返します
func NewEmail(value string) (Email, error) {
	var trimmed_value string

	trimmed_value = strings.TrimSpace(value)
	if trimmed_value == "" {
		return Email{}, fmt.Errorf("email cannot be empty")
	}
	if strings.Contains(trimmed_value, " ") {
		return Email{}, fmt.Errorf("email cannot contain spaces")
	}
	if !email_regex.MatchString(trimmed_value) {
		return Email{}, fmt.Errorf("invalid email format: %s", trimmed_value)
	}
	return Email{value: trimmed_value}, nil
}

// Value はメールアドレスの値を返します
func (e Email) Value() string {
	return e.value
}

// String はメールアドレスを文字列として返します
func (e Email) String() string {
	return e.value
}

// Equals は他のEmailと等しいかどうかを判定します
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
