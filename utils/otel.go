package utils

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var tracer = otel.Tracer("mux-server")
var Tracer = &tracer

func InitTracer() (*sdktrace.TracerProvider, error) {
	hostname, _ := os.Hostname()

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("expense-tracker"),
		semconv.HostArchKey.String(runtime.GOARCH),
		semconv.HostNameKey.String(hostname),
	)
	//  -- This prints out to terminal --
	// exporter, err := stdout.New(stdout.WithPrettyPrint())

	var (
		jaegerEndpoint string
		exists         bool
	)

	if jaegerEndpoint, exists = os.LookupEnv("JAEGER_ENDPOINT"); exists == false {
		if os.Getenv("ENV") == "local" {
			jaegerEndpoint = "localhost:4317"
		} else {
			panic("Missing JAEGER_ENDPOINT config")
		}

	}
	// -- This exports to jaeger --
	exporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithEndpoint(jaegerEndpoint), otlptracegrpc.WithInsecure())

	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

// Convert header to attributes
func HeaderToAttributes(data map[string][]string) []attribute.KeyValue {
	var dataString string

	for k, v := range data {
		dataString = dataString + k + " : " + strings.Join(v, ",") + "\n"
	}
	var concatenatedHeaders = make(map[string]interface{})

	concatenatedHeaders["http.request.headers"] = dataString

	return MapToAttributes(concatenatedHeaders)
}

// Convert Map to attribute.KeyValue
func MapToAttributes(data map[string]interface{}) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	for key, value := range data {
		switch v := value.(type) {
		case string:
			attrs = append(attrs, attribute.String(key, v))
		case int:
			attrs = append(attrs, attribute.Int(key, v))
		case bool:
			attrs = append(attrs, attribute.Bool(key, v))
		case float64:
			attrs = append(attrs, attribute.Float64(key, v))
		case []string:
			attrs = append(attrs, attribute.StringSlice(key, v))
		// Add more types as needed (float64, []string, etc.)
		default:
			// Handle unsupported types or convert to string representation
			GeneralLogger.Error("Warning: Unsupported type for key %s\n", key)
			attrs = append(attrs, attribute.String(key, fmt.Sprintf("%v", v)))
		}
	}
	return attrs
}
