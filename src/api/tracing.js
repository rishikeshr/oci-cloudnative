// EDOT (Elastic Distribution of OpenTelemetry) initialization
const { start } = require('@elastic/opentelemetry-node');

const serviceName = process.env.OTEL_SERVICE_NAME || 'mushop-api';

// Start EDOT with OTLP exporter configuration
start({
  serviceName,
  // EDOT will automatically use OTEL_EXPORTER_OTLP_* environment variables
  logLevel: process.env.OTEL_LOG_LEVEL || 'info',
});

console.log(`[EDOT] OpenTelemetry instrumentation started for service: ${serviceName}`);
