package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/britzc/go-kit_0dot12_fundamentals/current/service"
	"github.com/britzc/go-kit_0dot12_fundamentals/current/transport"
	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "priceservice01:8080,priceservice02:8080,priceservice03:8080", "List of URLs to proxy pricing requests")
	)
	flag.Parse()

	proxyList := strings.Split(*proxy, ",")
	for i := range proxyList {
		proxyList[i] = strings.TrimSpace(proxyList[i])
	}

	logger := log.NewLogfmtLogger(os.Stderr)

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
