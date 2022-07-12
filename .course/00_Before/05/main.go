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

	productRepo := NewProductRepo()
	pricingService := NewPricingService(productRepo)

	totalRetailPriceHandler := httptransport.NewServer(
		MakeTotalRetailPriceEndpoint(pricingService),
		DecodeTotalRetailPriceRequest,
		EncodeResponse,
	)

	http.Handle("/retail", totalRetailPriceHandler)

	http.ListenAndServe(*listen, nil)
}
