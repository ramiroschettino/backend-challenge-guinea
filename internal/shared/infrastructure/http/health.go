package http

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandlers struct {
	db *sql.DB
}

func NewHealthHandlers(db *sql.DB) *HealthHandlers {
	return &HealthHandlers{db: db}
}

func (h *HealthHandlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *HealthHandlers) Ready(c *gin.Context) {

	if err := h.db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"error":  "database not available",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

func (h *HealthHandlers) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.Health)
	router.GET("/ready", h.Ready)
}