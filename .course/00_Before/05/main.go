package main

import (
	"flag"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	productRepo := repos.NewProductRepo()
	pricingService := services.NewPricingService(productRepo)

	totalRetailPriceHandler := httptransport.NewServer(
		endpoints.MakeTotalRetailPriceEndpoint(pricingService),
		endpoints.DecodeTotalRetailPriceRequest,
		endpoints.EncodeResponse,
	)

	http.Handle("/retail", totalRetailPriceHandler)

	http.ListenAndServe(*listen, nil)
}
