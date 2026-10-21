package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"icpo-backend/internal/domain/token/dto"
	errs "icpo-backend/internal/errors"
)

const (
	issuer = "icpo-backend"
	// clockSkew — допустимое расхождение часов между инстансами при проверке exp/nbf.
	clockSkew = 5 * time.Second
)

type Config struct {
	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type Repository interface {
	Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Rotate(ctx context.Context, oldHash string, userID int64, newHash string, newExpiresAt, now time.Time) (bool, error)
	Revoke(ctx context.Context, tokenHash string) (bool, error)
	RevokeAllForUser(ctx context.Context, userID int64) error
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type service struct {
	config     Config
	repository Repository
	now        func() time.Time
}

func NewService(repository Repository, config Config) *service {
	return &service{repository: repository, config: config, now: time.Now}
}

func (s *service) generateAccessToken(userID int64, now time.Time) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.AccessTTL)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(s.config.JWTSecret))
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *service) newPair(access, refresh string) *dto.TokenPair {
	return &dto.TokenPair{
		AccessToken:      access,
		TokenType:        "Bearer",
		ExpiresIn:        int(s.config.AccessTTL.Seconds()),
		RefreshToken:     refresh,
		RefreshExpiresIn: int(s.config.RefreshTTL.Seconds()),
	}
}

func (s *service) IssuePair(ctx context.Context, userID int64) (*dto.TokenPair, error) {
	now := s.now()

	access, err := s.generateAccessToken(userID, now)
	if err != nil {
		return nil, err
	}
	refresh, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.repository.Create(ctx, userID, hashToken(refresh), now.Add(s.config.RefreshTTL)); err != nil {
		return nil, err
	}
	return s.newPair(access, refresh), nil
}

func (s *service) ParseAccessToken(tokenString string) (*Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims,
		func(*jwt.Token) (any, error) { return []byte(s.config.JWTSecret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(issuer),
		jwt.WithLeeway(clockSkew),
		jwt.WithTimeFunc(s.now),
	)
	if err != nil {
		// Причина (истёк, неверная подпись и т.д.) остаётся в цепочке ошибок для логов,
		// а errors.Is(err, ErrInvalidToken) продолжает работать.
		return nil, fmt.Errorf("%w: %v", errs.ErrInvalidToken, err)
	}
	if claims.UserID <= 0 {
		return nil, errs.ErrInvalidToken
	}

	return &claims, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshTokenString string) (*dto.TokenPair, error) {
	now := s.now()
	oldHash := hashToken(refreshTokenString)

	stored, err := s.repository.GetByHash(ctx, oldHash)
	if err != nil {
		return nil, err
	}

	if stored.RevokedAt != nil {
		if err := s.repository.RevokeAllForUser(ctx, stored.UserID); err != nil {
			return nil, err
		}
		return nil, errs.ErrInvalidToken
	}

	if !now.Before(stored.ExpiresAt) {
		return nil, errs.ErrInvalidToken
	}

	access, err := s.generateAccessToken(stored.UserID, now)
	if err != nil {
		return nil, err
	}
	newRefresh, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	rotated, err := s.repository.Rotate(ctx, oldHash, stored.UserID,
		hashToken(newRefresh), now.Add(s.config.RefreshTTL), now)
	if err != nil {
		return nil, err
	}
	if !rotated {
		return nil, errs.ErrInvalidToken
	}

	return s.newPair(access, newRefresh), nil
}

func (s *service) Logout(ctx context.Context, refreshTokenString string) error {
	_, err := s.repository.Revoke(ctx, hashToken(refreshTokenString))
	return err
}
