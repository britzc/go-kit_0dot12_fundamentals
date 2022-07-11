package main

import (
	"errors"
)

type ProductRepo interface {
	FetchProductDetails(code string) (retailPrice, wholesalePrice float64, found bool)
}

var ErrInvalidCode = errors.New("Invalid Code Requested")
var ErrInvalidQty = errors.New("Invalid Quantity Requested")
var ErrNotFound = errors.New("Product Not Found")

type PricingService struct {
	productRepo ProductRepo
}

func NewPricingService(pr ProductRepo) (ps *PricingService) {
	ps = &PricingService{
		productRepo: pr,
	}

	return ps
}

func (ps PricingService) GetTotalRetailPrice(code string, qty int) (total float64, err error) {
	if code == "" {
		return 0.0, ErrInvalidCode
	}

	if qty <= 0 {
		return 0.0, ErrInvalidQty
	}

	retailPrice, _, found := ps.productRepo.FetchProductDetails(code)
	if !found {
		return 0.0, ErrNotFound
	}

	return retailPrice * float64(qty), nil
}
