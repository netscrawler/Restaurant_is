package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/netscrawler/Restaurant_is/menu_service/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config содержит конфигурацию телеметрии

// Telemetry содержит все компоненты телеметрии
type Telemetry struct {
	Config         *config.TelemertyConfig
	Logger         *slog.Logger
	Tracer         trace.Tracer
	Meter          metric.Meter
	MeterProvider  *sdkmetric.MeterProvider
	TracerProvider *sdktrace.TracerProvider
	CustomMetrics  *CustomMetrics
}

// New создает новый экземпляр телеметрии
func New(cfg *config.TelemertyConfig, logger *slog.Logger) (*Telemetry, error) {
	// Создаем meter provider с Prometheus exporter
	prometheusExporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(prometheusExporter),
	)
	otel.SetMeterProvider(meterProvider)

	// Создаем OTLP exporter для трейсов
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.TraceEndpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Создаем tracer provider с OTLP exporter
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(createResource(cfg)),
	)
	otel.SetTracerProvider(tracerProvider)

	// Создаем tracer
	tracer := tracerProvider.Tracer(
		cfg.ServiceName,
		trace.WithInstrumentationVersion(cfg.ServiceVersion),
	)

	// Создаем meter
	meter := meterProvider.Meter(cfg.ServiceName)

	telemetry := &Telemetry{
		Config:         cfg,
		Logger:         logger,
		Tracer:         tracer,
		Meter:          meter,
		MeterProvider:  meterProvider,
		TracerProvider: tracerProvider,
	}

	// Создаем кастомные метрики
	telemetry.CustomMetrics = telemetry.NewCustomMetrics()

	return telemetry, nil
}

// createResource создает ресурс для трейсинга
func createResource(cfg *config.TelemertyConfig) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.ServiceName),
		semconv.ServiceVersion(cfg.ServiceVersion),
		semconv.ServiceNamespace("menu-service"),
	)
}

// Shutdown корректно завершает работу телеметрии
func (t *Telemetry) Shutdown(ctx context.Context) error {
	if err := t.MeterProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown meter provider: %w", err)
	}
	if err := t.TracerProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown tracer provider: %w", err)
	}
	return nil
}

// StartSpan создает новый span для трейсинга
func (t *Telemetry) StartSpan(
	ctx context.Context,
	name string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return t.Tracer.Start(ctx, name, opts...)
}

// RecordMetric записывает метрику
func (t *Telemetry) RecordMetric(name string, value float64, attrs ...string) {
	// Здесь можно добавить логику для записи метрик
	t.Logger.Debug("Recording metric",
		slog.String("name", name),
		slog.Float64("value", value),
		slog.Any("attributes", attrs),
	)
}
