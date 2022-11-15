package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type PricingService interface {
	GetRetailTotal(code string, qty int) (total float64, err error)
}

func MakeTotalRetailPriceEndpoint(ps PricingService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(totalRetailPriceRequest)
		total, err := ps.GetRetailTotal(req.Code, req.Qty)
		if err != nil {
			return totalRetailPriceResponse{0.0, err.Error()}, nil
		}

		return totalRetailPriceResponse{total, ""}, nil
	}
}
