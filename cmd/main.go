package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	otelglobal "go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

var echoTracer trace.Tracer

func newTracerProvider(exporter export.SpanExporter, name string) *sdktrace.TracerProvider {
	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	cfg := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithConfig(cfg),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.New(semconv.ServiceNameKey.String(name))))
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

	tp := newTracerProvider(exporter, "EchoTracerProvider")
	echoTracer = tp.Tracer("echo-tracer")

	if err := exporter.Start(); err != nil {
		log.Fatal(err)
	}

	r := echo.New()
	r.Use(otelecho.Middleware("server-name", otelecho.WithTracerProvider(tp)))
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.GET("/hello", hello)

	r.Logger.Fatal(r.Start(":9999"))
}

func hello(c echo.Context) error {
	ctx := c.Request().Context()
	var err error
	_, span1 := echoTracer.Start(ctx, "Span1")
	span1.RecordError(ctx, err)
	span1.End()

	_, span2 := echoTracer.Start(ctx, "Span2")
	span2.RecordError(ctx, err)
	time.Sleep(5 * time.Second)
	span2.End()

	_, span3 := echoTracer.Start(ctx, "Span3")
	err = errors.New("Dummy error")
	span3.RecordError(ctx, err)
	span3.SetStatus(codes.Internal, err.Error())
	span3.SetAttributes(label.Key("error").String(err.Error()))
	defer span3.End()
	if err != nil {
		return c.String(http.StatusInternalServerError, "deu ruim")
	}

	return c.String(http.StatusOK, "hello world!\n")
}
