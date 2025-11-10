/**
 * Copyright © 2020, Oracle and/or its affiliates. All rights reserved.
 * The Universal Permissive License (UPL), Version 1.0
 */
package events

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

// Middleware decorates a service.
type Middleware func(Service) Service

type Service interface {
	PostEvents(source string, track string, events []Event) (EventsReceived, error) // POST /events
	Health() []Health                                                               // GET /health
}

type EventsReceived struct {
	Success bool `json:"success"`
	Length  int  `json:"events"`
}

type Event struct {
	Time   string      `json:"time"`
	Type   string      `json:"type"`
	Detail interface{} `json:"detail"`
}

type EventRecord struct {
	Event
	Source string `json:"source"`
	Track  string `json:"track"`
}

type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

// NewEventsService returns a simple implementation of the Service interface
func NewEventsService(
	ctx context.Context,
	producer sarama.SyncProducer,
	topic string,
	logger log.Logger) Service {

	return &service{
		ctx:      ctx,
		producer: producer,
		topic:    topic,
		logger:   logger,
	}
}

type service struct {
	ctx      context.Context
	producer sarama.SyncProducer
	topic    string
	logger   log.Logger
}

func (s *service) PostEvents(source string, track string, events []Event) (EventsReceived, error) {

	numEvents := len(events)
	s.logger.Log(
		"source", source,
		"track", track,
		"length", numEvents,
	)

	var err error
	success := false

	if numEvents == 0 {
		err = errors.New("no events received")
		return EventsReceived{
			Success: success,
			Length:  numEvents,
		}, err
	}

	// Send messages to Kafka
	var messages []*sarama.ProducerMessage

	for _, evt := range events {
		msg := EventRecord{
			Source: source,
			Track:  track,
		}
		msg.Time = evt.Time
		msg.Type = evt.Type
		msg.Detail = evt.Detail

		data, e := json.Marshal(msg)
		if e == nil {
			// Create Kafka message
			messages = append(messages, &sarama.ProducerMessage{
				Topic: s.topic,
				Value: sarama.ByteEncoder(data),
			})
		}
	}

	// Send all messages
	successCount := 0
	for _, msg := range messages {
		_, _, err = s.producer.SendMessage(msg)
		if err != nil {
			s.logger.Log("error", "failed to send message", "err", err)
		} else {
			successCount++
		}
	}

	if successCount == len(messages) {
		success = true
		err = nil
	} else if successCount > 0 {
		success = true
		err = errors.New("some messages failed to send")
	} else {
		err = errors.New("all messages failed to send")
	}

	return EventsReceived{
		Success: success,
		Length:  numEvents,
	}, err
}

func (s *service) Health() []Health {
	var health []Health
	app := Health{"events", "OK", time.Now().String()}

	// Check Kafka connection
	kafkaStatus := "OK"
	if s.producer == nil {
		kafkaStatus = "ERROR"
	}
	kafka := Health{"kafka:events-stream", kafkaStatus, time.Now().String()}

	health = append(health, app)
	health = append(health, kafka)
	return health
}
