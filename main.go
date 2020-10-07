package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	otelglobal "go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/otlp"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

var echoTracer = otelglobal.Tracer("echo-tracer")

func newTracerProvider(exporter export.SpanExporter) *sdktrace.TracerProvider {
	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	cfg := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithConfig(cfg),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.New(semconv.ServiceNameKey.String("EchoTracer"))))
}

func newExporter() *otlp.Exporter {
	return otlp.NewUnstartedExporter(otlp.WithInsecure())
}

func shutDownExporter(exporter *otlp.Exporter) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := exporter.Shutdown(ctx); err != nil {
		otelglobal.Handle(err)
	}
}

func main() {
	exporter := newExporter()
	defer shutDownExporter(exporter)

	tp := newTracerProvider(exporter)
	otelglobal.SetTracerProvider(tp)

	if err := exporter.Start(); err != nil {
		log.Fatal(err)
	}

	r := echo.New()
	r.Use(otelecho.Middleware("server-name"))
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.GET("/hello", hello)

	r.Logger.Fatal(r.Start(":9999"))
}

func hello(c echo.Context) error {
	_, span := echoTracer.Start(c.Request().Context(), "Hello World!")
	defer span.End()
	return c.String(http.StatusOK, "hello world!\n")
}
