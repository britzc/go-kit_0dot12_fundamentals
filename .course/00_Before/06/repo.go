package main

type productRepo struct {
}

func NewProductRepo() (pr *productRepo) {
	pr = &productRepo{}

	return pr
}

func (productRepo) FetchProduct(code string) (retailPrice, wholesalePrice float64, found bool) {
	return 0.0, 0.0, true
}
