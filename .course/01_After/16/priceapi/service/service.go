package service

import "context"

type PricingService interface {
	GetRetailTotal(ctx context.Context, code string, qty int) (total float64, err error)
	GetWholesaleTotal(ctx context.Context, partner, code string, qty int) (total float64, err error)
}
