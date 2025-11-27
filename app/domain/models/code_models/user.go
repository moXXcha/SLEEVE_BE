package code_models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User はユーザーを表すエンティティです
type User struct {
	id           uuid.UUID
	firebase_uid string
	email        Email
	created_at   time.Time
	updated_at   time.Time
	deleted_at   *time.Time
}

// NewUser は新しいUserエンティティを作成します
func NewUser(firebase_uid string, email Email) (*User, error) {
	var now time.Time

	if firebase_uid == "" {
		return nil, fmt.Errorf("firebase_uid cannot be empty")
	}
	now = time.Now()
	return &User{
		id:           uuid.New(),
		firebase_uid: firebase_uid,
		email:        email,
		created_at:   now,
		updated_at:   now,
		deleted_at:   nil,
	}, nil
}

// NewUserWithID は既存のIDを持つUserエンティティを作成します（DBからの復元用）
func NewUserWithID(
	id uuid.UUID,
	firebase_uid string,
	email Email,
	created_at time.Time,
	updated_at time.Time,
	deleted_at *time.Time,
) (*User, error) {
	if firebase_uid == "" {
		return nil, fmt.Errorf("firebase_uid cannot be empty")
	}
	return &User{
		id:           id,
		firebase_uid: firebase_uid,
		email:        email,
		created_at:   created_at,
		updated_at:   updated_at,
		deleted_at:   deleted_at,
	}, nil
}

// ID はユーザーIDを返します
func (u *User) ID() uuid.UUID {
	return u.id
}

// FirebaseUID はFirebase UIDを返します
func (u *User) FirebaseUID() string {
	return u.firebase_uid
}

// Email はメールアドレスを返します
func (u *User) Email() Email {
	return u.email
}

// CreatedAt は作成日時を返します
func (u *User) CreatedAt() time.Time {
	return u.created_at
}

// UpdatedAt は更新日時を返します
func (u *User) UpdatedAt() time.Time {
	return u.updated_at
}

// DeletedAt は削除日時を返します
func (u *User) DeletedAt() *time.Time {
	return u.deleted_at
}

// IsDeleted はユーザーが論理削除されているかどうかを返します
func (u *User) IsDeleted() bool {
	return u.deleted_at != nil
}
