package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	errs "icpo-backend/internal/errors"
	tokenDto "icpo-backend/internal/domain/token/dto"
	userDomain "icpo-backend/internal/domain/user"
	"icpo-backend/internal/utils"
)

type UserRepository interface {
	Create(ctx context.Context, email, name, passwordHash string) (*userDomain.User, error)
	GetByEmail(ctx context.Context, email string) (*userDomain.User, error)
}

type TokenService interface {
	IssuePair(ctx context.Context, userID int64) (*tokenDto.TokenPair, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*tokenDto.TokenPair, error)
	Logout(ctx context.Context, refreshTokenString string) error
}

type service struct {
	userRepo UserRepository
	tokenSvc TokenService
}

func NewService(userRepo UserRepository, tokenSvc TokenService) *service {
	return &service{userRepo: userRepo, tokenSvc: tokenSvc}
}

func (s *service) Register(ctx context.Context, email, name, password string) (*userDomain.User, error) {
	email = utils.NormalizeEmail(email)

	if existing, err := s.userRepo.GetByEmail(ctx, email); err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return nil, err
	} else if existing != nil {
		return nil, errs.ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.userRepo.Create(ctx, email, name, string(hash))
}

func (s *service) Login(ctx context.Context, email, password string) (*tokenDto.TokenPair, error) {
	email = utils.NormalizeEmail(email)

	user, err := s.userRepo.GetByEmail(ctx, email)
	if errors.Is(err, errs.ErrUserNotFound) {
		return nil, errs.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	return s.tokenSvc.IssuePair(ctx, user.ID)
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*tokenDto.TokenPair, error) {
	return s.tokenSvc.RefreshToken(ctx, refreshToken)
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenSvc.Logout(ctx, refreshToken)
}
