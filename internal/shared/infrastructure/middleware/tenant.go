package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TenantMiddleware extrae el tenant ID del header X-Tenant-Id
// Todas las operaciones deben incluir el tenant para multi-tenancy
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-Id")
		
		// Si no hay tenant ID, rechazamos la request
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "X-Tenant-Id header is required",
			})
			c.Abort()
			return
		}

		c.Set("tenant_id", tenantID)
		c.Next()
	}
}

// tenant id del context
func GetTenantID(c *gin.Context) string {
	if tenantID, exists := c.Get("tenant_id"); exists {
		return tenantID.(string)
	}
	return ""
}