package repos

type ProductRepo struct {
}

func NewProductRepo() (pr *ProductRepo) {
	pr = &ProductRepo{}

	return pr
}

func (ProductRepo) FetchProductDetails(code string) (retailPrice, wholesalePrice float64, found bool) {
	return 0.0, 0.0, true
}
