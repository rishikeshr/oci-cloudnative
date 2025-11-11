// EDOT (Elastic Distribution of OpenTelemetry) initialization
import { start } from '@elastic/opentelemetry-node';

const serviceName = process.env.OTEL_SERVICE_NAME || 'mushop-user';
const otlpEndpoint = process.env.OTEL_EXPORTER_OTLP_ENDPOINT || 'http://otel-collector:4318';
const environment = process.env.OTEL_RESOURCE_ATTRIBUTES?.includes('deployment.environment')
  ? undefined
  : (process.env.DEPLOYMENT_ENVIRONMENT || 'production');

// Start EDOT with OTLP exporter configuration
start({
  serviceName,
  // EDOT will automatically use OTEL_EXPORTER_OTLP_* environment variables
  logLevel: process.env.OTEL_LOG_LEVEL || 'info',
});

console.log(`[EDOT] OpenTelemetry instrumentation started for service: ${serviceName}`);
