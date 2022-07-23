package main

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockProductRepo struct{}

func (MockProductRepo) FetchProduct(code string) (retailPrice, wholesalePrice float64, found bool) {
	data := []string{
		"aaa111,12.99,10.99",
		"bbb222,2.90,2.50",
		"ccc333,22.50,21.00",
	}

	for _, line := range data {
		parts := strings.Split(line, ",")
		if parts[0] != code {
			continue
		}

		retailPrice, _ = strconv.ParseFloat(parts[1], 64)
		wholesalePrice, _ = strconv.ParseFloat(parts[2], 64)

		return retailPrice, wholesalePrice, true
	}

	return 0, 0, false
}

func Test_GetRetailTotal(t *testing.T) {
	tests := []struct {
		code  string
		qty   int
		err   error
		total float64
	}{
		{
			code:  "",
			qty:   0,
			err:   ErrInvalidCode,
			total: 0.0,
		},
		{
			code:  "aaa111",
			qty:   0,
			err:   ErrInvalidQty,
			total: 0.0,
		},
		{
			code:  "aaa111",
			qty:   15,
			err:   nil,
			total: 194.85,
		},
		{
			code:  "fff000",
			qty:   10,
			err:   ErrNotFound,
			total: 0.0,
		},
	}

	mockProductRepo := new(MockProductRepo)

	priceService := NewPricingService(mockProductRepo)

	for id, test := range tests {
		total, err := priceService.GetRetailTotal(test.code, test.qty)
		assert.True(t, test.err == err, "~2|Test #%d expected error: %s, not error %s~", id, test.err, err)
		assert.True(t, test.total == total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.total, total)
	}
}

/*
func Test_GetWholesaleTotal(t *testing.T) {
	tests := []struct {
		code  string
		qty   int
		err   error
		total float64
	}{
		{
			code:  "",
			qty:   0,
			err:   ErrInvalidCode,
			total: 0.0,
		},
		{
			code:  "bbb222",
			qty:   0,
			err:   ErrInvalidQty,
			total: 0.0,
		},
		{
			code:  "bbb222",
			qty:   15,
			err:   nil,
			total: 37.50,
		},
		{
			code:  "xxx999",
			qty:   10,
			err:   ErrNotFound,
			total: 0.0,
		},
	}

	mockProductRepo := new(MockProductRepo)

	priceService := NewPricingService(mockProductRepo)

	for id, test := range tests {
		total, err := priceService.GetWholesaleTotal(test.code, test.qty)
		assert.True(t, test.err == err, "~2|Test #%d expected error: %s, not error %s~", id, test.err, err)
		assert.True(t, test.total == total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.total, total)
	}
}
*/
