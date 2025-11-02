package domain

import (
	"time"

	"github.com/google/uuid"
	vo "backend-challenge-guinea/internal/shared/domain/value_objects"
)

type User struct {
	id          string      
	name        string      
	email       vo.Email    
	password    vo.Password 
	displayName *string     
	tenantID    string      
	createdAt   time.Time   
	updatedAt   time.Time   
}


func NewUser(name string, email vo.Email, password vo.Password, tenantID string, displayName *string) (*User, error) {

	if name == "" {
		return nil, ErrInvalidUserName
	}

	now := time.Now().UTC()
	
	return &User{
		id:          uuid.New().String(), 
		name:        name,
		email:       email,
		password:    password,
		displayName: displayName,
		tenantID:    tenantID,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func Reconstitute(id, name string, email vo.Email, password vo.Password, tenantID string, displayName *string, createdAt, updatedAt time.Time) *User {
	return &User{
		id:          id,
		name:        name,
		email:       email,
		password:    password,
		displayName: displayName,
		tenantID:    tenantID,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (u *User) ID() string            { return u.id }
func (u *User) Name() string          { return u.name }
func (u *User) Email() vo.Email       { return u.email }
func (u *User) Password() vo.Password { return u.password }
func (u *User) DisplayName() *string  { return u.displayName }
func (u *User) TenantID() string      { return u.tenantID }
func (u *User) CreatedAt() time.Time  { return u.createdAt }
func (u *User) UpdatedAt() time.Time  { return u.updatedAt }