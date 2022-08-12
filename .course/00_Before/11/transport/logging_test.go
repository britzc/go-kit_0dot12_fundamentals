package transport

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-kit/log"
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

func Test_LogTotalRetailPriceEndpoint(t *testing.T) {
	tests := []struct {
		service  string
		request  TotalRetailPriceRequest
		expected string
	}{
		{
			service:  "endpointRetailTest",
			request:  TotalRetailPriceRequest{"aaa11", 10},
			expected: "service,endpointRetailTest,endpoint,TotalRetailPriceEndpoint,msg,Called endpoint",
		},
		{
			service:  "endpointRetailTest",
			request:  TotalRetailPriceRequest{"bbb11", 20},
			expected: "service,endpointRetailTest,endpoint,TotalRetailPriceEndpoint,msg,Called endpoint",
		},
	}

	endpoint := func(_ context.Context, request interface{}) (interface{}, error) {
		return request, nil
	}

	logger := &MockLogger{}

	for id, test := range tests {
		lmw := LogTotalRetailPriceEndpoint(log.With(logger, "service", test.service))(endpoint)

		lmw(context.Background(), test.request)

		actual := logger.Result()

		assert.True(t, test.expected == actual, "~2|Test #%d logger expected: \"%s\", not: \"%s\"~", id, test.expected, actual)
	}
}

func Test_LogTotalWholesalePriceEndpoint(t *testing.T) {
	tests := []struct {
		service  string
		request  TotalWholesalePriceRequest
		expected string
	}{
		{
			service:  "endpointWholesaleTest",
			request:  TotalWholesalePriceRequest{"testpartner", "aaa11", 10},
			expected: "service,endpointWholesaleTest,endpoint,TotalWholesalePriceEndpoint,msg,Called endpoint",
		},
		{
			service:  "endpointWholesaleTest",
			request:  TotalWholesalePriceRequest{"testpartner", "bbb11", 20},
			expected: "service,endpointWholesaleTest,endpoint,TotalWholesalePriceEndpoint,msg,Called endpoint",
		},
	}

	endpoint := func(_ context.Context, request interface{}) (interface{}, error) {
		return request, nil
	}

	logger := &MockLogger{}

	for id, test := range tests {
		lmw := LogTotalWholesalePriceEndpoint(log.With(logger, "service", test.service))(endpoint)

		lmw(context.Background(), test.request)

		actual := logger.Result()

		assert.True(t, test.expected == actual, "~2|Test #%d logger expected: \"%s\", not: \"%s\"~", id, test.expected, actual)
	}
}
