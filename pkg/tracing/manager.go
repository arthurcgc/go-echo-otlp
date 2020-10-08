package tracing

import (
	"context"
	"sync"
	"time"

	otelglobal "go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/otlp"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

var (
	once    sync.Once
	manager *TraceManager
)

// TraceManager manages all tracing spans and connections to the Open Telemetry Collector
type TraceManager struct {
	exporter       *otlp.Exporter
	TraceInstances map[string]TraceInstance
}

type TraceInstance struct {
	Tracer         trace.Tracer
	TracerProvider *sdktrace.TracerProvider
}

func new() *TraceManager {
	exp := newDefaultExporter()

	return &TraceManager{
		exporter:       exp,
		TraceInstances: make(map[string]TraceInstance),
	}
}

// Singleton
func GetManager() *TraceManager {
	if manager == nil {
		once.Do(func() { manager = new() })
	}

	return manager
}

func GetTracer(name string) trace.Tracer {
	return manager.TraceInstances[name].Tracer
}

func (tm *TraceManager) AddTrace(name string) {
	tm.TraceInstances[name] = newTraceInstance(name, tm.exporter)
}

func (tm *TraceManager) TracerProvider(name string) *sdktrace.TracerProvider {
	return tm.TraceInstances[name].TracerProvider
}

func newTraceInstance(name string, exporter *otlp.Exporter) TraceInstance {
	tp := newTracerProvider(exporter, name+"-tracer-provider")
	tracer := tp.Tracer(name + "-tracer")

	return TraceInstance{
		TracerProvider: tp,
		Tracer:         tracer,
	}
}

func newTracerProvider(exporter export.SpanExporter, name string) *sdktrace.TracerProvider {
	cfg := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithConfig(cfg),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.New(semconv.ServiceNameKey.String(name))))
}

func newDefaultExporter() *otlp.Exporter {
	return otlp.NewUnstartedExporter(otlp.WithInsecure())
}

func (tm *TraceManager) ShutdownExporter() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := tm.exporter.Shutdown(ctx); err != nil {
		otelglobal.Handle(err)
	}
}

func (tm *TraceManager) StartExporter() error {
	return tm.exporter.Start()
}
