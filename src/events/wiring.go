/**
 * Copyright © 2020, Oracle and/or its affiliates. All rights reserved.
 * The Universal Permissive License (UPL), Version 1.0
 */
package events

import (
	"log"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
	kitlog "github.com/go-kit/kit/log"
	"golang.org/x/net/context"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Time (in seconds) spent serving HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "status_code", "isWS"})
)

func init() {
	prometheus.MustRegister(HTTPLatency)
}

// WireUp the service to the provided context
func WireUp(
	ctx context.Context,
	tracer stdopentracing.Tracer,
	serviceName string) (http.Handler, kitlog.Logger) {
	// Logging
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	// Kafka configurations
	config := GetKafkaConfig()
	brokers := GetKafkaBrokers()
	topic := GetKafkaTopic()

	// Create Kafka producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	logger.Log("msg", "Connected to Kafka", "brokers", brokers, "topic", topic)

	// Service domain
	var service Service
	{
		service = NewEventsService(ctx, producer, topic, logger)
		service = LoggingMiddleware(logger)(service)
	}

	// Endpoint domain
	endpoints := MakeEndpoints(service, tracer)

	router := MakeHTTPHandler(ctx, endpoints, logger, tracer)

	return router, logger
}
