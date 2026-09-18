package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/victoryus84/gorders/internal/config"
	"github.com/victoryus84/gorders/internal/models"
	"github.com/victoryus84/gorders/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Interfața publică a serviciului
type UserService interface {
	Signup(email, password, role string) error
	Login(email, password string) (string, error)
}

// Structura privată cu câmpurile necesare
type userService struct {
	rep repository.UserRepository // l-am numit 'repo' pentru că așa îl folosești în metode
	cfg  *config.Config            // am adăugat cfg aici ca să poată fi salvat
}

// Constructorul
func NewUserService(rep repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		rep: rep,
		cfg:  cfg,
	}
}

// METODELE - toate folosesc receiver-ul (s *userService) cu "u" mic!

func (svc *userService) Signup(email, password, role string) error {
	if !svc.cfg.AllowSignup {
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
	
	// Folosim svc.rep
	return svc.rep.CreateUser(user)
}

func (svc *userService) Login(email, password string) (string, error) {
	user, err := svc.rep.FindUserByEmail(email)
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

	// Am presupus că JWTSecret se află în config-ul tău.
	// Dacă se numește altfel în config.Config, modifică 'JWTSecret' cu numele real.
	return token.SignedString([]byte(svc.cfg.JWTSecret))
}