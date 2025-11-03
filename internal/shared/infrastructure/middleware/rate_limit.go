package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// rate limiting por tenant
// usa un algoritmo de token bucket simple en memoria
type RateLimiter struct {
	requests int           // número de requests permitidos
	window   time.Duration // ventana de tiempo
	buckets  map[string]*bucket
	mu       sync.RWMutex
}

type bucket struct {
	count   int
	resetAt time.Time
	mu      sync.Mutex
}

// crea un rate limiter. ejemplo: 100 requests por minuto

func NewRateLimiter(requests int, window time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		requests: requests,
		window:   window,
		buckets:  make(map[string]*bucket),
	}

	// goroutine para limpiar buckets viejos
	go limiter.cleanup()

	return limiter
}

// Middleware retorna el middleware de Gin
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := GetTenantID(c)
		if tenantID == "" {
			c.Next()
			return
		}

		// Clave: tenant_id + path (ej: "tenant-1:/users")
		key := fmt.Sprintf("%s:%s", tenantID, c.FullPath())
		
		if !rl.allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// allow verifica si la request puede pasar
func (rl *RateLimiter) allow(key string) bool {
	rl.mu.RLock()
	b, exists := rl.buckets[key]
	rl.mu.RUnlock()

	// Si no existe el bucket, crearlo
	if !exists {
		rl.mu.Lock()
		b = &bucket{
			count:   0,
			resetAt: time.Now().Add(rl.window),
		}
		rl.buckets[key] = b
		rl.mu.Unlock()
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	
	// Si pasó la ventana, resetear contador
	if now.After(b.resetAt) {
		b.count = 0
		b.resetAt = now.Add(rl.window)
	}

	// Verificar si llegó al límite
	if b.count >= rl.requests {
		return false
	}

	b.count++
	return true
}

// cleanup limpia buckets viejos cada 5 minutos
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, b := range rl.buckets {
			b.mu.Lock()
			if now.After(b.resetAt.Add(rl.window)) {
				delete(rl.buckets, key)
			}
			b.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}