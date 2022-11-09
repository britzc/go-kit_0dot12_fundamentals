package service

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockLogger struct {
	result []string
}

func (ml *MockLogger) Log(keyvals ...interface{}) (err error) {
	var result []string
	for _, val := range keyvals {
		result = append(result, fmt.Sprint(val))
	}

	ml.result = result

	return nil
}

func (ml *MockLogger) Result() string {
	return strings.Join(ml.result[:], ",")
}

type MockPricingService struct{}

func (MockPricingService) GetRetailTotal(code string, qty int) (total float64, err error) {
	if code == "" {
		return 0.0, ErrInvalidCode
	}
	if qty <= 0 {
		return 0.0, ErrInvalidQty
	}

	data := []string{
		"aaa111,12.99,10.99",
		"bbb222,2.90,2.50",
		"ccc333,22.50,21.00",
	}

	for _, line := range data {
		parts := strings.Split(line, ",")
		if parts[0] == code {
			price, _ := strconv.ParseFloat(parts[1], 64)

			return (price * float64(qty)), nil
		}
	}

	return 0.0, ErrCodeNotFound
}

func (MockPricingService) GetWholesaleTotal(partner string, code string, qty int) (total float64, err error) {
	if partner == "" {
		return 0.0, ErrInvalidPartner
	}
	if code == "" {
		return 0.0, ErrInvalidCode
	}
	if qty <= 0 {
		return 0.0, ErrInvalidQty
	}

	prices := []string{
		"aaa111,12.99",
		"bbb222,2.90",
		"ccc333,22.50",
	}

	price := 0.0
	priceFound := false
	for _, line := range prices {
		parts := strings.Split(line, ",")
		if parts[0] == code {
			priceFound = true
			price, _ = strconv.ParseFloat(parts[1], 64)
		}
	}

	if !priceFound {
		return 0.0, ErrCodeNotFound
	}

	partners := []string{
		"superstore,0.15",
		"joesdiscount,0.05",
	}

	discount := 0.0
	discountFound := false
	for _, line := range partners {
		parts := strings.Split(line, ",")
		if parts[0] == partner {
			discountFound = true
			discount, _ = strconv.ParseFloat(parts[1], 64)
		}
	}

	if !discountFound {
		return 0.0, ErrPartnerNotFound
	}

	saved := (price * discount)
	total = (price - saved) * float64(qty)

	return math.Round(total*100) / 100, nil
}

func Test_Logging_GetRetailTotal(t *testing.T) {
	tests := []struct {
		code string
		qty  int
		msg  string
	}{
		{
			code: "aaa111",
			qty:  15,
			msg:  "method,GetRetailTotal,code,aaa111,quantity,15,total,194.85,error,<nil>,duration",
		},
		{
			code: "fff000",
			qty:  10,
			msg:  "method,GetRetailTotal,code,fff000,quantity,10,total,0,error,Code Not Found,duration",
		},
	}

	logger := new(MockLogger)
	var svc PricingService
	svc = new(MockPricingService)
	svc = NewLoggingMiddleware(logger, svc)

	for id, test := range tests {
		svc.GetRetailTotal(test.code, test.qty)

		actual := logger.Result()

		assert.True(t, strings.HasPrefix(actual, test.msg), "~2|Test #%d logger expected: \"%s\", not: \"%s\"~", id, test.msg, actual)
	}
}

func Test_Logging_GetWholesaleTotal(t *testing.T) {
	tests := []struct {
		partner string
		code    string
		qty     int
		msg     string
	}{
		{
			partner: "superstore",
			code:    "aaa111",
			qty:     15,
			msg:     "method,GetWholesaleTotal,partner,superstore,code,aaa111,quantity,15,total,165.62,error,<nil>,duration",
		},
		{
			partner: "smiles",
			code:    "fff000",
			qty:     10,
			msg:     "method,GetWholesaleTotal,partner,smiles,code,fff000,quantity,10,total,0,error,Code Not Found,duration",
		},
	}

	logger := new(MockLogger)
	var svc PricingService
	svc = new(MockPricingService)
	svc = NewLoggingMiddleware(logger, svc)

	for id, test := range tests {
		svc.GetWholesaleTotal(test.partner, test.code, test.qty)

		actual := logger.Result()

		assert.True(t, strings.HasPrefix(actual, test.msg), "~2|Test #%d logging expected: \"%s\", not: \"%s\"~", id, test.msg, actual)

	}
}
