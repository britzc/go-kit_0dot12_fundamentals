package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

type price struct {
	code      string
	retail    float64
	wholesale float64
}

type productRepo struct {
	prices map[string]*price
}

func NewProductRepo(path string) (pr *productRepo, err error) {
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

	prices := make(map[string]*price, 0)
	for _, record := range records {
		retail, _ := strconv.ParseFloat(record[1], 64)
		wholesale, _ := strconv.ParseFloat(record[2], 64)

		p := &price{
			code:      record[0],
			retail:    retail,
			wholesale: wholesale,
		}

		prices[p.code] = p
	}

	pr = &productRepo{
		prices: prices,
	}

	return pr, nil
}

func (pr *productRepo) FetchProduct(code string) (retailPrice, wholesalePrice float64, found bool) {
	p, ok := pr.prices[code]
	if !ok {
		return 0.0, 0.0, false
	}

	return p.retail, p.wholesale, true
}
