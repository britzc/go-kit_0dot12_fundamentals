package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func LogTotalRetailPriceEndpoint(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("endpoint", "TotalRetailPriceEndpoint", "msg", "Calling endpoint")
			defer logger.Log("endpoint", "TotalRetailPriceEndpoint", "msg", "Called endpoint")

			return next(ctx, request)
		}
	}
}

func LogTotalWholesalePriceEndpoint(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("endpoint", "TotalWholesalePriceEndpoint", "msg", "Calling endpoint")
			defer logger.Log("endpoint", "TotalWholesalePriceEndpoint", "msg", "Called endpoint")

			return next(ctx, request)
		}
	}
}
