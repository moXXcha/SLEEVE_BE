package utils

import (
	"fmt"
	"time"

	domain_errors "sleeve/domain/errors"

	"github.com/golang-jwt/jwt/v5"
)

// トークンタイプの定義
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// トークンの有効期限
const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

// JWTClaims はJWTのカスタムクレームです
type JWTClaims struct {
	UserID      string `json:"user_id"`
	FirebaseUID string `json:"firebase_uid"`
	TokenType   string `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenPair はアクセストークンとリフレッシュトークンのペアです
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// JWTService はJWTの生成・検証を行うサービスです
type JWTService struct {
	secret_key []byte
}

// NewJWTService は新しいJWTServiceを作成します
func NewJWTService(secret_key string) *JWTService {
	return &JWTService{
		secret_key: []byte(secret_key),
	}
}

// GenerateAccessToken はアクセストークンを生成します
func (s *JWTService) GenerateAccessToken(user_id, firebase_uid string) (string, error) {
	return s.generate_token(user_id, firebase_uid, TokenTypeAccess, AccessTokenExpiry)
}

// GenerateRefreshToken はリフレッシュトークンを生成します
func (s *JWTService) GenerateRefreshToken(user_id, firebase_uid string) (string, error) {
	return s.generate_token(user_id, firebase_uid, TokenTypeRefresh, RefreshTokenExpiry)
}

// GenerateTokenPair はアクセストークンとリフレッシュトークンのペアを生成します
func (s *JWTService) GenerateTokenPair(user_id, firebase_uid string) (*TokenPair, error) {
	var access_token string
	var refresh_token string
	var err error

	access_token, err = s.GenerateAccessToken(user_id, firebase_uid)
	if err != nil {
		return nil, err
	}
	refresh_token, err = s.GenerateRefreshToken(user_id, firebase_uid)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

// ValidateToken はトークンを検証し、クレームを返します
func (s *JWTService) ValidateToken(token_string string) (*JWTClaims, error) {
	var token *jwt.Token
	var claims *JWTClaims
	var claims_ok bool
	var err error

	token, err = jwt.ParseWithClaims(token_string, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		var is_hmac bool

		_, is_hmac = token.Method.(*jwt.SigningMethodHMAC)
		if !is_hmac {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret_key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain_errors.ErrJWTGenerationFailed, err)
	}
	claims, claims_ok = token.Claims.(*JWTClaims)
	if !claims_ok || !token.Valid {
		return nil, fmt.Errorf("%w: invalid token claims", domain_errors.ErrJWTGenerationFailed)
	}
	return claims, nil
}

// generate_token はトークンを生成する内部関数です
func (s *JWTService) generate_token(user_id, firebase_uid, token_type string, expiry time.Duration) (string, error) {
	var now time.Time
	var claims *JWTClaims
	var token *jwt.Token
	var token_string string
	var err error

	now = time.Now()
	claims = &JWTClaims{
		UserID:      user_id,
		FirebaseUID: firebase_uid,
		TokenType:   token_type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "sleeve",
			Subject:   user_id,
		},
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err = token.SignedString(s.secret_key)
	if err != nil {
		return "", fmt.Errorf("%w: %w", domain_errors.ErrJWTGenerationFailed, err)
	}
	return token_string, nil
}
