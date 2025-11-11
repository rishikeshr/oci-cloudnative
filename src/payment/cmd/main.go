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

	"github.com/go-kit/kit/log"
	"github.com/junior/mushop/src/payment"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"

	// OpenTelemetry imports
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	ServiceName = "payment"
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
		port          = flag.String("port", "8080", "Port to bind HTTP listener")
		zip           = flag.String("zipkin", os.Getenv("ZIPKIN"), "Zipkin address")
		declineAmount = flag.Float64("decline", 105, "Decline payments over certain amount")
	)
	flag.Parse()

	// Log domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Mechanical stuff.
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
		if *zip == "" {
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
			logger.Log("addr", zip)
			collector, err := zipkin.NewHTTPCollector(
				*zip,
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

	handler, logger := payment.WireUp(ctx, float32(*declineAmount), tracer, ServiceName)

	// Wrap with OpenTelemetry HTTP instrumentation
	otelHandler := otelhttp.NewHandler(handler, "mushop-payment",
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
