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

			// diagnostic functionality

			return next(ctx, request)
		}
	}
}
