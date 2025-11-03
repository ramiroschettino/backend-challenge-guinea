package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"backend-challenge-guinea/internal/contexts/users/application/commands"
	"backend-challenge-guinea/internal/contexts/users/application/queries"
	"backend-challenge-guinea/internal/shared/infrastructure/middleware"
)

// Agrupa todos los handlers de usuarios
type UserHandlers struct {
	createUserHandler *commands.CreateUserCommandHandler
	getUserHandler    *queries.GetUserQueryHandler
	featureFlags      *middleware.FeatureFlags
}

func NewUserHandlers(
	createUserHandler *commands.CreateUserCommandHandler,
	getUserHandler *queries.GetUserQueryHandler,
	featureFlags *middleware.FeatureFlags,
) *UserHandlers {
	return &UserHandlers{
		createUserHandler: createUserHandler,
		getUserHandler:    getUserHandler,
		featureFlags:      featureFlags,
	}
}

type CreateUserRequest struct {
	Name        string  `json:"name" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required"`
	DisplayName *string `json:"display_name,omitempty"`
}

type CreateUserResponse struct {
	ID            string `json:"id"`
	CorrelationID string `json:"correlation_id"`
}

func (h *UserHandlers) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tenantID := middleware.GetTenantID(c)
	correlationID := middleware.GetCorrelationID(c)

	if req.DisplayName != nil && !h.featureFlags.IsEnabled(tenantID, "user_display_name") {
		req.DisplayName = nil
	}

	idempotencyKey := c.GetHeader("X-Idempotency-Key")

	cmd := commands.CreateUserCommand{
		Name:           req.Name,
		Email:          req.Email,
		Password:       req.Password,
		DisplayName:    req.DisplayName,
		TenantID:       tenantID,
		CorrelationID:  correlationID,
		IdempotencyKey: idempotencyKey,
	}

	userID, err := h.createUserHandler.Handle(c.Request.Context(), cmd)
	if err != nil {

		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "user with this email already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, CreateUserResponse{
		ID:            userID,
		CorrelationID: correlationID,
	})
}

func (h *UserHandlers) GetUser(c *gin.Context) {

	userID := c.Param("id")

	tenantID := middleware.GetTenantID(c)

	query := queries.GetUserQuery{
		UserID:   userID,
		TenantID: tenantID,
	}

	user, err := h.getUserHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// registra las rutas en el router de Gin
func (h *UserHandlers) RegisterRoutes(router *gin.Engine, rateLimiter *middleware.RateLimiter) {

	users := router.Group("/api/v1/users")
	
	users.Use(middleware.TenantMiddleware())
	users.Use(middleware.CorrelationIDMiddleware())
	
	// Rutas
	users.POST("", rateLimiter.Middleware(), h.CreateUser)  
	users.GET("/:id", h.GetUser)                             
}