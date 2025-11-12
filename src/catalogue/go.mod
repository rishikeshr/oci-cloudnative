module mushop/catalogue

go 1.20

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/go-kit/kit v0.9.0
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.9
	github.com/opentracing/opentracing-go v1.1.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5
	github.com/prometheus/client_golang v1.7.1
	github.com/sony/gobreaker v0.4.1
	github.com/weaveworks/common v0.0.0-20200206153930-760e36ae819a
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.24.0
	go.opentelemetry.io/otel/sdk v1.24.0
	golang.org/x/net v0.19.0
)

require (
	github.com/Shopify/sarama v1.19.0 // indirect
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5 // indirect
	github.com/apache/thrift v0.13.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.1.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/googleapis v1.1.0 // indirect
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/gogo/status v1.0.3 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/opentracing-contrib/go-observer v0.0.0-20170622124052-a52f23424492 // indirect
	github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9 // indirect
	github.com/opentracing/basictracer-go v1.0.0 // indirect
	github.com/openzipkin/zipkin-go v0.2.2 // indirect
	github.com/pierrec/lz4 v1.0.2-0.20190131084431-473cd7ce01a1 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.10.0 // indirect
	github.com/prometheus/procfs v0.1.3 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271 // indirect
	github.com/streadway/handy v0.0.0-20200128134331-0f66f006fb2e // indirect
	github.com/uber/jaeger-client-go v2.15.0+incompatible // indirect
	github.com/uber/jaeger-lib v1.5.1-0.20181102163054-1fc5c315e03c // indirect
	github.com/weaveworks/promrus v1.2.0 // indirect
	go.opencensus.io v0.22.2 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.opentelemetry.io/proto/otlp v1.1.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/grpc v1.61.1 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

// Replace directive to fix hdrhistogram module path issue
replace github.com/codahale/hdrhistogram => github.com/HdrHistogram/hdrhistogram-go v1.1.2

// Replace directive to pin jaeger-lib to version with metrics/testutils package
replace github.com/uber/jaeger-lib => github.com/uber/jaeger-lib v1.5.1-0.20181102163054-1fc5c315e03c
