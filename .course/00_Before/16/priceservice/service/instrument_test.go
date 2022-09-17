package service

import (
	"context"
	"testing"

	"github.com/go-kit/kit/metrics"
	"github.com/stretchr/testify/assert"
)

type MockCounter struct {
	result float64
}

func (mc *MockCounter) Add(val float64) {
	mc.result += val
}

func (mc *MockCounter) With(lvs ...string) metrics.Counter {
	return mc
}

func (mc *MockCounter) Result() float64 {
	return mc.result
}

type MockLatency struct {
	result float64
}

func (ml *MockLatency) Observe(val float64) {
	ml.result += val
}

func (ml *MockLatency) With(lvs ...string) metrics.Histogram {
	return ml
}

func (ml *MockLatency) Result() float64 {
	return ml.result
}

func Test_Instrumenting_GetRetailTotal(t *testing.T) {
	ctx := context.Background()

	counter := new(MockCounter)
	latency := new(MockLatency)

	var svc PricingService
	svc = new(MockPricingService)
	svc = NewInstrumentingMiddleware(counter, latency, svc)

	svc.GetRetailTotal(ctx, "aaa111", 5)
	svc.GetRetailTotal(ctx, "bbb222", 10)
	svc.GetRetailTotal(ctx, "ccc333", 15)

	counterActual := counter.Result()
	latencyActual := latency.Result()

	assert.True(t, counterActual == 3.0, "~2|Test counter expected: 3, not: \"%.1f\"~", counterActual)
	assert.True(t, latencyActual > 0.0, "~2|Test latency expected value greater than 0.0~")
}

func Test_Instrumenting_GetWholesaleTotal(t *testing.T) {
	ctx := context.Background()

	counter := new(MockCounter)
	latency := new(MockLatency)

	var svc PricingService
	svc = new(MockPricingService)
	svc = NewInstrumentingMiddleware(counter, latency, svc)

	svc.GetWholesaleTotal(ctx, "superstore", "aaa111", 5)
	svc.GetWholesaleTotal(ctx, "superstore", "bbb222", 10)
	svc.GetWholesaleTotal(ctx, "superstore", "ccc333", 15)

	counterActual := counter.Result()
	latencyActual := latency.Result()

	assert.True(t, counterActual == 3.0, "~2|Test counter expected: 3, not: \"%.1f\"~", counterActual)
	assert.True(t, latencyActual > 0.0, "~2|Test latency expected value greater than 0.0~")
}
