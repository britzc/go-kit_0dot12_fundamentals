package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type PricingService interface {
	GetRetailTotal(code string, qty int) (total float64, err error)
	GetWholesaleTotal(partner string, code string, qty int) (total float64, err error)
}

func MakeTotalRetailPriceEndpoint(svc PricingService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(TotalRetailPriceRequest)
		total, err := svc.GetRetailTotal(req.Code, req.Qty)
		if err != nil {
			return TotalRetailPriceResponse{0.0, err.Error()}, nil
		}

		return TotalRetailPriceResponse{total, ""}, nil
	}
}

func MakeTotalWholesalePriceEndpoint(svc PricingService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(TotalWholesalePriceRequest)
		total, err := svc.GetWholesaleTotal(req.Partner, req.Code, req.Qty)
		if err != nil {
			return TotalWholesalePriceResponse{0.0, err.Error()}, nil
		}

		return TotalWholesalePriceResponse{total, ""}, nil
	}
}
