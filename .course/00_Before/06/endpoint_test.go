package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/stretchr/testify/assert"
)

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

	return 0.0, ErrNotFound
}

func Test_MakeTotalRetailPriceEndpoint(t *testing.T) {
	tests := []struct {
		request  totalRetailPriceRequest
		response totalRetailPriceResponse
	}{
		{
			request:  totalRetailPriceRequest{Code: "", Qty: 0},
			response: totalRetailPriceResponse{Err: "Invalid Code Requested"},
		},
		{
			request:  totalRetailPriceRequest{Code: "aaa111", Qty: 0},
			response: totalRetailPriceResponse{Err: "Invalid Quantity Requested"},
		},
		{
			request:  totalRetailPriceRequest{Code: "aaa111", Qty: 15},
			response: totalRetailPriceResponse{Total: 194.85},
		},
		{
			request:  totalRetailPriceRequest{Code: "fff000", Qty: 10},
			response: totalRetailPriceResponse{Err: "Product Not Found"},
		},
	}

	mockPricingService := new(MockPricingService)

	totalRetailPriceHandler := httptransport.NewServer(
		MakeTotalRetailPriceEndpoint(mockPricingService),
		decodeTotalRetailPriceRequest,
		encodeResponse,
	)

	server := httptest.NewServer(totalRetailPriceHandler)
	defer server.Close()

	for id, test := range tests {
		postBody, _ := json.Marshal(test.request)

		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(server.URL, "application/json", responseBody)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}

		var actualResponse totalRetailPriceResponse
		json.NewDecoder(resp.Body).Decode(&actualResponse)

		assert.True(t, test.response.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, test.response.Err, actualResponse.Err)
		assert.True(t, test.response.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.response.Total, actualResponse.Total)
	}
}
