package service

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           PricingService
}

func NewInstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram, next PricingService) (imw *instrumentingMiddleware) {
	imw = &instrumentingMiddleware{
		requestCount:   requestCount,
		requestLatency: requestLatency,
		next:           next,
	}

	return
}

func (mw instrumentingMiddleware) GetRetailTotal(code string, qty int) (total float64, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetRetailTotal", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	total, err = mw.next.GetRetailTotal(code, qty)

	return
}

func (mw instrumentingMiddleware) GetWholesaleTotal(partner string, code string, qty int) (total float64, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetWholesaleTotal", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	total, err = mw.next.GetWholesaleTotal(partner, code, qty)

	return
}
