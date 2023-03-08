package tracer

import (
	"io"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(os.Getenv("OTEL_SERVICE_NAME")),
			semconv.ServiceVersion("v1.0.0"),
			attribute.String("environment", "dev"),
		),
	)
	return r
}

// Init Tracer Provider
func NewTracerProvider(file io.Writer) error {
	exp, err := newExporter(file)
	if err != nil {
		return err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)

	otel.SetTracerProvider(tp)

	return nil
}

// func (tp *TracerProvider) Shutdown() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()

// 	if err := tp.Provider.Shutdown(ctx); err != nil {
// 		log.Fatal(err)
// 	}

// 	return tp.Provider.Shutdown(ctx)
// }
