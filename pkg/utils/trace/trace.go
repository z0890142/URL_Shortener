package trace

import (
	"URL_Shortener/pkg/utils/logger"
	"context"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	OPENTELEMETRYAPINameDefault      = "unknown_api"
	OPENTELEMETRYPackageNameDefault  = "unknown_package"
	OPENTELEMETRYFunctionNameDefault = "UnkownFunction"
)

func NewSpan(ctx context.Context, endpoint string) (context.Context, trace.Span) {
	var (
		pkgName  = OPENTELEMETRYPackageNameDefault
		funcName = OPENTELEMETRYFunctionNameDefault
	)

	if pc, _, _, ok := runtime.Caller(1); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			fullName := fn.Name()
			dotIndex := strings.LastIndexByte(fullName, '.')
			slashIndex := strings.LastIndexByte(fullName, '/')

			pkgName = fullName[slashIndex+1 : dotIndex]
			funcName = fullName[dotIndex+1:]
		} else {
			logger.Warn("apm: can't get function by PC")
		}
	} else {
		logger.Warn("apm: can't get runtime PC")
	}
	return otel.Tracer(pkgName).Start(ctx, funcName)
}

func NewTracerProvider(endpoint, service string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
