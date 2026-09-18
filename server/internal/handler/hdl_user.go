package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/service"
)

// UserService definește doar metodele de care are nevoie acest handler.
// Asta ajută la decuplare (Interface Segregation).
type UserService interface {
	Signup(email, password, role string) error
	Login(email, password string) (string, error)
}

type UserHandler struct {
	service service.UserService
}

// NewUserHandler creează o instanță nouă a handler-ului
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

// Signup gestionează înregistrarea utilizatorilor
func (hdl *UserHandler) Signup(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := hdl.service.Signup(req.Email, req.Password, req.Role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User created"})
}

// Login gestionează autentificarea
func (hdl *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email" xml:"email" binding:"required"`
		Password string `json:"password" xml:"password" binding:"required"`
	}

	// Folosim ParseBody (presupunând că e un utilitar global în pkg/utils sau similar)
	// Dacă nu ai ParseBody definit încă aici, poți folosi c.ShouldBind
	var req LoginReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Date invalide"})
		return
	}

	token, err := hdl.service.Login(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func RegisterUserRoutes(rte *gin.Engine, hdl *UserHandler) {
	rte.POST("/signup", hdl.Signup)
	rte.POST("/login", hdl.Login)
}