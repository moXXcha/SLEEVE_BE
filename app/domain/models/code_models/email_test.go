package code_models

import (
	"testing"
)

func TestNewEmail_Success(t *testing.T) {
	// Arrange
	var valid_emails []string
	var email Email
	var err error

	valid_emails = []string{
		"test@example.com",
		"user.name@domain.co.jp",
		"user+tag@example.org",
		"test123@subdomain.example.com",
	}
	// Act & Assert
	for _, valid_email := range valid_emails {
		email, err = NewEmail(valid_email)
		if err != nil {
			t.Errorf("expected no error for email %s, got %v", valid_email, err)
		}
		if email.Value() != valid_email {
			t.Errorf("expected email value %s, got %s", valid_email, email.Value())
		}
	}
}

func TestNewEmail_InvalidFormat(t *testing.T) {
	// Arrange
	var invalid_emails []string
	var err error

	invalid_emails = []string{
		"",
		"invalid",
		"invalid@",
		"@domain.com",
		"invalid@domain",
		"invalid@@domain.com",
		"invalid @domain.com",
		"invalid@ domain.com",
	}
	// Act & Assert
	for _, invalid_email := range invalid_emails {
		_, err = NewEmail(invalid_email)
		if err == nil {
			t.Errorf("expected error for invalid email %s, got nil", invalid_email)
		}
	}
}

func TestEmail_Equals(t *testing.T) {
	// Arrange
	var email1 Email
	var email2 Email
	var email3 Email
	var err error

	email1, err = NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create email1: %v", err)
	}
	email2, err = NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create email2: %v", err)
	}
	email3, err = NewEmail("other@example.com")
	if err != nil {
		t.Fatalf("failed to create email3: %v", err)
	}
	// Act & Assert
	if !email1.Equals(email2) {
		t.Error("expected email1 and email2 to be equal")
	}
	if email1.Equals(email3) {
		t.Error("expected email1 and email3 to not be equal")
	}
}

func TestEmail_String(t *testing.T) {
	// Arrange
	var email Email
	var expected_value string
	var err error

	expected_value = "test@example.com"
	email, err = NewEmail(expected_value)
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}
	// Act & Assert
	if email.String() != expected_value {
		t.Errorf("expected string %s, got %s", expected_value, email.String())
	}
}
