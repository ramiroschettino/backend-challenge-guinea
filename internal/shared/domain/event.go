package domain

import (
	"time"

	"github.com/google/uuid"
)

type DomainEvent interface {
	EventID() string        
	EventType() string      
	OccurredOn() time.Time  
	AggregateID() string    
	TenantID() string       
	CorrelationID() string 
}

type BaseEvent struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	AggregateId   string    `json:"aggregate_id"`
	TenantId      string    `json:"tenant_id"`
	CorrelationId string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
}

func NewBaseEvent(eventType, aggregateID, tenantID, correlationID string) BaseEvent {
	return BaseEvent{
		ID:            uuid.New().String(),
		Type:          eventType,
		AggregateId:   aggregateID,
		TenantId:      tenantID,
		CorrelationId: correlationID,
		Timestamp:     time.Now().UTC(),
	}
}

func (e BaseEvent) EventID() string       { return e.ID }
func (e BaseEvent) EventType() string     { return e.Type }
func (e BaseEvent) OccurredOn() time.Time { return e.Timestamp }
func (e BaseEvent) AggregateID() string   { return e.AggregateId }
func (e BaseEvent) TenantID() string      { return e.TenantId }
func (e BaseEvent) CorrelationID() string { return e.CorrelationId }