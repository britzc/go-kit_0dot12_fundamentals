package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"go.opentelemetry.io/otel"
)

type PricingService interface {
	GetRetailTotal(ctx context.Context, code string, qty int) (total float64, err error)
	GetWholesaleTotal(ctx context.Context, partner, code string, qty int) (total float64, err error)
}

func MakeTotalRetailPriceEndpoint(svc PricingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ctx, span := otel.Tracer("Transport.Endpoint").Start(ctx, "GetRetailTotal")
		defer span.End()

		req := request.(TotalRetailPriceRequest)
		total, err := svc.GetRetailTotal(ctx, req.Code, req.Qty)
		if err != nil {
			return TotalRetailPriceResponse{total, err.Error()}, nil
		}

		return TotalRetailPriceResponse{total, ""}, nil
	}
}

func MakeTotalWholesalePriceEndpoint(svc PricingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ctx, span := otel.Tracer("Transport.Endpoint").Start(ctx, "GetWholesaleTotal")
		defer span.End()

		req := request.(TotalWholesalePriceRequest)
		total, err := svc.GetWholesaleTotal(ctx, req.Partner, req.Code, req.Qty)
		if err != nil {
			return TotalWholesalePriceResponse{total, err.Error()}, nil
		}

		return TotalWholesalePriceResponse{total, ""}, nil
	}
}
