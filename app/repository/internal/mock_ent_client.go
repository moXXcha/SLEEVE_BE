package internal

import (
	"context"
	"fmt"
	"time"

	"sleeve/ent"

	"github.com/google/uuid"
)

// MockEntClient はテスト用のモックEntクライアントです
type MockEntClient struct {
	should_return_duplicate_error bool
	mock_user                     *ent.User
}

// NewMockEntClient は新しいMockEntClientを作成します
func NewMockEntClient() *MockEntClient {
	return &MockEntClient{
		should_return_duplicate_error: false,
		mock_user:                     nil,
	}
}

// NewMockEntClientWithDuplicateError は重複エラーを返すモッククライアントを作成します
func NewMockEntClientWithDuplicateError() *MockEntClient {
	return &MockEntClient{
		should_return_duplicate_error: true,
		mock_user:                     nil,
	}
}

// NewMockEntClientWithUser は指定されたpublic_idを持つユーザーを返すモッククライアントを作成します
func NewMockEntClientWithUser(public_id uuid.UUID) *MockEntClient {
	var now time.Time

	now = time.Now()
	return &MockEntClient{
		should_return_duplicate_error: false,
		mock_user: &ent.User{
			ID:          1,
			PublicID:    public_id,
			FirebaseUID: "firebase_uid_123",
			Email:       "test@example.com",
			CreatedAt:   now,
			UpdatedAt:   now,
			DeletedAt:   nil,
		},
	}
}

// NewMockEntClientWithUserByFirebaseUID は指定されたFirebase UIDを持つユーザーを返すモッククライアントを作成します
func NewMockEntClientWithUserByFirebaseUID(firebase_uid string) *MockEntClient {
	var now time.Time

	now = time.Now()
	return &MockEntClient{
		should_return_duplicate_error: false,
		mock_user: &ent.User{
			ID:          1,
			PublicID:    uuid.New(),
			FirebaseUID: firebase_uid,
			Email:       "test@example.com",
			CreatedAt:   now,
			UpdatedAt:   now,
			DeletedAt:   nil,
		},
	}
}

// NewMockEntClientWithUserByEmail は指定されたEmailを持つユーザーを返すモッククライアントを作成します
func NewMockEntClientWithUserByEmail(email string) *MockEntClient {
	var now time.Time

	now = time.Now()
	return &MockEntClient{
		should_return_duplicate_error: false,
		mock_user: &ent.User{
			ID:          1,
			PublicID:    uuid.New(),
			FirebaseUID: "firebase_uid_123",
			Email:       email,
			CreatedAt:   now,
			UpdatedAt:   now,
			DeletedAt:   nil,
		},
	}
}

// GetUserClient はモックのUserClientを返します
func (m *MockEntClient) GetUserClient() UserClientInterface {
	return &MockUserClient{
		should_return_duplicate_error: m.should_return_duplicate_error,
		mock_user:                     m.mock_user,
	}
}

// MockUserClient はモックのUserClientです
type MockUserClient struct {
	should_return_duplicate_error bool
	mock_user                     *ent.User
}

// Create はモックのUserCreate Builderを返します
func (m *MockUserClient) Create() UserCreateInterface {
	return &MockUserCreate{
		should_return_duplicate_error: m.should_return_duplicate_error,
	}
}

// Query はモックのUserQuery Builderを返します
func (m *MockUserClient) Query() UserQueryInterface {
	return &MockUserQuery{
		mock_user: m.mock_user,
	}
}

// MockUserCreate はモックのUserCreate Builderです
type MockUserCreate struct {
	should_return_duplicate_error bool
	public_id                     uuid.UUID
	firebase_uid                  string
	email                         string
}

// SetPublicID は公開IDを設定します
func (m *MockUserCreate) SetPublicID(id uuid.UUID) UserCreateInterface {
	m.public_id = id
	return m
}

// SetFirebaseUID はFirebase UIDを設定します
func (m *MockUserCreate) SetFirebaseUID(uid string) UserCreateInterface {
	m.firebase_uid = uid
	return m
}

// SetEmail はメールアドレスを設定します
func (m *MockUserCreate) SetEmail(email string) UserCreateInterface {
	m.email = email
	return m
}

// Save はユーザーを保存します
func (m *MockUserCreate) Save(ctx context.Context) (*ent.User, error) {
	var now time.Time

	if m.should_return_duplicate_error {
		return nil, fmt.Errorf("UNIQUE constraint failed")
	}
	now = time.Now()
	return &ent.User{
		ID:          1,
		PublicID:    m.public_id,
		FirebaseUID: m.firebase_uid,
		Email:       m.email,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   nil,
	}, nil
}

// MockUserQuery はモックのUserQuery Builderです
type MockUserQuery struct {
	mock_user *ent.User
}

// Where は条件を追加します
func (m *MockUserQuery) Where(predicates ...any) UserQueryInterface {
	return m
}

// Only は単一のユーザーを返します
func (m *MockUserQuery) Only(ctx context.Context) (*ent.User, error) {
	if m.mock_user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return m.mock_user, nil
}

// Exist はユーザーが存在するかを返します
func (m *MockUserQuery) Exist(ctx context.Context) (bool, error) {
	return m.mock_user != nil, nil
}
