package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/britzc/go-kit_0dot12_fundamentals/current/service"
	"github.com/britzc/go-kit_0dot12_fundamentals/current/transport"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "localhost:8081,localhost:8082,localhost:8083", "List of URLs to proxy pricing requests")
	)
	flag.Parse()

	proxyList := strings.Split(*proxy, ",")
	for i := range proxyList {
		proxyList[i] = strings.TrimSpace(proxyList[i])
	}

	fmt.Println("Logging and tracing: In progress")

	logger := log.NewLogfmtLogger(os.Stderr)

	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	if err != nil {
		logger.Log("err", err)
		return
	}
	defer f.Close()

	exp, err := newExporter(f)
	if err != nil {
		logger.Log("err", err)
		return
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Log("err", err)
			return
		}
	}()
	otel.SetTracerProvider(tp)

	fmt.Println("Logging and tracing: Ready")

	fmt.Println("Endpoints and handlers: In progress")

	var svc service.PricingService
	svc = transport.NewPricingServiceProxy(context.Background(), proxyList, logger)

	rtr := mux.NewRouter().StrictSlash(true)

	totalRetailPriceHandler := transport.MakeTotalRetailPriceHttpHandler(logger, svc)
	rtr.Handle("/retail", totalRetailPriceHandler).Methods(http.MethodPost)

	totalWholesalePriceHandler := transport.MakeTotalWholesalePriceHttpHandler(logger, svc)
	rtr.Handle("/wholesale", totalWholesalePriceHandler).Methods(http.MethodPost)

	rtr.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	fmt.Println("Endpoints and handlers: Ready")

	fmt.Printf("Hosting on %s\n", *listen)

	http.ListenAndServe(*listen, rtr)
}

func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("PriceAPI"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "pluralsight"),
		),
	)
	return r
}
