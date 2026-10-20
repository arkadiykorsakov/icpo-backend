package profile

import (
	"context"
	userDomain "icpo-backend/internal/domain/user"
	errs "icpo-backend/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*userDomain.User, error)
	Update(ctx context.Context, userID int64, user userDomain.User) (*userDomain.User, error)
	UpdatePassword(ctx context.Context, userID int64, passwordHash string) error
}

type TokenRepository interface {
	RevokeAllForUser(ctx context.Context, userID int64) error
}
type service struct {
	userRepo  UserRepository
	tokenRepo TokenRepository
}

func NewService(userRepo UserRepository, tokenRepo TokenRepository) *service {
	return &service{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *service) Profile(ctx context.Context, userID int64) (*userDomain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *service) Update(ctx context.Context, userID int64, email, name string) (*userDomain.User, error) {
	candidate := userDomain.NewUser(email, name)
	return s.userRepo.Update(ctx, userID, candidate)
}

func (s *service) UpdatePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errs.ErrWrongPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, userID, string(hash)); err != nil {
		return err
	}

	if err := s.tokenRepo.RevokeAllForUser(ctx, userID); err != nil {
		return err
	}

	return nil
}
