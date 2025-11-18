package firebase

import (
	"os"
	"testing"

	"firebase.google.com/go/v4/auth"
)

func TestInitializeFirebase_Success(t *testing.T) {
	// Arrange
	var err error

	// 環境変数を設定（テスト用）
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../../config/firebase.secret.json")
	defer os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
	defer func() {
		firebase_app = nil
		auth_client = nil
	}()
	// Act
	err = initialize_firebase()
	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if firebase_app == nil {
		t.Error("expected firebase_app to be initialized, got nil")
	}
	if auth_client == nil {
		t.Error("expected auth_client to be initialized, got nil")
	}
}

func TestInitializeFirebase_MissingCredentials(t *testing.T) {
	// Arrange
	var err error
	var original_credentials string
	var exists bool

	original_credentials, exists = os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if exists {
		os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
		defer os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", original_credentials)
	}
	// Act
	err = initialize_firebase()
	// Assert
	if err == nil {
		t.Error("expected error when credentials are missing, got nil")
	}
}

func TestGetAuthClient_Success(t *testing.T) {
	// Arrange
	var client *auth.Client
	var err error

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../../config/firebase.secret.json")
	defer os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
	defer func() {
		firebase_app = nil
		auth_client = nil
	}()
	err = initialize_firebase()
	if err != nil {
		t.Fatalf("failed to initialize firebase: %v", err)
	}
	// Act
	client, err = get_auth_client()
	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if client == nil {
		t.Error("expected auth client, got nil")
	}
}

func TestGetAuthClient_NotInitialized(t *testing.T) {
	// Arrange
	var client *auth.Client
	var err error

	// Firebase未初期化状態を作る（テスト開始時に確実にnilにする）
	firebase_app = nil
	auth_client = nil
	defer func() {
		firebase_app = nil
		auth_client = nil
	}()
	// Act
	client, err = get_auth_client()
	// Assert
	if err == nil {
		t.Error("expected error when firebase is not initialized, got nil")
	}
	if client != nil {
		t.Error("expected nil client, got non-nil")
	}
}
