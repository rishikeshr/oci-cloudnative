/**
 * Copyright © 2020, Oracle and/or its affiliates. All rights reserved.
 * The Universal Permissive License (UPL), Version 1.0
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"mushop/events"

	"github.com/go-kit/kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"

	// OpenTelemetry imports
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	ServiceName = "events"
)

// initOpenTelemetry initializes the OpenTelemetry SDK with OTLP exporter
func initOpenTelemetry(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	// Get OTLP endpoint from environment or use default
	otlpEndpoint := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://otel-collector:4318")

	// Create OTLP HTTP exporter
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(strings.TrimPrefix(strings.TrimPrefix(otlpEndpoint, "http://"), "https://")),
		otlptracehttp.WithInsecure(), // Use insecure for internal communication
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceNamespaceKey.String("mushop"),
			semconv.DeploymentEnvironmentKey.String(getEnv("DEPLOYMENT_ENVIRONMENT", "production")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	return tp, nil
}

// getEnv reads an environment variable value and returns a default value if environment variable does not exist
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func main() {
	var (
		port = flag.String("port", "8080", "Port to bind HTTP listener")
		zipk = flag.String("zipkin", os.Getenv("ZIPKIN"), "Zipkin address")
	)
	flag.Parse()

	// Log domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// context and error channel
	errc := make(chan error)
	ctx := context.Background()

	// Initialize OpenTelemetry
	tp, err := initOpenTelemetry(ctx, fmt.Sprintf("mushop-%s", ServiceName))
	if err != nil {
		logger.Log("err", "Failed to initialize OpenTelemetry", "error", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(shutdownCtx); err != nil {
			logger.Log("err", "Failed to shutdown tracer provider", "error", err)
		}
	}()
	logger.Log("msg", "OpenTelemetry initialized", "service", ServiceName)

	// Legacy Zipkin tracer for compatibility with existing middleware
	var tracer stdopentracing.Tracer
	{
		if *zipk == "" {
			tracer = stdopentracing.NoopTracer{}
		} else {
			// Find service local IP.
			conn, err := net.Dial("udp", "8.8.8.8:80")
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
			localAddr := conn.LocalAddr().(*net.UDPAddr)
			host := strings.Split(localAddr.String(), ":")[0]
			defer conn.Close()
			logger := log.With(logger, "tracer", "Zipkin")
			logger.Log("addr", zipk)
			collector, err := zipkin.NewHTTPCollector(
				*zipk,
				zipkin.HTTPLogger(logger),
			)
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, false, fmt.Sprintf("%v:%v", host, port), ServiceName),
			)
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
		}
		stdopentracing.InitGlobalTracer(tracer)
	}

	// Wire up the service with Kafka
	handler, logger := events.WireUp(ctx, tracer, ServiceName)

	// Wrap with OpenTelemetry HTTP instrumentation
	otelHandler := otelhttp.NewHandler(handler, "mushop-events",
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		}),
	)

	// Create and launch the HTTP server.
	go func() {
		logger.Log("transport", "HTTP", "port", *port)
		errc <- http.ListenAndServe(":"+*port, otelHandler)
	}()

	// Capture interrupts.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("exit", <-errc)
}
