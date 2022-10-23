package transport

import (
	"context"
	"encoding/json"
	"net/http"

	gkendpoint "github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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
	tracer := otel.Tracer("Transport.Transport")

	var retailEndpoint gkendpoint.Endpoint
	retailEndpoint = MakeTotalRetailPriceEndpoint(svc)
	retailEndpoint = LogTotalRetailPriceEndpoint(log.With(logger, "service", "PricingService"))(retailEndpoint)

	return httptransport.NewServer(
		retailEndpoint,
		decodeTotalRetailPriceRequest,
		encodeResponse,
		httptransport.ServerBefore(startTrace(tracer, logger)),
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
	tracer := otel.Tracer("Transport.Transport")

	var wholesaleEndpoint gkendpoint.Endpoint
	wholesaleEndpoint = MakeTotalWholesalePriceEndpoint(svc)
	wholesaleEndpoint = LogTotalWholesalePriceEndpoint(log.With(logger, "service", "PricingService"))(wholesaleEndpoint)

	return httptransport.NewServer(
		wholesaleEndpoint,
		decodeTotalWholesalePriceRequest,
		encodeResponse,
		httptransport.ServerBefore(startTrace(tracer, logger)),
	)
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func startTrace(tracer trace.Tracer, logger log.Logger) httptransport.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		ctx, span := otel.Tracer("Transport.PricingProxy").Start(req.Context(), "StartTrace")
		span.End()
		return ctx
	}
}
