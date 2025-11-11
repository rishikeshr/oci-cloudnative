# MuShop Observability with EDOT

This document describes the OpenTelemetry instrumentation implemented in MuShop using EDOT (Elastic Distribution of OpenTelemetry).

## Overview

MuShop is instrumented with OpenTelemetry to provide comprehensive observability through:
- **Distributed Tracing**: Track requests across all microservices
- **Metrics**: Monitor service performance and resource usage
- **Logs**: Centralized logging with trace correlation

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌──────────────┐
│   Services  │────▶│   OTel      │────▶│   Jaeger     │
│  (EDOT)     │     │  Collector  │     │  (Traces)    │
└─────────────┘     └─────────────┘     └──────────────┘
                          │
                          ├──────────────▶ Prometheus
                          │                (Metrics)
                          │
                          └──────────────▶ Logging
                                          (Future: Elastic)
```

## Instrumented Services

### Java Services (EDOT Java Agent)
- **mushop-orders** - Spring Boot application
- **mushop-carts** - Spring Data JPA application
- **mushop-fulfillment** - Micronaut application

**Configuration**:
- Agent: `elastic-otel-javaagent-1.7.0.jar`
- Auto-instruments: HTTP, JDBC, JPA, Spring Framework
- Location: Downloaded from Maven Central during Docker build

### Node.js Services (EDOT Node)
- **mushop-user** - NestJS application
- **mushop-api** - Express application
- **mushop-assets** - Static file server
- **mushop-newsletter** - Express application

**Configuration**:
- Package: `@elastic/opentelemetry-node@0.3.0`
- Auto-instruments: HTTP, Express, PostgreSQL, etc.
- Initialization: `require('./tracing')` at application startup

### Go Services (OpenTelemetry SDK)
- **mushop-catalogue** - Product catalog service
- **mushop-payment** - Payment processing service
- **mushop-events** - Event streaming service

**Status**: OpenTelemetry packages added to go.mod
**Note**: Manual instrumentation required (future enhancement)

## Configuration

### Environment Variables

All services are configured via standard OpenTelemetry environment variables:

```yaml
# Service identification
OTEL_SERVICE_NAME: mushop-<service-name>

# OTLP Exporter configuration
OTEL_EXPORTER_OTLP_ENDPOINT: http://otel-collector:4318
OTEL_EXPORTER_OTLP_PROTOCOL: http/protobuf

# Signal exporters
OTEL_TRACES_EXPORTER: otlp
OTEL_METRICS_EXPORTER: otlp
OTEL_LOGS_EXPORTER: otlp

# Resource attributes
OTEL_RESOURCE_ATTRIBUTES: deployment.environment=production,service.version=<version>
```

### OpenTelemetry Collector

The OTel Collector acts as a central telemetry hub:
- **Receives**: OTLP over gRPC (4317) and HTTP (4318)
- **Processes**: Batching, resource detection, memory limiting
- **Exports**: Jaeger (traces), Prometheus (metrics), Logging

Configuration file: `otel-collector-config.yaml`

## Usage

### Starting the Services

```bash
# Start all services including observability stack
docker-compose up -d

# View logs
docker-compose logs -f otel-collector
```

### Accessing Observability Tools

#### Jaeger UI (Distributed Tracing)
- URL: http://localhost:16686
- View traces, service dependencies, latency analysis
- Search by service, operation, tags, duration

#### Prometheus Metrics
- OTel Collector metrics: http://localhost:8888/metrics
- Service metrics: http://localhost:8889/metrics

#### OTel Collector Health
- Health check: http://localhost:13133
- zPages diagnostics: http://localhost:55679/debug/tracez

## Example Queries

### Jaeger Traces

1. **Find slow requests**:
   - Service: `mushop-orders`
   - Min Duration: `500ms`

2. **Trace a complete checkout flow**:
   - Service: `mushop-api`
   - Operation: `POST /orders`
   - Follow spans through: API → Carts → Orders → Fulfillment → Events

3. **Find database queries**:
   - Service: `mushop-catalogue`
   - Tags: `db.system=postgresql`

### Service Dependencies

View the Service Graph in Jaeger to see:
- Request rates between services
- Error rates
- Latency percentiles (P50, P95, P99)

## Troubleshooting

### No traces appearing in Jaeger

1. Check OTel Collector is running:
   ```bash
   docker-compose ps otel-collector
   docker-compose logs otel-collector
   ```

2. Verify services can reach the collector:
   ```bash
   docker-compose exec orders curl http://otel-collector:4318/v1/traces
   ```

3. Check service logs for EDOT initialization:
   ```bash
   docker-compose logs user | grep EDOT
   docker-compose logs orders | grep otel
   ```

### High memory usage

Adjust OTel Collector memory limit in `otel-collector-config.yaml`:
```yaml
processors:
  memory_limiter:
    limit_mib: 1024  # Increase from 512
```

### Missing spans for certain operations

Java services use automatic instrumentation, but custom operations may need manual spans:

```java
// Add to Java services
Span span = tracer.spanBuilder("custom-operation").startSpan();
try (Scope scope = span.makeCurrent()) {
    // Your code
} finally {
    span.end();
}
```

## Integration with Elastic Observability (Optional)

To export telemetry to Elastic APM Server:

1. Uncomment the Elastic exporter in `otel-collector-config.yaml`
2. Set environment variables:
   ```bash
   export ELASTIC_APM_SERVER_URL=https://your-apm-server:8200
   export ELASTIC_APM_SECRET_TOKEN=your-secret-token
   ```
3. Add to exporters in traces/metrics/logs pipelines

## Performance Impact

EDOT instrumentation has minimal performance impact:
- **Java**: ~2-5% CPU overhead, ~50MB additional memory
- **Node.js**: ~1-3% CPU overhead, ~30MB additional memory
- **Network**: ~1-2KB per request for telemetry export

Batching and sampling can further reduce impact in production.

## Best Practices

1. **Use meaningful span names**: Operations should be identifiable
2. **Add custom attributes**: Business context (user_id, order_id, etc.)
3. **Sample in production**: For high-traffic services, use sampling
4. **Monitor the collector**: Watch for dropped spans/metrics
5. **Correlate logs with traces**: Use trace_id in application logs

## Future Enhancements

- [ ] Complete Go services instrumentation
- [ ] Add custom business metrics
- [ ] Implement sampling strategies
- [ ] Add Elastic APM Server integration
- [ ] Create Grafana dashboards
- [ ] Add alerting rules
- [ ] Implement log correlation

## References

- [Elastic OpenTelemetry Distribution](https://www.elastic.co/docs/reference/opentelemetry)
- [EDOT Java Setup](https://www.elastic.co/docs/reference/opentelemetry/edot-sdks/java/setup/)
- [EDOT Node.js Setup](https://www.elastic.co/docs/reference/opentelemetry/edot-sdks/nodejs/setup/)
- [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
