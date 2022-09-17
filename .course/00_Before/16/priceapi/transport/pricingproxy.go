package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"
)

const (
	INVALID_RESPONSE = "Invalid Response"
)

func NewPricingServiceProxy(ctx context.Context, instanceList []string, logger log.Logger) PricingService {
	tracer := otel.Tracer("Transport.PricingProxy")

	getRetailTotal := makeRetailTotalEndpoint("RetailTotal", instanceList, tracer, logger)
	getWholesaleTotal := makeWholesaleTotalEndpoint("WholesaleTotal", instanceList, tracer, logger)

	return proxyMiddleware{ctx, getRetailTotal, getWholesaleTotal}
}

func makeRetailTotalEndpoint(name string, instanceList []string, tracer trace.Tracer, logger log.Logger) endpoint.Endpoint {
	var (
		qps         = 100
		maxAttempts = 3
		maxTime     = 250 * time.Millisecond
	)

	var endpointer sd.FixedEndpointer
	for _, instance := range instanceList {
		path := fmt.Sprintf("http://%s/retail", instance)
		u, _ := url.Parse(path)

		var e endpoint.Endpoint
		e = httptransport.NewClient(
			"POST",
			u,
			encodeRequest,
			decodeTotalRetailPriceResponse,
			httptransport.ClientBefore(startTrace(tracer, logger)),
			httptransport.ClientAfter(stopTrace(tracer, logger)),
		).Endpoint()
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	balancer := lb.NewRoundRobin(endpointer)
	return lb.Retry(maxAttempts, maxTime, balancer)
}

func makeWholesaleTotalEndpoint(name string, instanceList []string, tracer trace.Tracer, logger log.Logger) endpoint.Endpoint {
	var (
		qps         = 100
		maxAttempts = 3
		maxTime     = 250 * time.Millisecond
	)

	var endpointer sd.FixedEndpointer
	for _, instance := range instanceList {
		path := fmt.Sprintf("http://%s/wholesale", instance)
		u, _ := url.Parse(path)

		var e endpoint.Endpoint
		e = httptransport.NewClient(
			"POST",
			u,
			encodeRequest,
			decodeTotalWholesalePriceResponse,
			httptransport.ClientBefore(startTrace(tracer, logger)),
			httptransport.ClientAfter(stopTrace(tracer, logger)),
		).Endpoint()
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	balancer := lb.NewRoundRobin(endpointer)
	return lb.Retry(maxAttempts, maxTime, balancer)
}

type proxyMiddleware struct {
	ctx               context.Context
	getRetailTotal    endpoint.Endpoint
	getWholesaleTotal endpoint.Endpoint
}

func (mw proxyMiddleware) GetRetailTotal(ctx context.Context, code string, qty int) (total float64, err error) {
	ctx, span := otel.Tracer("Transport.PricingProxy").Start(ctx, "GetRetailTotal")
	defer span.End()

	req := TotalRetailPriceRequest{Code: code, Qty: qty}

	response, err := mw.getRetailTotal(ctx, req)
	if err != nil {
		return 0.0, err
	}

	resp := response.(TotalRetailPriceResponse)
	if resp.Err != "" {
		return 0.0, errors.New(resp.Err)
	}

	return resp.Total, nil
}

func (mw proxyMiddleware) GetWholesaleTotal(ctx context.Context, partner, code string, qty int) (total float64, err error) {
	ctx, span := otel.Tracer("Transport.PricingProxy").Start(ctx, "GetWholesaleTotal")
	defer span.End()

	req := TotalWholesalePriceRequest{Partner: partner, Code: code, Qty: qty}

	response, err := mw.getWholesaleTotal(ctx, req)
	if err != nil {
		return 0.0, err
	}

	resp := response.(TotalWholesalePriceResponse)
	if resp.Err != "" {
		return 0.0, errors.New(resp.Err)
	}

	return resp.Total, nil
}

func startTrace(tracer trace.Tracer, logger log.Logger) httptransport.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		ctx, span := otel.Tracer("Transport.PricingProxy").Start(ctx, "StartTrace")
		span.End()
		return ctx
	}
}

func stopTrace(tracer trace.Tracer, logger log.Logger) httptransport.ClientResponseFunc {
	return func(ctx context.Context, res *http.Response) context.Context {
		ctx, span := otel.Tracer("Transport.PricingProxy").Start(ctx, "StopTrace")
		span.End()
		return ctx
	}
}
