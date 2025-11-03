package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	id        string
	userID    string
	tenantID  string
	token     string
	expiresAt time.Time
	createdAt time.Time
}

func NewSession(userID, tenantID string, duration time.Duration) *Session {
	now := time.Now().UTC()
	return &Session{
		id:        uuid.New().String(),
		userID:    userID,
		tenantID:  tenantID,
		token:     generateToken(), 
		expiresAt: now.Add(duration),
		createdAt: now,
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().UTC().After(s.expiresAt)
}

func (s *Session) ID() string           { return s.id }
func (s *Session) UserID() string       { return s.userID }
func (s *Session) TenantID() string     { return s.tenantID }
func (s *Session) Token() string        { return s.token }
func (s *Session) ExpiresAt() time.Time { return s.expiresAt }

func generateToken() string {
	return uuid.New().String()
}