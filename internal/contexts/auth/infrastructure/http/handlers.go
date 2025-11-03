package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"backend-challenge-guinea/internal/contexts/auth/application/commands"
	"backend-challenge-guinea/internal/shared/infrastructure/middleware"
)

type AuthHandlers struct {
	authenticateHandler *commands.AuthenticateCommandHandler
}

func NewAuthHandlers(authenticateHandler *commands.AuthenticateCommandHandler) *AuthHandlers {
	return &AuthHandlers{
		authenticateHandler: authenticateHandler,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandlers) Login(c *gin.Context) {

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tenantID := middleware.GetTenantID(c)

	cmd := commands.AuthenticateCommand{
		Email:    req.Email,
		Password: req.Password,
		TenantID: tenantID,
	}

	response, err := h.authenticateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandlers) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/api/v1/auth")
	
	auth.Use(middleware.TenantMiddleware())
	
	auth.POST("/login", h.Login)
}