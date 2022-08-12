package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/britzc/go-kit_0dot12_fundamentals/current/endpoint"
	"github.com/britzc/go-kit_0dot12_fundamentals/current/payload"
	gkendpoint "github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

const (
	INVALID_REQUEST = "Invalid Request"
)

type PricingService interface {
	GetRetailTotal(code string, qty int) (total float64, err error)
	GetWholesaleTotal(partner, code string, qty int) (total float64, err error)
}

func decodeTotalRetailPriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request payload.TotalRetailPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, &payload.ErrorResponse{Err: INVALID_REQUEST}
	}

	return request, nil
}

func MakeTotalRetailPriceHttpHandler(logger log.Logger, svc PricingService) *httptransport.Server {
	var retailEndpoint gkendpoint.Endpoint
	retailEndpoint = endpoint.MakeTotalRetailPriceEndpoint(svc)
	retailEndpoint = LoggingMiddleware(log.With(logger, "method", "MakeTotalRetailPriceHandler"))(retailEndpoint)

	return httptransport.NewServer(
		retailEndpoint,
		decodeTotalRetailPriceRequest,
		encodeResponse,
	)
}

func decodeTotalWholesalePriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request payload.TotalWholesalePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, &payload.ErrorResponse{Err: INVALID_REQUEST}
	}

	return request, nil
}

func MakeTotalWholesalePriceHttpHandler(logger log.Logger, svc PricingService) *httptransport.Server {
	var wholesaleEndpoint gkendpoint.Endpoint
	wholesaleEndpoint = endpoint.MakeTotalWholesalePriceEndpoint(svc)
	wholesaleEndpoint = LoggingMiddleware(log.With(logger, "method", "MakeTotalWholesalePriceHandler"))(wholesaleEndpoint)

	return httptransport.NewServer(
		wholesaleEndpoint,
		decodeTotalWholesalePriceRequest,
		encodeResponse,
	)
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
