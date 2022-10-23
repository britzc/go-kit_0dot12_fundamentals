package main

import (
	"errors"
	"math"
)

type ProductRepo interface {
	FetchPrice(code string) (price float64, found bool)
	FetchDiscount(partner string) (discount float64, found bool)
}

var ErrInvalidPartner = errors.New("Invalid Partner Requested")
var ErrPartnerNotFound = errors.New("Partner Not Found")
var ErrInvalidCode = errors.New("Invalid Code Requested")
var ErrCodeNotFound = errors.New("Code Not Found")
var ErrInvalidQty = errors.New("Invalid Quantity Requested")

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

	price, found := ps.repo.FetchPrice(code)
	if !found {
		return 0.0, ErrCodeNotFound
	}

	total = price * float64(qty)

	return math.Round(total*100) / 100, nil
}

/*
Discount Calculation:
saved = (price x discount)
total = (price - saved) x quantity
*/
