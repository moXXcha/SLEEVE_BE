package code_models

import (
	"testing"
)

func TestNewPassword_Success(t *testing.T) {
	// Arrange
	var valid_passwords []string
	var password Password
	var err error

	valid_passwords = []string{
		"Password1!",
		"Abcdefg1@",
		"TestPass123#",
		"Complex$Pass99",
		"My_Secure!Pass1",
	}
	// Act & Assert
	for _, valid_password := range valid_passwords {
		password, err = NewPassword(valid_password)
		if err != nil {
			t.Errorf("expected no error for password %s, got %v", valid_password, err)
		}
		if password.Value() == "" {
			t.Error("expected password value to be not empty")
		}
	}
}

func TestNewPassword_TooShort(t *testing.T) {
	// Arrange
	var short_passwords []string
	var err error

	short_passwords = []string{
		"",
		"Pass1!",
		"Abc123!",
	}
	// Act & Assert
	for _, short_password := range short_passwords {
		_, err = NewPassword(short_password)
		if err == nil {
			t.Errorf("expected error for short password %s, got nil", short_password)
		}
	}
}

func TestNewPassword_NoUppercase(t *testing.T) {
	// Arrange
	var password_no_upper string
	var err error

	password_no_upper = "password123!"
	// Act & Assert
	_, err = NewPassword(password_no_upper)
	if err == nil {
		t.Error("expected error for password without uppercase, got nil")
	}
}

func TestNewPassword_NoLowercase(t *testing.T) {
	// Arrange
	var password_no_lower string
	var err error

	password_no_lower = "PASSWORD123!"
	// Act & Assert
	_, err = NewPassword(password_no_lower)
	if err == nil {
		t.Error("expected error for password without lowercase, got nil")
	}
}

func TestNewPassword_NoDigit(t *testing.T) {
	// Arrange
	var password_no_digit string
	var err error

	password_no_digit = "PasswordTest!"
	// Act & Assert
	_, err = NewPassword(password_no_digit)
	if err == nil {
		t.Error("expected error for password without digit, got nil")
	}
}

func TestNewPassword_NoSymbol(t *testing.T) {
	// Arrange
	var password_no_symbol string
	var err error

	password_no_symbol = "Password123"
	// Act & Assert
	_, err = NewPassword(password_no_symbol)
	if err == nil {
		t.Error("expected error for password without symbol, got nil")
	}
}

func TestPassword_Equals(t *testing.T) {
	// Arrange
	var password1 Password
	var password2 Password
	var password3 Password
	var err error

	password1, err = NewPassword("Password1!")
	if err != nil {
		t.Fatalf("failed to create password1: %v", err)
	}
	password2, err = NewPassword("Password1!")
	if err != nil {
		t.Fatalf("failed to create password2: %v", err)
	}
	password3, err = NewPassword("Different1!")
	if err != nil {
		t.Fatalf("failed to create password3: %v", err)
	}
	// Act & Assert
	if !password1.Equals(password2) {
		t.Error("expected password1 and password2 to be equal")
	}
	if password1.Equals(password3) {
		t.Error("expected password1 and password3 to not be equal")
	}
}
