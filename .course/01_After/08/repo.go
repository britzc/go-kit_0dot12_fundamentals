package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

type product struct {
	code  string
	price float64
}

type partner struct {
	name     string
	discount float64
}

type productRepo struct {
	products map[string]*product
	partners map[string]*partner
}

func NewProductRepo(productsPath string, partnersPath string) (pr *productRepo, err error) {
	productRecords, err := readCSV(productsPath)
	if err != nil {
		return nil, err
	}

	products := make(map[string]*product, 0)
	for _, record := range productRecords {
		price, _ := strconv.ParseFloat(record[1], 64)

		p := &product{
			code:  record[0],
			price: price,
		}

		products[p.code] = p
	}

	partnerRecords, err := readCSV(partnersPath)
	if err != nil {
		return nil, err
	}

	partners := make(map[string]*partner, 0)
	for _, record := range partnerRecords {
		discount, _ := strconv.ParseFloat(record[1], 64)

		p := &partner{
			name:     record[0],
			discount: discount,
		}

		partners[p.name] = p
	}

	pr = &productRepo{
		products: products,
		partners: partners,
	}

	return pr, nil
}

func readCSV(path string) (lines [][]string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (pr *productRepo) FetchPrice(code string) (price float64, found bool) {
	p, ok := pr.products[code]
	if !ok {
		return 0.0, false
	}

	return p.price, true
}

func (pr *productRepo) FetchDiscount(partner string) (discount float64, found bool) {
	p, ok := pr.partners[partner]
	if !ok {
		return 0.0, false
	}

	return p.discount, true
}
