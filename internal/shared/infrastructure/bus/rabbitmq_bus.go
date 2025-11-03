package bus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBus struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string                     
	handlers map[string][]EventHandler   
	log      Logger
}

func NewRabbitMQBus(url, exchange string, log Logger) (*RabbitMQBus, error) {

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = channel.ExchangeDeclare(
		exchange,
		"topic", 
		true,    
		false,   
		false,   
		false,   
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQBus{
		conn:     conn,
		channel:  channel,
		exchange: exchange,
		handlers: make(map[string][]EventHandler),
		log:      log,
	}, nil
}


func (b *RabbitMQBus) Publish(ctx context.Context, event interface{}) error {

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	eventType := extractEventType(event)

	correlationID := extractCorrelationID(ctx)

	err = b.channel.PublishWithContext(
		ctx,
		b.exchange,
		eventType,  
		false,      
		false,      
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			DeliveryMode:  amqp.Persistent, 
			Timestamp:     time.Now(),
			CorrelationId: correlationID,
		},
	)

	if err != nil {
		b.log.Error("failed to publish event", map[string]interface{}{
			"error":          err.Error(),
			"event_type":     eventType,
			"correlation_id": correlationID,
		})
		return err
	}

	b.log.Info("event published", map[string]interface{}{
		"event_type":     eventType,
		"correlation_id": correlationID,
	})

	return nil
}

func (b *RabbitMQBus) Subscribe(eventType string, handler EventHandler) error {
	if b.handlers[eventType] == nil {
		b.handlers[eventType] = make([]EventHandler, 0)
	}
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

func (b *RabbitMQBus) Start(ctx context.Context) error {
	for eventType := range b.handlers {
		queueName := fmt.Sprintf("%s_queue", eventType)
		
		queue, err := b.channel.QueueDeclare(
			queueName,
			true,  
			false, 
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue: %w", err)
		}

		err = b.channel.QueueBind(
			queue.Name,
			eventType,
			b.exchange,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue: %w", err)
		}


		msgs, err := b.channel.Consume(
			queue.Name,
			"",    
			false, 
			false, 
			false, 
			false, 
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to register consumer: %w", err)
		}

		go b.handleMessages(ctx, eventType, msgs)
		
		b.log.Info("started consuming", map[string]interface{}{
			"event_type": eventType,
			"queue":      queue.Name,
		})
	}

	return nil
}

func (b *RabbitMQBus) handleMessages(ctx context.Context, eventType string, msgs <-chan amqp.Delivery) {
	for msg := range msgs {

		msgCtx := context.WithValue(ctx, "correlation_id", msg.CorrelationId)

		var event map[string]interface{}
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			b.log.Error("failed to unmarshal event", map[string]interface{}{
				"error":          err.Error(),
				"correlation_id": msg.CorrelationId,
			})
			msg.Nack(false, false) 
			continue
		}

		handlers := b.handlers[eventType]
		success := true
		
		for _, handler := range handlers {
			if err := handler(msgCtx, event); err != nil {
				b.log.Error("handler failed", map[string]interface{}{
					"error":          err.Error(),
					"event_type":     eventType,
					"correlation_id": msg.CorrelationId,
				})
				success = false
				break
			}
		}


		if success {
			msg.Ack(false) 
			b.log.Debug("message processed", map[string]interface{}{
				"event_type":     eventType,
				"correlation_id": msg.CorrelationId,
			})
		} else {
			msg.Nack(false, true) 
		}
	}
}

func (b *RabbitMQBus) Close() error {
	if err := b.channel.Close(); err != nil {
		return err
	}
	return b.conn.Close()
}

func extractEventType(event interface{}) string {
	if e, ok := event.(interface{ EventType() string }); ok {
		return e.EventType()
	}
	return "unknown"
}

func extractCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value("correlation_id").(string); ok {
		return id
	}
	return ""
}

type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
}