package services

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gofiber-expense-tracker/internal/models"
)

// mockUserRepo implements UserRepo for testing
type mockUserRepo struct {
	users map[string]*models.User
	nextID int64
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:  make(map[string]*models.User),
		nextID: 1,
	}
}

func (m *mockUserRepo) Create(user *models.User) error {
	user.ID = m.nextID
	m.users[user.Email] = user
	m.nextID++
	return nil
}

func (m *mockUserRepo) GetByEmail(email string) (*models.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func newTestAuthService() *AuthService {
	return NewAuthService(newMockUserRepo(), "test-secret", 24)
}

func TestAuthService_Register(t *testing.T) {
	svc := newTestAuthService()

	t.Run("success", func(t *testing.T) {
		resp, err := svc.Register(models.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, "test@example.com", resp.User.Email)
		assert.NotEmpty(t, resp.User.PasswordHash) // bcrypt hash stored in DB, hidden from JSON via json:"-"
	})

	t.Run("duplicate email", func(t *testing.T) {
		_, err := svc.Register(models.RegisterRequest{Email: "dup@example.com", Password: "pass"})
		assert.NoError(t, err)
		_, err = svc.Register(models.RegisterRequest{Email: "dup@example.com", Password: "pass2"})
		assert.ErrorIs(t, err, ErrEmailTaken)
	})

	t.Run("empty fields", func(t *testing.T) {
		_, err := svc.Register(models.RegisterRequest{Email: "", Password: ""})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})
}

func TestAuthService_Login(t *testing.T) {
	svc := newTestAuthService()

	// Seed a user
	_, err := svc.Register(models.RegisterRequest{
		Email:    "user@example.com",
		Password: "secret123",
	})
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		resp, err := svc.Login(models.LoginRequest{
			Email:    "user@example.com",
			Password: "secret123",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("wrong password", func(t *testing.T) {
		_, err := svc.Login(models.LoginRequest{
			Email:    "user@example.com",
			Password: "wrongpassword",
		})
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := svc.Login(models.LoginRequest{
			Email:    "nobody@example.com",
			Password: "secret123",
		})
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("empty fields", func(t *testing.T) {
		_, err := svc.Login(models.LoginRequest{Email: "", Password: ""})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})
}
