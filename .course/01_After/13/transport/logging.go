package transport

import (
	"context"

	gkendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func LoggingMiddleware(logger log.Logger) gkendpoint.Middleware {
	return func(next gkendpoint.Endpoint) gkendpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			defer logger.Log("msg", "called endpoint")
			return next(ctx, request)
		}
	}
}
