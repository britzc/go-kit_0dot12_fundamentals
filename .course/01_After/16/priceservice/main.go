package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/britzc/go-kit_0dot12_fundamentals/current/repo"
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
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		listen = flag.String("listen", ":8081", "HTTP listen address")
	)
	flag.Parse()

	fmt.Println("Logging and tracing: In progress")

	logger := log.NewLogfmtLogger(os.Stderr)

	f, err := os.Create("traces.txt")
	if err != nil {
		logger.Log("error", err)
		return
	}
	defer f.Close()

	exp, err := newExporter(f)
	if err != nil {
		logger.Log("error", err)
		return
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Log("error", err)
			return
		}
	}()
	otel.SetTracerProvider(tp)

	fmt.Println("Logging and tracing: Ready")

	fmt.Println("Repository: In progress")

	productRepo, _ := repo.NewProductRepo("products.csv", "partners.csv")

	fmt.Println("Repository: Ready")

	fmt.Println("Endpoints and handlers: In progress")

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "gokitfundamentals",
		Subsystem: "pricing_service",
		Name:      "request_count",
		Help:      "Number of request received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "gokitfundamentals",
		Subsystem: "pricing_service",
		Name:      "request_latency",
		Help:      "Total duration of requests.",
	}, fieldKeys)

	var svc service.PricingService
	svc = service.NewPricingService(productRepo)
	svc = service.NewInstrumentingMiddleware(requestCount, requestLatency, svc)
	svc = service.NewLoggingMiddleware(logger, svc)

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
	)
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("PriceService"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "pluralsight"),
		),
	)
	return r
}
