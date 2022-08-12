package endpoint

import (
	"context"

	"github.com/britzc/go-kit_0dot12_fundamentals/current/payload"
	"github.com/go-kit/kit/endpoint"
)

type PricingService interface {
	GetRetailTotal(code string, qty int) (total float64, err error)
	GetWholesaleTotal(partner, code string, qty int) (total float64, err error)
}

func MakeTotalRetailPriceEndpoint(svc PricingService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(payload.TotalRetailPriceRequest)
		total, err := svc.GetRetailTotal(req.Code, req.Qty)
		if err != nil {
			return payload.TotalRetailPriceResponse{total, err.Error()}, nil
		}

		return payload.TotalRetailPriceResponse{total, ""}, nil
	}
}

func MakeTotalWholesalePriceEndpoint(svc PricingService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(payload.TotalWholesalePriceRequest)
		total, err := svc.GetWholesaleTotal(req.Partner, req.Code, req.Qty)
		if err != nil {
			return payload.TotalWholesalePriceResponse{total, err.Error()}, nil
		}

		return payload.TotalWholesalePriceResponse{total, ""}, nil
	}
}
