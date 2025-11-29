package firebase

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
)

// MockFirebaseAuthClient はテスト用のモックFirebase Authクライアントです
type MockFirebaseAuthClient struct {
	should_return_duplicate_error bool
	should_return_existing_email  bool
}

// NewMockFirebaseAuthClient は新しいMockFirebaseAuthClientを作成します
func NewMockFirebaseAuthClient() *MockFirebaseAuthClient {
	return &MockFirebaseAuthClient{
		should_return_duplicate_error: false,
		should_return_existing_email:  false,
	}
}

// NewMockFirebaseAuthClientWithDuplicateEmail は重複エラーを返すモッククライアントを作成します
func NewMockFirebaseAuthClientWithDuplicateEmail() *MockFirebaseAuthClient {
	return &MockFirebaseAuthClient{
		should_return_duplicate_error: true,
		should_return_existing_email:  false,
	}
}

// NewMockFirebaseAuthClientWithExistingEmail は既存メールを返すモッククライアントを作成します
func NewMockFirebaseAuthClientWithExistingEmail() *MockFirebaseAuthClient {
	return &MockFirebaseAuthClient{
		should_return_duplicate_error: false,
		should_return_existing_email:  true,
	}
}

// CreateUser はモックのユーザー作成処理です
func (m *MockFirebaseAuthClient) CreateUser(ctx context.Context, params *auth.UserToCreate) (*auth.UserRecord, error) {
	if m.should_return_duplicate_error {
		return nil, fmt.Errorf("EMAIL_EXISTS")
	}
	return &auth.UserRecord{
		UserInfo: &auth.UserInfo{
			UID: "mock_firebase_uid_123",
		},
	}, nil
}

// GetUserByEmail はモックのメールでユーザー取得処理です
func (m *MockFirebaseAuthClient) GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {
	if m.should_return_existing_email {
		return &auth.UserRecord{
			UserInfo: &auth.UserInfo{
				UID:   "existing_uid",
				Email: email,
			},
		}, nil
	}
	return nil, fmt.Errorf("USER_NOT_FOUND")
}

// DeleteUser はモックのユーザー削除処理です
func (m *MockFirebaseAuthClient) DeleteUser(ctx context.Context, uid string) error {
	return nil
}
