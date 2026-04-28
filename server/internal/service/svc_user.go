package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/victoryus84/gorders/internal/config"
	"github.com/victoryus84/gorders/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
}

type UserService struct {
	repo      UserRepository
	cfg       *config.Config
	jwtSecret string
}

func NewUserService(repo UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		repo:      repo,
		cfg:       cfg,
		jwtSecret: cfg.JWTSecret,
	}
}

func (s *UserService) Signup(email, password, role string) error {
	if !s.cfg.AllowSignup {
		return fmt.Errorf("înregistrarea utilizatorilor este dezactivată")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     map[bool]string{true: "admin", false: "user"}[strings.ToLower(role) == "trueadmin"],
	}
	return s.repo.CreateUser(user)
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("parolă incorectă")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}