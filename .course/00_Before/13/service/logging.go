package service

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   PricingService
}

func NewLoggingMiddleware(logger log.Logger, next PricingService) (lmw *loggingMiddleware) {
	lmw = &loggingMiddleware{
		logger: logger,
		next:   next,
	}

	return
}

func (mw loggingMiddleware) GetRetailTotal(code string, qty int) (total float64, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetRetailTotal",
			"code", code,
			"quantity", qty,
			"total", total,
			"error", err,
			"duration", time.Since(begin),
		)
	}(time.Now())

	total, err = mw.next.GetRetailTotal(code, qty)

	return
}

func (mw loggingMiddleware) GetWholesaleTotal(partner, code string, qty int) (total float64, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetWholesaleTotal",
			"partner", partner,
			"code", code,
			"quantity", qty,
			"total", total,
			"error", err,
			"duration", time.Since(begin),
		)
	}(time.Now())
	total, err = mw.next.GetWholesaleTotal(partner, code, qty)

	return
}
