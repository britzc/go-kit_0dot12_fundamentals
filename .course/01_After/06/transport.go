package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

type errorResponse struct {
	Err string `json:"err, omitEmpty"`
}

func (e *errorResponse) Error() string {
	return e.Err
}

func (e *errorResponse) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"err":"%s"}`, e.Err)), nil
}

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
		decodeTotalRetailPriceRequest,
		encodeResponse,
	)
}

func decodeTotalRetailPriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request totalRetailPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, &errorResponse{Err: "Invalid Request"}
	}

	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
