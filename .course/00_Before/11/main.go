package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)

	fmt.Println("Repository: In progress")

	productRepo, _ := NewProductRepo("products.csv", "partners.csv")

	fmt.Println("Repository: Ready")

	fmt.Println("Endpoints and handlers: In progress")

	var svc PricingService
	svc = NewPricingService(productRepo)
	svc = NewLoggingMiddleware(logger, svc)

	rtr := mux.NewRouter().StrictSlash(true)

	totalRetailPriceHandler := MakeTotalRetailPriceHttpHandler(svc)

	rtr.Handle("/retail", totalRetailPriceHandler).Methods(http.MethodPost)

	totalWholesalePriceHandler := MakeTotalWholesalePriceHttpHandler(svc)
	rtr.Handle("/wholesale", totalWholesalePriceHandler).Methods(http.MethodPost)

	fmt.Println("Endpoints and handlers: Ready")

	fmt.Printf("Hosting on %s\n", *listen)

	http.ListenAndServe(*listen, rtr)
}
