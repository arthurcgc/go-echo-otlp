module github.com/arthurcgc/dummy_server

go 1.13

require (
	github.com/gorilla/mux v1.7.3
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/echo/v4 v4.1.17
	github.com/labstack/gommon v0.3.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.12.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.12.0
	go.opentelemetry.io/otel v0.12.0
	go.opentelemetry.io/otel/exporters/otlp v0.12.0
	go.opentelemetry.io/otel/exporters/stdout v0.12.0
	go.opentelemetry.io/otel/sdk v0.12.0
	google.golang.org/grpc v1.32.0
)
