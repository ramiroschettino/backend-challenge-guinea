package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//generamos o extraemos correlation ID para trazabilidad
//permite seguir una request a través de toda la arquitectura
func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//intentamos get del header
		correlationID := c.GetHeader("X-Correlation-Id")
		
		// si no existe, generamos uno
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// guardamos en contexto y agregamos a la respuesta
		c.Set("correlation_id", correlationID)
		c.Header("X-Correlation-Id", correlationID)
		
		c.Next()
	}
}

// obtiene el correlation ID del contexto
func GetCorrelationID(c *gin.Context) string {
	if correlationID, exists := c.Get("correlation_id"); exists {
		return correlationID.(string)
	}
	return ""
}