package main

import (
	"errors"
)

type ProductRepo interface {
	FetchProduct(code string) (retailPrice, wholesalePrice float64, found bool)
}

var ErrInvalidCode = errors.New("Invalid Code Requested")
var ErrInvalidQty = errors.New("Invalid Quantity Requested")
var ErrNotFound = errors.New("Product Not Found")

type service struct {
	repo ProductRepo
}

func NewPricingService(pr ProductRepo) (ps *service) {
	ps = &service{
		repo: pr,
	}

	return ps
}

func (ps *service) GetRetailTotal(code string, qty int) (total float64, err error) {
	if code == "" {
		return 0.0, ErrInvalidCode
	}
	if qty <= 0 {
		return 0.0, ErrInvalidQty
	}

	price, _, found := ps.repo.FetchProduct(code)
	if !found {
		return 0.0, ErrNotFound
	}

	return price * float64(qty), nil
}
