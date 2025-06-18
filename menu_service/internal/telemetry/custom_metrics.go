package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// CustomMetrics содержит кастомные метрики для бизнес-логики
type CustomMetrics struct {
	// Счетчики
	dishCreatedCounter metric.Int64Counter
	dishUpdatedCounter metric.Int64Counter
	dishDeletedCounter metric.Int64Counter
	dishListedCounter  metric.Int64Counter

	// Гистограммы для времени выполнения
	dishCreateDuration metric.Float64Histogram
	dishUpdateDuration metric.Float64Histogram
	dishGetDuration    metric.Float64Histogram
	dishListDuration   metric.Float64Histogram

	// Gauge для активных блюд
	activeDishesGauge metric.Int64UpDownCounter
}

// NewCustomMetrics создает новые кастомные метрики
func (t *Telemetry) NewCustomMetrics() *CustomMetrics {
	// Счетчики
	dishCreatedCounter, _ := t.Meter.Int64Counter(
		"dish_created_total",
		metric.WithDescription("Total number of dishes created"),
	)

	dishUpdatedCounter, _ := t.Meter.Int64Counter(
		"dish_updated_total",
		metric.WithDescription("Total number of dishes updated"),
	)

	dishDeletedCounter, _ := t.Meter.Int64Counter(
		"dish_deleted_total",
		metric.WithDescription("Total number of dishes deleted"),
	)

	dishListedCounter, _ := t.Meter.Int64Counter(
		"dish_listed_total",
		metric.WithDescription("Total number of dish list requests"),
	)

	// Гистограммы для времени выполнения
	dishCreateDuration, _ := t.Meter.Float64Histogram(
		"dish_create_duration_seconds",
		metric.WithDescription("Time spent creating dishes"),
		metric.WithUnit("s"),
	)

	dishUpdateDuration, _ := t.Meter.Float64Histogram(
		"dish_update_duration_seconds",
		metric.WithDescription("Time spent updating dishes"),
		metric.WithUnit("s"),
	)

	dishGetDuration, _ := t.Meter.Float64Histogram(
		"dish_get_duration_seconds",
		metric.WithDescription("Time spent getting dishes"),
		metric.WithUnit("s"),
	)

	dishListDuration, _ := t.Meter.Float64Histogram(
		"dish_list_duration_seconds",
		metric.WithDescription("Time spent listing dishes"),
		metric.WithUnit("s"),
	)

	// Gauge для активных блюд
	activeDishesGauge, _ := t.Meter.Int64UpDownCounter(
		"active_dishes_count",
		metric.WithDescription("Current number of active dishes"),
	)

	return &CustomMetrics{
		dishCreatedCounter: dishCreatedCounter,
		dishUpdatedCounter: dishUpdatedCounter,
		dishDeletedCounter: dishDeletedCounter,
		dishListedCounter:  dishListedCounter,
		dishCreateDuration: dishCreateDuration,
		dishUpdateDuration: dishUpdateDuration,
		dishGetDuration:    dishGetDuration,
		dishListDuration:   dishListDuration,
		activeDishesGauge:  activeDishesGauge,
	}
}

// RecordDishCreated записывает метрику создания блюда
func (cm *CustomMetrics) RecordDishCreated(ctx context.Context, categoryID int32) {
	cm.dishCreatedCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.Int("category_id", int(categoryID)),
	))
}

// RecordDishUpdated записывает метрику обновления блюда
func (cm *CustomMetrics) RecordDishUpdated(ctx context.Context, categoryID int32) {
	cm.dishUpdatedCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.Int("category_id", int(categoryID)),
	))
}

// RecordDishDeleted записывает метрику удаления блюда
func (cm *CustomMetrics) RecordDishDeleted(ctx context.Context, categoryID int32) {
	cm.dishDeletedCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.Int("category_id", int(categoryID)),
	))
}

// RecordDishListed записывает метрику запроса списка блюд
func (cm *CustomMetrics) RecordDishListed(ctx context.Context, categoryID *int32, onlyAvailable bool) {
	attrs := []attribute.KeyValue{
		attribute.Bool("only_available", onlyAvailable),
	}
	if categoryID != nil {
		attrs = append(attrs, attribute.Int("category_id", int(*categoryID)))
	}

	cm.dishListedCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordDishCreateDuration записывает время создания блюда
func (cm *CustomMetrics) RecordDishCreateDuration(ctx context.Context, duration float64, categoryID int32) {
	cm.dishCreateDuration.Record(ctx, duration, metric.WithAttributes(
		attribute.Int("category_id", int(categoryID)),
	))
}

// RecordDishUpdateDuration записывает время обновления блюда
func (cm *CustomMetrics) RecordDishUpdateDuration(ctx context.Context, duration float64, categoryID int32) {
	cm.dishUpdateDuration.Record(ctx, duration, metric.WithAttributes(
		attribute.Int("category_id", int(categoryID)),
	))
}

// RecordDishGetDuration записывает время получения блюда
func (cm *CustomMetrics) RecordDishGetDuration(ctx context.Context, duration float64) {
	cm.dishGetDuration.Record(ctx, duration)
}

// RecordDishListDuration записывает время получения списка блюд
func (cm *CustomMetrics) RecordDishListDuration(ctx context.Context, duration float64, categoryID *int32) {
	attrs := []attribute.KeyValue{}
	if categoryID != nil {
		attrs = append(attrs, attribute.Int("category_id", int(*categoryID)))
	}

	cm.dishListDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
}

// SetActiveDishesCount устанавливает количество активных блюд
func (cm *CustomMetrics) SetActiveDishesCount(ctx context.Context, count int64) {
	cm.activeDishesGauge.Add(ctx, count)
}
