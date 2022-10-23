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

	rtr := mux.NewRouter().StrictSlash(true)

	totalRetailPriceHandler := MakeTotalRetailPriceHttpHandler(pricingService)
	rtr.Handle("/retail", totalRetailPriceHandler).Methods(http.MethodPost)
<<<<<<< HEAD
=======

	// totalWholesalePriceHandler := MakeTotalWholesalePriceHttpHandler(pricingService)
	// rtr.Handle("/wholesale", totalWholesalePriceHandler).Methods("POST")
>>>>>>> c3d1822cb45983a78c3b60edc34931d376136d1f

	fmt.Println("Endpoints and handlers: Ready")

	fmt.Printf("Hosting on %s\n", *listen)

	http.ListenAndServe(*listen, rtr)
}
