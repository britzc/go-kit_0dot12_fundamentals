package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/britzc/go-kit_0dot12_fundamentals/current/repo"
	"github.com/britzc/go-kit_0dot12_fundamentals/current/service"
	"github.com/britzc/go-kit_0dot12_fundamentals/current/transport"
	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)

	fmt.Println("Repository: In progress")

	productRepo, _ := repo.NewProductRepo("products.csv", "partners.csv")

	fmt.Println("Repository: Ready")

	fmt.Println("Endpoints and handlers: In progress")

	var svc service.PricingService
	svc = service.NewPricingService(productRepo)

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
