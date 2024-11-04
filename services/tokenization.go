package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"ticket-booking/configs/logs"
	"ticket-booking/dtos/responses"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Tokenization interface {
	GenerateToken(id string) (*responses.TokenResponse, error)
	ValidateToken(token string) (bool, error)
	GenerateRefreshToken(key string) (string, error)
	VerifyRefreshToken(key, token string) (bool, error)
	GetAccountID(token string) (uuid.UUID, error)
}

type tokenization struct {
	secret            string
	issuer            string
	audience          string
	expiry            time.Duration
	refreshTokenCache map[string]refreshTokenModel
	mu                sync.Mutex
}

type refreshTokenModel struct {
	token  string
	expiry time.Time
}

func NewTokenization() *tokenization {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret"
		logs.Warn("JWT_SECRET not set, using default secret")
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "ticket-booking"
		logs.Warn("JWT_ISSUER not set, using default issuer")
	}

	audience := os.Getenv("JWT_AUDIENCE")
	if audience == "" {
		audience = "ticket-booking"
		logs.Warn("JWT_AUDIENCE not set, using default audience")
	}

	expiryStr := os.Getenv("JWT_EXPIRY")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		expiry = time.Hour * 24
		logs.Warn("JWT_EXPIRY not set or invalid, using default expiry of 24 hours")
	}

	return &tokenization{
		secret:            secret,
		issuer:            issuer,
		audience:          audience,
		expiry:            expiry,
		refreshTokenCache: make(map[string]refreshTokenModel),
	}
}

func (t *tokenization) GenerateToken(id string) (*responses.TokenResponse, error) {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(t.expiry).Unix(),
		"aud": t.audience,
		"iss": t.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(t.secret))
	if err != nil {
		logs.Error("Error signing token", err)
		return nil, err
	}

	refreshToken, err := t.GenerateRefreshToken(id)
	if err != nil {
		logs.Error("Error generating refresh token", err)
		return nil, err
	}

	expiryDate := time.Now().Add(t.expiry)

	return responses.NewTokenResponse(tokenString, refreshToken, expiryDate), nil
}

func (t *tokenization) GenerateRefreshToken(key string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		logs.Error("Error generating salt for refresh token", err)
		return "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(salt)
	expiry := time.Now().Add(t.expiry)

	t.mu.Lock()
	t.refreshTokenCache[key] = refreshTokenModel{
		token:  refreshToken,
		expiry: expiry,
	}
	t.mu.Unlock()

	return refreshToken, nil
}

func (t *tokenization) ValidateToken(token string) (bool, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logs.Error("Unexpected signing method", fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.secret), nil
	})

	if err != nil {
		logs.Error("Error parsing token", err)
		return false, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		expiry := int64(claims["exp"].(float64))
		isValid := expiry >= time.Now().Unix()
		if !isValid {
			logs.Error("Token has expired", fmt.Errorf("expiry: %d", expiry))
		}
		return isValid, nil
	}

	logs.Error("Invalid token claims", errors.New("invalid token"))
	return false, errors.New("invalid token")
}

func (t *tokenization) VerifyRefreshToken(key, token string) (bool, error) {
	t.mu.Lock()
	model, exists := t.refreshTokenCache[key]
	t.mu.Unlock()

	if !exists {
		logs.Error("Refresh token not found", fmt.Errorf("key: %s", key))
		return false, nil
	}

	if model.token != token || model.expiry.Before(time.Now()) {
		logs.Error("Invalid or expired refresh token", fmt.Errorf("key: %s", key))
		return false, nil
	}

	t.mu.Lock()
	delete(t.refreshTokenCache, key)
	t.mu.Unlock()

	return true, nil
}

func (t *tokenization) GetAccountID(token string) (uuid.UUID, error) {
	// Parse the token with the signing method validation and secret key.
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			logs.Error("Unexpected signing method", err)
			return nil, err
		}
		return []byte(t.secret), nil
	})

	// Handle parsing errors
	if err != nil {
		logs.Error("Error parsing token", err)
		return uuid.Nil, fmt.Errorf("token parsing error: %w", err)
	}

	// Ensure claims are in the expected format and token is valid
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		logs.Error("Invalid token claims or token not valid", nil)
		return uuid.Nil, errors.New("invalid token: claims not valid")
	}

	// Extract account ID from claims and ensure it’s a string
	accountId, ok := claims["id"].(string)
	if !ok {
		err := errors.New("account ID claim missing or not a string")
		logs.Error("Invalid token claims", err)
		return uuid.Nil, err
	}

	// Convert account ID to UUID format
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		logs.Error("Invalid account ID format", err)
		return uuid.Nil, fmt.Errorf("account ID format error: %w", err)
	}

	return accountUUID, nil
}
