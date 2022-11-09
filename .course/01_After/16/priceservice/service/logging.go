package service

import (
	"context"
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

func (mw loggingMiddleware) GetRetailTotal(ctx context.Context, code string, qty int) (total float64, err error) {
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

	total, err = mw.next.GetRetailTotal(ctx, code, qty)

	return
}

func (mw loggingMiddleware) GetWholesaleTotal(ctx context.Context, partner string, code string, qty int) (total float64, err error) {
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
	total, err = mw.next.GetWholesaleTotal(ctx, partner, code, qty)

	return
}
