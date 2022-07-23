package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

var ErrInvalidRequest = errors.New("Invalid Request")

type totalRetailPriceRequest struct {
	Code string `json:"code"`
	Qty  int    `json:"qty"`
}

type totalRetailPriceResponse struct {
	Total float64 `json:"total"`
	Err   string  `json:"err,omitempty"`
}

func MakeTotalRetailPriceHttpHandler(pricingService PricingService) *httptransport.Server {
	return httptransport.NewServer(
		MakeTotalRetailPriceEndpoint(pricingService),
		DecodeTotalRetailPriceRequest,
		EncodeResponse,
	)
}

func DecodeTotalRetailPriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request totalRetailPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, ErrInvalidRequest
	}

	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
