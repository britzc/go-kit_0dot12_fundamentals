package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	gkendpoint "github.com/go-kit/kit/endpoint"
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

func decodeTotalRetailPriceResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response TotalRetailPriceResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, &ErrorResponse{Err: INVALID_RESPONSE}
	}
	return response, nil
}

func MakeTotalRetailPriceHttpHandler(logger log.Logger, svc PricingService) *httptransport.Server {
	var retailEndpoint gkendpoint.Endpoint
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

func decodeTotalWholesalePriceResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response TotalWholesalePriceResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, &ErrorResponse{Err: INVALID_RESPONSE}
	}
	return response, nil
}

func MakeTotalWholesalePriceHttpHandler(logger log.Logger, svc PricingService) *httptransport.Server {
	var wholesaleEndpoint gkendpoint.Endpoint
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

func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}
