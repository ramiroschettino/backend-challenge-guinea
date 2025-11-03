package middleware

import (
	"sync"
)

// feature flags por tenant. Habilitar/deshabilitar funcionalidades sin deployar. tenant_id -> feature -> enabled
type FeatureFlags struct {
	flags map[string]map[string]bool
	mu    sync.RWMutex
}

//crea gestor de feature flags
func NewFeatureFlags() *FeatureFlags {
	ff := &FeatureFlags{
		flags: make(map[string]map[string]bool),
	}
	
	// config inicial (en producción vendría de DB)
	ff.SetFeature("tenant-1", "user_display_name", true)
	ff.SetFeature("tenant-2", "user_display_name", false)
	
	return ff
}

// IsEnabled verifica si un feature está habilitado para un tenant
func (ff *FeatureFlags) IsEnabled(tenantID, feature string) bool {
	ff.mu.RLock()
	defer ff.mu.RUnlock()

	if tenantFlags, ok := ff.flags[tenantID]; ok {
		return tenantFlags[feature]
	}

	// Por defecto, todos los features están habilitados
	return true
}

// SetFeature configura un feature flag
func (ff *FeatureFlags) SetFeature(tenantID, feature string, enabled bool) {
	ff.mu.Lock()
	defer ff.mu.Unlock()

	if ff.flags[tenantID] == nil {
		ff.flags[tenantID] = make(map[string]bool)
	}

	ff.flags[tenantID][feature] = enabled
}