package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type PricingService interface {
	GetRetailTotal(code string, qty int) (total float64, err error)
	GetWholesaleTotal(partner, code string, qty int) (total float64, err error)
}

func MakeTotalRetailPriceEndpoint(svc PricingService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(totalRetailPriceRequest)
		total, err := svc.GetRetailTotal(req.Code, req.Qty)
		if err != nil {
			return totalRetailPriceResponse{total, err.Error()}, nil
		}

		return totalRetailPriceResponse{total, ""}, nil
	}
}

func MakeTotalWholesalePriceEndpoint(svc PricingService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(totalWholesalePriceRequest)
		total, err := svc.GetWholesaleTotal(req.Partner, req.Code, req.Qty)
		if err != nil {
			return totalWholesalePriceResponse{total, err.Error()}, nil
		}

		return totalWholesalePriceResponse{total, ""}, nil
	}
}
