package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"gofiber-expense-tracker/internal/models"
)

var (
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserRepo interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}

type AuthService struct {
	repo        UserRepo
	jwtSecret   string
	expireHours int
}

func NewAuthService(repo UserRepo, jwtSecret string, expireHours int) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret, expireHours: expireHours}
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, ErrInvalidInput
	}

	existing, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) generateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(s.expireHours) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
