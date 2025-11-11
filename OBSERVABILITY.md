# MuShop Observability with EDOT

This document describes the OpenTelemetry instrumentation implemented in MuShop using EDOT (Elastic Distribution of OpenTelemetry).

## Overview

MuShop is instrumented with OpenTelemetry to provide comprehensive observability through:
- **Distributed Tracing**: Track requests across all microservices
- **Metrics**: Monitor service performance and resource usage
- **Logs**: Centralized logging with trace correlation

## Architecture

```
┌─────────────────────┐     ┌─────────────────┐     ┌─────────────────────┐
│   MuShop Services   │────▶│   OTel          │────▶│   Elasticsearch     │
│   (EDOT Agents)     │     │   Collector     │     │   Managed OTLP      │
│                     │     │                 │     │   (Cloud)           │
│  • orders (Java)    │     │  Receives:      │     │                     │
│  • carts (Java)     │     │   - Traces      │     │  • APM UI           │
│  • fulfillment      │     │   - Metrics     │     │  • Logs Explorer    │
│  • user (Node.js)   │     │   - Logs        │     │  • Service Map      │
│  • api (Node.js)    │     │                 │     │  • Dashboards       │
│  • assets           │     │  Processes:     │     │  • Alerts           │
│  • newsletter       │     │   - Batching    │     │                     │
│  • catalogue (Go)*  │     │   - Resources   │     │                     │
│  • payment (Go)*    │     │   - Memory Mgmt │     │                     │
│  • events (Go)*     │     │                 │     │                     │
└─────────────────────┘     └─────────────────┘     └─────────────────────┘

* Go services have OTel packages added, manual instrumentation pending
```

All telemetry data (traces, metrics, logs) is sent to **Elasticsearch Managed OpenTelemetry (MOTLP)** via the OpenTelemetry Collector.

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
- **Receives**: OTLP over gRPC (4317) and HTTP (4318) from instrumented services
- **Processes**: Batching, resource detection, memory limiting, enrichment
- **Exports**: All signals to **Elasticsearch Managed OpenTelemetry (MOTLP)**

Configuration file: `otel-collector-config.yaml`

### Elasticsearch Managed OpenTelemetry (MOTLP)

MuShop sends all telemetry to Elasticsearch Cloud via MOTLP:

**Endpoint Format**: `https://<deployment-id>.apm.<region>.cloud.es.io`

**Authentication**: Bearer token (API Key or Secret Token)

**Signals Exported**:
- **Traces** → Elastic APM (distributed tracing)
- **Metrics** → Elasticsearch (time-series metrics)
- **Logs** → Elasticsearch (structured logging with trace correlation)

## Setup

### Prerequisites

1. **Elasticsearch Cloud Deployment**:
   - Sign up at https://cloud.elastic.co
   - Create a new deployment
   - Note your deployment ID and region

2. **API Key or Secret Token**:
   - In Kibana: Stack Management → API Keys → Create API Key
   - Or use APM Server secret token from deployment settings

### Configuration

1. **Copy the environment template**:
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your Elasticsearch credentials**:
   ```bash
   # Elasticsearch Managed OpenTelemetry endpoint
   ELASTIC_OTLP_ENDPOINT=https://abc123def456.apm.us-central1.gcp.cloud.es.io

   # API Key or Secret Token
   ELASTIC_OTLP_TOKEN=your_api_key_or_secret_token

   # Deployment environment
   DEPLOYMENT_ENVIRONMENT=production
   ```

3. **Start the services**:
   ```bash
   docker-compose up -d
   ```

4. **Verify telemetry is flowing**:
   ```bash
   # Check OTel Collector logs
   docker-compose logs -f otel-collector

   # Look for successful exports to Elasticsearch
   # Should see: "Traces export: OK" or similar
   ```

### Accessing Observability in Kibana

Once telemetry is flowing, access your observability data in Kibana:

#### APM UI (Application Performance Monitoring)
- **URL**: `https://your-deployment.kb.<region>.cloud.es.io/app/apm`
- **Features**:
  - Service overview and dependencies
  - Distributed traces
  - Transaction breakdowns
  - Error tracking
  - Service maps
  - Latency percentiles (P50, P95, P99)

#### Logs Explorer
- **URL**: `https://your-deployment.kb.<region>.cloud.es.io/app/logs`
- **Features**:
  - Centralized log viewing
  - Trace correlation (jump from logs to traces)
  - Advanced filtering and search
  - Log anomalies

#### Metrics Explorer
- **URL**: `https://your-deployment.kb.<region>.cloud.es.io/app/metrics`
- **Features**:
  - Infrastructure metrics
  - Application metrics
  - Custom dashboards
  - Alerting

#### OTel Collector Health (Local)
- Health check: http://localhost:13133
- Collector metrics: http://localhost:8888/metrics
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
