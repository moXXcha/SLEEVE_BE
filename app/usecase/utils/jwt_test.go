package utils

import (
	"testing"
	"time"
)

// テスト用定数
const (
	testSecretKey   = "test_secret_key_for_testing_1234567890"
	testUserID      = "user_123"
	testFirebaseUID = "firebase_uid_456"
)

// TestJWTService_GenerateAccessToken_Success はアクセストークン生成が成功するケースをテストします
func TestJWTService_GenerateAccessToken_Success(t *testing.T) {
	var service *JWTService
	var token string
	var user_id string
	var firebase_uid string
	var err error

	service = NewJWTService(testSecretKey)
	user_id = testUserID
	firebase_uid = testFirebaseUID
	token, err = service.GenerateAccessToken(user_id, firebase_uid)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected token to be non-empty")
	}
}

// TestJWTService_GenerateRefreshToken_Success はリフレッシュトークン生成が成功するケースをテストします
func TestJWTService_GenerateRefreshToken_Success(t *testing.T) {
	var service *JWTService
	var token string
	var user_id string
	var firebase_uid string
	var err error

	service = NewJWTService(testSecretKey)
	user_id = testUserID
	firebase_uid = testFirebaseUID
	token, err = service.GenerateRefreshToken(user_id, firebase_uid)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected token to be non-empty")
	}
}

// TestJWTService_ValidateToken_ValidAccessToken は有効なアクセストークンの検証をテストします
func TestJWTService_ValidateToken_ValidAccessToken(t *testing.T) {
	var service *JWTService
	var token string
	var claims *JWTClaims
	var user_id string
	var firebase_uid string
	var err error

	service = NewJWTService(testSecretKey)
	user_id = testUserID
	firebase_uid = testFirebaseUID
	token, err = service.GenerateAccessToken(user_id, firebase_uid)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	claims, err = service.ValidateToken(token)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if claims == nil {
		t.Fatal("expected claims to be non-nil")
	}
	if claims.UserID != user_id {
		t.Errorf("expected user_id %s, got %s", user_id, claims.UserID)
	}
	if claims.FirebaseUID != firebase_uid {
		t.Errorf("expected firebase_uid %s, got %s", firebase_uid, claims.FirebaseUID)
	}
	if claims.TokenType != TokenTypeAccess {
		t.Errorf("expected token_type %s, got %s", TokenTypeAccess, claims.TokenType)
	}
}

// TestJWTService_ValidateToken_ValidRefreshToken は有効なリフレッシュトークンの検証をテストします
func TestJWTService_ValidateToken_ValidRefreshToken(t *testing.T) {
	var service *JWTService
	var token string
	var claims *JWTClaims
	var user_id string
	var firebase_uid string
	var err error

	service = NewJWTService(testSecretKey)
	user_id = testUserID
	firebase_uid = testFirebaseUID
	token, err = service.GenerateRefreshToken(user_id, firebase_uid)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	claims, err = service.ValidateToken(token)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if claims == nil {
		t.Fatal("expected claims to be non-nil")
	}
	if claims.UserID != user_id {
		t.Errorf("expected user_id %s, got %s", user_id, claims.UserID)
	}
	if claims.FirebaseUID != firebase_uid {
		t.Errorf("expected firebase_uid %s, got %s", firebase_uid, claims.FirebaseUID)
	}
	if claims.TokenType != TokenTypeRefresh {
		t.Errorf("expected token_type %s, got %s", TokenTypeRefresh, claims.TokenType)
	}
}

// TestJWTService_ValidateToken_InvalidToken は無効なトークンの検証をテストします
func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	var service *JWTService
	var err error

	service = NewJWTService(testSecretKey)
	_, err = service.ValidateToken("invalid_token")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

// TestJWTService_ValidateToken_WrongSecret は異なるシークレットで署名されたトークンの検証をテストします
func TestJWTService_ValidateToken_WrongSecret(t *testing.T) {
	var service1 *JWTService
	var service2 *JWTService
	var token string
	var err error

	service1 = NewJWTService("secret_key_1_for_testing_1234567890")
	service2 = NewJWTService("secret_key_2_for_testing_1234567890")
	token, err = service1.GenerateAccessToken(testUserID, testFirebaseUID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	_, err = service2.ValidateToken(token)
	if err == nil {
		t.Error("expected error for wrong secret, got nil")
	}
}

// TestJWTService_AccessTokenExpiry はアクセストークンの有効期限をテストします
func TestJWTService_AccessTokenExpiry(t *testing.T) {
	var service *JWTService
	var token string
	var claims *JWTClaims
	var expected_expiry time.Duration
	var actual_expiry time.Duration
	var err error

	service = NewJWTService(testSecretKey)
	token, err = service.GenerateAccessToken(testUserID, testFirebaseUID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	claims, err = service.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}
	expected_expiry = AccessTokenExpiry
	actual_expiry = time.Until(claims.ExpiresAt.Time)
	// 許容誤差は1分
	if actual_expiry < expected_expiry-time.Minute || actual_expiry > expected_expiry+time.Minute {
		t.Errorf("expected expiry around %v, got %v", expected_expiry, actual_expiry)
	}
}

// TestJWTService_RefreshTokenExpiry はリフレッシュトークンの有効期限をテストします
func TestJWTService_RefreshTokenExpiry(t *testing.T) {
	var service *JWTService
	var token string
	var claims *JWTClaims
	var expected_expiry time.Duration
	var actual_expiry time.Duration
	var err error

	service = NewJWTService(testSecretKey)
	token, err = service.GenerateRefreshToken(testUserID, testFirebaseUID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	claims, err = service.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}
	expected_expiry = RefreshTokenExpiry
	actual_expiry = time.Until(claims.ExpiresAt.Time)
	// 許容誤差は1分
	if actual_expiry < expected_expiry-time.Minute || actual_expiry > expected_expiry+time.Minute {
		t.Errorf("expected expiry around %v, got %v", expected_expiry, actual_expiry)
	}
}

// TestJWTService_GenerateTokenPair_Success はトークンペア生成が成功するケースをテストします
func TestJWTService_GenerateTokenPair_Success(t *testing.T) {
	var service *JWTService
	var token_pair *TokenPair
	var user_id string
	var firebase_uid string
	var err error

	service = NewJWTService(testSecretKey)
	user_id = testUserID
	firebase_uid = testFirebaseUID
	token_pair, err = service.GenerateTokenPair(user_id, firebase_uid)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token_pair == nil {
		t.Fatal("expected token_pair to be non-nil")
	}
	if token_pair.AccessToken == "" {
		t.Error("expected access_token to be non-empty")
	}
	if token_pair.RefreshToken == "" {
		t.Error("expected refresh_token to be non-empty")
	}
	if token_pair.AccessToken == token_pair.RefreshToken {
		t.Error("access_token and refresh_token should be different")
	}
}
