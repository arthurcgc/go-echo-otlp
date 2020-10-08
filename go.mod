module github.com/arthurcgc/go-otel-example

go 1.13

require (
	github.com/labstack/echo/v4 v4.1.17
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.12.0
	go.opentelemetry.io/otel v0.12.0
	go.opentelemetry.io/otel/exporters/otlp v0.12.0
	go.opentelemetry.io/otel/sdk v0.12.0
)
