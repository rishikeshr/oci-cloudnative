module mushop/events

go 1.20

require (
	github.com/Shopify/sarama v1.38.1
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5 // indirect
	github.com/apache/thrift v0.12.0 // indirect
	github.com/go-kit/kit v0.9.0
	github.com/gogo/googleapis v1.2.0 // indirect
	github.com/gogo/status v1.1.0 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/microservices-demo/payment v0.0.0-20180427185901-3c33dbf915e2
	github.com/opentracing-contrib/go-observer v0.0.0-20170622124052-a52f23424492 // indirect
	github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5 // indirect
	github.com/openzipkin/zipkin-go-opentracing v0.3.5
	github.com/prometheus/client_golang v1.0.0
	github.com/sony/gobreaker v0.4.1 // indirect
	github.com/streadway/handy v0.0.0-20190108123426-d5acb3125c2a
	github.com/uber/jaeger-client-go v2.16.0+incompatible // indirect
	github.com/uber/jaeger-lib v1.5.1-0.20181102163054-1fc5c315e03c // indirect
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.24.0
	go.opentelemetry.io/otel/sdk v1.24.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80
	google.golang.org/grpc v1.22.1 // indirect
	gopkg.in/jcmturner/aescts.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/dnsutils.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/gokrb5.v7 v7.2.3 // indirect
	gopkg.in/jcmturner/rpc.v1 v1.1.0 // indirect
)

// Replace directive to fix hdrhistogram module path issue
replace github.com/codahale/hdrhistogram => github.com/HdrHistogram/hdrhistogram-go v1.1.2

// Replace directive to pin jaeger-lib to version with metrics/testutils package
replace github.com/uber/jaeger-lib => github.com/uber/jaeger-lib v1.5.1-0.20181102163054-1fc5c315e03c
