package main

import (
	"fmt"
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

/*
type MockPricingService struct{}

func (MockPricingService) GetRetailTotal(code string, qty int) (total float64, err error) {
	return 100.99, nil
}

func (MockPricingService) GetWholesaleTotal(partner, code string, qty int) (total float64, err error) {
	return 99.99, nil
}
*/

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
