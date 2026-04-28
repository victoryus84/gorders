package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// UserService definește doar metodele de care are nevoie acest handler.
// Asta ajută la decuplare (Interface Segregation).
type UserService interface {
	Signup(email, password, role string) error
	Login(email, password string) (string, error)
}

type UserHandler struct {
	service UserService
}

// NewUserHandler creează o instanță nouă a handler-ului
func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{service: s}
}

// RegisterRoutes înregistrează rutele de autentificare
func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/signup", h.Signup)
	router.POST("/login", h.Login)
}

// Signup gestionează înregistrarea utilizatorilor
func (h *UserHandler) Signup(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Signup(req.Email, req.Password, req.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

// Login gestionează autentificarea
func (h *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email" xml:"email" binding:"required"`
		Password string `json:"password" xml:"password" binding:"required"`
	}

	// Folosim ParseBody (presupunând că e un utilitar global în pkg/utils sau similar)
	// Dacă nu ai ParseBody definit încă aici, poți folosi c.ShouldBind
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date invalide"})
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}