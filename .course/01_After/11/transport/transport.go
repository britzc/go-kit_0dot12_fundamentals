package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

const (
	INVALID_REQUEST = "Invalid Request"
)

func decodeTotalRetailPriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request TotalRetailPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, &ErrorResponse{Err: INVALID_REQUEST}
	}

	return request, nil
}

func MakeTotalRetailPriceHttpHandler(logger log.Logger, svc PricingService) *httptransport.Server {
	var retailEndpoint endpoint.Endpoint
	retailEndpoint = MakeTotalRetailPriceEndpoint(svc)
	retailEndpoint = LogTotalRetailPriceEndpoint(log.With(logger, "service", "PricingService"))(retailEndpoint)

	return httptransport.NewServer(
		retailEndpoint,
		decodeTotalRetailPriceRequest,
		encodeResponse,
	)
}

func decodeTotalWholesalePriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request TotalWholesalePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, &ErrorResponse{Err: INVALID_REQUEST}
	}

	return request, nil
}

func MakeTotalWholesalePriceHttpHandler(logger log.Logger, svc PricingService) *httptransport.Server {
	var wholesaleEndpoint endpoint.Endpoint
	wholesaleEndpoint = MakeTotalWholesalePriceEndpoint(svc)
	wholesaleEndpoint = LogTotalWholesalePriceEndpoint(log.With(logger, "service", "PricingService"))(wholesaleEndpoint)

	return httptransport.NewServer(
		wholesaleEndpoint,
		decodeTotalWholesalePriceRequest,
		encodeResponse,
	)
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
