package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	traceManager "github.com/arthurcgc/go-otel-example/pkg/tracing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

func main() {
	tm := traceManager.GetManager()
	if err := tm.StartExporter(); err != nil {
		log.Print(err)
	}
	tm.AddTrace("echo")

	r := echo.New()
	r.Use(
		otelecho.Middleware("CWAPI-server-trace",
			otelecho.WithTracerProvider(tm.TracerProvider("echo"))),
	)
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.GET("/hello", hello)

	r.Logger.Fatal(r.Start(":9999"))
}

func hello(c echo.Context) error {
	var err error
	ctx := c.Request().Context()
	tracer := traceManager.GetTracer("echo")
	// checar se tracer é nil
	_, span1 := tracer.Start(ctx, "Span1")
	span1.RecordError(ctx, err)
	span1.End()

	_, span2 := tracer.Start(ctx, "Span2")
	span2.RecordError(ctx, err)
	time.Sleep(5 * time.Second)
	span2.End()

	_, span3 := tracer.Start(ctx, "Span3")
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
