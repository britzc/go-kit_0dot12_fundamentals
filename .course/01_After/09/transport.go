package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

const (
	INVALID_REQUEST = "Invalid Request"
)

// GENERIC ERROR

type errorResponse struct {
	Err string `json:"err, omitEmpty"`
}

func (e *errorResponse) Error() string {
	return e.Err
}

func (e *errorResponse) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"err":"%s"}`, e.Err)), nil
}

// RETAIL OBJECTS & HANDLERS

type totalRetailPriceRequest struct {
	Code string `json:"code"`
	Qty  int    `json:"qty"`
}

type totalRetailPriceResponse struct {
	Total float64 `json:"total"`
	Err   string  `json:"err,omitempty"`
}

func decodeTotalRetailPriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request totalRetailPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, &errorResponse{Err: INVALID_REQUEST}
	}

	return request, nil
}

func MakeTotalRetailPriceHttpHandler(pricingService PricingService) *httptransport.Server {
	return httptransport.NewServer(
		MakeTotalRetailPriceEndpoint(pricingService),
		decodeTotalRetailPriceRequest,
		encodeResponse,
	)
}

// WHOLESALE OBJECTS & HANDLERS

type totalWholesalePriceRequest struct {
	Partner string `json:"partner"`
	Code    string `json:"code"`
	Qty     int    `json:"qty"`
}

type totalWholesalePriceResponse struct {
	Total float64 `json:"total"`
	Err   string  `json:"err,omitempty"`
}

<<<<<<< HEAD
=======
func decodeTotalWholesalePriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request totalWholesalePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, &errorResponse{Err: INVALID_REQUEST}
	}

	return request, nil
}

// func MakeTotalWholesalePriceHttpHandler(pricingService PricingService) *httptransport.Server {
// 	return httptransport.NewServer(
// 		MakeTotalWholesalePriceEndpoint(pricingService),
// 		decodeTotalWholesalePriceRequest,
// 		encodeResponse,
// 	)
// }

>>>>>>> c3d1822cb45983a78c3b60edc34931d376136d1f
// GENERIC OBJECTS & HANDLERS

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
