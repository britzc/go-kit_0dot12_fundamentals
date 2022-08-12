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

func Test_LoggingMiddleware(t *testing.T) {
	tests := []struct {
		method   string
		expected string
	}{
		{
			method:   "function15",
			expected: "method,function15,msg,called endpoint",
		},
		{
			method:   "function55",
			expected: "method,function55,msg,called endpoint",
		},
	}

	endpoint := func(_ context.Context, request interface{}) (interface{}, error) {
		return request, nil
	}

	logger := &MockLogger{}

	for id, test := range tests {
		lmw := LoggingMiddleware(log.With(logger, "method", test.method))(endpoint)

		lmw(context.Background(), "Go kit!")

		actual := logger.Result()

		assert.True(t, test.expected == actual, "~2|Test #%d logger expected: \"%s\", not: \"%s\"~", id, test.expected, actual)
	}
}
