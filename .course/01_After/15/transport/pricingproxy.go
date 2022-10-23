package transport

import (
	"context"
	"errors"
	"fmt"
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
	"golang.org/x/time/rate"
)

const (
	INVALID_RESPONSE = "Invalid Response"
)

func NewPricingServiceProxy(ctx context.Context, instanceList []string, logger log.Logger) PricingService {
	getRetailTotal := makeEndpoint(ctx, instanceList, makeRetailTotalProxy)
	getWholesaleTotal := makeEndpoint(ctx, instanceList, makeWholesaleTotalProxy)

	return proxyMiddleware{ctx, getRetailTotal, getWholesaleTotal}
}

func makeEndpoint(ctx context.Context, instanceList []string, makeProxy func(context.Context, string) endpoint.Endpoint) endpoint.Endpoint {
	var (
		qps         = 100
		maxAttempts = 3
		maxTime     = 250 * time.Millisecond
	)

	var endpointer sd.FixedEndpointer
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeProxy(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	balancer := lb.NewRoundRobin(endpointer)
	return lb.Retry(maxAttempts, maxTime, balancer)
}

func makeRetailTotalProxy(ctx context.Context, instance string) endpoint.Endpoint {
	path := fmt.Sprintf("http://%s/retail", instance)
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	return httptransport.NewClient(
		"POST",
		u,
		encodeRequest,
		decodeTotalRetailPriceResponse,
	).Endpoint()
}

func makeWholesaleTotalProxy(ctx context.Context, instance string) endpoint.Endpoint {
	path := fmt.Sprintf("http://%s/wholesale", instance)
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	return httptransport.NewClient(
		"POST",
		u,
		encodeRequest,
		decodeTotalWholesalePriceResponse,
	).Endpoint()
}

type proxyMiddleware struct {
	ctx               context.Context
	getRetailTotal    endpoint.Endpoint
	getWholesaleTotal endpoint.Endpoint
}

func (mw proxyMiddleware) GetRetailTotal(code string, qty int) (total float64, err error) {
	req := TotalRetailPriceRequest{Code: code, Qty: qty}

	response, err := mw.getRetailTotal(mw.ctx, req)
	if err != nil {
		return 0.0, err
	}

	resp := response.(TotalRetailPriceResponse)
	if resp.Err != "" {
		return 0.0, errors.New(resp.Err)
	}

	return resp.Total, nil
}

func (mw proxyMiddleware) GetWholesaleTotal(partner, code string, qty int) (total float64, err error) {
	req := TotalWholesalePriceRequest{Partner: partner, Code: code, Qty: qty}

	response, err := mw.getWholesaleTotal(mw.ctx, req)
	if err != nil {
		return 0.0, err
	}

	resp := response.(TotalWholesalePriceResponse)
	if resp.Err != "" {
		return 0.0, errors.New(resp.Err)
	}

	return resp.Total, nil
}
