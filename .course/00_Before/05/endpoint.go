package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type PricingService interface {
	GetTotalRetailPrice(code string, qty int) (total float64, err error)
}

type totalRetailPriceRequest struct {
	Code string `json:"code"`
	Qty  int    `json:"qty"`
}

type totalRetailPriceResponse struct {
	Total float64 `json:"total"`
	Err   string  `json:"err,omitempty"`
}

func MakeTotalRetailPriceEndpoint(ps PricingService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(totalRetailPriceRequest)
		total, err := ps.GetTotalRetailPrice(req.Code, req.Qty)
		if err != nil {
			return totalRetailPriceResponse{total, err.Error()}, nil
		}

		return totalRetailPriceResponse{total, ""}, nil
	}

}

func DecodeTotalRetailPriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request totalRetailPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
