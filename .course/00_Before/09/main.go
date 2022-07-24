package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	fmt.Println("Repository: In progress")

	productRepo, err := NewProductRepo("products.csv", "partners.csv")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Repository: Ready")

	fmt.Println("Endpoints and handlers: In progress")

	pricingService := NewPricingService(productRepo)
	totalRetailPriceHandler := MakeTotalRetailPriceHttpHandler(pricingService)
	totalWholesalePriceHandler := MakeTotalWholesalePriceHttpHandler(pricingService)

	rtr := mux.NewRouter().StrictSlash(true)
	rtr.Handle("/retail", totalRetailPriceHandler).Methods("POST")
	rtr.Handle("/wholesale", totalWholesalePriceHandler).Methods("POST")

	fmt.Println("Endpoints and handlers: Ready")

	fmt.Printf("Hosting on %s\n", *listen)

	http.ListenAndServe(*listen, rtr)
}
