package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func Init(ctx context.Context, service string) func(context.Context) error {
	exp, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatalf("otel exporter: %v", err)
	}

	tp := oteltrace.NewTracerProvider(
		oteltrace.WithBatcher(exp),
		oteltrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp.Shutdown
}
