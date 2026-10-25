package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	tokenDto "icpo-backend/internal/domain/token/dto"
	userDomain "icpo-backend/internal/domain/user"
	errs "icpo-backend/internal/errors"
)

// dummyHash — валидный bcrypt-хеш (cost 10 = bcrypt.DefaultCost) случайного пароля.
// Используется в Login для несуществующего пользователя, чтобы время ответа не
// отличалось от случая «пользователь есть, пароль неверный».
const dummyHash = "$2b$10$GcGeDzte1Kj4pJA.3klmsuSbncEYzdX1PMfKtQQ6JW90Elh29VDAW"

type UserRepository interface {
	Create(ctx context.Context, user userDomain.User, passwordHash string) (userDomain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (userDomain.User, error)
}

type TokenService interface {
	IssuePair(ctx context.Context, userID int64) (tokenDto.TokenPair, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (tokenDto.TokenPair, error)
	Logout(ctx context.Context, refreshTokenString string) error
}

type service struct {
	userRepo UserRepository
	tokenSvc TokenService
}

func NewService(userRepo UserRepository, tokenSvc TokenService) *service {
	return &service{userRepo: userRepo, tokenSvc: tokenSvc}
}

func (s *service) Register(ctx context.Context, email, name, password string) (userDomain.User, error) {
	candidate := userDomain.NewUser(email, name)

	exists, err := s.userRepo.ExistsByEmail(ctx, candidate.Email)
	if err != nil {
		return userDomain.User{}, err
	}
	if exists {
		return userDomain.User{}, errs.ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return userDomain.User{}, err
	}

	return s.userRepo.Create(ctx, candidate, string(hash))
}

func (s *service) Login(ctx context.Context, email, password string) (tokenDto.TokenPair, error) {
	email = userDomain.NormalizeEmail(email)

	user, err := s.userRepo.GetByEmail(ctx, email)
	if errors.Is(err, errs.ErrUserNotFound) {
		_ = bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(password))
		return tokenDto.TokenPair{}, errs.ErrInvalidCredentials
	}
	if err != nil {
		return tokenDto.TokenPair{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return tokenDto.TokenPair{}, errs.ErrInvalidCredentials
	}

	return s.tokenSvc.IssuePair(ctx, user.ID)
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (tokenDto.TokenPair, error) {
	return s.tokenSvc.RefreshToken(ctx, refreshToken)
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenSvc.Logout(ctx, refreshToken)
}
