package service

type PricingService interface {
	GetRetailTotal(code string, qty int) (total float64, err error)
	GetWholesaleTotal(partner string, code string, qty int) (total float64, err error)
}
